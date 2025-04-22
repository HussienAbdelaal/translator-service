package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	OpenAI OpenAIConfig
}

var (
	cfg  *Config
	once sync.Once
	err  error
)

func loadConfig() (*Config, error) {
	err := godotenv.Load() // attempt to load .env file
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	// Load database configuration
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_NAME")
	if dbHost == "" || dbUsername == "" || dbPassword == "" || dbDatabase == "" {
		return nil, fmt.Errorf("DB_HOST, DB_USER, DB_PASSWORD, and DB_NAME are required but not found in environment variables")
	}

	// Load OpenAI configuration
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIModel := os.Getenv("OPENAI_MODEL")
	openAIBatchSize := os.Getenv("OPENAI_BATCH_SIZE")
	openAITemperature := os.Getenv("OPENAI_TEMPERATURE")
	if openAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required but not found in environment variables")
	}

	cfg = &Config{
		DB: DBConfig{
			Host:     dbHost,
			Username: dbUsername,
			Password: dbPassword,
			Database: dbDatabase,
		},
		OpenAI: OpenAIConfig{
			APIKey:      openAIAPIKey,
			Model:       openAIModel,
			BatchSize:   openAIBatchSize,
			Temperature: openAITemperature,
		},
	}
	return cfg, nil
}

// LoadConfig loads environment variables once and returns a shared config instance
func Load() (*Config, error) {
	once.Do(func() {
		cfg, err = loadConfig()
	})
	return cfg, err
}
