package data_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/natalie-o-perret/go-data"
)

func TestWalkOrderAndPaths(t *testing.T) {
	value := map[string]any{
		"b": []any{"x", "y"},
		"a": map[string]any{"d": 1, "c": 2},
	}
	var got []data.Path
	err := data.Walk(value, func(path data.Path, _ any) error {
		got = append(got, append(data.Path(nil), path...))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []data.Path{
		nil, {"a"}, {"a", "c"}, {"a", "d"},
		{"b"}, {"b", 0}, {"b", 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("paths = %#v, want %#v", got, want)
	}
}

func TestWalkReturnsVisitorError(t *testing.T) {
	want := errors.New("stop")
	calls := 0
	err := data.Walk(map[string]any{"a": 1, "b": 2}, func(path data.Path, _ any) error {
		calls++
		if reflect.DeepEqual(path, data.Path{"a"}) {
			return want
		}
		return nil
	})
	if !errors.Is(err, want) {
		t.Fatalf("Walk() error = %v, want %v", err, want)
	}
	if calls != 2 {
		t.Fatalf("Walk() made %d calls, want 2", calls)
	}
}
