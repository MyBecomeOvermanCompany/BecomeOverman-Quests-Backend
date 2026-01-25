package repositories

import (
	"BecomeOverMan/internal/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// QuestTreeRepository - репозиторий для работы с деревом квестов и пассивными баффами
type QuestTreeRepository struct {
	db *sqlx.DB
}

func NewQuestTreeRepository(db *sqlx.DB) *QuestTreeRepository {
	return &QuestTreeRepository{db: db}
}

// GetDevelopmentBranches - получить все ветки развития
func (r *QuestTreeRepository) GetDevelopmentBranches(ctx context.Context) ([]models.DevelopmentBranch, error) {
	var branches []models.DevelopmentBranch
	query := `SELECT id, name, display_name, description, parent_branch_id, level, icon, color 
	          FROM development_branches 
	          ORDER BY level, name`
	err := r.db.SelectContext(ctx, &branches, query)
	return branches, err
}

// GetQuestPrerequisites - получить зависимости квеста
func (r *QuestTreeRepository) GetQuestPrerequisites(ctx context.Context, questID int) ([]models.QuestPrerequisite, error) {
	var prerequisites []models.QuestPrerequisite
	query := `SELECT id, quest_id, prerequisite_quest_id, required_count 
	          FROM quest_prerequisites 
	          WHERE quest_id = $1`
	err := r.db.SelectContext(ctx, &prerequisites, query, questID)
	return prerequisites, err
}

// GetQuestBranches - получить ветки квеста
func (r *QuestTreeRepository) GetQuestBranches(ctx context.Context, questID int) ([]models.DevelopmentBranch, error) {
	var branches []models.DevelopmentBranch
	query := `SELECT b.id, b.name, b.display_name, b.description, b.parent_branch_id, b.level, b.icon, b.color
	          FROM development_branches b
	          INNER JOIN quest_branches qb ON b.id = qb.branch_id
	          WHERE qb.quest_id = $1`
	err := r.db.SelectContext(ctx, &branches, query, questID)
	return branches, err
}

// CheckQuestUnlocked - проверить, разблокирован ли квест для пользователя
func (r *QuestTreeRepository) CheckQuestUnlocked(ctx context.Context, userID, questID int) (bool, int, int, error) {
	// Получаем prerequisites
	prerequisites, err := r.GetQuestPrerequisites(ctx, questID)
	if err != nil {
		return false, 0, 0, err
	}

	// Если нет зависимостей - квест доступен
	if len(prerequisites) == 0 {
		return true, 0, 0, nil
	}

	// Группируем по required_count
	prerequisiteGroups := make(map[int][]int)
	for _, p := range prerequisites {
		prerequisiteGroups[p.RequiredCount] = append(prerequisiteGroups[p.RequiredCount], p.PrerequisiteQuestID)
	}

	// Проверяем каждую группу
	for requiredCount, questIDs := range prerequisiteGroups {
		// Проверяем сколько из этих квестов завершено
		var completedCount int
		placeholders := ""
		args := []interface{}{userID}
		for i, qID := range questIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += fmt.Sprintf("$%d", len(args)+1)
			args = append(args, qID)
		}

		query := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM user_quests 
			WHERE user_id = $1 
			AND quest_id IN (%s) 
			AND status = 'completed'`, placeholders)

		err := r.db.GetContext(ctx, &completedCount, query, args...)
		if err != nil {
			return false, 0, 0, err
		}

		// Если выполнено достаточно - квест разблокирован
		if completedCount >= requiredCount {
			return true, completedCount, requiredCount, nil
		}
	}

	// Если ни одна группа не выполнена - квест заблокирован
	// Возвращаем максимальное количество выполненных
	maxCompleted := 0
	maxRequired := 0
	for requiredCount, questIDs := range prerequisiteGroups {
		var completedCount int
		placeholders := ""
		args := []interface{}{userID}
		for i, qID := range questIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += fmt.Sprintf("$%d", len(args)+1)
			args = append(args, qID)
		}

		query := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM user_quests 
			WHERE user_id = $1 
			AND quest_id IN (%s) 
			AND status = 'completed'`, placeholders)

		r.db.GetContext(ctx, &completedCount, query, args...)
		if completedCount > maxCompleted {
			maxCompleted = completedCount
			maxRequired = requiredCount
		}
	}

	return false, maxCompleted, maxRequired, nil
}

// AddPassiveBuff - добавить пассивный бафф пользователю
func (r *QuestTreeRepository) AddPassiveBuff(ctx context.Context, userID, questID int, buffType string, buffData json.RawMessage) error {
	query := `INSERT INTO user_passive_buffs (user_id, quest_id, buff_type, buff_data, is_active)
	          VALUES ($1, $2, $3, $4, TRUE)
	          ON CONFLICT (user_id, quest_id) DO UPDATE 
	          SET buff_type = EXCLUDED.buff_type, 
	              buff_data = EXCLUDED.buff_data, 
	              is_active = TRUE`
	_, err := r.db.ExecContext(ctx, query, userID, questID, buffType, buffData)
	return err
}

// GetUserPassiveBuffs - получить все активные баффы пользователя
func (r *QuestTreeRepository) GetUserPassiveBuffs(ctx context.Context, userID int) ([]models.PassiveBuff, error) {
	var buffs []models.PassiveBuff
	query := `SELECT id, user_id, quest_id, buff_type, buff_data, is_active, created_at
	          FROM user_passive_buffs
	          WHERE user_id = $1 AND is_active = TRUE`
	err := r.db.SelectContext(ctx, &buffs, query, userID)
	return buffs, err
}

// GetRewardMultiplier - получить множитель награды для категории
func (r *QuestTreeRepository) GetRewardMultiplier(ctx context.Context, userID int, category string) (float64, error) {
	buffs, err := r.GetUserPassiveBuffs(ctx, userID)
	if err != nil {
		return 1.0, err
	}

	multiplier := 1.0
	for _, buff := range buffs {
		if buff.BuffType == "reward_multiplier" && buff.BuffData != nil {
			var buffData models.RewardMultiplierBuff
			if err := json.Unmarshal(*buff.BuffData, &buffData); err == nil {
				if buffData.Category == category {
					multiplier *= buffData.Multiplier
				}
			}
		}
	}

	return multiplier, nil
}

// LinkQuestToBranch - связать квест с веткой
func (r *QuestTreeRepository) LinkQuestToBranch(ctx context.Context, questID, branchID int, weight float64) error {
	query := `INSERT INTO quest_branches (quest_id, branch_id, weight)
	          VALUES ($1, $2, $3)
	          ON CONFLICT (quest_id, branch_id) DO UPDATE 
	          SET weight = EXCLUDED.weight`
	_, err := r.db.ExecContext(ctx, query, questID, branchID, weight)
	return err
}

// AddQuestPrerequisite - добавить зависимость квеста
func (r *QuestTreeRepository) AddQuestPrerequisite(ctx context.Context, questID, prerequisiteQuestID, requiredCount int) error {
	query := `INSERT INTO quest_prerequisites (quest_id, prerequisite_quest_id, required_count)
	          VALUES ($1, $2, $3)
	          ON CONFLICT (quest_id, prerequisite_quest_id) DO UPDATE 
	          SET required_count = EXCLUDED.required_count`
	_, err := r.db.ExecContext(ctx, query, questID, prerequisiteQuestID, requiredCount)
	return err
}

// GetAvailableQuestsWithTree - получить доступные квесты с информацией о дереве
func (r *QuestTreeRepository) GetAvailableQuestsWithTree(ctx context.Context, userID int) ([]models.QuestWithPrerequisites, error) {
	// Получаем все квесты
	var allQuests []models.Quest
	query := `SELECT id, title, description, category, rarity, difficulty, price, tasks_count, 
	          conditions_json, bonus_json, is_sequential, reward_xp, reward_coin, time_limit_hours
	          FROM quests`
	err := r.db.SelectContext(ctx, &allQuests, query)
	if err != nil {
		return nil, err
	}

	result := make([]models.QuestWithPrerequisites, 0, len(allQuests))
	for _, quest := range allQuests {
		// Получаем prerequisites
		prerequisites, _ := r.GetQuestPrerequisites(ctx, quest.ID)
		
		// Получаем ветки
		branches, _ := r.GetQuestBranches(ctx, quest.ID)
		
		// Проверяем разблокирован ли
		isUnlocked, unlockedCount, requiredCount, _ := r.CheckQuestUnlocked(ctx, userID, quest.ID)

		result = append(result, models.QuestWithPrerequisites{
			Quest:          quest,
			Prerequisites:  prerequisites,
			Branches:       branches,
			IsUnlocked:     isUnlocked,
			UnlockedCount:  unlockedCount,
			RequiredCount:  requiredCount,
		})
	}

	return result, nil
}
