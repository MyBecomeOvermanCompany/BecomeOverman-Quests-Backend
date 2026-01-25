package handlers

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/services"
	"BecomeOverMan/pkg/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type QuestHandler struct {
	questService *services.QuestService
}

func NewQuestHandler(questService *services.QuestService) *QuestHandler {
	return &QuestHandler{questService: questService}
}

// ==== Handlers ====

// GetQuestDetails возвращает детали квеста с задачами
func (h *QuestHandler) GetQuestDetails(c *gin.Context) {
	questIDStr := c.Param("questID")
	questID, err := strconv.Atoi(questIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questDetails, err := h.questService.GetQuestDetails(c.Request.Context(), questID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questDetails)
}

// ------------ GetMyAllQuestsWithDetails ---------------

func (h *QuestHandler) GetMyAllQuestsWithDetails(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetMyAllQuestsWithDetails(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// ------------ GetAvailableQuestsHandler ---------

func (h *QuestHandler) GetAvailableQuestsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetAvailableQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// <--------------------------------------------->
// ------------ GetQuestShopHandler --------------
// <--------------------------------------------->

func (h *QuestHandler) GetQuestShopHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetQuestShop(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

func (h *QuestHandler) GetMyActiveQuestsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetMyActiveQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

func (h *QuestHandler) GetMyCompletedQuests(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.GetMyCompletedQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

func (h *QuestHandler) PurchaseQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	if err := h.questService.PurchaseQuest(c.Request.Context(), userID, questID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *QuestHandler) StartQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	if err := h.questService.StartQuest(c.Request.Context(), userID, questID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *QuestHandler) CompleteTaskHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.questService.CompleteTask(c.Request.Context(), userID, questID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *QuestHandler) CompleteQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	if err := h.questService.CompleteQuest(c.Request.Context(), userID, questID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// ContinueQuestHandler - начать следующий уровень квеста
func (h *QuestHandler) ContinueQuestHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	nextQuest, err := h.questService.ContinueQuest(c.Request.Context(), userID, questID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nextQuest)
}

// MarkTaskHabitCompleteHandler - отметить задачу как выполненную сегодня (habit)
func (h *QuestHandler) MarkTaskHabitCompleteHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		CompletionTime *string `json:"completion_time"` // Опционально: время выполнения (HH:MM:SS)
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Если время не указано, используем текущее время
		now := time.Now().Format("15:04:05")
		req.CompletionTime = &now
	}

	err = h.questService.MarkTaskHabitComplete(c.Request.Context(), userID, questID, taskID, req.CompletionTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetTaskHabitProgressHandler - получить прогресс habit tracking для задачи
func (h *QuestHandler) GetTaskHabitProgressHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	progress, err := h.questService.GetTaskHabitProgress(c.Request.Context(), userID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": progress})
}

// FreezeMissedDayHandler - заморозить пропущенный день за монеты
func (h *QuestHandler) FreezeMissedDayHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questID, err := strconv.Atoi(c.Param("questID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Date string `json:"date" binding:"required"` // Формат: "2006-01-02"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.questService.FreezeMissedDay(c.Request.Context(), userID, questID, taskID, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Day frozen successfully"})
}

func (h *QuestHandler) CreateSharedQuest(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req models.CreateSharedQuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.questService.CreateSharedQuest(userID, req.FriendID, req.QuestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shared quest created successfully"})
}

// GetSharedQuestsHandler - получить совместные квесты пользователя
func (h *QuestHandler) GetSharedQuestsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sharedQuests, err := h.questService.GetSharedQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sharedQuests})
}

// Search relevant quests by title / description (поисковик - интеграция с Bert FastAPI-микросервисом)
// без аутентификаци-авторизации
func (h *QuestHandler) SearchQuests(c *gin.Context) {
	var req models.RecommendationService_SearchQuest_Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	req.Status = req.GetStatus() // default == "all"

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quests, err := h.questService.SearchQuests(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quests)
}

// Рекомендация друзей с помощью Recommendation Service
func (h *QuestHandler) RecommendFriends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req := models.RecommendationService_RecommendUsers_Request{
		UserID: userID,
	}

	resp, err := h.questService.RecommendFriends(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Рекомендация квестов с помощью Recommendation Service
func (h *QuestHandler) RecommendQuests(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.questService.RecommendQuests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetQuestTreeHandler - получить дерево квестов с зависимостями
func (h *QuestHandler) GetQuestTreeHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tree, err := h.questService.GetQuestTree(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GetDevelopmentBranchesHandler - получить все ветки развития
func (h *QuestHandler) GetDevelopmentBranchesHandler(c *gin.Context) {
	branches, err := h.questService.GetDevelopmentBranches(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// GetUserPassiveBuffsHandler - получить пассивные баффы пользователя
func (h *QuestHandler) GetUserPassiveBuffsHandler(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	buffs, err := h.questService.GetUserPassiveBuffs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, buffs)
}

// RegisterQuestRoutes sets up the routes for quest handling with Gin
func RegisterQuestRoutes(router *gin.Engine, questService *services.QuestService) {
	handler := NewQuestHandler(questService)

	questGroup := router.Group("/quests")

	questGroup.Use(middleware.JWTAuthMiddleware())
	{
		questGroup.GET("/:questID/details", handler.GetQuestDetails)
		questGroup.GET("/my-quests-with-details", handler.GetMyAllQuestsWithDetails)
		questGroup.GET("/available", handler.GetAvailableQuestsHandler)
		questGroup.GET("/shop", handler.GetQuestShopHandler)
		questGroup.GET("/active", handler.GetMyActiveQuestsHandler)
		questGroup.GET("/completed", handler.GetMyCompletedQuests)
		questGroup.POST("/:questID/purchase", handler.PurchaseQuestHandler)
		questGroup.POST("/:questID/start", handler.StartQuestHandler)
		questGroup.POST("/:questID/complete", handler.CompleteQuestHandler)
		questGroup.POST("/:questID/:taskID/complete", handler.CompleteTaskHandler)

		questGroup.GET("/shared", handler.GetSharedQuestsHandler)
		questGroup.POST("/shared", handler.CreateSharedQuest)

		questGroup.POST("/generate", handler.GenerateAIQuest)
		questGroup.POST("/schedule", handler.GenerateScheduleByAI)

		// recommendation service
		questGroup.POST("/search", handler.SearchQuests)
		questGroup.POST("/recommend/friends", handler.RecommendFriends)
		questGroup.POST("/recommend", handler.RecommendQuests)

		// New endpoints for quest tree system
		questGroup.GET("/tree", handler.GetQuestTreeHandler)
		questGroup.GET("/buffs", handler.GetUserPassiveBuffsHandler)

		// Habit tracking and quest levels
		questGroup.POST("/:questID/continue", handler.ContinueQuestHandler)
		questGroup.POST("/:questID/:taskID/habit-complete", handler.MarkTaskHabitCompleteHandler)
		questGroup.GET("/:questID/:taskID/habit-progress", handler.GetTaskHabitProgressHandler)
		questGroup.POST("/:questID/:taskID/freeze-day", handler.FreezeMissedDayHandler)

		// Enhanced AI features
		questGroup.POST("/ai/motivate", handler.GenerateMotivationHandler)
		questGroup.POST("/ai/analyze-progress", handler.AnalyzeProgressHandler)
		questGroup.POST("/ai/improve-quest", handler.ImproveQuestDescriptionHandler)
		questGroup.POST("/ai/personalized-quest", handler.GeneratePersonalizedQuestHandler)
	}

	// Public endpoint for branches (no auth required)
	branchesGroup := router.Group("/branches")
	{
		branchesGroup.GET("", handler.GetDevelopmentBranchesHandler)
	}
}
