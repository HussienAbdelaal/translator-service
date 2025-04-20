package service

import (
	"fmt"
	model "translator/models"
	"unicode/utf8"

	"golang.org/x/exp/slices"
)

var separators = []string{".", ",", "?", "!", ";", "\n"}

type BatchCollection struct {
	Batches       []model.Batch
	OriginMapping map[string][]string // Maps original IDs to their corresponding split IDs
}

// SplitBySeparator returns a slice of strings split by the specified separators.
// It doesn't remove the separators from the result.
func SplitBySeparator(text string) []string {
	var result []string
	runes := []rune(text)
	currentString := ""
	for i := 0; i < len(runes); i++ {
		// fmt.Printf("Rune %v is '%c'\n", i, runes[i])
		currentString += string(runes[i])
		for _, sep := range separators {
			if runes[i] == rune(sep[0]) {
				result = append(result, currentString)
				currentString = ""
				break
			}
		}
	}
	if currentString != "" {
		result = append(result, currentString)
	}
	return result
}

func NewBatchCollection(transcriptions []model.Transcription) BatchCollection {
	batches := BatchCollection{
		Batches:       []model.Batch{},
		OriginMapping: make(map[string][]string),
	}
	currentBatch := model.Batch{}
	maxSize := 50 // max size for each batch

	fmt.Println("Splitting bigger transcriptions first")
	splitTs := []model.Transcription{}
	for i, transcription := range transcriptions {
		fmt.Printf("Transcription %d: %v  Size: %d\n", i+1, transcription, utf8.RuneCountInString(transcription.Text))
		// check if transcription needs to be split
		if transcription.GetTextSize() > maxSize {
			// Split the transcription into smaller transcriptions and keep track of the IDs
			// NOTE: if the transcription can't be split, it will be added as is
			textSlices := SplitBySeparator(transcription.Text)
			for _, text := range textSlices {
				splitT := *model.NewTranscription(text, transcription.Speaker, transcription.Time)
				splitTs = append(splitTs, splitT)
				batches.OriginMapping[transcription.Hash] = append(batches.OriginMapping[transcription.Hash], splitT.Hash)
			}
		} else {
			// Add the original transcription to the split transcriptions
			splitTs = append(splitTs, transcription)
		}
	}

	// Now process the split transcriptions into batches
	fmt.Println("Processing split transcriptions into batches")
	for i, st := range splitTs {
		fmt.Printf("Edited Transcription %d: %v  New Size: %d\n", i+1, st, utf8.RuneCountInString(st.Text))

		// check if current batch + next transcription exceeds max size
		if currentBatch.GetTextSize() > 0 && currentBatch.GetTextSize()+st.GetTextSize() > maxSize {
			// batch is now considered done
			batches.Batches = append(batches.Batches, currentBatch)
			currentBatch = model.Batch{} // reset current batch
		}
		currentBatch.AddTranscription(st)
	}

	// Add the last batch if not empty
	if currentBatch.GetNumberOfTranscriptions() > 0 {
		batches.Batches = append(batches.Batches, currentBatch)
	}

	return batches
}

func (b *BatchCollection) ReconstructOriginalTranscriptions() []model.Transcription {
	// This function will reconstruct the original split transcriptions if any
	// by merging the split transcriptions back into their original form.

	splitHashes := []string{}
	for _, sHashes := range b.OriginMapping {
		splitHashes = append(splitHashes, sHashes...)
	}

	splitT := []model.Transcription{}
	unsplitT := []model.Transcription{}
	for _, batch := range b.Batches {
		// splitT = append(splitT, batch.Transcriptions...)
		for _, t := range batch.Transcriptions {
			// check if transcription is a split or not
			if slices.Contains(splitHashes, t.Hash) {
				splitT = append(splitT, t)
			} else {
				unsplitT = append(unsplitT, t)
			}
		}
	}

	// if splitT is empty, return the unsplit transcriptions as is
	if len(splitT) == 0 {
		return unsplitT
	}

	resultTs := []model.Transcription{}

	// Iterate over the OriginMapping to reconstruct the original transcriptions
	for originHash, splitHashes := range b.OriginMapping {
		text := ""
		translation := ""
		speaker := ""
		time := ""
		for _, splitHash := range splitHashes {
			for _, t := range splitT {
				if t.Hash == splitHash {
					text += t.Text
					translation += t.Translation
					// all split transcriptions should have the same speaker and time
					speaker = t.Speaker
					time = t.Time
					break
				}
			}
		}
		// create a new transcription record with the reconstructed text and translation
		newTranscription := *model.NewTranscription(text, speaker, time)
		newTranscription.Translation = translation
		newTranscription.Hash = originHash // This shouldn't be necessary, but just in case
		resultTs = append(resultTs, newTranscription)
	}

	// add the transcriptions that were not in OriginMapping
	resultTs = append(resultTs, unsplitT...)

	return resultTs
}
