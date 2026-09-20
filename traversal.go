package data

import (
	"errors"
	"reflect"
	"sort"
)

// Path identifies a value using string map keys and integer slice indexes.
// The root path is nil.
type Path []any

// VisitFunc visits a value during Walk.
type VisitFunc func(Path, any) error

// TransformFunc transforms a value during Transform.
type TransformFunc func(Path, any) (any, error)

// ErrRecursive reports a cycle in a map[string]any or []any being traversed.
var ErrRecursive = errors.New("data: recursive value")

type containerID struct {
	kind     byte
	pointer  uintptr
	length   int
	capacity int
}

func identify(value any) (containerID, bool) {
	switch value := value.(type) {
	case map[string]any:
		return containerID{kind: 'm', pointer: reflect.ValueOf(value).Pointer()}, true
	case []any:
		return containerID{
			kind: 's', pointer: reflect.ValueOf(value).Pointer(),
			length: len(value), capacity: cap(value),
		}, true
	default:
		return containerID{}, false
	}
}

func sortedKeys(value map[string]any) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func childPath(path Path, element any) Path {
	return append(path[:len(path):len(path)], element)
}
