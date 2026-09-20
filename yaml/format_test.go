package yaml

import (
	"strings"
	"testing"
)

func TestFormatStreamPreservesNodes(t *testing.T) {
	input := "# first\nz: &item \"quoted\"\na: *item\n---\nlist:\n    - value\n"
	formatted, err := Format([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	output := string(formatted)
	for _, fragment := range []string{
		"# first\n", `z: &item "quoted"`, "a: *item", "\n---\n", "list:\n  - value\n",
	} {
		if !strings.Contains(output, fragment) {
			t.Errorf("Format() output %q does not contain %q", output, fragment)
		}
	}
	if strings.Index(output, "z:") > strings.Index(output, "a:") {
		t.Errorf("Format() changed mapping order: %q", output)
	}
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Format() output is not newline-terminated: %q", output)
	}
}

func TestFormatRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{"", "ok: true\n---\n[", "ok: true\n\tbad"} {
		if output, err := Format([]byte(input)); err == nil || output != nil {
			t.Errorf("Format(%q) = %q, %v", input, output, err)
		}
	}
}

func TestFormatPreservesEmptyDocuments(t *testing.T) {
	formatted, err := Format([]byte("---\n---\n"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "null\n---\nnull\n"; string(formatted) != want {
		t.Fatalf("Format() = %q, want %q", formatted, want)
	}
	if !Valid(formatted) {
		t.Fatalf("Format() returned invalid YAML: %q", formatted)
	}
	reformatted, err := Format(formatted)
	if err != nil {
		t.Fatal(err)
	}
	if string(reformatted) != string(formatted) {
		t.Fatalf("Format() is not idempotent: %q, then %q", formatted, reformatted)
	}
}
