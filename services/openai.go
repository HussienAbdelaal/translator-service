package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"translator/config"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIService struct {
	config      *config.OpenAIConfig
	client      *openai.Client
	model       string
	temperature float64
	batchSize   int
}

func NewOpenAIService(cfg config.OpenAIConfig) (*OpenAIService, error) {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithMaxRetries(3),
	)

	// Set default values for model, temperature, and batch size
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	temperature := cfg.Temperature
	if temperature == "" {
		temperature = "0.3"
	}
	// Convert temperature to float64
	temperatureFloat, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid temperature: %s", temperature)
	}

	batchSize := cfg.BatchSize
	if batchSize == "" {
		batchSize = "3000"
	}
	// Convert batch size to int
	batchSizeInt, err := strconv.Atoi(batchSize)
	if err != nil {
		return nil, fmt.Errorf("invalid batch size: %s", batchSize)
	}

	return &OpenAIService{
		config:      &cfg,
		client:      &client,
		model:       model,
		temperature: temperatureFloat,
		batchSize:   batchSizeInt,
	}, nil
}

func (s *OpenAIService) Translate(ctx context.Context, text string) (string, error) {
	client := *s.client
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a translator. Translate only Arabic parts into English. Leave any English text unchanged. Maintain the same order"),
			openai.UserMessage(text),
		},
		Model:       s.model,
		Temperature: openai.Float(s.temperature),
	})
	if err != nil {
		return "", err
	}

	// print usage
	log.Printf("Usage total tokens: %v\n", chatCompletion.Usage.TotalTokens)

	return chatCompletion.Choices[0].Message.Content, nil
}
