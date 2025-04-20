package main

import (
	"net/http"

	model "translator/models"
	openai "translator/openai"
	repo "translator/repo"
	service "translator/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/transcriptions", getTranscriptionRecords)
	router.POST("/transcriptions", addTranscription)
	router.POST("/translate", translate)

	router.Run("0.0.0.0:8080")
}

func getTranscriptionRecords(c *gin.Context) {
	transcriptions := repo.GetAllTranscriptions()
	c.JSON(http.StatusOK, transcriptions)
}

func addTranscription(c *gin.Context) {
	var newTranscription model.TranscriptionRecord
	if err := c.ShouldBindJSON(&newTranscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo.AddTranscription(newTranscription)
	c.JSON(http.StatusCreated, newTranscription)
}

func translate(c *gin.Context) {
	var inputs []model.TranscriptionDTO
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transcriptionSet := model.TranscriptionSet{}
	for _, input := range inputs {
		transcription := *model.NewTranscription(input.Sentence, input.Speaker, input.Time)
		// check which transcriptions already exists
		existingTranscription, exists := repo.GetTranscriptionByHash(transcription.Hash)
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
			dto := service.MapTranscriptionToDTO(transcription)
			resultDTOs = append(resultDTOs, dto)
		}
	}

	// If there are new transcriptions, process them
	if len(transcriptionSet.New) > 0 {
		// Create a batch collection from the new transcriptions
		batchCollection := service.NewBatchCollection(transcriptionSet.New)

		// translate batches
		for _, batch := range batchCollection.Batches {
			prompt, _ := batch.BuildPrompt()
			// Simulate translation process
			translatedText := openai.Translate(prompt)
			decodedText := batch.UnmarshalResponse(translatedText)
			batch.MapTranslations(decodedText)
		}

		// reconstruct original transcriptions from the batches
		resultTranscription := batchCollection.ReconstructOriginalTranscriptions()
		for _, transcription := range resultTranscription {
			resultDTOs = append(resultDTOs, service.MapTranscriptionToDTO(transcription))
			// Add the new transcriptions to the repository
			record := service.MapTranscriptionToRecord(transcription)
			repo.AddTranscription(record)
		}
	}

	c.JSON(http.StatusOK, resultDTOs)
}
