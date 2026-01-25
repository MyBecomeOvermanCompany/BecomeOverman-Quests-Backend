package repositories

import (
	"context"

	"BecomeOverMan/internal/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrAlreadyFriends = errors.New("Эти пользователи уже друзья")
)

func (r *UserRepository) AddFriend(userID, friendID int) error {
	exists, err := r.isUserExists(friendID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Такого пользователя не существует")
	}

	return r.addFriend(userID, friendID)
}

func (r *UserRepository) AddFriendbyName(userID int, friendName string) error {
	friendID, err := r.getUserIdByUsername(friendName)
	if err != nil {
		return err
	}

	return r.addFriend(userID, friendID)
}

// Проверяем, что дружба не существует
func (r *UserRepository) isFriends(userID, friendID int) (bool, error) {
	var friendshipExists bool
	err := r.db.Get(&friendshipExists, `
		SELECT EXISTS(SELECT 1 FROM friends WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))`,
		userID, friendID)
	if err != nil {
		return false, err
	}

	return friendshipExists, nil
}

func (r *UserRepository) addFriend(userID, friendID int) error {
	// Используем SendFriendRequest для создания запроса
	return r.SendFriendRequest(userID, friendID)
}

func (r *UserRepository) GetAllAcceptedFriends(userID int) ([]int, error) {
	query := `
		-- Получить все ID друзей (только ID)
		SELECT 
			CASE 
				WHEN user_id = $1 THEN friend_id 
				ELSE user_id 
			END as friend_id
		FROM friends 
		WHERE (user_id = $1 OR friend_id = $1)
		AND status = 'accepted';
	  `
	var friendsIDS []int
	err := r.db.Select(&friendsIDS, query, userID)
	return friendsIDS, err
}

func (r *UserRepository) GetFriends(userID int) ([]models.Friend, error) {
	var friends []models.Friend
	query := `
		SELECT 
			f.id,
			$1 as user_id,
			CASE 
				WHEN f.user_id = $1 THEN f.friend_id 
				ELSE f.user_id 
			END as friend_id,
			f.status,
			CASE 
				WHEN f.user_id = $1 THEN u1.username 
				ELSE u2.username 
			END as username,
			f.created_at
		FROM friends f 
		LEFT JOIN users u1 ON f.friend_id = u1.id 
		LEFT JOIN users u2 ON f.user_id = u2.id 
		WHERE (f.user_id = $1 OR f.friend_id = $1) 
		AND f.status = 'accepted'
	`
	err := r.db.Select(&friends, query, userID)
	return friends, err
}

// Shared quest methods
func (r *QuestRepository) CreateSharedQuest(user1ID, user2ID, questID int) error {
	ctx := context.Background()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что пользователи друзья (проверяем оба направления)
	var areFriends bool
	err = tx.Get(&areFriends, `
		SELECT EXISTS(
			SELECT 1 FROM friends 
			WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
			AND status = 'accepted'
		)`, user1ID, user2ID)
	if err != nil {
		return err
	}
	if !areFriends {
		return errors.New("users are not friends")
	}

	// Создаем shared quest
	_, err = tx.Exec(`
		INSERT INTO shared_quests (user1_id, user2_id, quest_id, status) 
		VALUES ($1, $2, $3, 'active')`,
		user1ID, user2ID, questID)
	if err != nil {
		return err
	}

	// Стартуем квест для обоих пользователей
	if err := r.startQuestForUser(tx, user1ID, questID); err != nil {
		return err
	}
	if err := r.startQuestForUser(tx, user2ID, questID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *QuestRepository) startQuestForUser(tx *sqlx.Tx, userID, questID int) error {
	// Покупаем квест (если еще не куплен)
	var alreadyPurchased bool
	err := tx.Get(&alreadyPurchased, `
		SELECT EXISTS(SELECT 1 FROM user_quests WHERE user_id = $1 AND quest_id = $2)`,
		userID, questID)
	if err != nil {
		return err
	}

	if alreadyPurchased {
		return errors.New("quest already purchased")
	}

	// Получаем цену квеста
	var price int
	err = tx.Get(&price, "SELECT price FROM quests WHERE id = $1", questID)
	if err != nil {
		return err
	}

	// Проверяем баланс
	var balance int
	err = tx.Get(&balance, "SELECT coin_balance FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	if balance < price {
		return errors.New("not enough coins for shared quest")
	}

	// Покупаем квест
	_, err = tx.Exec(`
			INSERT INTO user_quests (user_id, quest_id, status) 
			VALUES ($1, $2, 'purchased')`,
		userID, questID)
	if err != nil {
		return err
	}

	// Списываем монеты
	_, err = tx.Exec(`
			UPDATE users SET coin_balance = coin_balance - $1 WHERE id = $2`,
		price, userID)
	if err != nil {
		return err
	}

	// Создаем user_tasks для всех задач квеста
	_, err = tx.Exec(`
		INSERT INTO user_tasks (user_id, task_id, quest_id, status)
		SELECT $1, qt.task_id, qt.quest_id, 'not_started'
		FROM quest_tasks qt
		WHERE qt.quest_id = $2
		ORDER BY qt.task_order
	`, userID, questID)
	if err != nil {
		return err
	}

	// Стартуем квест
	_, err = tx.Exec(`
			UPDATE user_quests 
			SET status = 'started', started_at = NOW(), expires_at = (
				SELECT NOW() + (time_limit_hours || ' hours')::interval 
				FROM quests WHERE id = $2
			)
			WHERE user_id = $1 AND quest_id = $2`,
		userID, questID)
	if err != nil {
		return err
	}

	// Активируем задачи
	_, err = tx.Exec(`
		UPDATE user_tasks
		SET status = 'active'
		WHERE user_id = $1 AND quest_id = $2 AND status = 'not_started'
	`, userID, questID)

	return err
}

// GetFriendRequests - получить входящие запросы в друзья
func (r *UserRepository) GetFriendRequests(userID int) ([]models.Friend, error) {
	var requests []models.Friend
	query := `
		SELECT 
			f.id,
			$1 as user_id,
			f.user_id as friend_id,
			f.status,
			u.username,
			f.created_at
		FROM friends f 
		INNER JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = $1 
		AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`
	err := r.db.Select(&requests, query, userID)
	return requests, err
}

// AcceptFriendRequest - принять запрос в друзья
func (r *UserRepository) AcceptFriendRequest(userID, friendID int) error {
	result, err := r.db.Exec(`
		UPDATE friends 
		SET status = 'accepted'
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`, friendID, userID)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("friend request not found")
	}
	
	return nil
}

// RejectFriendRequest - отклонить запрос в друзья
func (r *UserRepository) RejectFriendRequest(userID, friendID int) error {
	_, err := r.db.Exec(`
		DELETE FROM friends 
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`, friendID, userID)
	return err
}

// SendFriendRequest - отправить запрос в друзья (pending)
func (r *UserRepository) SendFriendRequest(userID, friendID int) error {
	exists, err := r.isUserExists(friendID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Такого пользователя не существует")
	}

	// Проверяем, что запрос не существует
	isFriends, err := r.isFriends(userID, friendID)
	if err != nil {
		return err
	}
	if isFriends {
		return ErrAlreadyFriends
	}

	// Проверяем, что нет обратного запроса
	var reverseRequest bool
	err = r.db.Get(&reverseRequest, `
		SELECT EXISTS(
			SELECT 1 FROM friends 
			WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
		)`, friendID, userID)
	if err != nil {
		return err
	}
	if reverseRequest {
		// Если есть обратный запрос, автоматически принимаем
		return r.AcceptFriendRequest(userID, friendID)
	}

	// Создаем запрос со статусом pending
	_, err = r.db.Exec(`
		INSERT INTO friends (user_id, friend_id, status) 
		VALUES ($1, $2, 'pending')
		ON CONFLICT (user_id, friend_id) DO NOTHING
	`, userID, friendID)
	return err
}

// GetFriendStats - получить статистику друга
type FriendStats struct {
	UserID              int    `json:"user_id" db:"user_id"`
	Username            string `json:"username" db:"username"`
	TotalQuestsCompleted int   `json:"total_quests_completed" db:"total_quests_completed"`
	TotalTasksCompleted  int   `json:"total_tasks_completed" db:"total_tasks_completed"`
	CurrentStreak        int   `json:"current_streak" db:"current_streak"`
	LongestStreak        int   `json:"longest_streak" db:"longest_streak"`
	Level                int   `json:"level" db:"level"`
}

func (r *UserRepository) GetFriendStats(friendID int) (*FriendStats, error) {
	var stats FriendStats
	query := `
		SELECT 
			u.id as user_id,
			u.username,
			(SELECT COUNT(*) FROM user_quests WHERE user_id = u.id AND status = 'completed') as total_quests_completed,
			(SELECT COUNT(*) FROM user_tasks WHERE user_id = u.id AND status = 'completed') as total_tasks_completed,
			COALESCE(u.current_streak, 0) as current_streak,
			COALESCE(u.longest_streak, 0) as longest_streak,
			COALESCE(u.level, 1) as level
		FROM users u
		WHERE u.id = $1
	`
	err := r.db.Get(&stats, query, friendID)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetSharedQuests - получить совместные квесты пользователя
func (r *QuestRepository) GetSharedQuests(ctx context.Context, userID int) ([]models.SharedQuest, error) {
	var sharedQuests []models.SharedQuest
	query := `
		SELECT 
			sq.id,
			sq.quest_id,
			sq.user1_id,
			sq.user2_id,
			sq.status,
			sq.created_at,
			q.title,
			q.description,
			CASE 
				WHEN sq.user1_id = $1 THEN u2.username
				ELSE u1.username
			END as friend_username
		FROM shared_quests sq
		INNER JOIN quests q ON sq.quest_id = q.id
		LEFT JOIN users u1 ON sq.user1_id = u1.id
		LEFT JOIN users u2 ON sq.user2_id = u2.id
		WHERE (sq.user1_id = $1 OR sq.user2_id = $1)
		ORDER BY sq.created_at DESC
	`
	err := r.db.SelectContext(ctx, &sharedQuests, query, userID)
	return sharedQuests, err
}
