package json_test

import (
	"strings"
	"testing"

	json "github.com/natalie-o-perret/go-data/json"
)

func TestFormat(t *testing.T) {
	input := []byte(`{"z":"<>&","n":9007199254740993,"fraction":1.2300e+45}`)
	want := "{\n" +
		"  \"fraction\": 1.2300e+45,\n" +
		"  \"n\": 9007199254740993,\n" +
		"  \"z\": \"<>&\"\n" +
		"}\n"

	got, err := json.Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("Format() =\n%s\nwant:\n%s", got, want)
	}
	if strings.Contains(string(got), `\u003`) {
		t.Fatalf("Format() HTML-escaped output: %s", got)
	}
}

func TestCompact(t *testing.T) {
	got, err := json.Compact([]byte(" { \"text\" : \"<>&\", \"items\" : [ 1, 2 ] } \n"))
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"text":"<>&","items":[1,2]}`; string(got) != want {
		t.Fatalf("Compact() = %q, want %q", got, want)
	}
	if _, err := json.Compact([]byte(`{"broken":`)); err == nil {
		t.Fatal("Compact() accepted invalid JSON")
	}
}

func TestFormatRejectsExtraInput(t *testing.T) {
	for _, input := range []string{"", `{"a":1} []`, `{"a":1} trailing`} {
		t.Run(input, func(t *testing.T) {
			if _, err := json.Format([]byte(input)); err == nil {
				t.Fatalf("Format(%q) succeeded", input)
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := map[string]bool{
		"null":          true,
		" \n [1, 2] \t": true,
		"":              false,
		"{} []":         false,
		"{":             false,
	}
	for input, want := range tests {
		if got := json.Valid([]byte(input)); got != want {
			t.Errorf("Valid(%q) = %v, want %v", input, got, want)
		}
	}
}
