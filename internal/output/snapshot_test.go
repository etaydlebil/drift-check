package output_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/drift-check/internal/compare"
	"github.com/example/drift-check/internal/output"
)

func makeResults(diffs ...string) []compare.DriftResult {
	results := make([]compare.DriftResult, len(diffs))
	for i, d := range diffs {
		results[i] = compare.DriftResult{
			Name:      "resource",
			Namespace: "default",
			Kind:      "Deployment",
			Diff:      d,
		}
	}
	return results
}

func TestNewSnapshot_FieldsSet(t *testing.T) {
	before := time.Now().UTC()
	s := output.NewSnapshot("my-release", "staging", makeResults("", "+ added"))
	after := time.Now().UTC()

	if s.Release != "my-release" {
		t.Errorf("expected release my-release, got %q", s.Release)
	}
	if s.Namespace != "staging" {
		t.Errorf("expected namespace staging, got %q", s.Namespace)
	}
	if len(s.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(s.Results))
	}
	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("captured_at %v outside expected range", s.CapturedAt)
	}
}

func TestSnapshot_DriftedCount(t *testing.T) {
	tests := []struct {
		name     string
		diffs    []string
		wantCount int
	}{
		{"no results", []string{}, 0},
		{"all clean", []string{"", ""}, 0},
		{"one drifted", []string{"", "+ field: new"}, 1},
		{"all drifted", []string{"- a", "+ b"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := output.NewSnapshot("r", "ns", makeResults(tt.diffs...))
			if got := s.DriftedCount(); got != tt.wantCount {
				t.Errorf("DriftedCount() = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := output.NewSnapshot("rel", "default", makeResults("+ x: 1"))
	if err := output.SaveSnapshot(path, orig); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := output.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if loaded.Release != orig.Release {
		t.Errorf("release mismatch: got %q want %q", loaded.Release, orig.Release)
	}
	if loaded.DriftedCount() != orig.DriftedCount() {
		t.Errorf("drifted count mismatch: got %d want %d", loaded.DriftedCount(), orig.DriftedCount())
	}
}

func TestLoadSnapshot_MissingFileReturnsError(t *testing.T) {
	_, err := output.LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSaveSnapshot_UnwritablePathReturnsError(t *testing.T) {
	s := output.NewSnapshot("r", "ns", nil)
	err := output.SaveSnapshot("/nonexistent/dir/snap.json", s)
	if err == nil {
		t.Error("expected error for unwritable path, got nil")
	}
	_ = os.Remove("/nonexistent/dir/snap.json")
}
