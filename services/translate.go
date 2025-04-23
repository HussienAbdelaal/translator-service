package service

import (
	"context"
	mapper "translator/mappers"
	model "translator/models"
	"translator/utils"
)

type ITranslateClient interface {
	Translate(ctx context.Context, prompt string) (string, error)
	GetBatchSize() int
}

type ITranslateRepo interface {
	Get(ctx context.Context, hash string) (*model.TranscriptionRecord, error)
	Create(ctx context.Context, t model.TranscriptionRecord) error
	GetAll(ctx context.Context) ([]model.TranscriptionRecord, error)
}

type TranslateService struct {
	translateRepo   ITranslateRepo
	translateClient ITranslateClient
}

func NewTranslateService(translateRepo ITranslateRepo, translateClient ITranslateClient) *TranslateService {
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
	transcriptionSet, inputOrderMap, err := s.classifyTranscriptions(ctx, inputs)
	if err != nil {
		return nil, err
	}

	resultDTOs := make([]model.TranscriptionDTO, len(inputs)) // Preallocate result slice

	s.insertExistingTranscriptions(transcriptionSet.Existing, resultDTOs, inputOrderMap)

	// If there are new transcriptions, process them
	if len(transcriptionSet.New) > 0 {
		if err := s.translateNewTranscriptions(ctx, transcriptionSet.New, resultDTOs, inputOrderMap); err != nil {
			return nil, err
		}
	}

	return resultDTOs, nil
}

func (s *TranslateService) classifyTranscriptions(ctx context.Context, inputs []model.TranscriptionDTO) (model.TranscriptionSet, map[string]int, error) {
	transcriptionSet := model.TranscriptionSet{}
	inputOrderMap := make(map[string]int)

	for index, input := range inputs {
		transcription := *model.NewTranscription(input.Sentence, input.Speaker, input.Time)
		inputOrderMap[transcription.Hash] = index

		existing, err := s.translateRepo.Get(ctx, transcription.Hash)
		if err != nil {
			return model.TranscriptionSet{}, nil, err
		}

		if existing != nil {
			transcription.Translation = existing.Translation
			transcriptionSet.Existing = append(transcriptionSet.Existing, transcription)
		} else {
			transcriptionSet.New = append(transcriptionSet.New, transcription)
		}
	}

	return transcriptionSet, inputOrderMap, nil
}

func (s *TranslateService) insertExistingTranscriptions(existing []model.Transcription, resultDTOs []model.TranscriptionDTO, inputOrderMap map[string]int) {
	for _, transcription := range existing {
		dto := mapper.MapTranscriptionToDTO(transcription)
		resultDTOs[inputOrderMap[transcription.Hash]] = dto
	}
}

func (s *TranslateService) translateNewTranscriptions(
	ctx context.Context,
	newTranscriptions []model.Transcription,
	resultDTOs []model.TranscriptionDTO,
	inputOrderMap map[string]int,
) error {
	batchSize := s.translateClient.GetBatchSize()
	batchCollection := NewBatchCollection(batchSize, newTranscriptions)

	prompts := make([]string, len(batchCollection.Batches))
	for i, batch := range batchCollection.Batches {
		prompt, _ := batch.BuildPrompt()
		prompts[i] = prompt
	}

	responses, err := utils.DoInParallelFailFast(ctx, prompts, s.translateClient.Translate)
	if err != nil {
		return err
	}

	for i, batch := range batchCollection.Batches {
		decoded, err := batch.UnmarshalResponse(responses[i])
		if err != nil {
			return err
		}
		batch.MapTranslations(decoded)
	}

	reconstructed := batchCollection.ReconstructOriginalTranscriptions()
	for _, t := range reconstructed {
		dto := mapper.MapTranscriptionToDTO(t)
		resultDTOs[inputOrderMap[t.Hash]] = dto

		record := mapper.MapTranscriptionToRecord(t)
		s.translateRepo.Create(ctx, record)
	}

	return nil
}
