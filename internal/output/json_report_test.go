package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func decode(t *testing.T, buf *bytes.Buffer) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.NewDecoder(buf).Decode(&m); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	return m
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	s := NewSummary("app", "default", "", nil)
	if err := WriteJSON(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, &buf)
	if m["drifted"].(bool) {
		t.Error("expected drifted=false")
	}
	if _, ok := m["diff"]; ok {
		t.Error("diff field should be omitted when empty")
	}
	if _, ok := m["error"]; ok {
		t.Error("error field should be omitted when nil")
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	s := NewSummary("app", "default", "-old\n+new", nil)
	if err := WriteJSON(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, &buf)
	if !m["drifted"].(bool) {
		t.Error("expected drifted=true")
	}
	if m["diff"] == "" {
		t.Error("expected diff to be present")
	}
}

func TestWriteJSON_WithError(t *testing.T) {
	var buf bytes.Buffer
	s := NewSummary("app", "default", "", errors.New("boom"))
	if err := WriteJSON(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, &buf)
	if m["error"] != "boom" {
		t.Errorf("unexpected error field: %v", m["error"])
	}
}

func TestWriteJSON_FieldsPresent(t *testing.T) {
	var buf bytes.Buffer
	s := NewSummary("my-release", "staging", "", nil)
	_ = WriteJSON(&buf, s)
	m := decode(t, &buf)
	for _, key := range []string{"release", "namespace", "drifted", "checked_at"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing expected field: %s", key)
		}
	}
}
