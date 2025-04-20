package openai

import (
	"encoding/json"
	"unicode/utf8"
)

type Batch struct {
	Messages []Transcription
}

type PromptMessage struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (b *Batch) AddMessage(msg Transcription) {
	b.Messages = append(b.Messages, msg)
}

func (b *Batch) GetTextSize() int {
	size := 0
	for _, msg := range b.Messages {
		size += utf8.RuneCountInString(msg.Text)
	}
	return size
}

func (b *Batch) BuildPrompt() (string, int) {
	// Create a slice of promptMessage
	promptMessages := make([]PromptMessage, len(b.Messages))
	for i, msg := range b.Messages {
		promptMessages[i] = PromptMessage{
			ID:   msg.ID,
			Text: msg.Text,
		}
	}
	// Marshal the promptMessages to JSON
	prompt, _ := json.Marshal(promptMessages)
	return string(prompt), utf8.RuneCountInString(string(prompt))
}

func (b *Batch) UnmarshalResponse(prompt string) []PromptMessage {
	var promptMessages []PromptMessage
	err := json.Unmarshal([]byte(prompt), &promptMessages)
	if err != nil {
		panic(err.Error())
	}
	return promptMessages
}

func (b *Batch) MapTranslationsToMessages(PromptMessages []PromptMessage) {
	// check if the length of PromptMessages and b.Messages are equal
	if len(PromptMessages) != len(b.Messages) {
		panic("Length of PromptMessages and b.Messages are not equal")
	}
	for i := 0; i < len(PromptMessages); i++ {
		b.Messages[i].Translation = PromptMessages[i].Text
	}
	// fmt.Printf("Mapped Translations: %v\n", b.Messages)
}

func (b *Batch) GetNumberOfMessages() int {
	return len(b.Messages)
}
