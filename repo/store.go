package repo

import model "translator/models"

// mimics a database index by Hash
var HashStore = make(map[string]model.TranscriptionRecord)

func AddTranscription(transcription model.TranscriptionRecord) {
	HashStore[transcription.Hash] = transcription
}

func GetTranscriptionByHash(hash string) (model.TranscriptionRecord, bool) {
	transcription, exists := HashStore[hash]
	return transcription, exists
}

func GetAllTranscriptions() []model.TranscriptionRecord {
	transcriptions := make([]model.TranscriptionRecord, 0, len(HashStore))
	for _, transcription := range HashStore {
		transcriptions = append(transcriptions, transcription)
	}
	return transcriptions
}
