package pretty

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func header(column string) string {
	return title(column)
}

func title(value string) string {
	runes := []rune(value)
	words := make([]string, 0, 4)
	start, uppercaseRun := -1, 0
	for i, current := range runes {
		if !unicode.IsLetter(current) && !unicode.IsDigit(current) {
			if start >= 0 {
				words = append(words, string(runes[start:i]))
			}
			start, uppercaseRun = -1, 0
			continue
		}
		if start < 0 {
			start = i
			if unicode.IsUpper(current) {
				uppercaseRun = 1
			}
			continue
		}

		previous := runes[i-1]
		nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
		boundary := unicode.IsUpper(current) && (unicode.IsLower(previous) || unicode.IsDigit(previous) ||
			unicode.IsUpper(previous) && nextIsLower && uppercaseRun > 1)
		if boundary {
			words = append(words, string(runes[start:i]))
			start, uppercaseRun = i, 1
			continue
		}
		if unicode.IsUpper(current) {
			uppercaseRun++
		} else {
			uppercaseRun = 0
		}
	}
	if start >= 0 {
		words = append(words, string(runes[start:]))
	}
	for i, word := range words {
		first, size := utf8.DecodeRuneInString(word)
		words[i] = string(unicode.ToUpper(first)) + word[size:]
	}
	return strings.Join(words, " ")
}
