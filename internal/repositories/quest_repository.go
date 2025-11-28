package repositories

import (
	"context"
	"errors"
	"time"

	"BecomeOverMan/internal/models"

	"github.com/jmoiron/sqlx"
)

type QuestRepository struct {
	db *sqlx.DB
}

func NewQuestRepository(db *sqlx.DB) *QuestRepository {
	return &QuestRepository{db: db}
}

func (r *QuestRepository) SaveQuestToDB(quest *models.Quest, tasks []models.Task) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Вставляем квест
	var questID int
	err = tx.QueryRow(`
		INSERT INTO quests (
			title, description, category, rarity, difficulty, price, tasks_count,
			reward_xp, reward_coin, time_limit_hours, is_sequential
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`,
		quest.Title, quest.Description, quest.Category, quest.Rarity,
		quest.Difficulty, quest.Price, quest.TasksCount, quest.RewardXP,
		quest.RewardCoin, quest.TimeLimitHours, true,
	).Scan(&questID)
	if err != nil {
		return 0, err
	}

	// Вставляем задачи
	for _, task := range tasks {
		var taskID int
		err = tx.QueryRow(`
			INSERT INTO tasks (
				title, description, difficulty, rarity, category, 
				base_xp_reward, base_coin_reward
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`,
			task.Title, task.Description, task.Difficulty, task.Rarity,
			task.Category, task.BaseXpReward, task.BaseCoinReward,
		).Scan(&taskID)
		if err != nil {
			return 0, err
		}

		// Связываем задачу с квестом
		_, err = tx.Exec(`
			INSERT INTO quest_tasks (quest_id, task_id, task_order)
			VALUES ($1, $2, $3)
		`, questID, taskID, task.TaskOrder)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return questID, nil
}

// GetQuestDetails возвращает детали квеста со всеми задачами
func (r *QuestRepository) GetQuestDetails(ctx context.Context, questID, userID int) (*models.Quest, error) {
	// Получаем основную информацию о квесте
	var quest models.Quest
	err := r.db.GetContext(ctx, &quest, `
        SELECT * FROM quests WHERE id = $1
    `, questID)
	if err != nil {
		return nil, err
	}

	// Получаем все задачи этого квеста
	var tasks []models.Task
	err = r.db.SelectContext(ctx, &tasks, `
        SELECT t.*, qt.task_order
        FROM tasks t
        INNER JOIN quest_tasks qt ON t.id = qt.task_id
        WHERE qt.quest_id = $1
        ORDER BY qt.task_order ASC
    `, questID)
	if err != nil {
		return nil, err
	}

	quest.Tasks = tasks

	return &quest, nil
}

// GetAvailableQuests возвращает квесты, доступные для пользователя
func (r *QuestRepository) GetAvailableQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	var quests []models.Quest

	// Получаем уровень пользователя
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT q.* FROM quests q 
		WHERE q.difficulty <= $1 + 1 AND q.price <= $2 AND NOT EXISTS (
			SELECT 1 FROM user_quests uq
			WHERE uq.quest_id = q.id AND uq.user_id = $3
		)
	`

	err = r.db.SelectContext(ctx, &quests, query, user.Level, user.CoinBalance, userID)
	if err != nil {
		return nil, err
	}

	return quests, nil
}

func (r *QuestRepository) GetQuestShop(ctx context.Context, userID int) ([]models.Quest, error) {
	var quests []models.Quest

	query := `
	SELECT q.* FROM quests q
	WHERE NOT EXISTS (
		SELECT 1 FROM user_quests uq
		WHERE uq.quest_id = q.id AND uq.user_id = $1
	)`
	// + TODO: conditions_json нужно проверить

	// Получаем все квесты, что у нас не куплены и не были пройдены
	if err := r.db.SelectContext(ctx, &quests, query, userID); err != nil {
		return nil, err
	}

	return quests, nil
}

func (r *QuestRepository) GetMyActiveQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	var quests []models.Quest

	query := `
		SELECT q.id, q.title, q.description, q.category, q.rarity, q.difficulty, 
		       q.price, q.tasks_count, q.conditions_json, q.bonus_json, q.is_sequential,
		       q.reward_xp, q.reward_coin, q.time_limit_hours
		FROM quests q
		INNER JOIN user_quests uq ON q.id = uq.quest_id
		WHERE uq.user_id = $1 AND uq.status IN ('purchased', 'started')`

	if err := r.db.SelectContext(ctx, &quests, query, userID); err != nil {
		return nil, err
	}

	return quests, nil
}

func (r *QuestRepository) GetMyCompletedQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	var quests []models.Quest

	query := `
	SELECT q.id, q.title, q.description, q.category, q.rarity, q.difficulty, 
		   q.price, q.tasks_count, q.conditions_json, q.bonus_json, q.is_sequential,
		   q.reward_xp, q.reward_coin, q.time_limit_hours
	FROM quests q
	INNER JOIN user_quests uq ON q.id = uq.quest_id
	WHERE uq.user_id = $1 AND uq.status = 'completed'`

	if err := r.db.SelectContext(ctx, &quests, query, userID); err != nil {
		return nil, err
	}

	return quests, nil
}

// PurchaseQuest покупает квест для пользователя
func (r *QuestRepository) PurchaseQuest(ctx context.Context, userID, questID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что квест существует
	var quest models.Quest
	err = tx.GetContext(ctx, &quest, "SELECT * FROM quests WHERE id = $1", questID)
	if err != nil {
		return err
	}

	// Проверяем что такой квест у нас не куплен и не был пройден
	var alreadyUsed bool
	err = tx.GetContext(ctx, &alreadyUsed, `
		SELECT EXISTS(
			SELECT 1 FROM user_quests WHERE user_id = $1 AND quest_id = $2
		)`, userID, questID,
	)

	if err != nil {
		return err
	}

	if alreadyUsed {
		return errors.New("quest already purchased or completed")
	}

	// Проверяем баланс пользователя
	var balance int
	err = tx.GetContext(ctx, &balance, "SELECT coin_balance FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	if balance < quest.Price {
		return errors.New("not enough currency")
	}

	// Списываем валюту
	_, err = tx.ExecContext(ctx,
		"UPDATE users SET coin_balance = coin_balance - $1 WHERE id = $2",
		quest.Price, userID)
	if err != nil {
		return err
	}

	// Добавляем квест пользователю
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_quests 
        (user_id, quest_id, status, started_at, expires_at)
        VALUES ($1, $2, 'purchased', NULL, NULL)`,
		userID, questID)
	if err != nil {
		return err
	}

	// Записываем транзакцию
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_coin_transactions 
        (user_id, amount, transaction_type, reference_type, reference_id, description)
        VALUES ($1, $2, 'spent', 'quest', $3, 'Purchased quest: ' || $4)`,
		userID, -quest.Price, quest.ID, quest.Title)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// StartQuest начинает выполнение квеста
func (r *QuestRepository) StartQuest(ctx context.Context, userID, questID int) error {
	// Устанавливаем время начала и завершения
	var timeLimitHours int
	err := r.db.GetContext(ctx, &timeLimitHours,
		"SELECT time_limit_hours FROM quests WHERE id = $1", questID)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Duration(timeLimitHours) * time.Hour)

	_, err = r.db.ExecContext(ctx, `
        UPDATE user_quests 
        SET status = 'started', started_at = NOW(), expires_at = $1
        WHERE user_id = $2 AND quest_id = $3 AND status = 'purchased'`,
		expiresAt, userID, questID)

	return err
}

// CompleteTask отмечает выполнение задачи
func (r *QuestRepository) CompleteTask(ctx context.Context, userID, questID, taskID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что такой квест и задача в нем - существуют и еще не выполнены
	var exists bool
	err = tx.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1 FROM user_quests WHERE user_id = $1 AND quest_id = $2 AND status = 'started'
		) AND EXISTS (
			SELECT 1 FROM quest_tasks WHERE quest_id = $2 AND task_id = $3
		)`, userID, questID, taskID,
	)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("quest or task not found or already completed")
	}

	// Добавляем задачу в буфер выполненных
	_, err = tx.ExecContext(ctx, `
        INSERT INTO user_completed_tasks 
        (user_id, task_id, completed_at, xp_gained, coin_gained)
        SELECT $1, $2, NOW(), t.base_xp_reward, t.base_coin_reward
        FROM tasks t WHERE t.id = $2`,
		userID, taskID)
	if err != nil {
		return err
	}

	// Увеличиваем счетчик выполненных задач
	_, err = tx.ExecContext(ctx, `
        UPDATE user_quests 
        SET tasks_done = tasks_done + 1
        WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CompleteQuest завершает квест
func (r *QuestRepository) CompleteQuest(ctx context.Context, userID, questID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что все задачи выполнены
	var totalTasks, completedTasks int
	err = tx.GetContext(ctx, &totalTasks, `
        SELECT COUNT(*) FROM quest_tasks WHERE quest_id = $1`, questID)
	if err != nil {
		return err
	}

	err = tx.GetContext(ctx, &completedTasks, `
        SELECT tasks_done FROM user_quests 
        WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)
	if err != nil {
		return err
	}

	if completedTasks < totalTasks {
		return errors.New("not all tasks completed")
	}

	// Проверяем, является ли квест совместным
	var sharedQuestID int
	var friendID int
	err = tx.QueryRowContext(ctx,
		`SELECT id, 
            CASE 
                WHEN user1_id = $1 THEN user2_id 
                ELSE user1_id 
            END as friend_id
        FROM shared_quests 
        WHERE quest_id = $2 AND (user1_id = $1 OR user2_id = $1) AND status = 'active'`,
		userID, questID).Scan(&sharedQuestID, &friendID)

	isSharedQuest := (err == nil)

	if isSharedQuest {
		// Это совместный квест - проверяем прогресс друга
		var friendTasksDone int
		err = tx.GetContext(ctx, &friendTasksDone, `
            SELECT tasks_done FROM user_quests 
            WHERE user_id = $1 AND quest_id = $2`,
			friendID, questID)
		if err != nil {
			return errors.New("friend quest progress not found")
		}

		if friendTasksDone < totalTasks {
			// Друг еще не завершил - НИКТО не получает награду
			return errors.New("friend has not completed all tasks yet")
		}

		// Оба завершили все задачи! Награждаем обоих и завершаем совместный квест
		if err := r.completeQuestForUsers(tx, ctx, []int{userID, friendID}, questID); err != nil {
			return err
		}

		// Отмечаем совместный квест как завершенный
		_, err = tx.ExecContext(ctx, `
            UPDATE shared_quests SET status = 'completed' 
            WHERE id = $1`, sharedQuestID)
		if err != nil {
			return err
		}
	} else {
		// Обычный квест - награждаем только текущего пользователя
		if err := r.completeQuestForUsers(tx, ctx, []int{userID}, questID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// completeQuestForUsers - упрощенная версия (если сложно с динамическими IN clause)
func (r *QuestRepository) completeQuestForUsers(tx *sqlx.Tx, ctx context.Context, userIDs []int, questID int) error {
	// Получаем награду за квест
	var rewardXP, rewardCoin int
	err := tx.QueryRowContext(ctx, `
        SELECT reward_xp, reward_coin FROM quests WHERE id = $1`, questID).
		Scan(&rewardXP, &rewardCoin)
	if err != nil {
		return err
	}

	// Для каждого пользователя выполняем операции
	for _, userID := range userIDs {
		// Начисляем награду
		_, err = tx.ExecContext(ctx, `
            UPDATE users 
            SET xp_points = xp_points + $1,
                coin_balance = coin_balance + $2
            WHERE id = $3`,
			rewardXP, rewardCoin, userID)
		if err != nil {
			return err
		}

		// Отмечаем квест как завершенный
		_, err = tx.ExecContext(ctx, `
            UPDATE user_quests 
            SET status = 'completed', completed_at = NOW(),
                xp_gained = $1, coin_gained = $2
            WHERE user_id = $3 AND quest_id = $4`,
			rewardXP, rewardCoin, userID, questID)
		if err != nil {
			return err
		}

		// Подтверждаем задачи
		_, err = tx.ExecContext(ctx, `
            UPDATE user_completed_tasks 
            SET is_confirmed = true
            WHERE user_id = $1 AND task_id IN (
                SELECT task_id FROM quest_tasks WHERE quest_id = $2
            )`,
			userID, questID)
		if err != nil {
			return err
		}
	}

	return nil
}
