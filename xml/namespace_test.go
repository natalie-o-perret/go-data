package xml

import (
	"bytes"
	stdxml "encoding/xml"
	"io"
	"testing"
)

func TestFormatPreservesNamespacesAndTokens(t *testing.T) {
	input := []byte(`<?work now?><!DOCTYPE root><root xmlns:p="urn:p" p:id="7" plain="ok"><!-- note --><p:item> x &amp; y </p:item></root>`)
	formatted, err := Format(input)
	if err != nil {
		t.Fatal(err)
	}
	want := "<?work now?>\n<!DOCTYPE root>\n<root xmlns:p=\"urn:p\" p:id=\"7\" plain=\"ok\">\n  <!-- note -->\n  <p:item> x &amp; y </p:item>\n</root>\n"
	if string(formatted) != want {
		t.Fatalf("Format() = %q, want %q", formatted, want)
	}

	decoder := stdxml.NewDecoder(bytes.NewReader(formatted))
	var comment, directive, instruction, text bool
	var namespacedElement, namespacedAttribute bool
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		switch value := token.(type) {
		case stdxml.Comment:
			comment = string(value) == " note "
		case stdxml.Directive:
			directive = string(value) == "DOCTYPE root"
		case stdxml.ProcInst:
			instruction = value.Target == "work" && string(value.Inst) == "now"
		case stdxml.CharData:
			text = text || string(value) == " x & y "
		case stdxml.StartElement:
			namespacedElement = namespacedElement || value.Name == (stdxml.Name{Space: "urn:p", Local: "item"})
			for _, attribute := range value.Attr {
				namespacedAttribute = namespacedAttribute ||
					(attribute.Name == (stdxml.Name{Space: "urn:p", Local: "id"}) && attribute.Value == "7")
			}
		}
	}
	if !comment || !directive || !instruction || !text || !namespacedElement || !namespacedAttribute {
		t.Fatalf("formatted tokens were not preserved: %s", formatted)
	}
}
