package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteTable_NoDrift(t *testing.T) {
	s := Summary{Release: "my-app", Namespace: "default", Drift: "", Error: ""}
	var buf bytes.Buffer
	if err := WriteTable(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no drift detected") {
		t.Errorf("expected 'no drift detected', got: %s", buf.String())
	}
}

func TestWriteTable_WithDrift(t *testing.T) {
	drift := "-replicas: 1\n+replicas: 3\n"
	s := Summary{Release: "my-app", Namespace: "default", Drift: drift, Error: ""}
	var buf bytes.Buffer
	if err := WriteTable(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FIELD") {
		t.Errorf("expected header row, got: %s", out)
	}
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected drift content, got: %s", out)
	}
}

func TestWriteTable_WithError(t *testing.T) {
	s := Summary{Release: "my-app", Namespace: "default", Drift: "- a: b", Error: "helm error"}
	var buf bytes.Buffer
	if err := WriteTable(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "helm error") {
		t.Errorf("expected error in output, got: %s", buf.String())
	}
}

func TestWriteTable_HeaderAlwaysPresent(t *testing.T) {
	s := Summary{}
	var buf bytes.Buffer
	_ = WriteTable(&buf, s)
	if !strings.Contains(buf.String(), "FIELD") {
		t.Errorf("expected header in all outputs, got: %s", buf.String())
	}
}

func TestWriteTable_ReleaseAndNamespaceInOutput(t *testing.T) {
	s := Summary{Release: "frontend", Namespace: "production", Drift: "", Error: ""}
	var buf bytes.Buffer
	if err := WriteTable(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "frontend") {
		t.Errorf("expected release name in output, got: %s", out)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected namespace in output, got: %s", out)
	}
}
