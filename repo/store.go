package repo

import (
	"context"
	"translator/db"
	model "translator/models"
)

func AddTranscription(transcription model.TranscriptionRecord) {
	// Insert the transcription into the database
	_, err := db.DB.Exec(
		context.TODO(),
		"INSERT INTO translation (hash, text, translation) VALUES ($1, $2, $3)",
		transcription.Hash, transcription.Text, transcription.Translation)
	if err != nil {
		panic(err.Error())
	}
}

func GetTranscriptionByHash(hash string) (model.TranscriptionRecord, bool) {
	row := db.DB.QueryRow(
		context.TODO(),
		"SELECT hash, text, translation FROM translation WHERE hash = $1", hash)

	// Scan the row into a TranscriptionRecord
	var transcription model.TranscriptionRecord
	err := row.Scan(&transcription.Hash, &transcription.Text, &transcription.Translation)

	if err != nil {
		if err.Error() == "no rows in result set" {
			// If no rows were found, return an empty TranscriptionRecord and false
			return model.TranscriptionRecord{}, false
		}
		panic(err.Error())
	}
	return transcription, true
}

func GetAllTranscriptions() []model.TranscriptionRecord {
	// Get all transcriptions from the database
	// Print a message indicating a database query is being attempted
	rows, err := db.DB.Query(
		context.TODO(),
		"SELECT hash, text, translation FROM translation")
	if err != nil {
		// Print an error message if the query fails
		panic(err.Error())
	}
	defer rows.Close()
	// Create a slice to hold the transcriptions
	transcriptions := []model.TranscriptionRecord{}
	// Iterate over the rows and append each transcription to the slice
	for rows.Next() {
		var transcription model.TranscriptionRecord
		err := rows.Scan(&transcription.Hash, &transcription.Text, &transcription.Translation)
		if err != nil {
			panic(err.Error())
		}
		transcriptions = append(transcriptions, transcription)
	}
	return transcriptions
}
