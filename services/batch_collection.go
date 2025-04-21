package service

import (
	"fmt"
	model "translator/models"
	"unicode/utf8"
)

type BatchCollection struct {
	Batches       []Batch
	OriginMapping map[string][]string // Maps original IDs to their corresponding split IDs
	MaxSize       int
}

func NewBatchCollection(maxSize int, transcriptions []model.Transcription) BatchCollection {
	batches := BatchCollection{
		Batches:       []Batch{},
		OriginMapping: make(map[string][]string),
		MaxSize:       maxSize,
	}
	currentBatch := Batch{}

	fmt.Println("Splitting bigger transcriptions first")
	normalizedT := batches.normalizeTranscriptions(transcriptions)

	// Now process the normalized transcriptions into batches
	fmt.Println("Processing normalizedT transcriptions into batches")
	for i, norm := range normalizedT {
		fmt.Printf("Normalized Transcription %d: %v  New Size: %d\n", i+1, norm, utf8.RuneCountInString(norm.Text))

		// check if current batch + next transcription exceeds max size
		if currentBatch.GetTextSize() > 0 && currentBatch.GetTextSize()+norm.GetTextSize() > maxSize {
			// batch is now considered done
			batches.Batches = append(batches.Batches, currentBatch)
			currentBatch = Batch{} // reset current batch
		}
		currentBatch.AddTranscription(norm)
	}

	// Add the last batch if not empty
	if currentBatch.GetNumberOfTranscriptions() > 0 {
		batches.Batches = append(batches.Batches, currentBatch)
	}

	return batches
}

// This function will reconstruct the original split transcriptions if any
// by merging the split transcriptions back into their original form.
func (b *BatchCollection) ReconstructOriginalTranscriptions() []model.Transcription {

	splitT, unsplitT := b.categorizeTranscriptions()

	// if splitT is empty, return the unsplit transcriptions as is
	if len(splitT) == 0 {
		return unsplitT
	}

	resultTs := []model.Transcription{}

	// reconstruct the split transcriptions
	resultTs = append(resultTs, b.reconstructParts(splitT)...)

	// add the transcriptions that were not in OriginMapping
	resultTs = append(resultTs, unsplitT...)

	return resultTs
}
