package handler

import (
	"context"
	"net/http"
	model "translator/models"

	"github.com/gin-gonic/gin"
)

type ITranslateService interface {
	GetAll(c context.Context) ([]model.TranscriptionRecord, error)
	Translate(c context.Context, inputs []model.TranscriptionDTO) ([]model.TranscriptionDTO, error)
}

type TranslateHandler struct {
	translateService ITranslateService
}

func NewTranslateHandler(translateService ITranslateService) *TranslateHandler {
	return &TranslateHandler{
		translateService: translateService,
	}
}

func (h *TranslateHandler) GetAllTranslations(c *gin.Context) {
	// Call the TranslateService to get all translations
	translations, err := h.translateService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if translations == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No translations found"})
		return
	}
	c.JSON(http.StatusOK, translations)
}

func (h *TranslateHandler) Translate(c *gin.Context) {
	var inputs []model.TranscriptionDTO
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the TranslateService to translate the inputs
	result, err := h.translateService.Translate(c.Request.Context(), inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
