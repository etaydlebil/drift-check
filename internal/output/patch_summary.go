package output

import (
	"fmt"
	"strings"
)

// PatchSummary holds counts of added, removed, and changed lines in a diff.
type PatchSummary struct {
	Added   int
	Removed int
	Changed int
}

// NewPatchSummary parses a unified diff string and returns a PatchSummary
// with counts of added and removed lines. Changed lines are counted as pairs
// where an add immediately follows a remove.
func NewPatchSummary(diff string) PatchSummary {
	if diff == "" {
		return PatchSummary{}
	}

	var added, removed int
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			added++
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			removed++
		}
	}

	changed := min(added, removed)
	pureAdded := added - changed
	pureRemoved := removed - changed

	return PatchSummary{
		Added:   pureAdded,
		Removed: pureRemoved,
		Changed: changed,
	}
}

// String returns a human-readable one-line summary of the patch.
func (p PatchSummary) String() string {
	if p.Added == 0 && p.Removed == 0 && p.Changed == 0 {
		return "no changes"
	}
	parts := make([]string, 0, 3)
	if p.Changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", p.Changed))
	}
	if p.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", p.Added))
	}
	if p.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", p.Removed))
	}
	return strings.Join(parts, ", ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
