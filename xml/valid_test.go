package xml

import "testing"

func TestValidDocuments(t *testing.T) {
	inputs := []string{
		`<root/>`,
		` <!--before--><?work now?><root></root><!--after--> `,
	}
	for _, input := range inputs {
		if !Valid([]byte(input)) {
			t.Errorf("Valid(%q) = false", input)
		}
	}
}

func TestInvalidDocuments(t *testing.T) {
	tests := []struct {
		input string
		error string
	}{
		{"", "xml: document must contain exactly one root element"},
		{" \n\t", "xml: document must contain exactly one root element"},
		{`<!--none-->`, "xml: document must contain exactly one root element"},
		{`<a/><b/>`, "xml: document must contain exactly one root element"},
		{`text<root/>`, "xml: character data outside the root element"},
		{`<root/>text`, "xml: character data outside the root element"},
		{` <?xml version="1.0"?><root/>`, "xml: XML declaration must be the first token"},
		{`<root>`, ""},
		{`</root>`, ""},
	}

	for _, test := range tests {
		if Valid([]byte(test.input)) {
			t.Errorf("Valid(%q) = true", test.input)
		}
		_, err := Format([]byte(test.input))
		if err == nil {
			t.Errorf("Format(%q) succeeded", test.input)
			continue
		}
		if test.error != "" && err.Error() != test.error {
			t.Errorf("Format(%q) error = %q, want %q", test.input, err, test.error)
		}
	}
}
