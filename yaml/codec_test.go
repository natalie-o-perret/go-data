package yaml

import (
	"errors"
	"reflect"
	"testing"
)

var errMarshal = errors.New("marshal failed")

type errorMarshaler struct{}

func (errorMarshaler) MarshalYAML() (any, error) { return nil, errMarshal }

func TestMarshalAndUnmarshal(t *testing.T) {
	input := map[string]any{"outer": map[string]any{"value": 1}}
	data, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "outer:\n  value: 1\n"; string(data) != want {
		t.Fatalf("Marshal() = %q, want %q", data, want)
	}

	var output map[string]any
	if err := Unmarshal(data, &output); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(output, input) {
		t.Fatalf("Unmarshal() = %#v, want %#v", output, input)
	}
}

func TestMarshalAndUnmarshalErrors(t *testing.T) {
	if _, err := Marshal(errorMarshaler{}); !errors.Is(err, errMarshal) {
		t.Fatalf("Marshal() error = %v, want %v", err, errMarshal)
	}
	var output any
	if err := Unmarshal([]byte("["), &output); err == nil {
		t.Fatal("Unmarshal() accepted malformed YAML")
	}
}

func TestUnmarshalDecodesFirstDocument(t *testing.T) {
	var output map[string]int
	if err := Unmarshal([]byte("one: 1\n---\n["), &output); err != nil {
		t.Fatal(err)
	}
	if want := map[string]int{"one": 1}; !reflect.DeepEqual(output, want) {
		t.Fatalf("Unmarshal() = %#v, want %#v", output, want)
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"document", "value: true\n", true},
		{"stream", "one: 1\n---\ntwo: 2\n", true},
		{"duplicate keys", "value: 1\nvalue: 2\n", true},
		{"recursive alias", "&self [*self]\n", true},
		{"empty", "", false},
		{"whitespace", "  \n", false},
		{"malformed", "[", false},
		{"malformed second document", "one: 1\n---\n[", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Valid([]byte(test.input)); got != test.want {
				t.Fatalf("Valid() = %t, want %t", got, test.want)
			}
		})
	}
}
