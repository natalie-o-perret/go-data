package xml

import (
	"bytes"
	stdxml "encoding/xml"
	"io"
	"reflect"
	"testing"
)

type codecDocument struct {
	XMLName stdxml.Name `xml:"document"`
	ID      string      `xml:"id,attr"`
	Items   []string    `xml:"item"`
}

var (
	_ *Encoder = (*stdxml.Encoder)(nil)
	_ *Decoder = (*stdxml.Decoder)(nil)
)

func TestMarshalAndUnmarshal(t *testing.T) {
	want := codecDocument{XMLName: stdxml.Name{Local: "document"}, ID: "7", Items: []string{"one", "two"}}
	data, err := Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `<document id="7"><item>one</item><item>two</item></document>` {
		t.Fatalf("Marshal() = %s", data)
	}

	var got codecDocument
	if err := Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unmarshal() = %#v, want %#v", got, want)
	}
	if _, err := Marshal(make(chan int)); err == nil {
		t.Fatal("Marshal() accepted an unsupported value")
	}
	if err := Unmarshal([]byte(`<document>`), &got); err == nil {
		t.Fatal("Unmarshal() accepted malformed XML")
	}
}

func TestStreamCodec(t *testing.T) {
	var output bytes.Buffer
	encoder := NewEncoder(&output)
	if err := encoder.Encode(codecDocument{ID: "1"}); err != nil {
		t.Fatal(err)
	}
	if output.String() != `<document id="1"></document>` {
		t.Fatalf("Encode() = %q", output.String())
	}

	decoder := NewDecoder(bytes.NewBufferString(`<item>a</item><item>b</item>`))
	for _, want := range []string{"a", "b"} {
		var got string
		if err := decoder.Decode(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("Decode() = %q, want %q", got, want)
		}
	}
	var extra string
	if err := decoder.Decode(&extra); err != io.EOF {
		t.Fatalf("final Decode() error = %v, want EOF", err)
	}
}
