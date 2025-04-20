package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	// Add other configs here, e.g.
	// DBUrl        string
	// Timeout      time.Duration
}

// LoadConfig loads environment variables once and returns a shared config instance
func LoadConfig() *Config {
	err := godotenv.Load() // attempt to load .env file
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is required but not found in environment variables")
	}

	return &Config{
		OpenAIAPIKey: apiKey,
	}
}
