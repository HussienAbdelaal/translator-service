package main

import (
	config "translator/config"
	db "translator/db"
	handler "translator/handlers"
	repo "translator/repo"
	service "translator/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	// Initialize database connection
	db := db.NewDBPool(cfg.DB)
	defer db.Close()
	// Initialize OpenAI client
	openAIService := service.NewOpenAIService(cfg.OpenAI)
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
