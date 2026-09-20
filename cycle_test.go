package data_test

import (
	"errors"
	"testing"

	"github.com/natalie-o-perret/go-data"
)

func TestRecursiveContainers(t *testing.T) {
	mapCycle := map[string]any{}
	mapCycle["self"] = mapCycle
	sliceCycle := make([]any, 1)
	sliceCycle[0] = sliceCycle
	mixedMap := map[string]any{}
	mixedSlice := []any{mixedMap}
	mixedMap["slice"] = mixedSlice

	for name, value := range map[string]any{
		"map": mapCycle, "slice": sliceCycle, "mixed": mixedMap,
	} {
		t.Run(name, func(t *testing.T) {
			err := data.Walk(value, func(data.Path, any) error { return nil })
			if !errors.Is(err, data.ErrRecursive) {
				t.Fatalf("Walk() error = %v, want ErrRecursive", err)
			}

			_, err = data.Transform(value,
				func(_ data.Path, value any) (any, error) { return value, nil })
			if !errors.Is(err, data.ErrRecursive) {
				t.Fatalf("Transform() error = %v, want ErrRecursive", err)
			}
		})
	}
}

func TestSharedContainersAreVisitedAndClonedAtEachPath(t *testing.T) {
	shared := map[string]any{"items": []any{1}}
	input := []any{shared, shared}
	visits := 0
	if err := data.Walk(input, func(_ data.Path, value any) error {
		if _, ok := value.(int); ok {
			visits++
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if visits != 2 {
		t.Fatalf("Walk() visited shared leaf %d times, want 2", visits)
	}

	result, err := data.Transform(input,
		func(_ data.Path, value any) (any, error) { return value, nil })
	if err != nil {
		t.Fatal(err)
	}
	clones := result.([]any)
	left := clones[0].(map[string]any)["items"].([]any)
	right := clones[1].(map[string]any)["items"].([]any)
	left[0] = 2
	if right[0] != 1 || shared["items"].([]any)[0] != 1 {
		t.Fatalf("Transform() retained shared containers: %#v", result)
	}
}
