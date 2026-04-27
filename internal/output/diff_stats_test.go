package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestComputeDiffStats_EmptyDiff(t *testing.T) {
	stats := ComputeDiffStats("")
	if stats.Added != 0 || stats.Removed != 0 || stats.Context != 0 {
		t.Fatalf("expected zero stats for empty diff, got %+v", stats)
	}
}

func TestComputeDiffStats_OnlyAdded(t *testing.T) {
	diff := "+++ b/file\n+added line one\n+added line two\n"
	stats := ComputeDiffStats(diff)
	if stats.Added != 2 {
		t.Errorf("expected 2 added, got %d", stats.Added)
	}
	if stats.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", stats.Removed)
	}
}

func TestComputeDiffStats_OnlyRemoved(t *testing.T) {
	diff := "--- a/file\n-removed line\n"
	stats := ComputeDiffStats(diff)
	if stats.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", stats.Removed)
	}
	if stats.Added != 0 {
		t.Errorf("expected 0 added, got %d", stats.Added)
	}
}

func TestComputeDiffStats_Mixed(t *testing.T) {
	diff := "--- a/f\n+++ b/f\n context\n+added\n-removed\n+another add\n"
	stats := ComputeDiffStats(diff)
	if stats.Added != 2 {
		t.Errorf("expected 2 added, got %d", stats.Added)
	}
	if stats.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", stats.Removed)
	}
	if stats.Total() != 3 {
		t.Errorf("expected total 3, got %d", stats.Total())
	}
}

func TestDiffStats_String(t *testing.T) {
	s := DiffStats{Added: 5, Removed: 2}
	if s.String() != "+5 -2" {
		t.Errorf("unexpected String(): %q", s.String())
	}
}

func TestWriteDiffStats_NoDrift(t *testing.T) {
	results := []DriftResult{{Kind: "Deployment", Name: "web", Diff: ""}}
	var buf bytes.Buffer
	if err := WriteDiffStats(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only the header line should be present
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d", len(lines))
	}
}

func TestWriteDiffStats_WithDrift(t *testing.T) {
	results := []DriftResult{
		{Kind: "Deployment", Name: "api", Diff: "+++ b/f\n+added\n-removed\n"},
		{Kind: "Service", Name: "svc", Diff: ""},
	}
	var buf bytes.Buffer
	if err := WriteDiffStats(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Deployment/api") {
		t.Errorf("expected resource key in output, got:\n%s", out)
	}
	if strings.Contains(out, "Service/svc") {
		t.Errorf("did not expect non-drifted resource in output")
	}
}
