package repo

import (
	"context"
	model "translator/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TranslationRepo struct {
	db *pgxpool.Pool
}

func NewTranslationRepo(db *pgxpool.Pool) *TranslationRepo {
	return &TranslationRepo{
		db: db,
	}
}

func (r *TranslationRepo) Create(ctx context.Context, t model.TranscriptionRecord) {
	// Insert the translation into the database
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO translation (hash, text, translation) VALUES ($1, $2, $3)",
		t.Hash, t.Text, t.Translation)
	if err != nil {
		panic(err.Error())
	}
}

func (r *TranslationRepo) Get(ctx context.Context, hash string) (model.TranscriptionRecord, bool) {
	row := r.db.QueryRow(
		ctx,
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

func (r *TranslationRepo) GetAll(ctx context.Context) []model.TranscriptionRecord {
	// Get all translations from the database
	// Print a message indicating a database query is being attempted
	rows, err := r.db.Query(
		ctx,
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
