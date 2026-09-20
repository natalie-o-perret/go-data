package toml

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	type config struct {
		Name    string `toml:"name"`
		Enabled bool   `toml:"enabled"`
		Ports   []int  `toml:"ports"`
	}
	want := config{Name: "api", Enabled: true, Ports: []int{80, 443}}

	data, err := Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got config
	if err := Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	if _, err := Marshal(make(chan int)); err == nil {
		t.Fatal("Marshal accepted an unsupported value")
	}
	if err := Unmarshal([]byte("value = ["), &got); err == nil {
		t.Fatal("Unmarshal accepted malformed TOML")
	}
}

func TestMarshalOrdersMapKeys(t *testing.T) {
	got, err := Marshal(map[string]any{"z": int64(2), "a": int64(1)})
	if err != nil {
		t.Fatal(err)
	}
	if want := "a = 1\nz = 2\n"; string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStreamConstructors(t *testing.T) {
	var output bytes.Buffer
	encoder := NewEncoder(&output)
	if err := encoder.SetArraysMultiline(false).Encode(map[string]any{"name": "api"}); err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	decoder := NewDecoder(strings.NewReader(output.String()))
	if err := decoder.Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got["name"] != "api" {
		t.Fatalf("got %#v", got)
	}
}
