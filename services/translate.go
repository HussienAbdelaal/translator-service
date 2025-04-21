package service

import (
	"context"
	model "translator/models"
	openai "translator/openai"
	repo "translator/repo"
)

type TranslateService struct {
	transRepo     repo.TranscriptionRepo
	openaiService openai.OpenAIService
}

func NewTranslateService(transRepo repo.TranscriptionRepo, openaiService openai.OpenAIService) *TranslateService {
	return &TranslateService{
		transRepo:     transRepo,
		openaiService: openaiService,
	}
}

func (s *TranslateService) GetAll(ctx context.Context) []model.TranscriptionRecord {
	// Get all transcriptions from the repository
	transcriptions := s.transRepo.GetAll(ctx)
	return transcriptions
}

func (s *TranslateService) Translate(ctx context.Context, inputs []model.TranscriptionDTO) ([]model.TranscriptionDTO, error) {
	transcriptionSet := model.TranscriptionSet{}
	for _, input := range inputs {
		transcription := *model.NewTranscription(input.Sentence, input.Speaker, input.Time)
		// check which transcriptions already exists
		existingTranscription, exists := s.transRepo.Get(ctx, transcription.Hash)
		if exists {
			// if exists, add the existing transcription
			transcription.Translation = existingTranscription.Translation
			transcriptionSet.Existing = append(transcriptionSet.Existing, transcription)
		} else {
			// if not, add the new transcription
			transcriptionSet.New = append(transcriptionSet.New, transcription)
		}
	}

	resultDTOs := []model.TranscriptionDTO{}

	if len(transcriptionSet.Existing) > 0 {
		// If there are existing transcriptions, add them to the result
		for _, transcription := range transcriptionSet.Existing {
			dto := MapTranscriptionToDTO(transcription)
			resultDTOs = append(resultDTOs, dto)
		}
	}

	// If there are new transcriptions, process them
	if len(transcriptionSet.New) > 0 {
		// Create a batch collection from the new transcriptions
		batchCollection := NewBatchCollection(transcriptionSet.New)

		// translate batches
		for _, batch := range batchCollection.Batches {
			prompt, _ := batch.BuildPrompt()
			// Simulate translation process
			translatedText := s.openaiService.Translate(ctx, prompt)
			decodedText := batch.UnmarshalResponse(translatedText)
			batch.MapTranslations(decodedText)
		}

		// reconstruct original transcriptions from the batches
		resultTranscription := batchCollection.ReconstructOriginalTranscriptions()
		for _, transcription := range resultTranscription {
			resultDTOs = append(resultDTOs, MapTranscriptionToDTO(transcription))
			// Add the new transcriptions to the repository
			record := MapTranscriptionToRecord(transcription)
			s.transRepo.Create(ctx, record)
		}
	}
	return resultDTOs, nil
}
