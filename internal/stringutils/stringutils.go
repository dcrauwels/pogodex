package stringutils

import (
	"strings"
)

func CleanInput(text string) []string {
	lc := strings.ToLower(text)   //lowercase
	trLc := strings.TrimSpace(lc) //trim spaces
	return strings.Fields(trLc)   //return as slice of strings (sep = whitespace)
}
