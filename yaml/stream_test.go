package yaml

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	yamlv3 "go.yaml.in/yaml/v3"
)

var errWrite = errors.New("write failed")

var (
	_ *Encoder = (*yamlv3.Encoder)(nil)
	_ *Decoder = (*yamlv3.Decoder)(nil)
)

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) { return 0, errWrite }

func TestStreamConstructors(t *testing.T) {
	var output bytes.Buffer
	encoder := NewEncoder(&output)
	if err := encoder.Encode(map[string]any{
		"outer": map[string]any{"value": true},
	}); err != nil {
		t.Fatal(err)
	}
	if err := encoder.Close(); err != nil {
		t.Fatal(err)
	}
	if want := "outer:\n  value: true\n"; output.String() != want {
		t.Fatalf("encoded stream = %q, want %q", output.String(), want)
	}

	decoder := NewDecoder(
		bytes.NewBufferString("one: 1\n---\ntwo: 2\n"),
	)
	for index, key := range []string{"one", "two"} {
		var document map[string]int
		if err := decoder.Decode(&document); err != nil {
			t.Fatal(err)
		}
		if document[key] != index+1 {
			t.Fatalf("decoded %q document as %#v", key, document)
		}
	}
	var document any
	if err := decoder.Decode(&document); !errors.Is(err, io.EOF) {
		t.Fatalf("final Decode() error = %v, want EOF", err)
	}
}

func TestEncoderReportsWriterError(t *testing.T) {
	encoder := NewEncoder(errorWriter{})
	if err := encoder.Encode("value"); err == nil || !strings.Contains(err.Error(), errWrite.Error()) {
		t.Fatalf("Encode() error = %v, want %v", err, errWrite)
	}
}
