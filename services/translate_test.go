package service

import (
	"context"
	"encoding/json"
	"testing"
	model "translator/models"
	"translator/test/mocks"
	"translator/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildPrompt(input []string) string {
	marshaled, _ := json.Marshal(input)
	return string(marshaled)
}

func TestGetTranscriptions(t *testing.T) {
	mockRepo := new(mocks.TranslateRepo)
	mockClient := new(mocks.TranslateClient)

	expected := []model.TranscriptionRecord{
		{Hash: "hash1", Text: "Hello", Translation: "Hola"},
		{Hash: "hash2", Text: "World", Translation: "Mundo"},
	}
	mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

	translateService := NewTranslateService(mockRepo, mockClient)
	transcriptions, err := translateService.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, transcriptions, 2)
	assert.Equal(t, "Hello", transcriptions[0].Text)
	assert.Equal(t, "World", transcriptions[1].Text)
}

func TestTranslateExistingTranscriptions(t *testing.T) {
	mockRepo := new(mocks.TranslateRepo)
	mockClient := new(mocks.TranslateClient)

	translateService := NewTranslateService(mockRepo, mockClient)

	// Mock the Get method to return a record for both of the calls
	hash1 := utils.GenerateHash("مرحبا")
	hash2 := utils.GenerateHash("اليوم")
	mockRepo.On("Get", mock.Anything, hash1).Return(&model.TranscriptionRecord{Hash: hash1, Text: "مرحبا", Translation: "Hello"}, nil)
	mockRepo.On("Get", mock.Anything, hash2).Return(&model.TranscriptionRecord{Hash: hash2, Text: "اليوم", Translation: "Today"}, nil)

	inputs := []model.TranscriptionDTO{
		{Sentence: "مرحبا", Speaker: "Speaker1", Time: "0"},
		{Sentence: "اليوم", Speaker: "Speaker2", Time: "1"},
	}

	results, err := translateService.Translate(context.TODO(), inputs)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Hello", results[0].Sentence)
	assert.Equal(t, "Today", results[1].Sentence)
}

func TestTranslateNewAndExistingTranscriptions(t *testing.T) {
	mockRepo := new(mocks.TranslateRepo)
	mockClient := new(mocks.TranslateClient)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	inputPrompt := buildPrompt([]string{"مرحبا"})
	outputPrompt := buildPrompt([]string{"Hello"})
	mockClient.On("Translate", mock.Anything, inputPrompt).Return(outputPrompt, nil)
	mockClient.On("GetBatchSize").Return(100)

	translateService := NewTranslateService(mockRepo, mockClient)

	// Mock the Get method to return nil for the first call and a record for the second call
	hash1 := utils.GenerateHash("مرحبا")
	hash2 := utils.GenerateHash("اليوم")
	mockRepo.On("Get", mock.Anything, hash1).Return(nil, nil)
	mockRepo.On("Get", mock.Anything, hash2).Return(&model.TranscriptionRecord{Hash: hash2, Text: "اليوم", Translation: "Today"}, nil)

	inputs := []model.TranscriptionDTO{
		{Sentence: "مرحبا", Speaker: "Speaker1", Time: "0"},
		{Sentence: "اليوم", Speaker: "Speaker2", Time: "1"},
	}

	results, err := translateService.Translate(context.TODO(), inputs)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Hello", results[0].Sentence)
	assert.Equal(t, "Today", results[1].Sentence)
}

func TestTranslateBiggerThanBatchSize(t *testing.T) {
	mockRepo := new(mocks.TranslateRepo)
	mockClient := new(mocks.TranslateClient)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// This proves that the input string is split into two prompts
	// and that the translation is done in two steps
	inputPrompt1 := buildPrompt([]string{"مرحبا."})
	outputPrompt1 := buildPrompt([]string{"Hello."})
	mockClient.On("Translate", mock.Anything, inputPrompt1).Return(outputPrompt1, nil)
	inputPrompt2 := buildPrompt([]string{" كيف حالك؟"})
	outputPrompt2 := buildPrompt([]string{" How are you?"})
	mockClient.On("Translate", mock.Anything, inputPrompt2).Return(outputPrompt2, nil)
	mockClient.On("GetBatchSize").Return(10)

	// Mock the Get method to return nil
	hash1 := utils.GenerateHash("مرحبا. كيف حالك؟")
	mockRepo.On("Get", mock.Anything, hash1).Return(nil, nil)

	inputs := []model.TranscriptionDTO{
		{Sentence: "مرحبا. كيف حالك؟", Speaker: "Speaker1", Time: "0"},
	}

	translateService := NewTranslateService(mockRepo, mockClient)

	results, err := translateService.Translate(context.TODO(), inputs)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Hello. How are you?", results[0].Sentence)
}
