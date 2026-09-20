package xml

import (
	"bytes"
	"testing"
)

func TestFormatIsDeterministic(t *testing.T) {
	input := []byte(`<root><group><item id="1">x</item><empty/></group></root>`)
	want := "<root>\n  <group>\n    <item id=\"1\">x</item>\n    <empty></empty>\n  </group>\n</root>\n"
	got, err := Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}

	again, err := Format(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(again, got) {
		t.Fatalf("Format(Format()) = %q, want %q", again, got)
	}
}

func TestFormatPreservesMixedContent(t *testing.T) {
	input := []byte(`<p>Hello <b>world</b> and <i>XML</i>!</p>`)
	want := "<p>Hello <b>world</b> and <i>XML</i>!</p>\n"
	got, err := Format(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}
}
