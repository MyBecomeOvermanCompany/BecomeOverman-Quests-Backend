package services

import (
	"BecomeOverMan/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// GenerateAIQuest - чистая Go реализация без Python
func GenerateAIQuest(userMessage string) (*models.AIQuestResponse, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not found in environment variables")
	}

	aiModel := "moonshotai/Kimi-K2-Thinking"
	url := "https://api.intelligence.io.solutions/api/v1/chat/completions"

	systemPrompt := `
	Ты помощник для генерации квестов в формате строгого JSON. 
	ВОЗВРАЩАЙ ТОЛЬКО JSON БЕЗ ЛЮБЫХ ДОПОЛНИТЕЛЬНЫХ ТЕКСТОВ И КОММЕНТАРИЕВ!

	Структура JSON должна быть такой:
	{
		"quest": {
			"title": "Название квеста",
			"description": "Описание квеста [GENERATED]",
			"category": "health/willpower/intelligence/creativity/social",
			"rarity": "common/rare/epic/legendary",
			"difficulty": 1-5,
			"price": 10-100,
			"tasks_count": 3-7,
			"reward_xp": 50-500,
			"reward_coin": 25-250,
			"time_limit_hours": 24-336
		},
		"tasks": [
			{
				"title": "Название задачи 1",
				"description": "Описание задачи 1",
				"difficulty": 1-3,
				"rarity": "common/rare/epic",
				"category": "health/willpower/intelligence/creativity/social",
				"base_xp_reward": 10-50,
				"base_coin_reward": 5-25,
				"task_order": 1
			}
		]
	}

	Правила:
	- difficulty квеста должен быть средним от difficulty задач
	- price = reward_coin * 1.5 (округлить)
	- tasks_count должно соответствовать количеству задач в массиве
	- time_limit_hours: 24-168 (1-7 дней)
	- reward_xp = сумма base_xp_reward всех задач * 1.5
	- reward_coin = сумма base_coin_reward всех задач * 1.5
	`

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

	client := &http.Client{}
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

	// Парсим финальный JSON
	var aiResponse models.AIQuestResponse
	err = json.Unmarshal([]byte(content), &aiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing AI quest response: %v", err)
	}

	return &aiResponse, nil
}
