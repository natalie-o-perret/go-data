package anonymize

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

func (a *identifierAnonymizer) randomLike(original string) (string, error) {
	var result strings.Builder
	for _, current := range original {
		var alphabet string
		switch {
		case current >= 'a' && current <= 'z', unicode.IsLower(current):
			alphabet = "abcdefghijklmnopqrstuvwxyz"
		case current >= 'A' && current <= 'Z', unicode.IsUpper(current):
			alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		case current >= '0' && current <= '9', unicode.IsDigit(current):
			alphabet = "0123456789"
		default:
			result.WriteRune(current)
			continue
		}
		index, err := a.randomIndex(len(alphabet))
		if err != nil {
			return "", fmt.Errorf("anonymize: generate identifier: %w", err)
		}
		result.WriteByte(alphabet[index])
	}
	return result.String(), nil
}

func (a *identifierAnonymizer) randomDecimal(original string) (string, error) {
	sign := ""
	digits := original
	if strings.HasPrefix(digits, "-") {
		sign, digits = "-", digits[1:]
	}
	result := make([]byte, len(digits))
	for i := range result {
		alphabet := "0123456789"
		if i == 0 && len(result) > 1 {
			alphabet = "123456789"
		}
		index, err := a.randomIndex(len(alphabet))
		if err != nil {
			return "", fmt.Errorf("anonymize: generate numeric identifier: %w", err)
		}
		result[i] = alphabet[index]
	}
	return sign + string(result), nil
}

func (a *identifierAnonymizer) randomIndex(size int) (int, error) {
	limit := 256 - 256%size
	var data [1]byte
	for {
		if _, err := io.ReadFull(a.random, data[:]); err != nil {
			return 0, err
		}
		if int(data[0]) < limit {
			return int(data[0]) % size, nil
		}
	}
}
