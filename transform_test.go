package data_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/natalie-o-perret/go-data"
)

func TestTransformPostOrderWithoutMutation(t *testing.T) {
	original := map[string]any{"b": 1, "a": []any{2}}
	var paths []data.Path
	result, err := data.Transform(original, func(path data.Path, value any) (any, error) {
		paths = append(paths, append(data.Path(nil), path...))
		if number, ok := value.(int); ok {
			return number * 10, nil
		}
		return value, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	wantPaths := []data.Path{{"a", 0}, {"a"}, {"b"}, nil}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("paths = %#v, want %#v", paths, wantPaths)
	}
	if want := (map[string]any{"b": 10, "a": []any{20}}); !reflect.DeepEqual(result, want) {
		t.Fatalf("Transform() = %#v, want %#v", result, want)
	}

	result.(map[string]any)["a"].([]any)[0] = 30
	if want := (map[string]any{"b": 1, "a": []any{2}}); !reflect.DeepEqual(original, want) {
		t.Fatalf("Transform() changed original: %#v", original)
	}
}

func TestTransformReturnsVisitorError(t *testing.T) {
	want := errors.New("stop")
	result, err := data.Transform(map[string]any{"a": 1},
		func(_ data.Path, _ any) (any, error) { return nil, want })
	if !errors.Is(err, want) {
		t.Fatalf("Transform() error = %v, want %v", err, want)
	}
	if result != nil {
		t.Fatalf("Transform() result = %#v, want nil", result)
	}
}

func TestTransformPreservesNilContainers(t *testing.T) {
	value := map[string]any{"map": map[string]any(nil), "slice": []any(nil)}
	result, err := data.Transform(value,
		func(_ data.Path, value any) (any, error) { return value, nil })
	if err != nil {
		t.Fatal(err)
	}
	got := result.(map[string]any)
	if !reflect.ValueOf(got["map"]).IsNil() || !reflect.ValueOf(got["slice"]).IsNil() {
		t.Fatalf("Transform() changed nil containers: %#v", got)
	}
}
