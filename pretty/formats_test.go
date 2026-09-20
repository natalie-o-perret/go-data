package pretty

import (
	"os"
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-data/json"
	"github.com/natalie-o-perret/go-data/toml"
	"github.com/natalie-o-perret/go-data/xml"
	"github.com/natalie-o-perret/go-data/yaml"
)

type xmlFixture struct {
	Catalog struct {
		Name string `xml:"name"`
	} `xml:"catalog"`
	Publications []struct {
		ID        string `xml:"id"`
		Title     string `xml:"title"`
		FileBytes int64  `xml:"file_bytes"`
	} `xml:"publications"`
}

func decodeXML(data []byte) (any, error) {
	var decoded xmlFixture
	if err := xml.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}
	publications := make([]any, len(decoded.Publications))
	for i, publication := range decoded.Publications {
		publications[i] = map[string]any{
			"id": publication.ID, "title": publication.Title, "file_bytes": publication.FileBytes,
		}
	}
	return map[string]any{
		"catalog":      map[string]any{"name": decoded.Catalog.Name},
		"publications": publications,
	}, nil
}

func TestDecodedFormatsRenderIdentically(t *testing.T) {
	tests := map[string]func([]byte) (any, error){
		"json": func(data []byte) (any, error) {
			var value any
			err := json.Unmarshal(data, &value)
			return value, err
		},
		"yaml": func(data []byte) (any, error) {
			var value any
			err := yaml.Unmarshal(data, &value)
			return value, err
		},
		"toml": func(data []byte) (any, error) {
			var value any
			err := toml.Unmarshal(data, &value)
			return value, err
		},
		"xml": decodeXML,
	}

	want := "├── Catalog\n" +
		"│   └── Name: Example Press\n" +
		"├──────────────────────────────────────╮\n" +
		"│ Publications                         │\n" +
		"├────────────┬─────────┬───────────────┤\n" +
		"│ File Bytes │ Id      │ Title         │\n" +
		"├────────────┼─────────┼───────────────┤\n" +
		"│      8 MiB │ pub-001 │ Field Notes   │\n" +
		"├────────────┼─────────┼───────────────┤\n" +
		"│      4 MiB │ pub-002 │ Short Stories │\n" +
		"╰────────────┴─────────┴───────────────╯\n"

	for format, decode := range tests {
		t.Run(format, func(t *testing.T) {
			data, err := os.ReadFile("examples/catalog." + format)
			if err != nil {
				t.Fatal(err)
			}
			value, err := decode(data)
			if err != nil {
				t.Fatal(err)
			}
			var out strings.Builder
			if err := Render(&out, value); err != nil {
				t.Fatal(err)
			}
			if got := out.String(); got != want {
				t.Fatalf("output mismatch:\nwant:\n%s\ngot:\n%s", want, got)
			}
		})
	}
}
