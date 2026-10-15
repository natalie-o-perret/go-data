package anonymize

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestIdentifiersPreservesTypesAndShapes(t *testing.T) {
	const uppercaseUUID = "AAAAAAAA-BBBB-4CCC-8DDD-EEEEEEEEEEEE"
	input := map[string]any{
		"signed_id": int8(12), "unsigned_id": uint8(34), "floating_id": float64(56),
		"fractional_id": 1.5, "boolean_id": true, "empty_id": "",
		"uuid": uppercaseUUID, "checksum": "ABCDEF0123456789",
		"ids": []any{"alpha-1", "alpha-1"},
	}
	result, err := identifiers(input, bytes.NewReader(bytes.Repeat([]byte{0x42}, 4096)))
	if err != nil {
		t.Fatal(err)
	}
	output := result.(map[string]any)
	for name, original := range map[string]any{
		"signed_id": int8(12), "unsigned_id": uint8(34), "floating_id": float64(56),
	} {
		if got := output[name]; got == original || reflect.TypeOf(got) != reflect.TypeOf(original) {
			t.Errorf("%s = %v (%T), original %v (%T)", name, got, got, original, original)
		}
	}
	for name, expected := range map[string]any{"fractional_id": 1.5, "boolean_id": true, "empty_id": ""} {
		if got := output[name]; got != expected {
			t.Errorf("%s = %v, want %v", name, got, expected)
		}
	}
	if got := output["uuid"].(string); got == uppercaseUUID || got != strings.ToUpper(got) || !uuidPattern.MatchString(got) {
		t.Errorf("uuid = %q", got)
	}
	if got := output["checksum"].(string); got == "ABCDEF0123456789" || !hexPattern.MatchString(got) {
		t.Errorf("checksum = %q", got)
	}
	ids := output["ids"].([]any)
	if ids[0] == "alpha-1" || ids[0] != ids[1] {
		t.Errorf("ids = %#v", ids)
	}
}

func TestIdentifiersReportsRandomSourceErrors(t *testing.T) {
	_, err := identifiers(map[string]any{"id": "value"}, strings.NewReader(""))
	if err == nil || !strings.Contains(err.Error(), "generate identifier") {
		t.Fatalf("error = %v", err)
	}
}

func TestIdentifiersPublicAPI(t *testing.T) {
	const original = "11111111-2222-4333-8444-555555555555"
	result, err := Identifiers(map[string]any{"id": original})
	if err != nil {
		t.Fatal(err)
	}
	got := result.(map[string]any)["id"].(string)
	if got == original || !uuidPattern.MatchString(got) {
		t.Fatalf("anonymized id = %q", got)
	}
}
