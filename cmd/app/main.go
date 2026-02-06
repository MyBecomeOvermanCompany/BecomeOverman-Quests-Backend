package main

import (
	"BecomeOverMan/internal/handlers"
	"log"
	"log/slog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"BecomeOverMan/internal/config"
	grpcclient "BecomeOverMan/internal/grpc"
	_ "BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"BecomeOverMan/internal/services"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug) // Включаем DEBUG-логирование

	db, err := sqlx.Connect("postgres", config.Cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize gRPC client for recommendation service
	grpcClient, err := grpcclient.NewRecommendationClient(config.Cfg.RecommendationGRPCAddress)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to connect to recommendation gRPC service: %v", err)
		log.Printf("   Recommendation features will not be available")
		// Don't exit - allow the app to run without recommendations
	} else {
		defer grpcClient.Close()
		log.Println("✓ Recommendation gRPC client initialized")
	}

	techRepo := repositories.NewTechRepository(db)
	techService := services.NewTechService(techRepo)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	questRepo := repositories.NewQuestRepository(db)
	questService := services.NewQuestService(questRepo, userRepo, grpcClient)

	r := gin.Default()
	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	{
		handlers.RegisterTechRoutes(r, techService)

		handlers.RegisterUserRoutes(r, userService)
		handlers.RegisterQuestRoutes(r, questService)
	}

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
