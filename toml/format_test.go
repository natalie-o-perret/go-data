package toml

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFormatIsDeterministicAndDropsComments(t *testing.T) {
	input := []byte("# remove this comment\nz = \"last\"\na=1\n")
	got, err := Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "a = 1\nz = 'last'\n"; string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if bytes.Contains(got, []byte("remove this comment")) {
		t.Fatalf("comment remains in %q", got)
	}
	again, err := Format(got)
	if err != nil || !bytes.Equal(again, got) {
		t.Fatalf("second format = %q, %v", again, err)
	}
}

func TestFormatDateTimesAndArraysOfTables(t *testing.T) {
	input := []byte(`offset = 1979-05-27T07:32:00Z
local-date = 1979-05-27
local-time = 07:32:00
local-date-time = 1979-05-27T07:32:00

[[products]]
name = "hammer"

[[products]]
name = "nail"
`)
	formatted, err := Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if !Valid(formatted) || formatted[len(formatted)-1] != '\n' {
		t.Fatalf("invalid formatted document: %q", formatted)
	}
	for _, literal := range []string{"1979-05-27", "07:32:00", "[[products]]"} {
		if !strings.Contains(string(formatted), literal) {
			t.Errorf("formatted output does not contain %q", literal)
		}
	}
	var got struct {
		Offset   time.Time `toml:"offset"`
		Products []struct {
			Name string `toml:"name"`
		} `toml:"products"`
	}
	if err := Unmarshal(formatted, &got); err != nil {
		t.Fatal(err)
	}
	if got.Offset.IsZero() || len(got.Products) != 2 || got.Products[1].Name != "nail" {
		t.Fatalf("got %#v", got)
	}
}

func TestFormatRejectsInvalidDocuments(t *testing.T) {
	for _, input := range [][]byte{[]byte("a = ["), []byte("a=1\na=2\n")} {
		if _, err := Format(input); err == nil {
			t.Errorf("Format(%q) succeeded", input)
		}
	}

	for _, input := range [][]byte{nil, []byte(" \n\t"), []byte("# comment-only document")} {
		got, err := Format(input)
		if err != nil || string(got) != "\n" {
			t.Errorf("Format(%q) = %q, %v", input, got, err)
		}
		again, err := Format(got)
		if err != nil || !bytes.Equal(again, got) {
			t.Errorf("second Format(%q) = %q, %v", input, again, err)
		}
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", true}, {" \n", true}, {"value = [", false},
		{"value = 1\n", true}, {"# comment-only", true},
		{"[[items]]\nname='one'\n", true},
	}
	for _, test := range tests {
		if got := Valid([]byte(test.input)); got != test.want {
			t.Errorf("Valid(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}
