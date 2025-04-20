package openai

import (
	"context"
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
			openai.SystemMessage("You are a translator. Translate only Arabic parts into English. Leave any English text unchanged. Maintain the same order"),
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
