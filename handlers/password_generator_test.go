package handlers

import (
	"fmt"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

var caps = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func TestGenerateRandomHash_16CharWithCapitalsWithSymbols(t *testing.T) {
	length := 16
	uppercase := true
	symbols := true

	result, err := generateRandomHash(length, uppercase, symbols)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Len(t, result, 16)

	if !containsUppercase(result) {
		assert.FailNow(t, "no capital in result", fmt.Sprint("result: ", result))
	}

	if !containsNumberOrSymbol(result) {
		assert.FailNow(t, "no number or symbol in result", fmt.Sprint("result: ", result))
	}
}

func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func containsNumberOrSymbol(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) || unicode.IsSymbol(r) || unicode.IsPunct(r) {
			return true
		}
	}
	return false
}
