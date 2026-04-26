package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/drift-check/internal/compare"
)

func makeRemediationResults(drifted bool) []compare.DriftResult {
	if !drifted {
		return []compare.DriftResult{
			{Name: "my-app", Namespace: "default", Kind: "Deployment", Diff: ""},
		}
	}
	return []compare.DriftResult{
		{Name: "my-app", Namespace: "default", Kind: "Deployment", Diff: "-replicas: 1\n+replicas: 3"},
		{Name: "my-svc", Namespace: "staging", Kind: "Service", Diff: "-port: 80\n+port: 8080"},
	}
}

func TestNewRemediationHints_NoDrift(t *testing.T) {
	hints := NewRemediationHints(makeRemediationResults(false))
	if len(hints) != 0 {
		t.Fatalf("expected 0 hints, got %d", len(hints))
	}
}

func TestNewRemediationHints_WithDrift(t *testing.T) {
	hints := NewRemediationHints(makeRemediationResults(true))
	if len(hints) != 2 {
		t.Fatalf("expected 2 hints, got %d", len(hints))
	}
	if hints[0].Resource != "my-app" {
		t.Errorf("expected resource my-app, got %s", hints[0].Resource)
	}
	if hints[1].Namespace != "staging" {
		t.Errorf("expected namespace staging, got %s", hints[1].Namespace)
	}
}

func TestNewRemediationHints_CommandContainsNamespace(t *testing.T) {
	hints := NewRemediationHints(makeRemediationResults(true))
	for _, h := range hints {
		if !strings.Contains(h.Command, h.Namespace) {
			t.Errorf("command %q does not contain namespace %q", h.Command, h.Namespace)
		}
	}
}

func TestWriteRemediationReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRemediationReport(&buf, makeRemediationResults(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No remediation required") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteRemediationReport_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRemediationReport(&buf, makeRemediationResults(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Remediation suggestions") {
		t.Errorf("expected header, got: %s", out)
	}
	if !strings.Contains(out, "Deployment/my-app") {
		t.Errorf("expected Deployment/my-app in output, got: %s", out)
	}
	if !strings.Contains(out, "Service/my-svc") {
		t.Errorf("expected Service/my-svc in output, got: %s", out)
	}
}

func TestWriteRemediationReport_WriterError(t *testing.T) {
	err := WriteRemediationReport(errWriter{}, makeRemediationResults(false))
	if err == nil {
		t.Fatal("expected error from failing writer")
	}
}
