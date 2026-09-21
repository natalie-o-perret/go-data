# go-data

`go-data` provides small, consistent packages for structured data formats and
deterministic traversal of decoded Go values.

## Packages

| Package | Purpose |
| --- | --- |
| `data` | Shared deterministic traversal and non-mutating transformation |
| `data/json` | JSON codec and formatting |
| `data/yaml` | YAML v3 codec and formatting |
| `data/toml` | TOML codec and deterministic formatting |
| `data/xml` | Standard-library XML codec and formatting |
| `data/pretty` | Deterministic tree and table rendering for decoded values |
| `data/barcode` | Common 1D and 2D barcode encoding and decoding |

Package import paths use `github.com/natalie-o-perret/go-data` as their base.

## Installation

Install the shared traversal package:

```sh
go get github.com/natalie-o-perret/go-data
```

Or install the format package you need:

```sh
go get github.com/natalie-o-perret/go-data/json
go get github.com/natalie-o-perret/go-data/yaml
go get github.com/natalie-o-perret/go-data/toml
go get github.com/natalie-o-perret/go-data/xml
go get github.com/natalie-o-perret/go-data/pretty
go get github.com/natalie-o-perret/go-data/barcode
```

All packages live in one Go module. The module lists dependencies needed across
all packages, but Go does not link unused packages into your program.

## Codec API

The four codec packages expose the same core API:

```go
func Marshal(value any) ([]byte, error)
func Unmarshal(input []byte, value any) error
func Format(input []byte) ([]byte, error)
func Valid(input []byte) bool
func NewEncoder(writer io.Writer) *Encoder
func NewDecoder(reader io.Reader) *Decoder
```

JSON additionally provides `Compact`. Each package aliases the encoder and
decoder types from its standard or upstream implementation.

Formatting policies:

- `json`: two spaces, stable object keys, preserved number tokens, and no HTML
  escaping.
- `yaml`: two spaces and complete multi-document streams. Comments and node
  styles are retained where YAML v3 supports them.
- `toml`: deterministic value-based output. Comments and lexical choices are
  not preserved.
- `xml`: two spaces and one root element. Mixed content and lexical namespace
  prefixes are preserved.

XML needs typed structs for useful marshaling and unmarshaling because
`encoding/xml` has no generic `map[string]any` representation.

## Barcodes

`barcode.Encode` renders common one-dimensional formats, QR codes, and Data
Matrix codes as `image.Image` values. `barcode.Decode` detects those formats,
plus Aztec and RSS-14, in any `image.Image`:

```go
img, err := barcode.Encode("hello", barcode.QRCode, 256, 256)
result, err := barcode.Decode(img)
```

## Traversal

The root package exposes the shared traversal API:

```go
type Path []any
type VisitFunc func(Path, any) error
type TransformFunc func(Path, any) (any, error)

var ErrRecursive error

func Walk(any, VisitFunc) error
func Transform(any, TransformFunc) (any, error)
```

`Walk` visits the root and then its descendants. `Transform` clones supported
containers, transforms their descendants, and then visits each clone. Both
descend only into `map[string]any` and `[]any`. Map keys use lexical order,
slice elements use index order, and the root path is nil.

Shared non-cyclic containers are visited at every occurrence. `Transform`
clones each occurrence independently and does not mutate the input. A container
reached again on its active traversal path returns `ErrRecursive`. Callback
errors are returned unchanged.

## Pretty Rendering

`pretty.Render` writes already-decoded values as deterministic trees and
tables, falling back to indented JSON for other values:

```go
err := pretty.Render(os.Stdout, value)
```

It handles nested `map[string]any` and `[]any` values, preserves array order,
and returns an error wrapping `data.ErrRecursive` before writing output when a
recursive container is found.

## Go Versions

The module requires Go 1.25.13 or newer. XML users must use Go 1.25.13, Go
1.26.6, or a newer release to include the minimum safe `encoding/xml`
toolchain fixes.

## License

MIT
