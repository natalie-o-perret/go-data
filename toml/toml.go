package toml

import (
	"io"

	pelletier "github.com/pelletier/go-toml/v2"
)

// Encoder is Pelletier's TOML encoder.
type Encoder = pelletier.Encoder

// Decoder is Pelletier's TOML decoder.
type Decoder = pelletier.Decoder

// Marshal serializes v as one TOML document using Pelletier's encoder.
func Marshal(v any) ([]byte, error) {
	return pelletier.Marshal(v)
}

// Unmarshal deserializes one TOML document into v.
func Unmarshal(data []byte, v any) error {
	return pelletier.Unmarshal(data, v)
}

// NewEncoder returns a TOML encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return pelletier.NewEncoder(w)
}

// NewDecoder returns a TOML decoder that reads all remaining bytes in r as one document.
func NewDecoder(r io.Reader) *Decoder {
	return pelletier.NewDecoder(r)
}

// Format validates one TOML document and returns deterministic, newline-terminated TOML.
//
// Formatting is value-based and therefore drops comments and lexical choices such as
// whitespace, quoting style, and numeric notation. The stable decoder does not expose
// a lossless syntax tree. An empty document formats as a single newline.
func Format(data []byte) ([]byte, error) {
	var value map[string]any
	if err := Unmarshal(data, &value); err != nil {
		return nil, err
	}
	out, err := Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	return out, nil
}

// Valid reports whether data is a well-formed TOML document.
// Empty and whitespace-only documents are valid TOML.
func Valid(data []byte) bool {
	var value map[string]any
	return Unmarshal(data, &value) == nil
}
