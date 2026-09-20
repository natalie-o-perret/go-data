package json_test

import (
	"bytes"
	"strings"
	"testing"

	json "github.com/natalie-o-perret/go-data/json"
)

func TestMarshal(t *testing.T) {
	got, err := json.Marshal(map[string]any{"text": "<café & tea>"})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"text":"<café & tea>"}`; string(got) != want {
		t.Fatalf("Marshal() = %q, want %q", got, want)
	}
}

func TestMarshalError(t *testing.T) {
	if _, err := json.Marshal(make(chan int)); err == nil {
		t.Fatal("Marshal() accepted an unsupported value")
	}
}

func TestUnmarshal(t *testing.T) {
	var got struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(`{"name":"Ada"}`), &got); err != nil {
		t.Fatal(err)
	}
	if got.Name != "Ada" {
		t.Fatalf("Unmarshal() name = %q, want Ada", got.Name)
	}
	if err := json.Unmarshal([]byte(`{`), &got); err == nil {
		t.Fatal("Unmarshal() accepted invalid JSON")
	}
}

func TestStreamConstructors(t *testing.T) {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	if err := encoder.Encode(map[string]string{"text": "<>&"}); err != nil {
		t.Fatal(err)
	}
	if want := "{\"text\":\"<>&\"}\n"; output.String() != want {
		t.Fatalf("Encode() = %q, want %q", output.String(), want)
	}

	decoder := json.NewDecoder(strings.NewReader("1 2"))
	for _, want := range []float64{1, 2} {
		var got float64
		if err := decoder.Decode(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("Decode() = %v, want %v", got, want)
		}
	}
}
