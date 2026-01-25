package services

import (
	"BecomeOverMan/internal/models"
	"context"
	"fmt"
	"time"
)

// ContinueQuest - начать следующий уровень квеста
func (s *QuestService) ContinueQuest(ctx context.Context, userID, completedQuestID int) (*models.Quest, error) {
	// Получаем следующий уровень квеста
	nextQuest, err := s.questRepo.GetNextQuestLevel(ctx, completedQuestID)
	if err != nil {
		return nil, fmt.Errorf("error getting next quest level: %w", err)
	}
	if nextQuest == nil {
		return nil, fmt.Errorf("no next level available for this quest")
	}

	// Покупаем и начинаем следующий уровень
	err = s.PurchaseQuest(ctx, userID, nextQuest.ID)
	if err != nil {
		return nil, fmt.Errorf("error purchasing next quest level: %w", err)
	}

	err = s.StartQuest(ctx, userID, nextQuest.ID)
	if err != nil {
		return nil, fmt.Errorf("error starting next quest level: %w", err)
	}

	return nextQuest, nil
}

// MarkTaskHabitComplete - отметить задачу как выполненную сегодня (habit tracking)
func (s *QuestService) MarkTaskHabitComplete(ctx context.Context, userID, questID, taskID int, completionTime *string) error {
	// Получаем требования для habit tracking
	requirement, err := s.habitTrackingRepo.GetTaskHabitRequirement(ctx, taskID)
	if err != nil {
		return fmt.Errorf("error getting habit requirement: %w", err)
	}

	// Если нет требований, просто отмечаем как выполненную
	if requirement == nil {
		// Обычное выполнение задачи
		return s.CompleteTask(ctx, userID, questID, taskID)
	}

	// Проверяем время дня, если требуется
	if requirement.DaytimeRequired != nil && completionTime != nil {
		if !s.habitTrackingRepo.CheckDaytimeRequirement(completionTime, requirement.DaytimeRequired) {
			return fmt.Errorf("task must be completed during %s", *requirement.DaytimeRequired)
		}
	}

	// Отмечаем выполнение сегодня
	err = s.habitTrackingRepo.MarkTaskCompletedToday(ctx, userID, questID, taskID, completionTime)
	if err != nil {
		return fmt.Errorf("error marking task as completed: %w", err)
	}

	// Проверяем, выполнено ли требование (N дней подряд)
	consecutiveDays, err := s.habitTrackingRepo.GetTaskHabitProgress(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("error getting habit progress: %w", err)
	}

	// Если выполнено достаточно дней подряд, завершаем задачу
	if consecutiveDays >= requirement.ConsecutiveDays {
		return s.CompleteTask(ctx, userID, questID, taskID)
	}

	return nil
}

// GetTaskHabitProgress - получить прогресс habit tracking для задачи
func (s *QuestService) GetTaskHabitProgress(ctx context.Context, userID, taskID int) (map[string]interface{}, error) {
	// Получаем требования
	requirement, err := s.habitTrackingRepo.GetTaskHabitRequirement(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("error getting habit requirement: %w", err)
	}

	if requirement == nil {
		return map[string]interface{}{
			"has_requirement": false,
		}, nil
	}

	// Получаем прогресс
	consecutiveDays, err := s.habitTrackingRepo.GetTaskHabitProgress(ctx, userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("error getting habit progress: %w", err)
	}

	// Получаем выполнения за период
	completions, err := s.habitTrackingRepo.GetTaskCompletionsForPeriod(ctx, userID, taskID, requirement.ConsecutiveDays)
	if err != nil {
		return nil, fmt.Errorf("error getting completions: %w", err)
	}

	// Получаем пропущенные дни
	missedDays, err := s.habitTrackingRepo.GetMissedDays(ctx, userID, taskID, requirement.ConsecutiveDays)
	if err != nil {
		return nil, fmt.Errorf("error getting missed days: %w", err)
	}

	// Форматируем пропущенные дни для JSON
	missedDaysStr := make([]string, len(missedDays))
	for i, day := range missedDays {
		missedDaysStr[i] = day.Format("2006-01-02")
	}

	return map[string]interface{}{
		"has_requirement":   true,
		"consecutive_days":   consecutiveDays,
		"required_days":      requirement.ConsecutiveDays,
		"daytime_required":  requirement.DaytimeRequired,
		"completions":       completions,
		"is_completed":      consecutiveDays >= requirement.ConsecutiveDays,
		"missed_days":       missedDaysStr,
		"has_missed_days":   len(missedDays) > 0,
	}, nil
}

// FreezeMissedDay - заморозить пропущенный день за монеты
func (s *QuestService) FreezeMissedDay(ctx context.Context, userID, questID, taskID int, dateStr string) error {
	// Парсим дату
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	// Проверяем, что дата в прошлом (нельзя заморозить будущее)
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	if date.After(today) {
		return fmt.Errorf("cannot freeze future date")
	}

	// Проверяем баланс пользователя
	var balance int
	err = s.questRepo.GetDB().GetContext(ctx, &balance, "SELECT coin_balance FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("error getting user balance: %w", err)
	}

	// Стоимость заморозки: 50 монет за день
	freezeCost := 50
	if balance < freezeCost {
		return fmt.Errorf("not enough coins: need %d, have %d", freezeCost, balance)
	}

	// Замораживаем день
	err = s.habitTrackingRepo.FreezeMissedDay(ctx, userID, questID, taskID, date)
	if err != nil {
		return fmt.Errorf("error freezing day: %w", err)
	}

	// Списываем монеты
	_, err = s.questRepo.GetDB().ExecContext(ctx, 
		"UPDATE users SET coin_balance = coin_balance - $1 WHERE id = $2", 
		freezeCost, userID)
	if err != nil {
		return fmt.Errorf("error updating coin balance: %w", err)
	}

	// Записываем транзакцию (если таблица существует)
	_, err = s.questRepo.GetDB().ExecContext(ctx, `
		INSERT INTO user_coin_transactions (user_id, amount, transaction_type, description)
		VALUES ($1, $2, $3, $4)
	`, userID, -freezeCost, "freeze_missed_day", fmt.Sprintf("Заморозка пропущенного дня для задачи %d", taskID))
	if err != nil {
		// Не критично, просто логируем
		fmt.Printf("Warning: failed to record coin transaction: %v\n", err)
	}

	return nil
}
