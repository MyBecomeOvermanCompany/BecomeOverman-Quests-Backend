package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL                     string
	JWTSecret                       string
	TokenExpireHours                int
	APIKeyIntelligenceIO            string
	Recommendation_Service_BASE_URL string // Legacy HTTP endpoint (deprecated)
	RecommendationGRPCAddress       string // gRPC server address
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewConfig() Config {
	if err := godotenv.Load("../../.env"); err != nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	return Config{
		DatabaseURL:                     os.Getenv("DATABASE_URL"),
		JWTSecret:                       os.Getenv("JWT_SECRET"),
		TokenExpireHours:                24,
		APIKeyIntelligenceIO:            os.Getenv("API_KEY_INTELLIGENCE_IO"),
		Recommendation_Service_BASE_URL: "http://localhost:8000/api",
		RecommendationGRPCAddress:       getEnvOrDefault("RECOMMENDATION_GRPC_ADDRESS", "localhost:50051"),
	}
}

var Cfg = NewConfig()
