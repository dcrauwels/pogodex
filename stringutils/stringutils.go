package stringutils

import (
	"strings"
)

func CleanInput(text string) []string {
	lc := strings.ToLower(text)
	trLc := strings.TrimSpace(lc)
	return strings.Fields(trLc)
}
