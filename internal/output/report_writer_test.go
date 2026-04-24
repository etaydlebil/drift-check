package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewReportWriter_InvalidFormat(t *testing.T) {
	_, err := NewReportWriter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestNewReportWriter_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatText, FormatJSON} {
		_, err := NewReportWriter(f, &bytes.Buffer{})
		if err != nil {
			t.Errorf("format %q: unexpected error: %v", f, err)
		}
	}
}

func TestReportWriter_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	rw, _ := NewReportWriter(FormatText, &buf)

	s := NewSummary("my-release", "", nil)
	if err := rw.Write(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "my-release") {
		t.Errorf("output missing release name: %q", buf.String())
	}
}

func TestReportWriter_JSONWithDrift(t *testing.T) {
	var buf bytes.Buffer
	rw, _ := NewReportWriter(FormatJSON, &buf)

	s := NewSummary("rel", "--- a\n+++ b\n", nil)
	if err := rw.Write(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["release"] != "rel" {
		t.Errorf("release field: got %v", out["release"])
	}
	if out["drift"] != true {
		t.Errorf("drift field: expected true, got %v", out["drift"])
	}
}

func TestReportWriter_Format(t *testing.T) {
	rw, _ := NewReportWriter(FormatJSON, &bytes.Buffer{})
	if rw.Format() != FormatJSON {
		t.Errorf("got %q, want %q", rw.Format(), FormatJSON)
	}
}
