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
	transcriptionSet := model.TranscriptionSet{}
	inputOrderMap := make(map[string]int) // Map to track the input order
	for index, input := range inputs {
		transcription := *model.NewTranscription(input.Sentence, input.Speaker, input.Time)
		inputOrderMap[transcription.Hash] = index // Store the input order using the hash
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

	resultDTOs := make([]model.TranscriptionDTO, len(inputs)) // Preallocate result slice

	if len(transcriptionSet.Existing) > 0 {
		// If there are existing transcriptions, add them to the result
		for _, transcription := range transcriptionSet.Existing {
			dto := mapper.MapTranscriptionToDTO(transcription)
			// Place the DTO in the correct order
			resultDTOs[inputOrderMap[transcription.Hash]] = dto
		}
	}

	// If there are new transcriptions, process them
	if len(transcriptionSet.New) > 0 {
		// Create a batch collection from the new transcriptions
		batchSize := s.translateClient.GetBatchSize()
		batchCollection := NewBatchCollection(batchSize, transcriptionSet.New)

		prompts := []string{}
		for _, batch := range batchCollection.Batches {
			prompt, _ := batch.BuildPrompt()
			prompts = append(prompts, prompt)
		}
		// translate all batches in parallel
		promptResponses, err := utils.DoInParallelFailFast(ctx, prompts, s.translateClient.Translate)
		if err != nil {
			// error returned if any translation fails
			return nil, err
		}
		// map the responses to the batches. They are in the same order as the prompts
		// so we can use the index to map them to the correct batch
		for i, batch := range batchCollection.Batches {
			decodedText, err := batch.UnmarshalResponse(promptResponses[i])
			if err != nil {
				// fail fast if any unmarshaling fails
				return nil, err
			}
			batch.MapTranslations(decodedText)
		}

		// reconstruct original transcriptions from the batches
		resultTranscription := batchCollection.ReconstructOriginalTranscriptions()
		for _, transcription := range resultTranscription {
			dto := mapper.MapTranscriptionToDTO(transcription)
			// Place the DTO in the correct order
			resultDTOs[inputOrderMap[transcription.Hash]] = dto
			// Add the new transcriptions to the repository
			record := mapper.MapTranscriptionToRecord(transcription)
			s.translateRepo.Create(ctx, record)
		}
	}

	return resultDTOs, nil
}
