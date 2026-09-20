package xml

import "errors"

var (
	errDeclaration = errors.New("xml: XML declaration must be the first token")
	errRoot        = errors.New("xml: document must contain exactly one root element")
	errText        = errors.New("xml: character data outside the root element")
)

// Format validates data, indents nested markup by two spaces without reflowing
// mixed content, and adds one trailing newline.
func Format(data []byte) ([]byte, error) {
	document, err := parseDocument(data)
	if err != nil {
		return nil, err
	}
	return renderDocument(document)
}

// Valid reports whether Format accepts data as a single XML document.
func Valid(data []byte) bool {
	_, err := Format(data)
	return err == nil
}
