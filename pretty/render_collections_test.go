package pretty

import (
	"strings"
	"testing"
)

func TestRenderCollections(t *testing.T) {
	value := map[string]any{
		"source": "fixture",
		"publications": []any{
			map[string]any{"id": "pub-b", "title": "Field Notes", "active": false, "pages": 320, "size_bytes": uint64(32 << 20), "formats": []any{"print", "ebook"}, "featured": true},
			map[string]any{"id": "pub-a", "title": "Short Stories", "active": true, "pages": 144, "size_bytes": uint64(8 << 20), "formats": []any{"ebook"}, "featured": true},
		},
		"warnings": []any{
			map[string]any{"code": "stale", "message": "One source is delayed"},
		},
	}

	var out strings.Builder
	if err := Render(&out, value); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Source=fixture", "Featured: true", "│ Publications", "Active", "Pages",
		"32 MiB", "print, ebook", "│ Warnings", "One source is delayed",
	} {
		if !strings.Contains(out.String(), expected) {
			t.Fatalf("output omitted %q:\n%s", expected, out.String())
		}
	}
	if strings.Index(out.String(), "pub-b") > strings.Index(out.String(), "pub-a") {
		t.Fatalf("table did not preserve array order:\n%s", out.String())
	}

	ids := []string{
		"11111111-1111-1111-1111-111111111111",
		"22222222-2222-2222-2222-222222222222",
		"33333333-3333-3333-3333-333333333333",
	}
	formats := []any{"hardcover", "paperback", "ebook", "audiobook", "large-print", "braille", "serial", "archive"}
	var wrapped strings.Builder
	if err := Render(&wrapped, []any{
		map[string]any{"id": ids[0], "sequence": 1, "formats": formats},
		map[string]any{"id": ids[1], "sequence": 2},
		map[string]any{"id": ids[2], "sequence": 3},
	}); err != nil {
		t.Fatal(err)
	}
	for _, expected := range append(ids, "hardcover", "paperback", "ebook", "audiobook", "large-print", "braille", "serial", "archive") {
		if !strings.Contains(wrapped.String(), expected) {
			t.Fatalf("wrapped table omitted %q:\n%s", expected, wrapped.String())
		}
	}
}
