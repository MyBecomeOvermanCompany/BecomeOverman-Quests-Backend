package services

import (
	"BecomeOverMan/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LLMProvider - тип провайдера LLM
type LLMProvider string

const (
	ProviderOpenAI  LLMProvider = "openai"
	ProviderClaude  LLMProvider = "claude"
	ProviderKimi    LLMProvider = "kimi"
	ProviderDefault LLMProvider = "default"
)

// EnhancedLLMService - улучшенный сервис для работы с LLM
type EnhancedLLMService struct {
	provider LLMProvider
	apiKey   string
	baseURL  string
}

// NewEnhancedLLMService создает новый сервис LLM
func NewEnhancedLLMService() *EnhancedLLMService {
	provider := ProviderDefault
	if config.Cfg.OpenAIAPIKey != "" {
		provider = ProviderOpenAI
	} else if config.Cfg.ClaudeAPIKey != "" {
		provider = ProviderClaude
	} else if config.Cfg.APIKeyIntelligenceIO != "" {
		provider = ProviderKimi
	}

	return &EnhancedLLMService{
		provider: provider,
		apiKey:   getAPIKey(provider),
		baseURL:  getBaseURL(provider),
	}
}

func getAPIKey(provider LLMProvider) string {
	switch provider {
	case ProviderOpenAI:
		return config.Cfg.OpenAIAPIKey
	case ProviderClaude:
		return config.Cfg.ClaudeAPIKey
	case ProviderKimi:
		return config.Cfg.APIKeyIntelligenceIO
	default:
		return config.Cfg.APIKeyIntelligenceIO
	}
}

func getBaseURL(provider LLMProvider) string {
	switch provider {
	case ProviderOpenAI:
		return "https://api.openai.com/v1/chat/completions"
	case ProviderClaude:
		return "https://api.anthropic.com/v1/messages"
	case ProviderKimi:
		return "https://api.intelligence.io.solutions/api/v1/chat/completions"
	default:
		return "https://api.intelligence.io.solutions/api/v1/chat/completions"
	}
}

// GenerateMotivation - генерирует мотивационное сообщение для пользователя
func (s *EnhancedLLMService) GenerateMotivation(ctx context.Context, userProgress map[string]interface{}) (string, error) {
	systemPrompt := `Ты - мотивационный коуч, который помогает людям достигать целей в системе саморазвития.
Твоя задача - создать короткое, вдохновляющее и персонализированное мотивационное сообщение на основе прогресса пользователя.
Сообщение должно быть:
- Позитивным и поддерживающим
- Конкретным (упоминать достижения пользователя)
- Мотивирующим к дальнейшим действиям
- Коротким (2-3 предложения)
Ответь только текстом сообщения, без дополнительных комментариев.`

	userMessage := fmt.Sprintf(`Прогресс пользователя:
- Уровень: %v
- Опыт: %v XP
- Завершенных квестов: %v
- Активных квестов: %v
- Текущий streak: %v дней

Создай мотивационное сообщение для этого пользователя.`,
		userProgress["level"],
		userProgress["xp"],
		userProgress["completed_quests"],
		userProgress["active_quests"],
		userProgress["streak"],
	)

	return s.callLLM(ctx, systemPrompt, userMessage, "gpt-4o-mini")
}

// AnalyzeProgress - анализирует прогресс пользователя и дает рекомендации
func (s *EnhancedLLMService) AnalyzeProgress(ctx context.Context, userData map[string]interface{}) (map[string]interface{}, error) {
	systemPrompt := `Ты - эксперт по анализу прогресса в системах саморазвития.
Твоя задача - проанализировать данные пользователя и дать конкретные, полезные рекомендации.
Ответь в формате JSON:
{
  "strengths": ["сильная сторона 1", "сильная сторона 2"],
  "weaknesses": ["слабая сторона 1", "слабая сторона 2"],
  "recommendations": ["рекомендация 1", "рекомендация 2"],
  "next_steps": ["шаг 1", "шаг 2"]
}`

	userMessage := fmt.Sprintf(`Данные пользователя:
%+v

Проанализируй прогресс и дай рекомендации.`, userData)

	response, err := s.callLLM(ctx, systemPrompt, userMessage, "gpt-4o-mini")
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Если не JSON, возвращаем как текст
		return map[string]interface{}{
			"analysis": response,
		}, nil
	}

	return result, nil
}

// ImproveQuestDescription - улучшает описание квеста
func (s *EnhancedLLMService) ImproveQuestDescription(ctx context.Context, questTitle, questDescription string) (string, error) {
	systemPrompt := `Ты - эксперт по написанию описаний для квестов саморазвития.
Твоя задача - улучшить описание квеста, сделав его:
- Более мотивирующим и вдохновляющим
- Конкретным и понятным
- Научно обоснованным (если возможно)
- Структурированным
Ответь только улучшенным описанием, без дополнительных комментариев.`

	userMessage := fmt.Sprintf(`Улучши описание квеста:

Название: %s
Текущее описание: %s`, questTitle, questDescription)

	return s.callLLM(ctx, systemPrompt, userMessage, "gpt-4o-mini")
}

// GeneratePersonalizedQuest - генерирует персонализированный квест на основе истории пользователя
func (s *EnhancedLLMService) GeneratePersonalizedQuest(ctx context.Context, userHistory map[string]interface{}) (map[string]interface{}, error) {
	systemPrompt := `Ты - эксперт по созданию квестов саморазвития.
На основе истории пользователя создай персонализированный квест, который:
- Учитывает предыдущие достижения
- Подходит к текущему уровню
- Реалистичен и выполним
- Мотивирует к росту

Ответь в формате JSON:
{
  "title": "Название квеста",
  "description": "Подробное описание",
  "difficulty": 5,
  "category": "health",
  "tasks": [
    {"title": "Задача 1", "description": "Описание задачи 1"},
    {"title": "Задача 2", "description": "Описание задачи 2"}
  ]
}`

	userMessage := fmt.Sprintf(`История пользователя:
%+v

Создай персонализированный квест.`, userHistory)

	response, err := s.callLLM(ctx, systemPrompt, userMessage, "gpt-4o-mini")
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return result, nil
}

// callLLM - универсальный метод для вызова LLM
func (s *EnhancedLLMService) callLLM(ctx context.Context, systemPrompt, userMessage, model string) (string, error) {
	switch s.provider {
	case ProviderOpenAI:
		return s.callOpenAI(ctx, systemPrompt, userMessage, model)
	case ProviderClaude:
		return s.callClaude(ctx, systemPrompt, userMessage, model)
	case ProviderKimi, ProviderDefault:
		result, err := requestAI(userMessage, systemPrompt, model)
		if err != nil {
			return "", err
		}
		return string(result), nil
	default:
		result, err := requestAI(userMessage, systemPrompt, model)
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
}

// callOpenAI - вызов OpenAI API
func (s *EnhancedLLMService) callOpenAI(ctx context.Context, systemPrompt, userMessage, model string) (string, error) {
	if model == "" {
		model = "gpt-4o-mini"
	}

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userMessage},
		},
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// callClaude - вызов Claude API
func (s *EnhancedLLMService) callClaude(ctx context.Context, systemPrompt, userMessage, model string) (string, error) {
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	reqBody := map[string]interface{}{
		"model": model,
		"max_tokens": 1024,
		"system": systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userMessage},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Claude API error: %s", string(body))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return strings.TrimSpace(result.Content[0].Text), nil
}
