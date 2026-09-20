package pretty

import (
	"strings"
	"testing"
)

func TestRenderNestedStructures(t *testing.T) {
	value := map[string]any{
		"catalog": map[string]any{
			"name":  "demo",
			"owner": map[string]any{"name": "Alex"},
			"sections": []any{
				map[string]any{
					"name": "fiction",
					"books": []any{
						map[string]any{"id": "book-1", "status": "available"},
						map[string]any{"id": "book-2", "status": "borrowed"},
					},
				},
			},
		},
	}

	var out strings.Builder
	if err := Render(&out, value); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"└── Catalog", "    ├── Name: demo", "    ├── Owner",
		"    │   └── Name: Alex", "    └── Sections", "        └── Item 1",
		"            ├── Name: fiction", "│ Books", "book-1", "book-2",
	} {
		if !strings.Contains(out.String(), expected) {
			t.Fatalf("output omitted %q:\n%s", expected, out.String())
		}
	}
}

func TestRenderEmptyContainersInline(t *testing.T) {
	value := map[string]any{
		"alerts":   []any{},
		"metadata": map[string]any{},
	}

	var out strings.Builder
	if err := Render(&out, value); err != nil {
		t.Fatal(err)
	}
	want := "├── Alerts: []\n└── Metadata: {}\n"
	if got := out.String(); got != want {
		t.Fatalf("output mismatch:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
