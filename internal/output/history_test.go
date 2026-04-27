package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func makeEntry(release string, drifted bool, count int) HistoryEntry {
	return HistoryEntry{
		Timestamp:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Release:    release,
		Namespace:  "default",
		Drifted:    drifted,
		DriftCount: count,
	}
}

func TestHistory_AddAndLast(t *testing.T) {
	h := NewHistory()
	h.Add(makeEntry("app", true, 3))
	h.Add(makeEntry("app", false, 0))

	if got := len(h.Entries); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
	last := h.Last()
	if last == nil {
		t.Fatal("expected non-nil last entry")
	}
	if last.DriftCount != 0 {
		t.Errorf("expected last drift count 0, got %d", last.DriftCount)
	}
}

func TestHistory_LastOnEmpty(t *testing.T) {
	h := NewHistory()
	if h.Last() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestSaveAndLoadHistory_RoundTrip(t *testing.T) {
	h := NewHistory()
	h.Add(makeEntry("svc", true, 5))

	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	if err := SaveHistory(path, h); err != nil {
		t.Fatalf("SaveHistory: %v", err)
	}
	loaded, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Release != "svc" {
		t.Errorf("unexpected release: %s", loaded.Entries[0].Release)
	}
}

func TestLoadHistory_MissingFileReturnsEmpty(t *testing.T) {
	h, err := LoadHistory(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestLoadHistory_CorruptFileReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	if _, err := LoadHistory(path); err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestWriteHistoryReport_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteHistoryReport(&buf, NewHistory()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no history") {
		t.Errorf("expected 'no history' message, got: %s", buf.String())
	}
}

func TestWriteHistoryReport_WithEntries(t *testing.T) {
	h := NewHistory()
	h.Add(makeEntry("frontend", true, 2))
	h.Add(makeEntry("backend", false, 0))

	var buf bytes.Buffer
	if err := WriteHistoryReport(&buf, h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"frontend", "backend", "yes", "no", "TIMESTAMP", "DRIFTED"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}
