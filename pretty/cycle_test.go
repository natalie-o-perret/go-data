package pretty

import (
	"errors"
	"strings"
	"testing"

	data "github.com/natalie-o-perret/go-data"
)

func TestRenderRejectsRecursiveData(t *testing.T) {
	object := map[string]any{}
	object["self"] = object

	items := make([]any, 1)
	items[0] = items

	mixedObject := map[string]any{}
	mixedItems := []any{mixedObject}
	mixedObject["items"] = mixedItems

	for name, value := range map[string]any{
		"map":   object,
		"slice": items,
		"mixed": mixedObject,
	} {
		t.Run(name, func(t *testing.T) {
			var out strings.Builder
			err := Render(&out, value)
			if !errors.Is(err, data.ErrRecursive) {
				t.Fatalf("Render() error = %v, want wrapped data.ErrRecursive", err)
			}
			if !strings.Contains(err.Error(), "pretty: recursive data structures are not supported") {
				t.Fatalf("Render() error = %q", err)
			}
			if out.Len() != 0 {
				t.Fatalf("Render() wrote %q before returning the error", out.String())
			}
		})
	}
}

func TestRenderAllowsSharedData(t *testing.T) {
	shared := map[string]any{"name": "same"}
	var out strings.Builder
	if err := Render(&out, map[string]any{"left": shared, "right": shared}); err != nil {
		t.Fatal(err)
	}
	if strings.Count(out.String(), "Name: same") != 2 {
		t.Fatalf("Render() output = %q, want shared value at both paths", out.String())
	}
}
