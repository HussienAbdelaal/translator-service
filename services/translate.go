package service

import (
	"context"
	mapper "translator/mappers"
	model "translator/models"
)

type TranslateClient interface {
	Translate(ctx context.Context, prompt string) (string, error)
	GetBatchSize() int
}

type TranslateRepo interface {
	Get(ctx context.Context, hash string) (*model.TranscriptionRecord, error)
	Create(ctx context.Context, t model.TranscriptionRecord) error
	GetAll(ctx context.Context) ([]model.TranscriptionRecord, error)
}

type TranslateService struct {
	translateRepo   TranslateRepo
	translateClient TranslateClient
}

func NewTranslateService(translateRepo TranslateRepo, translateClient TranslateClient) *TranslateService {
	return &TranslateService{
		translateRepo:   translateRepo,
		translateClient: translateClient,
	}
}

func (s *TranslateService) GetAll(ctx context.Context) ([]model.TranscriptionRecord, error) {
	// Get all transcriptions from the repository
	transcriptions, err := s.translateRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return transcriptions, nil
}

func (s *TranslateService) Translate(ctx context.Context, inputs []model.TranscriptionDTO) ([]model.TranscriptionDTO, error) {
	transcriptionSet := model.TranscriptionSet{}
	for _, input := range inputs {
		transcription := *model.NewTranscription(input.Sentence, input.Speaker, input.Time)
		// check which transcriptions already exists
		existingTranscription, err := s.translateRepo.Get(ctx, transcription.Hash)
		if err != nil {
			return nil, err
		}
		if existingTranscription != nil {
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
			dto := mapper.MapTranscriptionToDTO(transcription)
			resultDTOs = append(resultDTOs, dto)
		}
	}

	// If there are new transcriptions, process them
	if len(transcriptionSet.New) > 0 {
		// Create a batch collection from the new transcriptions
		batchSize := s.translateClient.GetBatchSize()
		batchCollection := NewBatchCollection(batchSize, transcriptionSet.New)

		// translate batches
		for _, batch := range batchCollection.Batches {
			prompt, _ := batch.BuildPrompt()
			translatedText, err := s.translateClient.Translate(ctx, prompt)
			if err != nil {
				// fail fast if any translation fails
				return nil, err
			}
			decodedText, err := batch.UnmarshalResponse(translatedText)
			if err != nil {
				// fail fast if any unmarshaling fails
				return nil, err
			}
			batch.MapTranslations(decodedText)
		}

		// reconstruct original transcriptions from the batches
		resultTranscription := batchCollection.ReconstructOriginalTranscriptions()
		for _, transcription := range resultTranscription {
			resultDTOs = append(resultDTOs, mapper.MapTranscriptionToDTO(transcription))
			// Add the new transcriptions to the repository
			record := mapper.MapTranscriptionToRecord(transcription)
			s.translateRepo.Create(ctx, record)
		}
	}
	return resultDTOs, nil
}
