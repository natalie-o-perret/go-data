package json

import (
	"bytes"
	stdjson "encoding/json"
	"io"
)

// Encoder is the standard library JSON encoder.
type Encoder = stdjson.Encoder

// Decoder is the standard library JSON decoder.
type Decoder = stdjson.Decoder

// Marshal returns the compact JSON encoding of v without HTML escaping.
func Marshal(v any) ([]byte, error) {
	var output bytes.Buffer
	if err := NewEncoder(&output).Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(output.Bytes(), []byte("\n")), nil
}

// Unmarshal parses JSON-encoded data into v.
func Unmarshal(data []byte, v any) error {
	return stdjson.Unmarshal(data, v)
}

// Valid reports whether data contains exactly one valid JSON value.
func Valid(data []byte) bool {
	return stdjson.Valid(data)
}

// NewEncoder returns a standard library encoder with HTML escaping disabled.
func NewEncoder(w io.Writer) *Encoder {
	encoder := stdjson.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder
}

// NewDecoder returns a standard library JSON decoder.
func NewDecoder(r io.Reader) *Decoder {
	return stdjson.NewDecoder(r)
}
