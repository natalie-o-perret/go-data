package pretty

import (
	"strings"
	"testing"
)

func TestRenderFallsBackToJSON(t *testing.T) {
	var out strings.Builder
	if err := Render(&out, []any{"one", "two"}); err != nil {
		t.Fatal(err)
	}
	want := "[\n  \"one\",\n  \"two\"\n]\n"
	if got := out.String(); got != want {
		t.Fatalf("output mismatch:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
