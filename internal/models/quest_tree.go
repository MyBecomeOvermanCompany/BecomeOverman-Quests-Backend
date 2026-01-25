package models

import "encoding/json"

// DevelopmentBranch - ветка развития
type DevelopmentBranch struct {
	ID             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	DisplayName    string `json:"display_name" db:"display_name"`
	Description    string `json:"description" db:"description"`
	ParentBranchID *int   `json:"parent_branch_id" db:"parent_branch_id"`
	Level          int    `json:"level" db:"level"`
	Icon           string `json:"icon" db:"icon"`
	Color          string `json:"color" db:"color"`
}

// QuestPrerequisite - зависимость квеста
type QuestPrerequisite struct {
	ID                  int `json:"id" db:"id"`
	QuestID             int `json:"quest_id" db:"quest_id"`
	PrerequisiteQuestID int `json:"prerequisite_quest_id" db:"prerequisite_quest_id"`
	RequiredCount       int `json:"required_count" db:"required_count"`
}

// QuestBranch - связь квеста с веткой
type QuestBranch struct {
	ID       int     `json:"id" db:"id"`
	QuestID  int     `json:"quest_id" db:"quest_id"`
	BranchID int     `json:"branch_id" db:"branch_id"`
	Weight   float64 `json:"weight" db:"weight"`
}

// PassiveBuff - пассивный бафф от завершенного квеста
type PassiveBuff struct {
	ID        int              `json:"id" db:"id"`
	UserID    int              `json:"user_id" db:"user_id"`
	QuestID   int              `json:"quest_id" db:"quest_id"`
	BuffType  string           `json:"buff_type" db:"buff_type"` // 'reward_multiplier', 'stat_boost', 'unlock_quest'
	BuffData  *json.RawMessage `json:"buff_data" db:"buff_data"`
	IsActive  bool             `json:"is_active" db:"is_active"`
	CreatedAt string           `json:"created_at" db:"created_at"`
}

// BuffData структуры для разных типов баффов
type RewardMultiplierBuff struct {
	Type        string  `json:"type"`       // "reward_multiplier"
	Category    string  `json:"category"`   // "health", "charisma", etc.
	Multiplier  float64 `json:"multiplier"` // 1.1 = +10%
	Description string  `json:"description"`
}

type StatBoostBuff struct {
	Type        string `json:"type"`         // "stat_boost"
	StatName    string `json:"stat_name"`    // "health_level", "charisma_level", etc.
	BoostAmount int    `json:"boost_amount"` // +1, +2, etc.
	Description string `json:"description"`
}

type UnlockQuestBuff struct {
	Type        string `json:"type"`     // "unlock_quest"
	QuestID     int    `json:"quest_id"` // ID квеста который разблокируется
	Description string `json:"description"`
}

// QuestWithPrerequisites - квест с информацией о зависимостях
type QuestWithPrerequisites struct {
	Quest         Quest               `json:"quest"`
	Prerequisites []QuestPrerequisite `json:"prerequisites"`
	Branches      []DevelopmentBranch `json:"branches"`
	IsUnlocked    bool                `json:"is_unlocked"`    // Доступен ли квест пользователю
	UnlockedCount int                 `json:"unlocked_count"` // Сколько prerequisites выполнено
	RequiredCount int                 `json:"required_count"` // Сколько нужно выполнить
}
