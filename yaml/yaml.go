package yaml

import (
	"bytes"
	"errors"
	"io"

	yamlv3 "go.yaml.in/yaml/v3"
)

// Encoder is the upstream streaming YAML encoder.
type Encoder = yamlv3.Encoder

// Decoder is the upstream streaming YAML decoder.
type Decoder = yamlv3.Decoder

var errNoDocuments = errors.New("yaml: input contains no documents")

// Marshal encodes v as two-space-indented YAML.
func Marshal(v any) ([]byte, error) {
	var output bytes.Buffer
	encoder := NewEncoder(&output)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// Unmarshal decodes the first YAML document into v.
func Unmarshal(data []byte, v any) error {
	return yamlv3.Unmarshal(data, v)
}

// Format parses and formats every document in a YAML stream.
func Format(data []byte) ([]byte, error) {
	documents, err := decodeDocuments(data)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	encoder := NewEncoder(&output)
	for i := range documents {
		root := documents[i].Content[0]
		// Without content, the encoder emits no first document marker.
		if root.Kind == yamlv3.ScalarNode && root.Style == 0 && root.Tag == "!!null" && root.Value == "" {
			root.Value = "null"
		}
		if err := encoder.Encode(&documents[i]); err != nil {
			return nil, err
		}
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// Valid reports whether data is a parseable stream with at least one YAML document.
// It checks syntax only. Recursive aliases and duplicate keys may fail during decoding.
func Valid(data []byte) bool {
	_, err := decodeDocuments(data)
	return err == nil
}

// NewEncoder returns a YAML encoder configured with two-space indentation.
func NewEncoder(writer io.Writer) *Encoder {
	encoder := yamlv3.NewEncoder(writer)
	encoder.SetIndent(2)
	return encoder
}

// NewDecoder returns a YAML decoder reading from reader.
func NewDecoder(reader io.Reader) *Decoder {
	return yamlv3.NewDecoder(reader)
}

func decodeDocuments(data []byte) ([]yamlv3.Node, error) {
	decoder := yamlv3.NewDecoder(bytes.NewReader(data))
	var documents []yamlv3.Node
	for {
		var document yamlv3.Node
		err := decoder.Decode(&document)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	if len(documents) == 0 {
		return nil, errNoDocuments
	}
	return documents, nil
}
