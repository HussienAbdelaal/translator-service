package main

import (
	"fmt"
	config "translator/config"
	db "translator/db"
	handler "translator/handlers"
	repo "translator/repo"
	service "translator/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	// Initialize database connection
	db, err := db.NewDBPool(cfg.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer db.Close()
	// Initialize OpenAI client
	openAIService, err := service.NewOpenAIService(cfg.OpenAI)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize OpenAI service: %v", err))
	}
	// Initialize transcription repository
	repo := repo.NewTranslationRepo(db)
	// Initialize transcription service
	transcriptionService := service.NewTranslateService(*repo, *openAIService)
	// Initialize transcription handler
	handler := handler.NewTranslateHandler(*transcriptionService)

	router := gin.Default()
	router.GET("/translations", handler.GetAllTranslations)
	router.POST("/translate", handler.Translate)

	router.Run("0.0.0.0:8080")
}
