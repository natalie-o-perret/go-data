package anonymize

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	data "github.com/natalie-o-perret/go-data"
)

func TestIdentifiersKeepsReferencesConsistent(t *testing.T) {
	const sharedID = "11111111-2222-4333-8444-555555555555"
	input := map[string]any{
		"id":    sharedID,
		"name":  "example-volume",
		"owner": map[string]any{"id": sharedID},
		"items": []any{
			map[string]any{"instanceId": "instance-7", "numeric_id": uint64(42)},
			map[string]any{"instanceId": "instance-7", "numeric_id": uint64(42)},
		},
		"by_id": map[string]any{sharedID: "present"},
	}

	result, err := identifiers(input, bytes.NewReader(bytes.Repeat([]byte{0x42}, 4096)))
	if err != nil {
		t.Fatal(err)
	}
	output := result.(map[string]any)
	anonymousID := output["id"].(string)
	if anonymousID == sharedID || !uuidPattern.MatchString(anonymousID) {
		t.Fatalf("anonymized id = %q", anonymousID)
	}
	if got := output["owner"].(map[string]any)["id"]; got != anonymousID {
		t.Fatalf("repeated id = %q, want %q", got, anonymousID)
	}
	items := output["items"].([]any)
	first, second := items[0].(map[string]any), items[1].(map[string]any)
	if first["instanceId"] == "instance-7" || first["instanceId"] != second["instanceId"] {
		t.Fatalf("instance ids = %q and %q", first["instanceId"], second["instanceId"])
	}
	if first["numeric_id"] == uint64(42) || first["numeric_id"] != second["numeric_id"] {
		t.Fatalf("numeric ids = %v and %v", first["numeric_id"], second["numeric_id"])
	}
	if output["by_id"].(map[string]any)[anonymousID] != "present" {
		t.Fatal("identifier object key was not mapped consistently")
	}
	if output["name"] != "example-volume" || input["id"] != sharedID {
		t.Fatal("non-identifier changed or input was mutated")
	}
}

func TestIdentifiersRejectsCyclesAndMissingRandomSource(t *testing.T) {
	cycle := map[string]any{}
	cycle["self"] = cycle
	if _, err := identifiers(cycle, strings.NewReader("random")); !errors.Is(err, data.ErrRecursive) {
		t.Fatalf("cycle error = %v", err)
	}
	if _, err := identifiers(map[string]any{}, nil); err == nil {
		t.Fatal("nil random source was accepted")
	}
}

func TestIdentifierName(t *testing.T) {
	for _, name := range []string{"id", "IDs", "organization_id", "instance-id", "request.id", "instanceId", "resourceUUID"} {
		if !identifierName(name) {
			t.Errorf("identifierName(%q) = false", name)
		}
	}
	for _, name := range []string{"valid", "grid", "openid", "identity", "name"} {
		if identifierName(name) {
			t.Errorf("identifierName(%q) = true", name)
		}
	}
}
