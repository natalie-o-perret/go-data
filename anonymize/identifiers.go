package anonymize

import (
	"crypto/rand"
	"fmt"
	"io"
	"sort"

	data "github.com/natalie-o-perret/go-data"
)

type identifierAnonymizer struct {
	random       io.Reader
	replacements map[string]string
}

// Identifiers returns a non-mutating copy of value with identifier-shaped
// values and object keys replaced. Equal identifiers receive equal replacements.
func Identifiers(value any) (any, error) {
	return identifiers(value, rand.Reader)
}

func identifiers(value any, random io.Reader) (any, error) {
	if random == nil {
		return nil, fmt.Errorf("anonymize: random source is nil")
	}
	a := &identifierAnonymizer{random: random, replacements: map[string]string{}}
	return data.Transform(value, a.transform)
}

func (a *identifierAnonymizer) transform(path data.Path, value any) (any, error) {
	name := pathName(path)
	switch value := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		result := make(map[string]any, len(value))
		for _, key := range keys {
			anonymousKey, err := a.string(key, false)
			if err != nil {
				return nil, err
			}
			if _, exists := result[anonymousKey]; exists {
				return nil, fmt.Errorf("anonymize: object key collision for %q", key)
			}
			result[anonymousKey] = value[key]
		}
		return result, nil
	case string:
		return a.string(value, identifierName(name))
	default:
		if identifierName(name) {
			return a.number(value)
		}
		return value, nil
	}
}

func pathName(path data.Path) string {
	for i := len(path) - 1; i >= 0; i-- {
		if name, ok := path[i].(string); ok {
			return name
		}
	}
	return ""
}
