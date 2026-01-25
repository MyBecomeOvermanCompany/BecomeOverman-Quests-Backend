package repositories

import (
	"BecomeOverMan/internal/models"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type HabitTrackingRepository struct {
	db *sqlx.DB
}

func NewHabitTrackingRepository(db *sqlx.DB) *HabitTrackingRepository {
	return &HabitTrackingRepository{db: db}
}

// MarkTaskCompletedToday - отметить задачу как выполненную сегодня
func (r *HabitTrackingRepository) MarkTaskCompletedToday(ctx context.Context, userID, questID, taskID int, completionTime *string) error {
	today := time.Now()
	
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO habit_tracking (user_id, quest_id, task_id, completion_date, completion_time, is_confirmed)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (user_id, task_id, completion_date) 
		DO UPDATE SET 
			completion_time = COALESCE(EXCLUDED.completion_time, habit_tracking.completion_time),
			is_confirmed = TRUE
	`, userID, questID, taskID, today.Format("2006-01-02"), completionTime)
	
	if err != nil {
		return err
	}
	
	// Обновляем streak в user_tasks
	consecutiveDays, err := r.GetTaskHabitProgress(ctx, userID, taskID)
	if err == nil {
		_, err = r.db.ExecContext(ctx, `
			UPDATE user_tasks
			SET current_streak = $1, last_completed_date = $2
			WHERE user_id = $3 AND quest_id = $4 AND task_id = $5
		`, consecutiveDays, today.Format("2006-01-02"), userID, questID, taskID)
	}
	
	return err
}

// GetTaskHabitProgress - получить прогресс выполнения задачи (сколько дней подряд выполнено)
func (r *HabitTrackingRepository) GetTaskHabitProgress(ctx context.Context, userID, taskID int) (int, error) {
	var consecutiveDays int
	
	// Получаем количество последовательных дней выполнения
	err := r.db.GetContext(ctx, &consecutiveDays, `
		WITH RECURSIVE consecutive_days AS (
			SELECT completion_date, 
				   ROW_NUMBER() OVER (ORDER BY completion_date DESC) as rn
			FROM habit_tracking
			WHERE user_id = $1 AND task_id = $2 AND is_confirmed = TRUE
			ORDER BY completion_date DESC
		),
		date_gaps AS (
			SELECT completion_date,
				   completion_date - INTERVAL '1 day' * (rn - 1) as expected_date,
				   rn
			FROM consecutive_days
		)
		SELECT COUNT(*) 
		FROM date_gaps
		WHERE completion_date = expected_date
		AND completion_date >= CURRENT_DATE - INTERVAL '30 days'
	`, userID, taskID)
	
	if err == sql.ErrNoRows {
		return 0, nil
	}
	
	return consecutiveDays, err
}

// GetTaskHabitRequirement - получить требования для habit tracking задачи
func (r *HabitTrackingRepository) GetTaskHabitRequirement(ctx context.Context, taskID int) (*models.TaskHabitRequirement, error) {
	var requirement models.TaskHabitRequirement
	err := r.db.GetContext(ctx, &requirement, `
		SELECT id, task_id, consecutive_days, daytime_required, created_at
		FROM task_habit_requirements
		WHERE task_id = $1
	`, taskID)
	
	if err == sql.ErrNoRows {
		return nil, nil // Нет требований
	}
	
	return &requirement, err
}

// GetTaskCompletionsForPeriod - получить выполнения задачи за период
func (r *HabitTrackingRepository) GetTaskCompletionsForPeriod(ctx context.Context, userID, taskID int, days int) ([]models.HabitTracking, error) {
	var completions []models.HabitTracking
	err := r.db.SelectContext(ctx, &completions, `
		SELECT id, user_id, quest_id, task_id, completion_date, completion_time, is_confirmed, created_at
		FROM habit_tracking
		WHERE user_id = $1 AND task_id = $2 
		AND completion_date >= CURRENT_DATE - INTERVAL '1 day' * $3
		AND is_confirmed = TRUE
		ORDER BY completion_date DESC
	`, userID, taskID, days)
	
	return completions, err
}

// CheckDaytimeRequirement - проверить, соответствует ли время выполнения требованиям
func (r *HabitTrackingRepository) CheckDaytimeRequirement(completionTime *string, requiredDaytime *string) bool {
	if requiredDaytime == nil || *requiredDaytime == "any" || completionTime == nil {
		return true
	}
	
	// Парсим время
	t, err := time.Parse("15:04:05", *completionTime)
	if err != nil {
		return false
	}
	
	hour := t.Hour()
	
	switch *requiredDaytime {
	case "morning":
		return hour >= 6 && hour < 12
	case "afternoon":
		return hour >= 12 && hour < 18
	case "evening":
		return hour >= 18 || hour < 6
	default:
		return true
	}
}

// GetMissedDays - получить список пропущенных дней в текущем периоде
func (r *HabitTrackingRepository) GetMissedDays(ctx context.Context, userID, taskID int, requiredDays int) ([]time.Time, error) {
	var missedDays []time.Time
	
	// Получаем последнее выполнение
	var lastCompletion *time.Time
	err := r.db.GetContext(ctx, &lastCompletion, `
		SELECT MAX(completion_date) 
		FROM habit_tracking
		WHERE user_id = $1 AND task_id = $2 AND is_confirmed = TRUE
	`, userID, taskID)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	
	// Если нет выполнений, все дни до сегодня пропущены
	if lastCompletion == nil {
		for i := requiredDays - 1; i >= 0; i-- {
			date := today.AddDate(0, 0, -i)
			if date.Before(today) {
				missedDays = append(missedDays, date)
			}
		}
		return missedDays, nil
	}
	
	// Проверяем дни от последнего выполнения до сегодня
	lastDate := time.Date(lastCompletion.Year(), lastCompletion.Month(), lastCompletion.Day(), 0, 0, 0, 0, lastCompletion.Location())
	
	for d := lastDate.AddDate(0, 0, 1); d.Before(today) || d.Equal(today); d = d.AddDate(0, 0, 1) {
		// Проверяем, было ли выполнение в этот день
		var exists bool
		err := r.db.GetContext(ctx, &exists, `
			SELECT EXISTS(
				SELECT 1 FROM habit_tracking
				WHERE user_id = $1 AND task_id = $2 
				AND completion_date = $3 AND is_confirmed = TRUE
			)
		`, userID, taskID, d.Format("2006-01-02"))
		
		if err != nil {
			return nil, err
		}
		
		if !exists {
			missedDays = append(missedDays, d)
		}
	}
	
	return missedDays, nil
}

// FreezeMissedDay - "заморозить" пропущенный день (отметить как выполненный за монеты)
func (r *HabitTrackingRepository) FreezeMissedDay(ctx context.Context, userID, questID, taskID int, date time.Time) error {
	dateStr := date.Format("2006-01-02")
	
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO habit_tracking (user_id, quest_id, task_id, completion_date, completion_time, is_confirmed)
		VALUES ($1, $2, $3, $4, NULL, TRUE)
		ON CONFLICT (user_id, task_id, completion_date) 
		DO UPDATE SET is_confirmed = TRUE
	`, userID, questID, taskID, dateStr)
	
	return err
}