package services

import (
	"BecomeOverMan/internal/models"
	"context"
	"fmt"
)

// GetUserProgress - получает прогресс пользователя для мотивации
func (s *QuestService) GetUserProgress(ctx context.Context, userID int) (map[string]interface{}, error) {
	user, err := s.userRepo.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %w", err)
	}

	activeQuests, err := s.GetMyActiveQuests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting active quests: %w", err)
	}

	completedQuests, err := s.GetMyCompletedQuests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting completed quests: %w", err)
	}

	// Получаем streak (можно добавить в User модель)
	streak := 0 // TODO: получить из user.streaks

	return map[string]interface{}{
		"level":           user.Level,
		"xp":              user.XpPoints,
		"completed_quests": len(completedQuests),
		"active_quests":   len(activeQuests),
		"streak":          streak,
	}, nil
}

// GetUserProgressData - получает детальные данные для анализа
func (s *QuestService) GetUserProgressData(ctx context.Context, userID int) (map[string]interface{}, error) {
	user, err := s.userRepo.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %w", err)
	}

	activeQuests, err := s.GetMyActiveQuests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting active quests: %w", err)
	}

	completedQuests, err := s.GetMyCompletedQuests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting completed quests: %w", err)
	}

	allQuests, err := s.GetMyAllQuestsWithDetails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting all quests: %w", err)
	}

	// Анализ категорий
	categoryStats := make(map[string]int)
	for _, q := range completedQuests {
		categoryStats[q.Category]++
	}

	return map[string]interface{}{
		"user": map[string]interface{}{
			"level":      user.Level,
			"xp":         user.XpPoints,
			"coins":      user.CoinBalance,
			"created_at": user.CreatedAt,
		},
		"quests": map[string]interface{}{
			"active":    len(activeQuests),
			"completed": len(completedQuests),
			"total":     len(allQuests),
		},
		"categories": categoryStats,
		"active_quests_details": activeQuests,
		"completed_quests_details": completedQuests,
	}, nil
}

// GetUserHistory - получает историю пользователя для персонализации
func (s *QuestService) GetUserHistory(ctx context.Context, userID int) (map[string]interface{}, error) {
	user, err := s.userRepo.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %w", err)
	}

	completedQuests, err := s.GetMyCompletedQuests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting completed quests: %w", err)
	}

	activeQuests, err := s.GetMyActiveQuests(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting active quests: %w", err)
	}

	// Получаем детали завершенных квестов
	completedDetails := make([]map[string]interface{}, 0, len(completedQuests))
	for _, q := range completedQuests {
		completedDetails = append(completedDetails, map[string]interface{}{
			"id":          q.ID,
			"title":       q.Title,
			"category":    q.Category,
			"difficulty":  q.Difficulty,
			"quest_level": q.QuestLevel,
		})
	}

	return map[string]interface{}{
		"user_level":      user.Level,
		"user_xp":         user.XpPoints,
		"completed_quests": completedDetails,
		"active_quests_count": len(activeQuests),
		"preferred_categories": getPreferredCategories(completedQuests),
	}, nil
}

func getPreferredCategories(quests []models.Quest) []string {
	categoryCount := make(map[string]int)
	for _, q := range quests {
		categoryCount[q.Category]++
	}

	// Сортируем по популярности
	type catCount struct {
		category string
		count    int
	}
	cats := make([]catCount, 0, len(categoryCount))
	for cat, count := range categoryCount {
		cats = append(cats, catCount{cat, count})
	}

	// Простая сортировка (можно улучшить)
	result := make([]string, 0, 3)
	for i := 0; i < len(cats) && i < 3; i++ {
		maxIdx := i
		for j := i + 1; j < len(cats); j++ {
			if cats[j].count > cats[maxIdx].count {
				maxIdx = j
			}
		}
		cats[i], cats[maxIdx] = cats[maxIdx], cats[i]
		result = append(result, cats[i].category)
	}

	return result
}
