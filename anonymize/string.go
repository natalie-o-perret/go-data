package anonymize

import (
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	hexPattern  = regexp.MustCompile(`(?i)^[0-9a-f]{16,}$`)
)

func (a *identifierAnonymizer) string(value string, force bool) (string, error) {
	if value == "" {
		return value, nil
	}
	switch {
	case uuidPattern.MatchString(value):
		return a.mapped("string:"+value, func() (string, error) { return a.randomUUID(value) })
	case hexPattern.MatchString(value):
		return a.mapped("string:"+value, func() (string, error) { return a.randomHex(value) })
	case force:
		return a.mapped("string:"+value, func() (string, error) { return a.randomLike(value) })
	default:
		return value, nil
	}
}

func (a *identifierAnonymizer) mapped(key string, generate func() (string, error)) (string, error) {
	if replacement, ok := a.replacements[key]; ok {
		return replacement, nil
	}
	replacement, err := generate()
	if err != nil {
		return "", err
	}
	a.replacements[key] = replacement
	return replacement, nil
}

func (a *identifierAnonymizer) randomUUID(original string) (string, error) {
	data := make([]byte, 16)
	if _, err := io.ReadFull(a.random, data); err != nil {
		return "", fmt.Errorf("anonymize: generate UUID: %w", err)
	}
	data[6] = data[6]&0x0f | 0x40
	data[8] = data[8]&0x3f | 0x80
	encoded := hex.EncodeToString(data)
	result := encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:]
	if original == strings.ToUpper(original) {
		result = strings.ToUpper(result)
	}
	return result, nil
}

func (a *identifierAnonymizer) randomHex(original string) (string, error) {
	result := make([]byte, len(original))
	const digits = "0123456789abcdef"
	for i := range result {
		index, err := a.randomIndex(len(digits))
		if err != nil {
			return "", fmt.Errorf("anonymize: generate hexadecimal value: %w", err)
		}
		result[i] = digits[index]
		if original[i] >= 'A' && original[i] <= 'F' && result[i] >= 'a' {
			result[i] -= 'a' - 'A'
		}
	}
	return string(result), nil
}
