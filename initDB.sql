-- Удаление таблиц, если они существуют (с правильным порядком и CASCADE)
DROP TABLE IF EXISTS task_habit_requirements CASCADE;
DROP TABLE IF EXISTS habit_tracking CASCADE;
DROP TABLE IF EXISTS user_passive_buffs CASCADE;
DROP TABLE IF EXISTS quest_branches CASCADE;
DROP TABLE IF EXISTS quest_prerequisites CASCADE;
DROP TABLE IF EXISTS development_branches CASCADE;
DROP TABLE IF EXISTS user_achievements CASCADE;
DROP TABLE IF EXISTS achievements CASCADE;
DROP TABLE IF EXISTS user_coin_transactions CASCADE;
DROP TABLE IF EXISTS user_daily_streaks CASCADE;
DROP TABLE IF EXISTS user_tasks CASCADE;
DROP TABLE IF EXISTS task_variants CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS user_completed_quests CASCADE;
DROP TABLE IF EXISTS user_current_quests CASCADE;
DROP TABLE IF EXISTS user_quests CASCADE;
DROP TABLE IF EXISTS quest_tasks CASCADE;
DROP TABLE IF EXISTS quests CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS friends CASCADE;
DROP TABLE IF EXISTS shared_quests CASCADE;

-- Удаление типов
DROP TYPE IF EXISTS category_name CASCADE;
DROP TYPE IF EXISTS difficulty_level CASCADE;
DROP TYPE IF EXISTS task_type CASCADE;
DROP TYPE IF EXISTS rarity CASCADE;

-- Создание ENUM типов
CREATE TYPE category_name AS ENUM ('health', 'intelligence', 'charisma', 'willpower');
CREATE TYPE rarity AS ENUM ('free', 'common', 'rare', 'epic', 'legendary');
CREATE TYPE task_type AS ENUM ('daily', 'weekly', 'special', 'user_generated');

-- Таблица пользователей с расширенными полями
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    xp_points INT DEFAULT 0,
    coin_balance INT DEFAULT 0,
    level INT DEFAULT 1,

    -- TODO: если появятся еще ветки то поправим
    health_level INT DEFAULT 0,
    mental_health_level INT DEFAULT 0,
    intelligence_level INT DEFAULT 0,
    charisma_level INT DEFAULT 0,
    willpower_level INT DEFAULT 0,

    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица задач
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,

    difficulty INT DEFAULT 0,
    rarity VARCHAR(255) NOT NULL DEFAULT 'free', -- 'free', 'common', 'rare', 'epic', 'legendary'
    category VARCHAR(255) NOT NULL,
    base_xp_reward INT NOT NULL DEFAULT 0,
    base_coin_reward INT NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    -- type task_type NOT NULL DEFAULT 'daily',
    -- cooldown_hours INT DEFAULT 24,
);

-- Транзакции валюты пользователя
CREATE TABLE user_coin_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    reference_id INT, -- ID связанной сущности (задача, покупка и т.д.)

    transaction_type VARCHAR(50) NOT NULL, -- 'earned', 'spent', 'bonus'
    amount INT NOT NULL,
    
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Достижения
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    criteria_json JSONB NOT NULL, -- например: {"tasks_completed": 100}
    bonus_json JSONB, -- например заморозка или пассивный бонус
    reward_xp INT DEFAULT 0,
    reward_coin INT DEFAULT 0,
    is_secret BOOLEAN DEFAULT FALSE
);

-- Достижения пользователей
CREATE TABLE user_achievements (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    achievement_id INT NOT NULL,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_achievement UNIQUE (user_id, achievement_id)
);

-- Таблица квестов
CREATE TABLE quests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255) NOT NULL, -- 'health', 'willpower', 'intelligence', 'charisma'
    rarity VARCHAR(255) NOT NULL, -- 'free', 'common', 'rare', 'epic', 'legendary'
    difficulty INT NOT NULL DEFAULT 0,
    price INT NOT NULL DEFAULT 0,
    tasks_count INT DEFAULT 1, -- Сколько задач в квесте
    conditions_json JSONB,
    bonus_json JSONB,
    is_sequential BOOLEAN DEFAULT FALSE,   -- Нужно ли выполнять по порядку
    reward_xp INT NOT NULL,
    reward_coin INT NOT NULL,
    time_limit_hours INT DEFAULT 0,          -- Ограничение по времени (опционально)
    quest_level INT DEFAULT 1,               -- Уровень квеста (1, 2, 3, ...)
    parent_quest_id INT REFERENCES quests(id) ON DELETE SET NULL, -- Связь с родительским квестом
    max_level INT DEFAULT 1                  -- Максимальный уровень для этого квеста
);


-- задачи пользователя
CREATE TABLE user_tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    quest_id INT REFERENCES quests(id) ON DELETE SET NULL,  -- опционально

    status VARCHAR(50) NOT NULL DEFAULT 'active',         -- not_started, active, completed, TODO: failed
    scheduled_start TIMESTAMP,
    scheduled_end TIMESTAMP,
    deadline TIMESTAMP,
    duration INT,                                        -- время выделенное на задачу в минутах (?)
    updated_by_ai BOOLEAN DEFAULT FALSE,                 -- была ли запланирована AI

    is_confirmed BOOL DEFAULT FALSE NOT NULL,            -- прежнее поле
    completed_at TIMESTAMP,                              -- прежнее поле
    xp_gained INT NOT NULL DEFAULT 0,
    coin_gained INT NOT NULL DEFAULT 0,
    current_streak INT DEFAULT 0,                        -- Текущая серия дней подряд (для habit tracking)
    last_completed_date DATE,                            -- Дата последнего выполнения (для habit tracking)

    CONSTRAINT unique_user_task UNIQUE (user_id, task_id)
);

-- Связь квестов и задач (какие задачи входят в квест)
CREATE TABLE quest_tasks (
    id SERIAL PRIMARY KEY,
    quest_id INT NOT NULL,
    task_id INT NOT NULL,

    task_order INT,                        -- Порядок (если is_sequential = TRUE)

    FOREIGN KEY (quest_id) REFERENCES quests(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- Прогресс пользователя по квестам
CREATE TABLE user_quests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    quest_id INT NOT NULL,

    status VARCHAR(255) NOT NULL DEFAULT 'purchased', -- "purchased", "started", "failed", "completed"

    xp_gained INT,
    coin_gained INT,

    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    expires_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (quest_id) REFERENCES quests(id) ON DELETE CASCADE
);

-- friends
-- Друзья
CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, friend_id)
);

-- Совместные квесты
CREATE TABLE shared_quests (
    id SERIAL PRIMARY KEY,
    quest_id INTEGER NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    user1_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- Система дерева квестов и пассивных баффов
-- ============================================

-- 1. Таблица для зависимостей квестов (prerequisites)
CREATE TABLE quest_prerequisites (
    id SERIAL PRIMARY KEY,
    quest_id INT NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    prerequisite_quest_id INT NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    required_count INT DEFAULT 1, -- Сколько из prerequisites нужно выполнить
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_quest_prerequisite UNIQUE (quest_id, prerequisite_quest_id),
    CONSTRAINT no_self_reference CHECK (quest_id != prerequisite_quest_id)
);

-- 2. Таблица для веток развития
CREATE TABLE development_branches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- 'health', 'charisma', 'intelligence', 'money'
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_branch_id INT REFERENCES development_branches(id) ON DELETE SET NULL,
    level INT DEFAULT 1, -- Уровень вложенности (1 = основная ветка)
    icon VARCHAR(50), -- Иконка для отображения
    color VARCHAR(7), -- Цвет ветки (#000000)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Связь квестов с ветками (многие-ко-многим)
CREATE TABLE quest_branches (
    id SERIAL PRIMARY KEY,
    quest_id INT NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    branch_id INT NOT NULL REFERENCES development_branches(id) ON DELETE CASCADE,
    weight DECIMAL(3,2) DEFAULT 1.0, -- Вес квеста в этой ветке (0.0-1.0)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_quest_branch UNIQUE (quest_id, branch_id)
);

-- 4. Таблица для пассивных баффов от завершенных квестов
CREATE TABLE user_passive_buffs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id INT NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    buff_type VARCHAR(50) NOT NULL, -- 'reward_multiplier', 'stat_boost', 'unlock_quest'
    buff_data JSONB NOT NULL, -- Данные баффа (multiplier, category, stat_name, etc.)
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_quest_buff UNIQUE (user_id, quest_id)
);

-- 5. Индексы для производительности
CREATE INDEX idx_quest_prerequisites_quest_id ON quest_prerequisites(quest_id);
CREATE INDEX idx_quest_prerequisites_prerequisite_id ON quest_prerequisites(prerequisite_quest_id);
CREATE INDEX idx_quest_branches_quest_id ON quest_branches(quest_id);
CREATE INDEX idx_quest_branches_branch_id ON quest_branches(branch_id);
CREATE INDEX idx_user_passive_buffs_user_id ON user_passive_buffs(user_id);
CREATE INDEX idx_user_passive_buffs_quest_id ON user_passive_buffs(quest_id);
CREATE INDEX idx_development_branches_parent ON development_branches(parent_branch_id);

-- 6. Вставка основных веток развития
INSERT INTO development_branches (name, display_name, description, level, icon, color) VALUES
('health', 'Здоровье', 'Физическое и ментальное здоровье, энергия, внешний вид', 1, 'favorite', '#ef4444'),
('charisma', 'Харизма', 'Социальные навыки, уверенность, лидерство, привлекательность', 1, 'people', '#8b5cf6'),
('intelligence', 'Интеллект', 'Когнитивные навыки, образование, креативность, мышление', 1, 'psychology', '#3b82f6'),
('money', 'Деньги', 'Заработок, финансовая грамотность, инвестиции, бизнес', 1, 'attach_money', '#10b981'),
('willpower', 'Сила воли', 'Дисциплина, самоконтроль, формирование привычек, продуктивность', 1, 'fitness_center', '#f59e0b');

-- 7. Подветки для Health
INSERT INTO development_branches (name, display_name, description, parent_branch_id, level, icon, color) VALUES
('health_physical', 'Физическое здоровье', 'Спорт, питание, сон, выносливость', 
 (SELECT id FROM development_branches WHERE name = 'health'), 2, 'fitness_center', '#ef4444'),
('health_mental', 'Ментальное здоровье', 'Медитация, стресс-менеджмент, психология', 
 (SELECT id FROM development_branches WHERE name = 'health'), 2, 'self_improvement', '#f59e0b'),
('health_appearance', 'Внешний вид', 'Стиль, уход за собой, гигиена', 
 (SELECT id FROM development_branches WHERE name = 'health'), 2, 'face', '#ec4899');

-- 8. Подветки для Charisma
INSERT INTO development_branches (name, display_name, description, parent_branch_id, level, icon, color) VALUES
('charisma_social', 'Социальные навыки', 'Коммуникация, networking, общение', 
 (SELECT id FROM development_branches WHERE name = 'charisma'), 2, 'groups', '#8b5cf6'),
('charisma_confidence', 'Уверенность', 'Самооценка, смелость, открытость', 
 (SELECT id FROM development_branches WHERE name = 'charisma'), 2, 'star', '#a855f7'),
('charisma_leadership', 'Лидерство', 'Управление, влияние, ответственность', 
 (SELECT id FROM development_branches WHERE name = 'charisma'), 2, 'military_tech', '#9333ea'),
('charisma_attraction', 'Привлекательность', 'Романтические отношения, обаяние', 
 (SELECT id FROM development_branches WHERE name = 'charisma'), 2, 'favorite', '#7c3aed');

-- 9. Подветки для Intelligence
INSERT INTO development_branches (name, display_name, description, parent_branch_id, level, icon, color) VALUES
('intelligence_cognitive', 'Когнитивные навыки', 'Память, концентрация, скорость мышления', 
 (SELECT id FROM development_branches WHERE name = 'intelligence'), 2, 'memory', '#3b82f6'),
('intelligence_education', 'Образование', 'Чтение, обучение, знания', 
 (SELECT id FROM development_branches WHERE name = 'intelligence'), 2, 'school', '#2563eb'),
('intelligence_creativity', 'Креативность', 'Творчество, инновации, идеи', 
 (SELECT id FROM development_branches WHERE name = 'intelligence'), 2, 'palette', '#1d4ed8'),
('intelligence_critical', 'Критическое мышление', 'Анализ, логика, решение проблем', 
 (SELECT id FROM development_branches WHERE name = 'intelligence'), 2, 'lightbulb', '#1e40af');

-- 10. Подветки для Money
INSERT INTO development_branches (name, display_name, description, parent_branch_id, level, icon, color) VALUES
('money_earning', 'Заработок', 'Карьера, работа, доход', 
 (SELECT id FROM development_branches WHERE name = 'money'), 2, 'work', '#10b981'),
('money_literacy', 'Финансовая грамотность', 'Бюджет, планирование, экономия', 
 (SELECT id FROM development_branches WHERE name = 'money'), 2, 'account_balance', '#059669'),
('money_investing', 'Инвестиции', 'Акции, недвижимость, пассивный доход', 
 (SELECT id FROM development_branches WHERE name = 'money'), 2, 'trending_up', '#047857'),
('money_business', 'Предпринимательство', 'Стартапы, бизнес, продажи', 
 (SELECT id FROM development_branches WHERE name = 'money'), 2, 'store', '#065f46');

-- ============================================
-- Система habit tracking (уровни квестов уже добавлены выше)
-- ============================================

-- Таблица для habit tracking (отслеживание выполнения задач по дням)
CREATE TABLE habit_tracking (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id INT NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    completion_date DATE NOT NULL, -- Дата выполнения
    completion_time TIME, -- Время выполнения (для проверки daytime)
    is_confirmed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_task_date UNIQUE (user_id, task_id, completion_date)
);

-- Таблица для хранения требований habit tracking (сколько дней подряд нужно выполнить)
CREATE TABLE task_habit_requirements (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    consecutive_days INT DEFAULT 1, -- Сколько дней подряд нужно выполнить
    daytime_required VARCHAR(50), -- 'morning', 'afternoon', 'evening', 'any' или NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_task_habit UNIQUE (task_id)
);

-- Индексы для habit tracking
CREATE INDEX idx_habit_tracking_user_quest ON habit_tracking(user_id, quest_id);
CREATE INDEX idx_habit_tracking_date ON habit_tracking(completion_date);
CREATE INDEX idx_task_habit_requirements_task ON task_habit_requirements(task_id);
CREATE INDEX idx_quests_parent_quest_id ON quests(parent_quest_id);
CREATE INDEX idx_quests_quest_level ON quests(quest_level);