package handlers

import (
	"BecomeOverMan/internal/integrations"
	"BecomeOverMan/internal/services"
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TechHandler struct {
	service    *services.TechService
	grpcClient *integrations.RecommendationGRPCClient
}

func NewTechHandler(service *services.TechService, grpcClient *integrations.RecommendationGRPCClient) *TechHandler {
	return &TechHandler{service: service, grpcClient: grpcClient}
}

func (h *TechHandler) CheckConnectionDB(c *gin.Context) {
	if err := h.service.CheckConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "database is available"})
}

func (h *TechHandler) RecommendationServiceHealth(c *gin.Context) {
	// gRPC health check вместо HTTP
	resp, err := h.grpcClient.HealthCheck(context.Background())
	if err != nil {
		slog.Error("Failed to check recommendation service health via gRPC", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Cannot connect to recommendation service",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       resp.Status,
		"model":        resp.Model,
		"device":       resp.Device,
		"quests_count": resp.QuestsCount,
	})
}

func RegisterTechRoutes(router *gin.Engine, techService *services.TechService, grpcClient *integrations.RecommendationGRPCClient) {
	handler := NewTechHandler(techService, grpcClient)

	g := router.Group("/tech")
	{
		g.GET("/ping-db", handler.CheckConnectionDB)

		g.GET("/recommendation-service/health", handler.RecommendationServiceHealth)
	}
}
