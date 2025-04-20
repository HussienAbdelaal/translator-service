package service

import (
	model "translator/models"
)

func MapTranscriptionToDTO(transcription model.Transcription) model.TranscriptionDTO {
	return model.TranscriptionDTO{
		Speaker:  transcription.Speaker,
		Time:     transcription.Time,
		Sentence: transcription.Translation, // Assuming Translation is the final text
	}
}

func MapTranscriptionToRecord(transcription model.Transcription) model.TranscriptionRecord {
	return model.TranscriptionRecord{
		Hash:        transcription.Hash,
		Text:        transcription.Text,
		Translation: transcription.Translation,
	}
}
