package json

import (
	"bytes"
	stdjson "encoding/json"
	"errors"
	"io"
)

// Compact validates data and removes insignificant JSON whitespace.
func Compact(data []byte) ([]byte, error) {
	var output bytes.Buffer
	if err := stdjson.Compact(&output, data); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// Format validates and returns deterministic, two-space-indented JSON.
func Format(data []byte) ([]byte, error) {
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, errors.New("json: multiple values")
		}
		return nil, err
	}

	var output bytes.Buffer
	encoder := NewEncoder(&output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
