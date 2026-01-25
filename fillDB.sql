-- ============================================
-- Тестовые данные: Квесты с иерархией и зависимостями
-- Структура: Базовые → Промежуточные → Продвинутые → Легендарные
-- ============================================

-- ============================================
-- УРОВЕНЬ 1: БАЗОВЫЕ КВЕСТЫ (без prerequisites)
-- ============================================

-- HEALTH: Базовый квест 1
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Утренняя зарядка: первые шаги', 
        'Базовый квест для формирования привычки утренней физической активности. 7 дней простых упражнений.', 
        'health', 
        'common', 
        2, 
        30, 
        5,
        200, 
        100, 
        168,  -- 7 дней
        true,
        '{"type": "reward_multiplier", "category": "health", "multiplier": 1.05, "description": "+5% к наградам за квесты здоровья"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Подготовка пространства', 'Подготовьте место для зарядки: коврик, удобная одежда, вода.', 1, 'common', 'health', 15, 8),
    ('Разминка (день 1-2)', 'Выполните 5-минутную разминку: вращения суставами, наклоны, потягивания.', 1, 'common', 'health', 20, 10),
    ('Базовые упражнения (день 3-5)', 'Выполните 10 приседаний, 10 отжиманий, планку 30 секунд.', 2, 'common', 'health', 30, 15),
    ('Увеличение нагрузки (день 6-7)', 'Увеличьте количество повторений на 50% или добавьте новые упражнения.', 2, 'common', 'health', 35, 18),
    ('Рефлексия', 'Оцените изменения: уровень энергии, настроение, самочувствие.', 1, 'common', 'health', 20, 10)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- HEALTH: Базовый квест 2
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Здоровое питание: основы', 
        'Базовый квест по формированию здоровых пищевых привычек. Отслеживание питания в течение недели.', 
        'health', 
        'common', 
        2, 
        40, 
        6,
        250, 
        125, 
        168,  -- 7 дней
        true,
        '{"type": "reward_multiplier", "category": "health", "multiplier": 1.05, "description": "+5% к наградам за квесты здоровья"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Анализ текущего питания', 'Запишите все приемы пищи в течение 3 дней. Оцените качество питания.', 2, 'common', 'health', 25, 12),
    ('Планирование меню', 'Составьте план здорового питания на неделю: завтрак, обед, ужин, перекусы.', 2, 'common', 'health', 30, 15),
    ('Приготовление здоровой еды', 'Приготовьте 3 здоровых блюда: завтрак, обед, ужин. Используйте свежие продукты.', 3, 'common', 'health', 40, 20),
    ('Отслеживание воды', 'Пейте минимум 2 литра воды в день. Отмечайте каждый стакан.', 2, 'common', 'health', 35, 18),
    ('Сокращение сахара', 'Сократите потребление добавленного сахара на 50% в течение 7 дней.', 3, 'rare', 'health', 50, 25),
    ('Итоговая оценка', 'Проанализируйте изменения: энергия, вес, самочувствие.', 2, 'common', 'health', 30, 15)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- INTELLIGENCE: Базовый квест 1
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Чтение: 30 минут в день', 
        'Базовый квест по формированию привычки ежедневного чтения. 14 дней по 30 минут.', 
        'intelligence', 
        'common', 
        2, 
        35, 
        4,
        220, 
        110, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "intelligence", "multiplier": 1.05, "description": "+5% к наградам за квесты интеллекта"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Выбор книги', 'Выберите книгу для чтения: художественная или нон-фикшн по интересующей теме.', 1, 'common', 'intelligence', 15, 8),
    ('Настройка времени', 'Определите время для чтения: утро, обед, вечер. Установите напоминание.', 1, 'common', 'intelligence', 20, 10),
    ('Ежедневное чтение (14 дней)', 'Читайте минимум 30 минут каждый день. Отмечайте прогресс.', 3, 'common', 'intelligence', 100, 50),
    ('Конспектирование', 'После каждого чтения записывайте 3 ключевых момента или цитаты.', 2, 'common', 'intelligence', 40, 20)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- INTELLIGENCE: Базовый квест 2
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Изучение нового навыка: основы', 
        'Базовый квест по изучению нового навыка. Выберите навык и изучайте его 10 дней.', 
        'intelligence', 
        'common', 
        3, 
        50, 
        5,
        300, 
        150, 
        240,  -- 10 дней
        true,
        '{"type": "reward_multiplier", "category": "intelligence", "multiplier": 1.05, "description": "+5% к наградам за квесты интеллекта"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Выбор навыка', 'Выберите навык для изучения: программирование, язык, музыка, рисование и т.д.', 1, 'common', 'intelligence', 20, 10),
    ('Поиск ресурсов', 'Найдите 3-5 качественных ресурсов: курсы, книги, видео, статьи.', 2, 'common', 'intelligence', 30, 15),
    ('План обучения', 'Составьте план обучения на 10 дней: что изучать каждый день.', 2, 'common', 'intelligence', 35, 18),
    ('Ежедневная практика (10 дней)', 'Изучайте и практикуйте навык минимум 1 час в день.', 4, 'common', 'intelligence', 120, 60),
    ('Применение на практике', 'Примените изученное: создайте проект, используйте навык в реальной ситуации.', 3, 'rare', 'intelligence', 50, 25)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- CHARISMA: Базовый квест 1
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Ежедневные комплименты', 
        'Базовый квест по развитию социальных навыков. Делайте 3 комплимента в день в течение 7 дней.', 
        'charisma', 
        'common', 
        2, 
        30, 
        4,
        200, 
        100, 
        168,  -- 7 дней
        true,
        '{"type": "reward_multiplier", "category": "charisma", "multiplier": 1.05, "description": "+5% к наградам за квесты харизмы"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Изучение техники комплиментов', 'Изучите, как делать искренние комплименты: конкретика, искренность, уместность.', 1, 'common', 'charisma', 20, 10),
    ('Практика комплиментов (7 дней)', 'Делайте минимум 3 искренних комплимента в день разным людям.', 3, 'common', 'charisma', 80, 40),
    ('Отслеживание реакций', 'Обращайте внимание на реакции людей. Записывайте наблюдения.', 2, 'common', 'charisma', 30, 15),
    ('Рефлексия', 'Проанализируйте изменения в ваших социальных взаимодействиях.', 1, 'common', 'charisma', 20, 10)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- CHARISMA: Базовый квест 2
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Развитие уверенности: первые шаги', 
        'Базовый квест по развитию уверенности в себе. 10 дней практики позитивного мышления.', 
        'charisma', 
        'common', 
        3, 
        45, 
        5,
        280, 
        140, 
        240,  -- 10 дней
        true,
        '{"type": "reward_multiplier", "category": "charisma", "multiplier": 1.05, "description": "+5% к наградам за квесты харизмы"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Анализ самооценки', 'Оцените текущий уровень уверенности по шкале 1-10. Запишите области для улучшения.', 2, 'common', 'charisma', 25, 12),
    ('Аффирмации', 'Создайте 5 персональных аффирмаций. Повторяйте их утром и вечером.', 2, 'common', 'charisma', 30, 15),
    ('Ежедневная практика (10 дней)', 'Каждый день записывайте 3 своих достижения и 3 сильные стороны.', 3, 'common', 'charisma', 100, 50),
    ('Выход из зоны комфорта', 'Каждый день делайте одно действие, которое выходит за пределы зоны комфорта.', 4, 'rare', 'charisma', 60, 30),
    ('Итоговая оценка', 'Оцените изменения в уровне уверенности. Запишите прогресс.', 2, 'common', 'charisma', 30, 15)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- MONEY: Базовый квест 1
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Отслеживание расходов', 
        'Базовый квест по финансовой грамотности. Отслеживайте все расходы в течение 14 дней.', 
        'money', 
        'common', 
        2, 
        35, 
        4,
        240, 
        120, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "money", "multiplier": 1.05, "description": "+5% к наградам за квесты денег"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Настройка системы учета', 'Выберите метод учета: приложение, таблица, блокнот. Настройте категории расходов.', 1, 'common', 'money', 20, 10),
    ('Ежедневное отслеживание (14 дней)', 'Записывайте каждую трату сразу после покупки. Категоризируйте расходы.', 3, 'common', 'money', 120, 60),
    ('Еженедельный анализ', 'В конце каждой недели анализируйте расходы: на что тратите больше всего.', 2, 'common', 'money', 40, 20),
    ('Итоговый отчет', 'Создайте итоговый отчет: общая сумма, категории, паттерны расходов.', 2, 'common', 'money', 30, 15)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- MONEY: Базовый квест 2
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Создание бюджета', 
        'Базовый квест по созданию личного бюджета. Правило 50/30/20.', 
        'money', 
        'common', 
        3, 
        50, 
        5,
        320, 
        160, 
        168,  -- 7 дней
        true,
        '{"type": "reward_multiplier", "category": "money", "multiplier": 1.05, "description": "+5% к наградам за квесты денег"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Анализ доходов', 'Рассчитайте средний ежемесячный доход за последние 3 месяца.', 2, 'common', 'money', 30, 15),
    ('Категоризация расходов', 'Разделите расходы на категории: нужды (50%), желания (30%), сбережения (20%).', 3, 'common', 'money', 50, 25),
    ('Создание бюджета', 'Создайте детальный бюджет на месяц по правилу 50/30/20.', 3, 'common', 'money', 60, 30),
    ('Планирование сбережений', 'Определите сумму для ежемесячных сбережений. Откройте отдельный счет.', 2, 'common', 'money', 40, 20),
    ('Тестирование бюджета', 'Следуйте бюджету в течение недели. Отслеживайте отклонения.', 3, 'rare', 'money', 50, 25)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- WILLPOWER: Базовый квест 1
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Формирование утреннего ритуала', 
        'Базовый квест по созданию продуктивного утреннего ритуала. 14 дней практики.', 
        'willpower', 
        'common', 
        2, 
        40, 
        5,
        260, 
        130, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "willpower", "multiplier": 1.05, "description": "+5% к наградам за квесты силы воли"}'::jsonb,
        1,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Планирование ритуала', 'Создайте план утреннего ритуала: 3-5 действий (медитация, зарядка, чтение и т.д.).', 2, 'common', 'willpower', 30, 15),
    ('Подготовка пространства', 'Подготовьте все необходимое для ритуала с вечера.', 1, 'common', 'willpower', 20, 10),
    ('Практика ритуала (14 дней)', 'Выполняйте утренний ритуал каждый день в одно и то же время.', 4, 'common', 'willpower', 140, 70),
    ('Отслеживание прогресса', 'Ведите дневник: что работает, что нужно изменить.', 2, 'common', 'willpower', 35, 18),
    ('Оптимизация', 'Скорректируйте ритуал на основе опыта. Упростите или добавьте элементы.', 2, 'common', 'willpower', 30, 15)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- ============================================
-- УРОВЕНЬ 2: ПРОМЕЖУТОЧНЫЕ КВЕСТЫ (требуют 1-2 базовых)
-- ============================================

-- HEALTH: Промежуточный квест 1 (требует: Утренняя зарядка + Здоровое питание)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Интегральный подход к здоровью', 
        'Продвинутый квест, объединяющий физическую активность и правильное питание. Требует завершения базовых квестов по зарядке и питанию.', 
        'health', 
        'rare', 
        4, 
        120, 
        6,
        500, 
        250, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "health", "multiplier": 1.15, "description": "+15% к наградам за квесты здоровья"}'::jsonb,
        2,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Планирование тренировок', 'Создайте план тренировок на 14 дней: сочетание кардио и силовых упражнений.', 3, 'common', 'health', 50, 25),
    ('Планирование питания', 'Составьте план питания, поддерживающий тренировки: белки, углеводы, витамины.', 3, 'common', 'health', 50, 25),
    ('Интеграция (14 дней)', 'Выполняйте тренировки и следуйте плану питания каждый день.', 5, 'rare', 'health', 200, 100),
    ('Отслеживание прогресса', 'Ведите дневник: тренировки, питание, самочувствие, энергия.', 3, 'common', 'health', 60, 30),
    ('Корректировка плана', 'На основе опыта скорректируйте план тренировок и питания.', 3, 'rare', 'health', 50, 25),
    ('Итоговая оценка', 'Проанализируйте результаты: изменения в весе, силе, выносливости, энергии.', 2, 'common', 'health', 40, 20)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- INTELLIGENCE: Промежуточный квест 1 (требует: Чтение + Изучение навыка)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Глубокое обучение: от теории к практике', 
        'Продвинутый квест по систематическому обучению. Требует завершения базовых квестов по чтению и изучению навыка.', 
        'intelligence', 
        'rare', 
        4, 
        130, 
        7,
        550, 
        275, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "intelligence", "multiplier": 1.15, "description": "+15% к наградам за квесты интеллекта"}'::jsonb,
        2,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Выбор темы для глубокого изучения', 'Выберите тему для глубокого изучения на 14 дней.', 2, 'common', 'intelligence', 30, 15),
    ('Создание учебного плана', 'Составьте детальный план: книги, курсы, практика, проекты.', 3, 'common', 'intelligence', 50, 25),
    ('Ежедневное чтение (14 дней)', 'Читайте минимум 1 час в день по выбранной теме.', 4, 'common', 'intelligence', 140, 70),
    ('Практические упражнения', 'Выполняйте практические упражнения каждый день.', 4, 'common', 'intelligence', 120, 60),
    ('Конспектирование и заметки', 'Ведите детальные конспекты: ключевые идеи, цитаты, вопросы.', 3, 'common', 'intelligence', 60, 30),
    ('Создание проекта', 'Создайте проект, демонстрирующий ваши знания: статья, презентация, код.', 5, 'rare', 'intelligence', 100, 50),
    ('Обучение других', 'Объясните изученное другому человеку или запишите видео-объяснение.', 4, 'epic', 'intelligence', 80, 40)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- CHARISMA: Промежуточный квест 1 (требует: Комплименты + Уверенность)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Мастер социальных взаимодействий', 
        'Продвинутый квест по развитию социальных навыков. Требует завершения базовых квестов по комплиментам и уверенности.', 
        'charisma', 
        'rare', 
        5, 
        150, 
        8,
        600, 
        300, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "charisma", "multiplier": 1.15, "description": "+15% к наградам за квесты харизмы"}'::jsonb,
        2,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Изучение техник общения', 'Изучите техники активного слушания, эмпатии, невербального общения.', 3, 'common', 'charisma', 50, 25),
    ('Практика в знакомой среде', 'Примените техники в общении с друзьями и семьей. Получите обратную связь.', 3, 'common', 'charisma', 60, 30),
    ('Разговоры с незнакомцами (14 дней)', 'Каждый день заводите разговор с незнакомым человеком минимум на 5 минут.', 5, 'rare', 'charisma', 200, 100),
    ('Публичные выступления', 'Выступите публично 3 раза: рассказ, презентация, мнение в группе.', 6, 'epic', 'charisma', 150, 75),
    ('Организация встречи', 'Организуйте встречу или мероприятие для группы людей (минимум 5 человек).', 5, 'rare', 'charisma', 100, 50),
    ('Менторство', 'Помогите кому-то развить навык или решить проблему. Проведите минимум 2 сессии.', 4, 'rare', 'charisma', 80, 40),
    ('Отслеживание прогресса', 'Ведите дневник: какие техники работают, что нужно улучшить.', 2, 'common', 'charisma', 40, 20),
    ('Итоговая оценка', 'Оцените изменения: уверенность, количество социальных связей, качество общения.', 2, 'common', 'charisma', 40, 20)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- MONEY: Промежуточный квест 1 (требует: Отслеживание расходов + Бюджет)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Оптимизация финансов', 
        'Продвинутый квест по оптимизации личных финансов. Требует завершения базовых квестов по отслеживанию и бюджету.', 
        'money', 
        'rare', 
        4, 
        140, 
        6,
        520, 
        260, 
        336,  -- 14 дней
        true,
        '{"type": "reward_multiplier", "category": "money", "multiplier": 1.15, "description": "+15% к наградам за квесты денег"}'::jsonb,
        2,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Анализ расходов', 'Проанализируйте расходы за последние 3 месяца. Найдите паттерны и возможности для экономии.', 3, 'common', 'money', 50, 25),
    ('Сокращение расходов', 'Найдите и реализуйте 5 способов сократить расходы на 20% без снижения качества жизни.', 4, 'rare', 'money', 80, 40),
    ('Оптимизация подписок', 'Проанализируйте все подписки и отмените неиспользуемые. Сэкономьте минимум 500 единиц в месяц.', 2, 'common', 'money', 40, 20),
    ('Создание фонда на черный день', 'Создайте резервный фонд: минимум 3 месячных дохода. Начните откладывать.', 3, 'common', 'money', 60, 30),
    ('Автоматизация сбережений', 'Настройте автоматические переводы на сберегательный счет: 20% от дохода.', 2, 'common', 'money', 50, 25),
    ('Итоговый финансовый план', 'Создайте долгосрочный финансовый план: цели, сроки, стратегия.', 3, 'rare', 'money', 70, 35)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- CROSS-CATEGORY: Промежуточный квест (требует: HEALTH + INTELLIGENCE базовые)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Здоровый ум в здоровом теле', 
        'Интегральный квест, объединяющий физическое здоровье и интеллектуальное развитие. Требует завершения базовых квестов здоровья и интеллекта.', 
        'health', 
        'rare', 
        5, 
        160, 
        7,
        650, 
        325, 
        504,  -- 21 день
        true,
        '[{"type": "stat_boost", "stat_name": "health_level", "boost_amount": 1, "description": "+1 к уровню здоровья"}, {"type": "stat_boost", "stat_name": "intelligence_level", "boost_amount": 1, "description": "+1 к уровню интеллекта"}]'::jsonb,
        2,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Планирование интеграции', 'Создайте план, объединяющий физическую активность и обучение на 21 день.', 3, 'common', 'health', 60, 30),
    ('Утренний ритуал', 'Создайте утренний ритуал: зарядка + чтение/обучение (минимум 1 час).', 4, 'common', 'health', 70, 35),
    ('Ежедневная практика (21 день)', 'Выполняйте ритуал каждый день: физическая активность + интеллектуальное развитие.', 6, 'rare', 'health', 300, 150),
    ('Изучение связи тела и ума', 'Изучите научные исследования о связи физической активности и когнитивных функций.', 3, 'common', 'intelligence', 50, 25),
    ('Применение знаний', 'Примените изученное: оптимизируйте тренировки для улучшения когнитивных функций.', 4, 'rare', 'health', 80, 40),
    ('Отслеживание прогресса', 'Ведите дневник: физические показатели, когнитивные тесты, энергия, настроение.', 3, 'common', 'health', 60, 30),
    ('Итоговая оценка', 'Проанализируйте результаты: изменения в физической форме, когнитивных способностях, продуктивности.', 3, 'rare', 'health', 70, 35)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- ============================================
-- УРОВЕНЬ 3: ПРОДВИНУТЫЕ КВЕСТЫ (требуют 2+ промежуточных)
-- ============================================

-- HEALTH: Продвинутый квест (требует: Интегральный подход к здоровью + Здоровый ум в здоровом теле)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Мастер здорового образа жизни', 
        'Эпический квест по созданию комплексной системы здорового образа жизни. Требует завершения промежуточных квестов здоровья.', 
        'health', 
        'epic', 
        6, 
        250, 
        8,
        1000, 
        500, 
        720,  -- 30 дней
        true,
        '[{"type": "stat_boost", "stat_name": "health_level", "boost_amount": 2, "description": "+2 к уровню здоровья"}, {"type": "reward_multiplier", "category": "health", "multiplier": 1.2, "description": "+20% к наградам за квесты здоровья"}]'::jsonb,
        3,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Создание системы', 'Создайте комплексную систему здорового образа жизни: тренировки, питание, сон, восстановление.', 4, 'rare', 'health', 100, 50),
    ('Планирование на месяц', 'Составьте детальный план на 30 дней: тренировки, питание, режим сна, активный отдых.', 4, 'rare', 'health', 120, 60),
    ('Ежедневное выполнение (30 дней)', 'Следуйте системе каждый день: тренировки, питание, сон, восстановление.', 7, 'epic', 'health', 500, 250),
    ('Мониторинг показателей', 'Отслеживайте ключевые показатели: вес, сила, выносливость, энергия, настроение, сон.', 4, 'rare', 'health', 100, 50),
    ('Оптимизация системы', 'Еженедельно анализируйте и оптимизируйте систему на основе результатов.', 4, 'rare', 'health', 80, 40),
    ('Интеграция с другими сферами', 'Интегрируйте здоровый образ жизни с работой, обучением, социальной жизнью.', 5, 'epic', 'health', 120, 60),
    ('Обучение других', 'Поделитесь опытом: помогите кому-то создать свою систему здорового образа жизни.', 5, 'epic', 'health', 100, 50),
    ('Итоговая оценка', 'Проанализируйте результаты: физические изменения, энергия, продуктивность, качество жизни.', 3, 'rare', 'health', 80, 40)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- CROSS-CATEGORY: Продвинутый квест (требует: Мастер социальных взаимодействий + Глубокое обучение)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Лидер и наставник', 
        'Эпический квест по развитию лидерских качеств через обучение других. Требует завершения промежуточных квестов харизмы и интеллекта.', 
        'charisma', 
        'epic', 
        7, 
        300, 
        9,
        1200, 
        600, 
        504,  -- 21 день
        true,
        '[{"type": "stat_boost", "stat_name": "charisma_level", "boost_amount": 2, "description": "+2 к уровню харизмы"}, {"type": "stat_boost", "stat_name": "intelligence_level", "boost_amount": 1, "description": "+1 к уровню интеллекта"}]'::jsonb,
        3,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Выбор области для обучения', 'Выберите область, в которой вы можете обучать других на основе ваших знаний.', 3, 'common', 'charisma', 60, 30),
    ('Создание учебной программы', 'Создайте структурированную программу обучения: цели, материалы, упражнения.', 5, 'rare', 'intelligence', 120, 60),
    ('Поиск учеников', 'Найдите минимум 3 человека, готовых обучаться у вас.', 4, 'rare', 'charisma', 100, 50),
    ('Проведение занятий (21 день)', 'Проводите регулярные занятия: минимум 2 раза в неделю, минимум 1 час каждое.', 7, 'epic', 'charisma', 400, 200),
    ('Адаптация под учеников', 'Адаптируйте программу под потребности каждого ученика. Получайте обратную связь.', 5, 'rare', 'charisma', 120, 60),
    ('Создание сообщества', 'Создайте сообщество учеников: группа, чат, регулярные встречи.', 5, 'epic', 'charisma', 150, 75),
    ('Развитие лидерских навыков', 'Примените лидерские техники: мотивация, делегирование, обратная связь.', 6, 'epic', 'charisma', 140, 70),
    ('Измерение результатов', 'Оцените результаты обучения: прогресс учеников, их достижения.', 4, 'rare', 'charisma', 80, 40),
    ('Итоговая оценка', 'Проанализируйте свой рост как лидера и наставника. Запишите ключевые уроки.', 3, 'rare', 'charisma', 60, 30)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- MONEY: Продвинутый квест (требует: Оптимизация финансов + любой промежуточный квест из другой категории)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Основы инвестирования', 
        'Эпический квест по изучению и началу инвестирования. Требует завершения квеста по оптимизации финансов и базовых знаний.', 
        'money', 
        'epic', 
        6, 
        280, 
        8,
        1100, 
        550, 
        504,  -- 21 день
        true,
        '[{"type": "stat_boost", "stat_name": "intelligence_level", "boost_amount": 1, "description": "+1 к уровню интеллекта"}, {"type": "reward_multiplier", "category": "money", "multiplier": 1.2, "description": "+20% к наградам за квесты денег"}]'::jsonb,
        3,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Изучение основ инвестирования', 'Изучите основы: акции, облигации, ETF, индексные фонды, диверсификация.', 4, 'common', 'money', 100, 50),
    ('Понимание рисков', 'Изучите типы рисков и определите свою толерантность к риску.', 3, 'common', 'money', 80, 40),
    ('Инвестиционные цели', 'Определите цели: срок, сумма, цель инвестирования. Запишите конкретные цифры.', 3, 'common', 'money', 70, 35),
    ('Выбор стратегии', 'Выберите инвестиционную стратегию: пассивное или активное инвестирование.', 4, 'rare', 'money', 100, 50),
    ('Выбор брокера', 'Исследуйте и выберите брокерскую платформу. Откройте счет (можно демо).', 3, 'common', 'money', 80, 40),
    ('Создание портфеля', 'Создайте диверсифицированный портфель: 60-80% акции (ETF), 20-40% облигации.', 5, 'rare', 'money', 150, 75),
    ('Первая инвестиция', 'Сделайте первую реальную инвестицию: минимум 1000 единиц валюты в выбранный ETF.', 6, 'epic', 'money', 200, 100),
    ('Настройка автоматических инвестиций', 'Настройте автоматическое ежемесячное инвестирование: минимум 10% от дохода.', 4, 'rare', 'money', 120, 60)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- ============================================
-- УРОВЕНЬ 4: ЛЕГЕНДАРНЫЕ КВЕСТЫ (требуют несколько продвинутых)
-- ============================================

-- LEGENDARY: Квест требует 2+ продвинутых квеста
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours, is_sequential, bonus_json, quest_level, max_level
    ) VALUES (
        'Путь к совершенству: интеграция всех сфер', 
        'Легендарный квест, объединяющий все сферы развития: здоровье, интеллект, харизма, деньги, сила воли. Требует завершения продвинутых квестов во всех категориях.', 
        'willpower', 
        'legendary', 
        10, 
        500, 
        10,
        2500, 
        1250, 
        1008,  -- 42 дня
        true,
        '[{"type": "stat_boost", "stat_name": "health_level", "boost_amount": 3, "description": "+3 к уровню здоровья"}, {"type": "stat_boost", "stat_name": "intelligence_level", "boost_amount": 3, "description": "+3 к уровню интеллекта"}, {"type": "stat_boost", "stat_name": "charisma_level", "boost_amount": 3, "description": "+3 к уровню харизмы"}, {"type": "stat_boost", "stat_name": "willpower_level", "boost_amount": 3, "description": "+3 к уровню силы воли"}]'::jsonb,
        4,
        1
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Создание интегральной системы', 'Создайте систему, объединяющую все сферы: здоровье, интеллект, харизма, деньги, сила воли.', 6, 'epic', 'willpower', 200, 100),
    ('Планирование на 42 дня', 'Составьте детальный план на 42 дня, интегрирующий все сферы развития.', 6, 'epic', 'willpower', 250, 125),
    ('Ежедневная практика (42 дня)', 'Следуйте системе каждый день: все сферы развития должны быть задействованы.', 10, 'legendary', 'willpower', 1200, 600),
    ('Мониторинг всех показателей', 'Отслеживайте показатели во всех сферах: физические, когнитивные, социальные, финансовые.', 5, 'epic', 'willpower', 150, 75),
    ('Еженедельная оптимизация', 'Еженедельно анализируйте и оптимизируйте систему на основе результатов.', 6, 'epic', 'willpower', 180, 90),
    ('Создание личной философии', 'Сформулируйте свою философию развития: принципы, ценности, цели.', 5, 'epic', 'willpower', 150, 75),
    ('Обучение и наставничество', 'Поделитесь опытом: обучите других своей системе или станьте наставником.', 7, 'legendary', 'charisma', 200, 100),
    ('Создание сообщества', 'Создайте или присоединитесь к сообществу единомышленников.', 6, 'epic', 'charisma', 180, 90),
    ('Измерение трансформации', 'Оцените трансформацию: сравните себя до и после. Запишите все изменения.', 5, 'epic', 'willpower', 150, 75),
    ('Создание наследия', 'Создайте что-то долгосрочное: книга, курс, сообщество, проект, который поможет другим.', 8, 'legendary', 'willpower', 250, 125)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order FROM numbered_tasks;

-- ============================================
-- СВЯЗЫВАНИЕ КВЕСТОВ С ВЕТКАМИ РАЗВИТИЯ
-- ============================================

WITH quest_ids AS (
    SELECT id, title FROM quests WHERE title IN (
        'Утренняя зарядка: первые шаги',
        'Здоровое питание: основы',
        'Чтение: 30 минут в день',
        'Изучение нового навыка: основы',
        'Ежедневные комплименты',
        'Развитие уверенности: первые шаги',
        'Отслеживание расходов',
        'Создание бюджета',
        'Формирование утреннего ритуала',
        'Интегральный подход к здоровью',
        'Глубокое обучение: от теории к практике',
        'Мастер социальных взаимодействий',
        'Оптимизация финансов',
        'Здоровый ум в здоровом теле',
        'Мастер здорового образа жизни',
        'Лидер и наставник',
        'Основы инвестирования',
        'Путь к совершенству: интеграция всех сфер'
    )
),
branch_ids AS (
    SELECT id, name FROM development_branches
)
INSERT INTO quest_branches (quest_id, branch_id, weight)
SELECT 
    qi.id as quest_id,
    bi.id as branch_id,
    CASE 
        -- HEALTH базовые
        WHEN qi.title = 'Утренняя зарядка: первые шаги' AND bi.name = 'health' THEN 1.0
        WHEN qi.title = 'Утренняя зарядка: первые шаги' AND bi.name = 'health_physical' THEN 0.8
        WHEN qi.title = 'Здоровое питание: основы' AND bi.name = 'health' THEN 1.0
        WHEN qi.title = 'Здоровое питание: основы' AND bi.name = 'health_physical' THEN 0.7
        -- INTELLIGENCE базовые
        WHEN qi.title = 'Чтение: 30 минут в день' AND bi.name = 'intelligence' THEN 1.0
        WHEN qi.title = 'Чтение: 30 минут в день' AND bi.name = 'intelligence_education' THEN 0.8
        WHEN qi.title = 'Изучение нового навыка: основы' AND bi.name = 'intelligence' THEN 1.0
        WHEN qi.title = 'Изучение нового навыка: основы' AND bi.name = 'intelligence_education' THEN 0.7
        WHEN qi.title = 'Изучение нового навыка: основы' AND bi.name = 'intelligence_cognitive' THEN 0.6
        -- CHARISMA базовые
        WHEN qi.title = 'Ежедневные комплименты' AND bi.name = 'charisma' THEN 1.0
        WHEN qi.title = 'Ежедневные комплименты' AND bi.name = 'charisma_social' THEN 0.8
        WHEN qi.title = 'Развитие уверенности: первые шаги' AND bi.name = 'charisma' THEN 1.0
        WHEN qi.title = 'Развитие уверенности: первые шаги' AND bi.name = 'charisma_confidence' THEN 0.9
        -- MONEY базовые
        WHEN qi.title = 'Отслеживание расходов' AND bi.name = 'money' THEN 1.0
        WHEN qi.title = 'Отслеживание расходов' AND bi.name = 'money_literacy' THEN 0.9
        WHEN qi.title = 'Создание бюджета' AND bi.name = 'money' THEN 1.0
        WHEN qi.title = 'Создание бюджета' AND bi.name = 'money_literacy' THEN 0.9
        -- WILLPOWER базовые
        WHEN qi.title = 'Формирование утреннего ритуала' AND bi.name = 'willpower' THEN 1.0
        -- HEALTH промежуточные
        WHEN qi.title = 'Интегральный подход к здоровью' AND bi.name = 'health' THEN 1.0
        WHEN qi.title = 'Интегральный подход к здоровью' AND bi.name = 'health_physical' THEN 0.9
        WHEN qi.title = 'Здоровый ум в здоровом теле' AND bi.name = 'health' THEN 0.8
        WHEN qi.title = 'Здоровый ум в здоровом теле' AND bi.name = 'intelligence' THEN 0.8
        WHEN qi.title = 'Здоровый ум в здоровом теле' AND bi.name = 'health_physical' THEN 0.6
        WHEN qi.title = 'Здоровый ум в здоровом теле' AND bi.name = 'intelligence_cognitive' THEN 0.6
        -- INTELLIGENCE промежуточные
        WHEN qi.title = 'Глубокое обучение: от теории к практике' AND bi.name = 'intelligence' THEN 1.0
        WHEN qi.title = 'Глубокое обучение: от теории к практике' AND bi.name = 'intelligence_education' THEN 0.9
        WHEN qi.title = 'Глубокое обучение: от теории к практике' AND bi.name = 'intelligence_cognitive' THEN 0.7
        -- CHARISMA промежуточные
        WHEN qi.title = 'Мастер социальных взаимодействий' AND bi.name = 'charisma' THEN 1.0
        WHEN qi.title = 'Мастер социальных взаимодействий' AND bi.name = 'charisma_social' THEN 0.9
        WHEN qi.title = 'Мастер социальных взаимодействий' AND bi.name = 'charisma_confidence' THEN 0.8
        WHEN qi.title = 'Мастер социальных взаимодействий' AND bi.name = 'charisma_leadership' THEN 0.7
        -- MONEY промежуточные
        WHEN qi.title = 'Оптимизация финансов' AND bi.name = 'money' THEN 1.0
        WHEN qi.title = 'Оптимизация финансов' AND bi.name = 'money_literacy' THEN 0.9
        -- HEALTH продвинутые
        WHEN qi.title = 'Мастер здорового образа жизни' AND bi.name = 'health' THEN 1.0
        WHEN qi.title = 'Мастер здорового образа жизни' AND bi.name = 'health_physical' THEN 0.9
        WHEN qi.title = 'Мастер здорового образа жизни' AND bi.name = 'health_mental' THEN 0.7
        -- CHARISMA продвинутые
        WHEN qi.title = 'Лидер и наставник' AND bi.name = 'charisma' THEN 1.0
        WHEN qi.title = 'Лидер и наставник' AND bi.name = 'charisma_leadership' THEN 0.9
        WHEN qi.title = 'Лидер и наставник' AND bi.name = 'intelligence' THEN 0.8
        WHEN qi.title = 'Лидер и наставник' AND bi.name = 'intelligence_education' THEN 0.7
        -- MONEY продвинутые
        WHEN qi.title = 'Основы инвестирования' AND bi.name = 'money' THEN 1.0
        WHEN qi.title = 'Основы инвестирования' AND bi.name = 'money_investing' THEN 0.9
        WHEN qi.title = 'Основы инвестирования' AND bi.name = 'intelligence' THEN 0.6
        WHEN qi.title = 'Основы инвестирования' AND bi.name = 'intelligence_critical' THEN 0.5
        -- LEGENDARY
        WHEN qi.title = 'Путь к совершенству: интеграция всех сфер' AND bi.name = 'willpower' THEN 1.0
        WHEN qi.title = 'Путь к совершенству: интеграция всех сфер' AND bi.name = 'health' THEN 0.8
        WHEN qi.title = 'Путь к совершенству: интеграция всех сфер' AND bi.name = 'intelligence' THEN 0.8
        WHEN qi.title = 'Путь к совершенству: интеграция всех сфер' AND bi.name = 'charisma' THEN 0.8
        WHEN qi.title = 'Путь к совершенству: интеграция всех сфер' AND bi.name = 'money' THEN 0.8
        ELSE NULL
    END as weight
FROM quest_ids qi
CROSS JOIN branch_ids bi
WHERE CASE 
    WHEN qi.title = 'Утренняя зарядка: первые шаги' AND bi.name IN ('health', 'health_physical') THEN true
    WHEN qi.title = 'Здоровое питание: основы' AND bi.name IN ('health', 'health_physical') THEN true
    WHEN qi.title = 'Чтение: 30 минут в день' AND bi.name IN ('intelligence', 'intelligence_education') THEN true
    WHEN qi.title = 'Изучение нового навыка: основы' AND bi.name IN ('intelligence', 'intelligence_education', 'intelligence_cognitive') THEN true
    WHEN qi.title = 'Ежедневные комплименты' AND bi.name IN ('charisma', 'charisma_social') THEN true
    WHEN qi.title = 'Развитие уверенности: первые шаги' AND bi.name IN ('charisma', 'charisma_confidence') THEN true
    WHEN qi.title = 'Отслеживание расходов' AND bi.name IN ('money', 'money_literacy') THEN true
    WHEN qi.title = 'Создание бюджета' AND bi.name IN ('money', 'money_literacy') THEN true
    WHEN qi.title = 'Формирование утреннего ритуала' AND bi.name = 'willpower' THEN true
    WHEN qi.title = 'Интегральный подход к здоровью' AND bi.name IN ('health', 'health_physical') THEN true
    WHEN qi.title = 'Здоровый ум в здоровом теле' AND bi.name IN ('health', 'intelligence', 'health_physical', 'intelligence_cognitive') THEN true
    WHEN qi.title = 'Глубокое обучение: от теории к практике' AND bi.name IN ('intelligence', 'intelligence_education', 'intelligence_cognitive') THEN true
    WHEN qi.title = 'Мастер социальных взаимодействий' AND bi.name IN ('charisma', 'charisma_social', 'charisma_confidence', 'charisma_leadership') THEN true
    WHEN qi.title = 'Оптимизация финансов' AND bi.name IN ('money', 'money_literacy') THEN true
    WHEN qi.title = 'Мастер здорового образа жизни' AND bi.name IN ('health', 'health_physical', 'health_mental') THEN true
    WHEN qi.title = 'Лидер и наставник' AND bi.name IN ('charisma', 'charisma_leadership', 'intelligence', 'intelligence_education') THEN true
    WHEN qi.title = 'Основы инвестирования' AND bi.name IN ('money', 'money_investing', 'intelligence', 'intelligence_critical') THEN true
    WHEN qi.title = 'Путь к совершенству: интеграция всех сфер' AND bi.name IN ('willpower', 'health', 'intelligence', 'charisma', 'money') THEN true
    ELSE false
END;

-- ============================================
-- СОЗДАНИЕ ЗАВИСИМОСТЕЙ МЕЖДУ КВЕСТАМИ
-- ============================================
-- Система: required_count определяет, сколько из группы prerequisites нужно выполнить
-- Например: required_count = 2 означает, что нужно выполнить 2 из перечисленных квестов

WITH quest_map AS (
    SELECT id, title FROM quests WHERE title IN (
        'Утренняя зарядка: первые шаги',
        'Здоровое питание: основы',
        'Чтение: 30 минут в день',
        'Изучение нового навыка: основы',
        'Ежедневные комплименты',
        'Развитие уверенности: первые шаги',
        'Отслеживание расходов',
        'Создание бюджета',
        'Формирование утреннего ритуала',
        'Интегральный подход к здоровью',
        'Глубокое обучение: от теории к практике',
        'Мастер социальных взаимодействий',
        'Оптимизация финансов',
        'Здоровый ум в здоровом теле',
        'Мастер здорового образа жизни',
        'Лидер и наставник',
        'Основы инвестирования',
        'Путь к совершенству: интеграция всех сфер'
    )
)
INSERT INTO quest_prerequisites (quest_id, prerequisite_quest_id, required_count)
SELECT 
    q2.id as quest_id,
    q1.id as prerequisite_quest_id,
    CASE
        -- Промежуточные квесты: требуют ОБА базовых (required_count = 2)
        WHEN q2.title = 'Интегральный подход к здоровью' THEN 2
        WHEN q2.title = 'Глубокое обучение: от теории к практике' THEN 2
        WHEN q2.title = 'Мастер социальных взаимодействий' THEN 2
        WHEN q2.title = 'Оптимизация финансов' THEN 2
        -- Кросс-категорийный: требует 2 из 4 базовых (health + intelligence)
        WHEN q2.title = 'Здоровый ум в здоровом теле' THEN 2
        -- Продвинутые квесты: требуют ОБА промежуточных (required_count = 2)
        WHEN q2.title = 'Мастер здорового образа жизни' THEN 2
        WHEN q2.title = 'Лидер и наставник' THEN 2
        -- Продвинутый квест: требует 1 промежуточный
        WHEN q2.title = 'Основы инвестирования' THEN 1
        -- Легендарный квест: требует ВСЕ 3 продвинутых (required_count = 3)
        WHEN q2.title = 'Путь к совершенству: интеграция всех сфер' THEN 3
        ELSE 1
    END as required_count
FROM quest_map q1
JOIN quest_map q2 ON (
    -- Промежуточные квесты требуют ОБА базовых
    (q2.title = 'Интегральный подход к здоровью' AND q1.title IN ('Утренняя зарядка: первые шаги', 'Здоровое питание: основы'))
    OR
    (q2.title = 'Глубокое обучение: от теории к практике' AND q1.title IN ('Чтение: 30 минут в день', 'Изучение нового навыка: основы'))
    OR
    (q2.title = 'Мастер социальных взаимодействий' AND q1.title IN ('Ежедневные комплименты', 'Развитие уверенности: первые шаги'))
    OR
    (q2.title = 'Оптимизация финансов' AND q1.title IN ('Отслеживание расходов', 'Создание бюджета'))
    OR
    -- Кросс-категорийный промежуточный квест: требует 2 из 4 (health + intelligence)
    (q2.title = 'Здоровый ум в здоровом теле' AND q1.title IN ('Утренняя зарядка: первые шаги', 'Здоровое питание: основы', 'Чтение: 30 минут в день', 'Изучение нового навыка: основы'))
    OR
    -- Продвинутые квесты требуют ОБА промежуточных
    (q2.title = 'Мастер здорового образа жизни' AND q1.title IN ('Интегральный подход к здоровью', 'Здоровый ум в здоровом теле'))
    OR
    (q2.title = 'Лидер и наставник' AND q1.title IN ('Мастер социальных взаимодействий', 'Глубокое обучение: от теории к практике'))
    OR
    -- Продвинутый квест: требует 1 промежуточный
    (q2.title = 'Основы инвестирования' AND q1.title = 'Оптимизация финансов')
    OR
    -- Легендарный квест требует ВСЕ 3 продвинутых
    (q2.title = 'Путь к совершенству: интеграция всех сфер' AND q1.title IN ('Мастер здорового образа жизни', 'Лидер и наставник', 'Основы инвестирования'))
);

-- ============================================
-- ДОБАВЛЕНИЕ HABIT REQUIREMENTS ДЛЯ ЗАДАЧ
-- ============================================

-- Базовые квесты: ежедневные задачи (7-14 дней)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 7, 'any'
FROM tasks
WHERE title IN ('Практика комплиментов (7 дней)', 'Ежедневное отслеживание (14 дней)', 'Практика ритуала (14 дней)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Базовые квесты: ежедневные задачи (10-14 дней)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 10, 'any'
FROM tasks
WHERE title IN ('Ежедневная практика (10 дней)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Базовые квесты: ежедневное чтение (14 дней)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 14, 'any'
FROM tasks
WHERE title IN ('Ежедневное чтение (14 дней)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Промежуточные квесты: более длительные привычки (14-21 день)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 14, 'any'
FROM tasks
WHERE title IN ('Интеграция (14 дней)', 'Разговоры с незнакомцами (14 дней)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Промежуточные квесты: 21 день
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 21, 'any'
FROM tasks
WHERE title IN ('Ежедневная практика (21 день)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Продвинутые квесты: длительные привычки (21-30 дней)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 21, 'any'
FROM tasks
WHERE title IN ('Проведение занятий (21 день)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Продвинутые квесты: очень длительные привычки (30 дней)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 30, 'any'
FROM tasks
WHERE title IN ('Ежедневное выполнение (30 дней)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Легендарный квест: очень длительная привычка (42 дня)
INSERT INTO task_habit_requirements (task_id, consecutive_days, daytime_required)
SELECT id, 42, 'any'
FROM tasks
WHERE title IN ('Ежедневная практика (42 дня)')
AND NOT EXISTS (SELECT 1 FROM task_habit_requirements WHERE task_id = tasks.id);

-- Проверка: сколько задач получили habit requirements
SELECT 
    COUNT(*) as total_requirements,
    SUM(CASE WHEN daytime_required = 'morning' THEN 1 ELSE 0 END) as morning_tasks,
    SUM(CASE WHEN daytime_required = 'any' OR daytime_required IS NULL THEN 1 ELSE 0 END) as any_time_tasks
FROM task_habit_requirements;
