package main

import (
	"fmt"
	config "translator/config"
	"translator/container"
	handler "translator/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	// Initialize dependency container
	container := container.NewContainer(cfg)
	defer container.Close()
	handler := handler.NewTranslateHandler(container.TranscriptionService)

	router := gin.Default()
	router.GET("/translations", handler.GetAllTranslations)
	router.POST("/translate", handler.Translate)

	router.Run("0.0.0.0:8080")
}
