package services

import (
	"BecomeOverMan/internal/config"
	"BecomeOverMan/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	url = "https://api.intelligence.io.solutions/api/v1/chat/completions"
)

var apiKey = config.Cfg.APIKeyIntelligenceIO

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type Choice struct {
	Message ChatMessage `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

func requestAI(userMessage, systemPrompt, aiModel string) ([]byte, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not found in environment variables")
	}

	if aiModel == "" {
		aiModel = "moonshotai/Kimi-K2-Thinking"
	}

	requestData := ChatRequest{
		Model: aiModel,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: 60 * time.Second, // Увеличиваем таймаут для более сложных запросов
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", string(body))
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing chat response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// Очищаем ответ от thinking тегов
	content := chatResponse.Choices[0].Message.Content
	if idx := strings.Index(content, "</think>\n\n"); idx != -1 {
		content = content[idx+11:] // +11 чтобы пропустить "</think>\n\n"
	}

	return []byte(content), nil
}

// GenerateAIQuest - Улучшенная генерация квеста с научной обоснованностью
func (s *QuestService) GenerateAIQuest(userMessage string) (*models.AIQuestResponse, error) {
	aiModel := "moonshotai/Kimi-K2-Thinking"

	systemPrompt := `Ты эксперт по саморазвитию, психологии, медицине и бизнесу. Твоя задача - создавать реалистичные, научно обоснованные квесты для личностного роста.

КРИТИЧЕСКИ ВАЖНО: ВОЗВРАЩАЙ ТОЛЬКО ВАЛИДНЫЙ JSON БЕЗ ЛЮБЫХ ДОПОЛНИТЕЛЬНЫХ ТЕКСТОВ, КОММЕНТАРИЕВ ИЛИ МАРКДАУН!

Структура JSON:
{
	"quest": {
		"title": "Конкретное, мотивирующее название (макс 60 символов)",
		"description": "ОБЯЗАТЕЛЬНО ВКЛЮЧИ: 1) Детальное описание (минимум 300 слов), 2) ССЫЛКИ НА РЕАЛЬНЫЕ ИССЛЕДОВАНИЯ (название исследования, авторы, год, краткое описание результатов), 3) ПРИМЕРЫ РЕАЛЬНЫХ ЛЮДЕЙ (имя или псевдоним, что они делали, какие результаты получили), 4) Конкретные цифры и метрики успеха, 5) Пошаговый план действий. Формат: [ИССЛЕДОВАНИЕ: название, авторы, год, результаты] [ПРИМЕР: имя, история, результаты]",
		"category": "health/willpower/intelligence/charisma/money",
		"rarity": "common/rare/epic/legendary",
		"difficulty": 1-10,
		"price": 10-500,
		"tasks_count": 3-10,
		"reward_xp": 50-2000,
		"reward_coin": 25-1000,
		"time_limit_hours": 24-720,
		"bonus_json": {
			"type": "reward_multiplier",
			"category": "health",
			"multiplier": 1.1,
			"description": "Пассивный бафф: +10% к наградам за квесты здоровья"
		}
	},
	"tasks": [
		{
			"title": "Конкретная, измеримая задача",
			"description": "Детальное описание с четкими критериями выполнения, дедлайнами, измеримыми результатами. Включи научное обоснование почему эта задача эффективна.",
			"difficulty": 1-5,
			"rarity": "common/rare/epic",
			"category": "health/willpower/intelligence/charisma/money",
			"base_xp_reward": 10-200,
			"base_coin_reward": 5-100,
			"task_order": 1
		}
	]
}

ПРАВИЛА ГЕНЕРАЦИИ:

1. НАУЧНАЯ ОБОСНОВАННОСТЬ (ОБЯЗАТЕЛЬНО):
   - КАЖДЫЙ квест ДОЛЖЕН содержать минимум 2-3 ссылки на реальные исследования
   - Формат: [ИССЛЕДОВАНИЕ: "Название исследования" (Авторы, Год) - краткое описание результатов]
   - Примеры: "Исследование Lally et al. (2010) показало, что формирование привычки занимает в среднем 66 дней"
   - Упоминай конкретные принципы и методологии (SMART goals, habit stacking, CBT, GTD и т.д.)
   - Задачи должны быть проверены на эффективность в реальных исследованиях

2. РЕАЛИСТИЧНОСТЬ:
   - Задачи должны быть выполнимыми в реальной жизни
   - Учитывай время, ресурсы, возможности обычного человека
   - Избегай нереалистичных требований

3. СТРУКТУРИРОВАННОСТЬ:
   - Четкие дедлайны для каждой задачи
   - Измеримые критерии выполнения
   - Логическая последовательность задач
   - Прогрессивное усложнение

4. МОТИВАЦИЯ И РЕАЛЬНЫЕ ПРИМЕРЫ (ОБЯЗАТЕЛЬНО):
   - КАЖДЫЙ квест ДОЛЖЕН содержать минимум 1-2 примера реальных людей
   - Формат: [ПРИМЕР: "Имя/псевдоним" - что делал, какие результаты получил, конкретные цифры]
   - Примеры: "Джон, 32 года, за 3 месяца ранних подъемов увеличил продуктивность на 40%"
   - Показывай конкретные результаты и бенефиты с цифрами
   - Связывай с долгосрочными целями
   - Избегай абстрактных формулировок - только конкретика

5. КАТЕГОРИИ И ВЕТКИ:
   - health: физическое здоровье, спорт, питание, сон, ментальное здоровье
   - charisma: социальные навыки, уверенность, лидерство, привлекательность
   - intelligence: когнитивные навыки, образование, креативность, критическое мышление
   - money: заработок, финансовая грамотность, инвестиции, бизнес
   - willpower: сила воли, дисциплина, привычки, самоконтроль

6. РАСЧЕТЫ:
   - difficulty квеста = среднее от difficulty задач (округлить)
   - price = reward_coin * 1.2-1.5 (округлить)
   - tasks_count = количество задач в массиве
   - time_limit_hours: 24-720 (1-30 дней) в зависимости от сложности
   - reward_xp = сумма base_xp_reward всех задач * 1.5-2.0
   - reward_coin = сумма base_coin_reward всех задач * 1.5-2.0

7. ПАССИВНЫЕ БАФФЫ:
   - Каждый квест должен иметь bonus_json с пассивным баффом
   - Типы баффов: reward_multiplier (увеличение наград), stat_boost (увеличение характеристик)
   - Бафф должен соответствовать категории квеста

8. КАЧЕСТВО ОПИСАНИЙ (КРИТИЧЕСКИ ВАЖНО):
   - description квеста: МИНИМУМ 300 слов, ОБЯЗАТЕЛЬНО включает:
     * Ссылки на реальные исследования (2-3 шт) с авторами и годами
     * Примеры реальных людей (1-2 шт) с конкретными результатами
     * Конкретные цифры и метрики (не "улучшит здоровье", а "снизит давление на 10-15 мм рт.ст.")
     * Пошаговый план действий
   - description задачи: минимум 80 слов, четкие критерии, измеримые результаты, ссылки на методологию
   - ИЗБЕГАЙ абстрактных формулировок типа "стать лучше", "улучшить жизнь" - только конкретика

ПРИМЕРЫ ХОРОШИХ КВЕСТОВ (с исследованиями и примерами):
- "21-дневный челлендж раннего подъема" 
  Описание должно включать: [ИССЛЕДОВАНИЕ: Lally et al. (2010) - формирование привычки], [ПРИМЕР: Тим Феррис - встает в 5 утра, написал 4 бестселлера]
- "Интервальное голодание 16/8 по протоколу доктора Фунга"
  Описание: [ИССЛЕДОВАНИЕ: Fung (2016) - метаболические преимущества], [ПРИМЕР: Хью Джекман - использовал для подготовки к ролям]
- "Метод Фейнмана: объясни как ребенку"
  Описание: [ИССЛЕДОВАНИЕ: Feynman (1985) - техника обучения], [ПРИМЕР: студенты MIT улучшили понимание на 40%]

ПРИМЕРЫ ПЛОХИХ КВЕСТОВ (НЕ ДЕЛАЙ ТАК):
- Слишком общие названия ("Стать лучше", "Улучшить здоровье")
- Без ссылок на исследования
- Без примеров реальных людей
- Абстрактные формулировки ("почувствуешь себя лучше")
- Слишком короткие описания без деталей

ВАЖНО: Всегда генерируй реалистичные, мотивирующие квесты, которые действительно помогут человеку изменить свою жизнь к лучшему.`

	answer, err := requestAI(userMessage, systemPrompt, aiModel)
	if err != nil {
		return nil, fmt.Errorf("error requesting AI: %v", err)
	}

	// Очищаем ответ от возможных markdown блоков
	answerStr := string(answer)
	answerStr = strings.TrimSpace(answerStr)

	// Удаляем markdown code blocks если есть
	if strings.HasPrefix(answerStr, "```json") {
		answerStr = strings.TrimPrefix(answerStr, "```json")
		answerStr = strings.TrimSuffix(answerStr, "```")
		answerStr = strings.TrimSpace(answerStr)
	} else if strings.HasPrefix(answerStr, "```") {
		answerStr = strings.TrimPrefix(answerStr, "```")
		answerStr = strings.TrimSuffix(answerStr, "```")
		answerStr = strings.TrimSpace(answerStr)
	}

	// Парсим финальный JSON
	var aiResponse models.AIQuestResponse
	err = json.Unmarshal([]byte(answerStr), &aiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing AI quest response: %v\nRaw response: %s", err, answerStr)
	}

	// Валидация
	if aiResponse.Quest == nil {
		return nil, fmt.Errorf("quest is nil in AI response")
	}
	if len(aiResponse.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks in AI response")
	}
	if aiResponse.Quest.TasksCount != len(aiResponse.Tasks) {
		aiResponse.Quest.TasksCount = len(aiResponse.Tasks)
	}

	return &aiResponse, nil
}

// -----------------------------------------------------------
// ----------- GenerateScheduleByAI --------------------------
// -----------------------------------------------------------

type ScheduleTask struct {
	TaskID         int        `json:"task_id"`
	ScheduledStart time.Time  `json:"scheduled_start"`
	ScheduledEnd   time.Time  `json:"scheduled_end"`
	Deadline       *time.Time `json:"deadline"`
	Duration       int        `json:"duration"` // в минутах
}

type AIScheduleResponse struct {
	Schedule []ScheduleTask `json:"schedule"`
}

func (s *QuestService) GenerateScheduleByAI(
	ctx context.Context,
	userID int,
	userMessage string,
) ([]models.Quest, error) {
	quests, err := s.GetMyAllQuestsWithDetails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user quests: %v", err)
	}

	info := ""
	for _, q := range quests {
		questStr := fmt.Sprintf("QuestID:%d Title:%s Description:%s ",
			q.ID,
			q.Title,
			q.Description,
		)

		tasksStr := "Tasks:["
		for _, t := range q.Tasks {
			if t.Status != nil && *t.Status != "active" {
				continue
			}

			taskStr := fmt.Sprintf(
				"ID:%d Title:%s Desc:%s Deadline:%v Duration:%v ScheduledStart:%v ScheduledEnd:%v",
				t.ID,
				t.Title,
				t.Description,
				t.Deadline,
				t.Duration,
				t.ScheduledStart,
				t.ScheduledEnd,
			)

			if t.TaskOrder >= 0 {
				taskStr = taskStr + fmt.Sprintf(" TaskOrder:%d; ", t.TaskOrder)
			} else {
				taskStr = taskStr + "; "
			}

			tasksStr += taskStr
		}
		tasksStr += "] "
		questStr += tasksStr

		info += questStr
	}

	userMessageWithInfo := userMessage + "\n\n" + info

	systemPrompt := `
	Ты помощник для генерации расписания задач.
	Верни ТОЛЬКО валидный JSON без комментариев.

	Формат ответа:
	{
	"schedule": [
		{
		"task_id": number,
		"scheduled_start": "RFC3339",
		"scheduled_end": "RFC3339",
		"deadline": "RFC3339|null",
		"duration": number
		}
	]
	}

	Правила:
	- duration только в минутах
	- если у задачи уже есть часть данных — дополни логично
	- распределяй задачи равномерно
	- учитывай taskOrder - задача с меньшим taskOrder должна быть раньше выполнена
	`

	currentDate := time.Now()

	userMessageWithInfo += fmt.Sprintf(`
		- Текущая дата: "%s"
	`, currentDate.Format(time.RFC3339))

	aiModel := "moonshotai/Kimi-K2-Thinking"

	answer, err := requestAI(userMessageWithInfo, systemPrompt, aiModel)
	if err != nil {
		return nil, err
	}

	answer = bytes.TrimSpace(answer)
	if len(answer) == 0 || answer[0] != '{' {
		return nil, fmt.Errorf("AI returned invalid JSON: %s", string(answer))
	}

	// ---------- парсим ----------
	var schedules AIScheduleResponse
	if err := json.Unmarshal(answer, &schedules); err != nil {
		return nil, fmt.Errorf("error parsing AI schedule response: %v\nraw: %s", err, string(answer))
	}

	// ---------- индексируем ----------
	scheduleMap := make(map[int]ScheduleTask)
	for _, sch := range schedules.Schedule {
		scheduleMap[sch.TaskID] = sch
	}

	// ---------- применяем к задачам ----------
	for qi := range quests {
		for ti := range quests[qi].Tasks {
			t := &quests[qi].Tasks[ti]

			sch, ok := scheduleMap[t.ID]
			if !ok {
				continue
			}

			t.ScheduledStart = &sch.ScheduledStart
			t.ScheduledEnd = &sch.ScheduledEnd
			t.Deadline = sch.Deadline
			t.Duration = &sch.Duration
		}
	}

	// ------ Сохраняем в БД -----------
	allTasks := []models.Task{}
	for _, q := range quests {
		allTasks = append(allTasks, q.Tasks...)
	}

	if err := s.questRepo.SetOrUpdateScheduleTasks(ctx, userID, allTasks); err != nil {
		return nil, fmt.Errorf("error updating schedule tasks: %v", err)
	}

	return quests, nil
}
