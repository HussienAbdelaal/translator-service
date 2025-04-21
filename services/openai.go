package service

import (
	"context"
	"fmt"
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

func NewOpenAIService(cfg config.OpenAIConfig) *OpenAIService {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)

	// Set default values for model, temperature, and batch size
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}
	// // TODO: Validate the model name against OpenAI's available models
	// // Use OpenAI SDK to validate the model name
	// modelsList, err := client.Models.List(context.TODO())
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to fetch models list: %v", err))
	// }

	// isValidModel := false
	// for _, m := range modelsList.Data {
	// 	if m.ID == model {
	// 		isValidModel = true
	// 		break
	// 	}
	// }

	// if !isValidModel {
	// 	panic(fmt.Sprintf("Invalid model name: %s", model))
	// }

	temperature := cfg.Temperature
	if temperature == "" {
		temperature = "0.3"
	}
	// Convert temperature to float64
	temperatureFloat, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		panic(fmt.Sprintf("Invalid temperature: %s", temperature))
	}

	batchSize := cfg.BatchSize
	if batchSize == "" {
		batchSize = "3000"
	}
	// Convert batch size to int
	batchSizeInt, err := strconv.Atoi(batchSize)
	if err != nil {
		panic(fmt.Sprintf("Invalid batch size: %s", batchSize))
	}

	return &OpenAIService{
		config:      &cfg,
		client:      &client,
		model:       model,
		temperature: temperatureFloat,
		batchSize:   batchSizeInt,
	}
}

func (s *OpenAIService) Translate(ctx context.Context, text string) string {
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
		panic(err.Error())
	}

	// print usage
	fmt.Printf("Usage total tokens: %v\n", chatCompletion.Usage.TotalTokens)

	return chatCompletion.Choices[0].Message.Content
}
