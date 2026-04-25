package output

import (
	"testing"
)

func TestNewPatchSummary_EmptyDiff(t *testing.T) {
	s := NewPatchSummary("")
	if s.Added != 0 || s.Removed != 0 || s.Changed != 0 {
		t.Fatalf("expected zero counts, got %+v", s)
	}
}

func TestNewPatchSummary_OnlyAdded(t *testing.T) {
	diff := `--- a\n+++ b\n+line1\n+line2`
	s := NewPatchSummary(diff)
	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 0 || s.Changed != 0 {
		t.Errorf("unexpected removed/changed: %+v", s)
	}
}

func TestNewPatchSummary_OnlyRemoved(t *testing.T) {
	diff := "--- a\n+++ b\n-line1\n-line2\n-line3"
	s := NewPatchSummary(diff)
	if s.Removed != 3 {
		t.Errorf("expected 3 removed, got %d", s.Removed)
	}
	if s.Added != 0 || s.Changed != 0 {
		t.Errorf("unexpected added/changed: %+v", s)
	}
}

func TestNewPatchSummary_ChangedPairs(t *testing.T) {
	// 2 removed + 2 added => 2 changed, 0 pure added, 0 pure removed
	diff := "--- a\n+++ b\n-old1\n-old2\n+new1\n+new2"
	s := NewPatchSummary(diff)
	if s.Changed != 2 {
		t.Errorf("expected 2 changed, got %d", s.Changed)
	}
	if s.Added != 0 || s.Removed != 0 {
		t.Errorf("expected no pure adds/removes: %+v", s)
	}
}

func TestNewPatchSummary_MixedChanges(t *testing.T) {
	// 1 removed + 3 added => 1 changed, 2 pure added
	diff := "--- a\n+++ b\n-old1\n+new1\n+new2\n+new3"
	s := NewPatchSummary(diff)
	if s.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", s.Changed)
	}
	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", s.Removed)
	}
}

func TestPatchSummary_String_NoChanges(t *testing.T) {
	s := PatchSummary{}
	if s.String() != "no changes" {
		t.Errorf("expected 'no changes', got %q", s.String())
	}
}

func TestPatchSummary_String_WithChanges(t *testing.T) {
	s := PatchSummary{Added: 2, Removed: 1, Changed: 3}
	out := s.String()
	for _, want := range []string{"3 changed", "2 added", "1 removed"} {
		if !contains(out, want) {
			t.Errorf("expected %q in output %q", want, out)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
