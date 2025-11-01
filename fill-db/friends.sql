-- Создаем первый квест с задачами
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Утренний марафон с другом', 
        'Совместный недельный челлендж для развития силы воли и здоровья', 
        'health', 
        'common', 
        2, 
        15, 
        5,
        150, 
        75, 
        168
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Пробуждение в 6 утра', 'Проснуться в 6 утра без будильника', 1, 'common', 'willpower', 20, 10),
    ('Утренняя зарядка', 'Сделать 15-минутную зарядку', 1, 'common', 'health', 25, 12),
    ('Совместная пробежка', 'Пробежать 3 км вместе с другом', 2, 'common', 'health', 30, 15),
    ('Медитация 10 минут', 'Помедитировать 10 минут утром', 1, 'common', 'mental_health', 20, 10),
    ('План на день', 'Составить план задач на день', 1, 'common', 'intelligence', 15, 8)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;

-------------------------------------------------------------------------------------------------------------------

-- Квест для развития интеллекта и общения
INSERT INTO quests (
    title, description, category, rarity, difficulty, price, tasks_count,
    reward_xp, reward_coin, time_limit_hours
) VALUES (
    'Интеллектуальный дуэт', 
    'Совместное развитие интеллектуальных навыков', 
    'intelligence', 
    'rare', 
    3, 
    25, 
    4,
    200, 
    100, 
    120  -- 5 дней
);

-- Задачи для этого квеста
INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
('Прочитать книгу', 'Прочитать 50 страниц научной литературы', 2, 'rare', 'intelligence', 35, 18),
('Обсуждение книги', 'Обсудить прочитанное с другом в течение 30 минут', 1, 'common', 'charisma', 25, 12),
('Парное программирование', 'Вместе написать сложную программу или часть программного проекта', 3, 'rare', 'intelligence', 40, 20),
('Поделиться знаниями', 'Рассказать другу о новой изученной теме', 2, 'common', 'charisma', 30, 15);

-- Связываем (предположим, квест id=2, задачи id=6,7,8,9)
INSERT INTO quest_tasks (quest_id, task_id, task_order) VALUES
(2, 6, 1),
(2, 7, 2),
(2, 8, 3),
(2, 9, 4);


