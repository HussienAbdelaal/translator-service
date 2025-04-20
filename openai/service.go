package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"translator/config"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Translate(text string) string {
	// Load the configuration
	cfg := config.LoadConfig()

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIAPIKey),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a translator. Translate only Arabic parts into English. Leave any English text unchanged."),
			openai.UserMessage(text),
		},
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0.3),
	})
	if err != nil {
		panic(err.Error())
	}

	// print usage
	fmt.Printf("Usage total tokens: %v\n", chatCompletion.Usage.TotalTokens)

	return chatCompletion.Choices[0].Message.Content
}

func TranslateBatch(batch Batch) Batch {
	// Load the configuration
	cfg := config.LoadConfig()

	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIAPIKey),
	)

	prompt, _ := batch.BuildPrompt()

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a translator. Translate only Arabic parts into English. Leave any English text unchanged."),
			openai.UserMessage(prompt),
		},
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0.3),
	})
	if err != nil {
		panic(err.Error())
	}

	// print usage
	fmt.Printf("Usage total tokens: %v\n", chatCompletion.Usage.TotalTokens)

	fmt.Printf("Translated text: %s\n", chatCompletion.Choices[0].Message.Content)
	// return chatCompletion.Choices[0].Message.Content
	responseBatch := Batch{}
	json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &responseBatch)
	fmt.Printf("Translated RESPONSE: %v\n", responseBatch)

	return Batch{
		Messages: []Transcription{},
	}
}
