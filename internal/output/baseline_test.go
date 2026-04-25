package output

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBaseline_Empty(t *testing.T) {
	b := NewBaseline()
	if b == nil {
		t.Fatal("expected non-nil baseline")
	}
	if len(b.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(b.Entries))
	}
}

func TestBaseline_SetAndGet(t *testing.T) {
	b := NewBaseline()
	b.Set("my-app", "default", "--- a\n+++ b\n")

	e, ok := b.Get("my-app", "default")
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if e.Release != "my-app" {
		t.Errorf("release: got %q, want %q", e.Release, "my-app")
	}
	if e.Namespace != "default" {
		t.Errorf("namespace: got %q, want %q", e.Namespace, "default")
	}
	if e.Diff != "--- a\n+++ b\n" {
		t.Errorf("diff mismatch: %q", e.Diff)
	}
	if e.Captured.IsZero() {
		t.Error("expected non-zero captured timestamp")
	}
}

func TestBaseline_GetMissing(t *testing.T) {
	b := NewBaseline()
	_, ok := b.Get("missing", "default")
	if ok {
		t.Error("expected missing entry to return false")
	}
}

func TestBaseline_SetOverwrites(t *testing.T) {
	b := NewBaseline()
	b.Set("app", "ns", "first")
	time.Sleep(time.Millisecond)
	b.Set("app", "ns", "second")

	e, _ := b.Get("app", "ns")
	if e.Diff != "second" {
		t.Errorf("expected overwritten diff %q, got %q", "second", e.Diff)
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	b := NewBaseline()
	b.Set("nginx", "production", "+replica: 3\n")

	tmp := filepath.Join(t.TempDir(), "baseline.json")
	if err := SaveBaseline(tmp, b); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(tmp)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	e, ok := loaded.Get("nginx", "production")
	if !ok {
		t.Fatal("expected entry after round-trip")
	}
	if e.Diff != "+replica: 3\n" {
		t.Errorf("diff mismatch after round-trip: %q", e.Diff)
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveBaseline_InvalidPath(t *testing.T) {
	b := NewBaseline()
	err := SaveBaseline(filepath.Join(t.TempDir(), "no", "such", "dir", "b.json"), b)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestLoadBaseline_NilEntries(t *testing.T) {
	// Write a JSON file with null entries to exercise the nil guard.
	tmp := filepath.Join(t.TempDir(), "null.json")
	if err := os.WriteFile(tmp, []byte(`{"entries":null}`), 0o644); err != nil {
		t.Fatal(err)
	}
	b, err := LoadBaseline(tmp)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if b.Entries == nil {
		t.Error("expected non-nil entries map after load")
	}
}
