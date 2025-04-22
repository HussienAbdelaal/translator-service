package container

import (
	"fmt"
	client "translator/clients"
	"translator/config"
	"translator/db"
	"translator/repo"
	service "translator/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	DB                   *pgxpool.Pool
	OpenAIClient         *client.OpenAIClient
	Repo                 *repo.TranslationRepo
	TranscriptionService *service.TranslateService
}

func NewContainer(cfg *config.Config) Container {
	// Initialize database connection
	db, err := db.NewDBPool(&cfg.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	// Initialize OpenAI client
	openAIClient, err := client.NewOpenAIClient(cfg.OpenAI)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize OpenAI service: %v", err))
	}
	// Initialize transcription repository
	repo := repo.NewTranslationRepo(db)
	// Initialize transcription service
	transcriptionService := service.NewTranslateService(repo, openAIClient)

	return Container{
		DB:                   db,
		OpenAIClient:         openAIClient,
		Repo:                 repo,
		TranscriptionService: transcriptionService,
	}
}

func (c *Container) Close() {
	c.DB.Close()
}
