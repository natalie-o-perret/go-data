package xml

import (
	stdxml "encoding/xml"
	"io"
)

// Encoder is the standard library XML encoder.
type Encoder = stdxml.Encoder

// Decoder is the standard library XML decoder.
type Decoder = stdxml.Decoder

// Marshal returns the compact XML encoding of v.
func Marshal(v any) ([]byte, error) {
	return stdxml.Marshal(v)
}

// Unmarshal parses XML data into v.
func Unmarshal(data []byte, v any) error {
	return stdxml.Unmarshal(data, v)
}

// NewEncoder returns a standard library XML encoder writing to w.
func NewEncoder(w io.Writer) *Encoder {
	return stdxml.NewEncoder(w)
}

// NewDecoder returns a standard library XML decoder reading from r.
func NewDecoder(r io.Reader) *Decoder {
	return stdxml.NewDecoder(r)
}
