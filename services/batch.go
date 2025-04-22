package service

import (
	"encoding/json"
	"fmt"
	model "translator/models"
	"unicode/utf8"
)

type Batch struct {
	Transcriptions []model.Transcription
}

func (b *Batch) AddTranscription(transcription model.Transcription) {
	b.Transcriptions = append(b.Transcriptions, transcription)
}

func (b *Batch) GetTextSize() int {
	size := 0
	for _, transcription := range b.Transcriptions {
		size += utf8.RuneCountInString(transcription.Text)
	}
	return size
}

func (b *Batch) BuildPrompt() (string, int) {
	// Create a slice of transcription texts
	payloads := make([]string, len(b.Transcriptions))
	for i, transcription := range b.Transcriptions {
		payloads[i] = transcription.Text
	}
	// Marshal the minimal payloads to JSON
	prompt, _ := json.Marshal(payloads)
	return string(prompt), utf8.RuneCountInString(string(prompt))
}

func (b *Batch) UnmarshalResponse(prompt string) ([]string, error) {
	payloads := make([]string, 0)
	err := json.Unmarshal([]byte(prompt), &payloads)
	if err != nil {
		return nil, err
	}
	return payloads, nil
}

func (b *Batch) MapTranslations(payloads []string) error {
	// check if the length of payloads and b.Transcriptions are equal
	if len(payloads) != len(b.Transcriptions) {
		return fmt.Errorf("length of Payloads and b.Transcriptions are not equal")
	}
	// payloads are in the same order as b.Transcriptions
	for i := 0; i < len(payloads); i++ {
		b.Transcriptions[i].Translation = payloads[i]
	}
	return nil
}

func (b *Batch) GetNumberOfTranscriptions() int {
	return len(b.Transcriptions)
}
