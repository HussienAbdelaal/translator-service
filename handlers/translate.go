package handler

import (
	"net/http"
	model "translator/models"
	service "translator/services"

	"github.com/gin-gonic/gin"
)

type TranslateHandler struct {
	translateService service.TranslateService
}

func NewTranslateHandler(translateService service.TranslateService) *TranslateHandler {
	return &TranslateHandler{
		translateService: translateService,
	}
}

func (h *TranslateHandler) GetAllTranscriptions(c *gin.Context) {
	// Call the TranslateService to get all transcriptions
	transcriptions := h.translateService.GetAll(c)
	if transcriptions == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No transcriptions found"})
		return
	}
	c.JSON(http.StatusOK, transcriptions)
}

func (h *TranslateHandler) Translate(c *gin.Context) {
	var inputs []model.TranscriptionDTO
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the TranslateService to translate the inputs
	result, err := h.translateService.Translate(c, inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
