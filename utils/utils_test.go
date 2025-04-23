package utils

import (
	"testing"
)

func TestGenerateHash(t *testing.T) {
	input := "hello"
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	actual := GenerateHash(input)

	if actual != expected {
		t.Errorf("hash mismatch:\nexpected: %s\ngot:      %s", expected, actual)
	}
}
