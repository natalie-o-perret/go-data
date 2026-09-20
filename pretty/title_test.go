package pretty

import "testing"

func TestTitleSplitsCommonIdentifierStyles(t *testing.T) {
	tests := map[string]string{
		"snake_case":   "Snake Case",
		"spinal-case":  "Spinal Case",
		"camelCase":    "Camel Case",
		"PascalCase":   "Pascal Case",
		"HTTPServerID": "HTTP Server ID",
		"IPv6Address":  "IPv6 Address",
	}
	for input, want := range tests {
		if got := title(input); got != want {
			t.Errorf("title(%q) = %q, want %q", input, got, want)
		}
	}
}
