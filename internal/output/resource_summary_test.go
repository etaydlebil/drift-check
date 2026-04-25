package output

import (
	"bytes"
	"strings"
	"testing"
)

func makeDriftResult(kind, name string, diff string) DriftResult {
	return DriftResult{Kind: kind, Name: name, Diff: diff}
}

func TestNewResourceSummary_Empty(t *testing.T) {
	rs := NewResourceSummary(nil)
	if rs.Total != 0 || rs.Drifted != 0 {
		t.Errorf("expected zero counts, got total=%d drifted=%d", rs.Total, rs.Drifted)
	}
}

func TestNewResourceSummary_NoDrift(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("Deployment", "api", ""),
		makeDriftResult("Service", "svc", ""),
	}
	rs := NewResourceSummary(results)
	if rs.Total != 2 {
		t.Errorf("expected total=2, got %d", rs.Total)
	}
	if rs.Drifted != 0 {
		t.Errorf("expected drifted=0, got %d", rs.Drifted)
	}
}

func TestNewResourceSummary_WithDrift(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("Deployment", "api", "- old\n+ new"),
		makeDriftResult("Deployment", "worker", ""),
		makeDriftResult("ConfigMap", "cfg", "- a\n+ b"),
	}
	rs := NewResourceSummary(results)
	if rs.Total != 3 {
		t.Errorf("expected total=3, got %d", rs.Total)
	}
	if rs.Drifted != 2 {
		t.Errorf("expected drifted=2, got %d", rs.Drifted)
	}
	if rs.KindCounts["Deployment"].Drifted != 1 {
		t.Errorf("expected Deployment drifted=1, got %d", rs.KindCounts["Deployment"].Drifted)
	}
	if rs.KindCounts["ConfigMap"].Drifted != 1 {
		t.Errorf("expected ConfigMap drifted=1, got %d", rs.KindCounts["ConfigMap"].Drifted)
	}
}

func TestWriteResourceSummary_ContainsKinds(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("Deployment", "api", "diff"),
		makeDriftResult("Service", "svc", ""),
	}
	rs := NewResourceSummary(results)
	var buf bytes.Buffer
	if err := WriteResourceSummary(&buf, rs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Deployment", "Service", "TOTAL", "KIND"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWriteResourceSummary_WriterError(t *testing.T) {
	rs := NewResourceSummary([]DriftResult{makeDriftResult("Pod", "p", "")})
	err := WriteResourceSummary(&errWriter{}, rs)
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}
