package handlers

import (
	"BecomeOverMan/internal/services"
	"BecomeOverMan/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenerateMotivationHandler - генерирует мотивационное сообщение
func (h *QuestHandler) GenerateMotivationHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем прогресс пользователя
	progress, err := h.questService.GetUserProgress(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	llmService := services.NewEnhancedLLMService()
	motivation, err := llmService.GenerateMotivation(c.Request.Context(), progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate motivation: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"motivation": motivation,
	})
}

// AnalyzeProgressHandler - анализирует прогресс пользователя
func (h *QuestHandler) AnalyzeProgressHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем данные пользователя
	userData, err := h.questService.GetUserProgressData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	llmService := services.NewEnhancedLLMService()
	analysis, err := llmService.AnalyzeProgress(c.Request.Context(), userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze progress: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// ImproveQuestDescriptionHandler - улучшает описание квеста
func (h *QuestHandler) ImproveQuestDescriptionHandler(c *gin.Context) {
	var req struct {
		QuestID int `json:"quest_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	quest, err := h.questService.GetQuestDetails(c.Request.Context(), req.QuestID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	llmService := services.NewEnhancedLLMService()
	improvedDescription, err := llmService.ImproveQuestDescription(
		c.Request.Context(),
		quest.Title,
		quest.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to improve description: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"improved_description": improvedDescription,
	})
}

// GeneratePersonalizedQuestHandler - генерирует персонализированный квест
func (h *QuestHandler) GeneratePersonalizedQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		UserMessage string `json:"user_message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Получаем историю пользователя
	userHistory, err := h.questService.GetUserHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	llmService := services.NewEnhancedLLMService()
	quest, err := llmService.GeneratePersonalizedQuest(c.Request.Context(), userHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate quest: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, quest)
}
