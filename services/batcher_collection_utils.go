package service

import (
	"slices"
	model "translator/models"
)

// Split transcriptions bigger than maxSize into smaller ones
// while keeping track of the original hashes
// also updates the originMapping
func (b *BatchCollection) normalizeTranscriptions(transcriptions []model.Transcription) []model.Transcription {
	normalized := []model.Transcription{}
	for _, transcription := range transcriptions {
		// check if transcription needs to be split
		if transcription.GetTextSize() > b.MaxSize {
			// Split the transcription into smaller transcriptions and keep track of the IDs
			// NOTE: if the transcription can't be split, it will be added as is
			textSlices := SplitBySeparator(transcription.Text)
			for _, text := range textSlices {
				part := *model.NewTranscription(text, transcription.Speaker, transcription.Time)
				normalized = append(normalized, part)
				b.OriginMapping[transcription.Hash] = append(b.OriginMapping[transcription.Hash], part.Hash)
			}
		} else {
			// Add the original transcription to the split transcriptions
			normalized = append(normalized, transcription)
		}
	}
	return normalized
}

// This function will reconstruct transcriptions from the given parts
// by merging parts based on OriginMapping and returning the reconstructed transcriptions
func (b *BatchCollection) reconstructParts(parts []model.Transcription) []model.Transcription {
	result := []model.Transcription{}
	for originHash, splitHashes := range b.OriginMapping {
		text := ""
		translation := ""
		speaker := ""
		time := ""
		for _, splitHash := range splitHashes {
			for _, part := range parts {
				if part.Hash == splitHash {
					text += part.Text
					translation += part.Translation
					// all split transcriptions should have the same speaker and time
					speaker = part.Speaker
					time = part.Time
					break
				}
			}
		}
		// create a new transcription record with the reconstructed text and translation
		newTranscription := *model.NewTranscription(text, speaker, time)
		newTranscription.Translation = translation
		newTranscription.Hash = originHash // This shouldn't be necessary, but just in case
		result = append(result, newTranscription)
	}
	return result
}

// This function will categorize the transcriptions into split and unsplit
// by checking if the transcription is in the OriginMapping or not
func (b *BatchCollection) categorizeTranscriptions() (splitT, unsplitT []model.Transcription) {
	splitHashes := []string{}
	for _, sHashes := range b.OriginMapping {
		splitHashes = append(splitHashes, sHashes...)
	}
	for _, batch := range b.Batches {
		for _, t := range batch.Transcriptions {
			// check if transcription is a split or not
			if slices.Contains(splitHashes, t.Hash) {
				splitT = append(splitT, t)
			} else {
				unsplitT = append(unsplitT, t)
			}
		}
	}
	return splitT, unsplitT
}
