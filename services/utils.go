package service

var separators = []string{".", ",", "?", "!", ";", "\n"}

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
