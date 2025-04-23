package repo

import (
	"context"
	"fmt"
	model "translator/models"

	pgxV5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IDBPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgxV5.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgxV5.Rows, error)
	Close()
}

type TranslationRepo struct {
	db IDBPool
}

func NewTranslationRepo(db IDBPool) *TranslationRepo {
	return &TranslationRepo{
		db: db,
	}
}

func (r *TranslationRepo) Create(ctx context.Context, t model.TranscriptionRecord) error {
	// Insert the translation into the database
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO translation (hash, text, translation) VALUES ($1, $2, $3)",
		t.Hash, t.Text, t.Translation)
	if err != nil {
		return fmt.Errorf("failed to insert translation: %w", err)
	}
	return nil
}

func (r *TranslationRepo) Get(ctx context.Context, hash string) (*model.TranscriptionRecord, error) {
	row := r.db.QueryRow(
		ctx,
		"SELECT hash, text, translation FROM translation WHERE hash = $1", hash)

	// Scan the row into a TranscriptionRecord
	var transcription model.TranscriptionRecord
	err := row.Scan(&transcription.Hash, &transcription.Text, &transcription.Translation)

	if err != nil {
		// TODO: catch this error properly
		if err.Error() == "no rows in result set" {
			// If no rows were found, return an empty TranscriptionRecord and false
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	return &transcription, nil
}

func (r *TranslationRepo) GetAll(ctx context.Context) ([]model.TranscriptionRecord, error) {
	// Get all translations from the database
	rows, err := r.db.Query(
		ctx,
		"SELECT hash, text, translation FROM translation")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Create a slice to hold the transcriptions
	transcriptions := []model.TranscriptionRecord{}
	// Iterate over the rows and append each transcription to the slice
	for rows.Next() {
		var transcription model.TranscriptionRecord
		err := rows.Scan(&transcription.Hash, &transcription.Text, &transcription.Translation)
		if err != nil {
			return nil, err
		}
		transcriptions = append(transcriptions, transcription)
	}
	return transcriptions, nil
}
