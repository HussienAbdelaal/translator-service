package openai

import (
	"fmt"
	"unicode/utf8"
)

var separators = []string{".", ",", "?", "!", ";", "\n"}

type Transcription struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

func (t *Transcription) GetTextSize() int {
	return utf8.RuneCountInString(t.Text)
}

type BatchCollection struct {
	Batches       []Batch
	OriginMapping map[string][]string // Maps original message IDs to their corresponding split message IDs
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

func Reorg(messages []Transcription) BatchCollection {
	batches := BatchCollection{
		Batches:       []Batch{},
		OriginMapping: make(map[string][]string),
	}
	currentBatch := Batch{}
	maxSize := 50 // max size for each batch

	fmt.Println("Splitting bigger messages first")
	splitMessages := []Transcription{}
	for i, msg := range messages {
		fmt.Printf("Message %d: %v  Size: %d\n", i+1, msg, utf8.RuneCountInString(msg.Text))
		// check if message needs to be split
		if msg.GetTextSize() > maxSize {
			// Split the message into smaller messages and keep track of the IDs
			// NOTE: if the message can't be split, it will be added as is
			textSlices := SplitBySeparator(msg.Text)
			for j, text := range textSlices {
				splitMsg := Transcription{ID: fmt.Sprintf("%s.%d", msg.ID, j), Text: text}
				splitMessages = append(splitMessages, splitMsg)
				batches.OriginMapping[msg.ID] = append(batches.OriginMapping[msg.ID], splitMsg.ID)
			}
		} else {
			// Add the original message to the split messages
			splitMessages = append(splitMessages, msg)
		}
	}

	// Now process the split messages into batches
	fmt.Println("Processing split messages into batches")
	for i, msg := range splitMessages {
		fmt.Printf("Edited Message %d: %v  New Size: %d\n", i+1, msg, utf8.RuneCountInString(msg.Text))

		// check if current batch + next message exceeds max size
		if currentBatch.GetTextSize() > 0 && currentBatch.GetTextSize()+msg.GetTextSize() > maxSize {
			// batch is now considered done
			batches.Batches = append(batches.Batches, currentBatch)
			currentBatch = Batch{} // reset current batch
		}
		currentBatch.AddMessage(msg)
	}

	// Add the last batch if not empty
	if currentBatch.GetNumberOfMessages() > 0 {
		batches.Batches = append(batches.Batches, currentBatch)
	}

	// // Print the batches
	// for i, batch := range batches.Batches {
	// 	prompt, size := batch.BuildPrompt()
	// 	fmt.Printf("Batch %d (size=%d):\n%s\n", i+1, size, prompt)
	// }
	// // print originMapping
	// for k, v := range batches.OriginMapping {
	// 	fmt.Printf("Original ID: %s, Split IDs: %v\n", k, v)
	// }

	// // reconstruct the original messages
	// fmt.Println("Reconstructing original messages")
	// for originalMsgId, splitMsgIds := range batches.OriginMapping {
	// 	// reconstruct the original message
	// 	reconstructedMessage := ""
	// 	for _, splitID := range splitMsgIds {
	// 		for _, msg := range splitMessages {
	// 			if msg.ID == splitID {
	// 				reconstructedMessage += msg.Text
	// 				break
	// 			}
	// 		}
	// 	}
	// 	fmt.Printf("Reconstructed Message for ID %s: %s\n", originalMsgId, reconstructedMessage)
	// }

	return batches
}

/*

sudo code:
1. take as input []Transcription
2. currentBatch => empty batch
3. loop over input messages
4. if input message size > max size => split message into smaller messages
5. reconstruct smaller messages into moderate messages
6. keep track of IDs belonging to the same message
7. if batch size + next message > max size => batch is considered done
8. send batch to open ai for translation
9. receive response => regenerate the same input messages

*/
