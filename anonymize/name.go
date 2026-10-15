package anonymize

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func identifierName(name string) bool {
	lower := strings.ToLower(name)
	for _, suffix := range []string{"id", "ids", "uuid", "uuids", "identifier", "identifiers"} {
		if lower == suffix {
			return true
		}
		if !strings.HasSuffix(lower, suffix) {
			continue
		}
		start := len(name) - len(suffix)
		first, _ := utf8.DecodeRuneInString(name[start:])
		previous, _ := utf8.DecodeLastRuneInString(name[:start])
		if !unicode.IsLetter(previous) && !unicode.IsDigit(previous) || unicode.IsUpper(first) {
			return true
		}
	}
	return false
}
