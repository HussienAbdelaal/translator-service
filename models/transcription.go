package model

import (
	"crypto/sha256"
	"fmt"
	"unicode/utf8"
)

type TranscriptionDTO struct {
	Speaker  string `json:"speaker"`
	Time     string `json:"time"`
	Sentence string `json:"sentence"`
}

type Transcription struct {
	Hash        string `json:"hash"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
	Speaker     string `json:"speaker"`
	Time        string `json:"time"`
}

type TranscriptionRecord struct {
	Hash        string `json:"hash"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

func NewTranscription(text string, speaker string, time string) *Transcription {
	t := &Transcription{}
	t.Text = text
	t.Speaker = speaker
	t.Time = time
	t.Translation = ""
	h := sha256.New()
	h.Write([]byte(t.Text))
	t.Hash = fmt.Sprintf("%x", h.Sum(nil))
	return t
}

func (t *Transcription) GenerateHash() {
	h := sha256.New()
	h.Write([]byte(t.Text))
	t.Hash = fmt.Sprintf("%x", h.Sum(nil))
}

func (t *Transcription) GetTextSize() int {
	return utf8.RuneCountInString(t.Text)
}

type TranscriptionSet struct {
	Existing []Transcription
	New      []Transcription
}
