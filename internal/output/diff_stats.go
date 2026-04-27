package output

import (
	"fmt"
	"io"
	"strings"
)

// DiffStats holds line-level statistics for a unified diff.
type DiffStats struct {
	Added   int
	Removed int
	Context int
}

// Total returns the total number of changed lines (added + removed).
func (d DiffStats) Total() int {
	return d.Added + d.Removed
}

// String returns a compact human-readable summary, e.g. "+3 -1".
func (d DiffStats) String() string {
	return fmt.Sprintf("+%d -%d", d.Added, d.Removed)
}

// ComputeDiffStats parses a unified diff string and counts added, removed,
// and context lines. Lines beginning with '+' (but not '+++') are added;
// lines beginning with '-' (but not '---') are removed; all other non-header
// lines are treated as context.
func ComputeDiffStats(diff string) DiffStats {
	var stats DiffStats
	if diff == "" {
		return stats
	}
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			// file header – skip
		case strings.HasPrefix(line, "+"):
			stats.Added++
		case strings.HasPrefix(line, "-"):
			stats.Removed++
		case strings.HasPrefix(line, " "):
			stats.Context++
		}
	}
	return stats
}

// WriteDiffStats writes a one-line stats summary for each drifted result to w.
func WriteDiffStats(w io.Writer, results []DriftResult) error {
	_, err := fmt.Fprintf(w, "%-40s %-8s %-8s %s\n", "RESOURCE", "ADDED", "REMOVED", "TOTAL")
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Diff == "" {
			continue
		}
		stats := ComputeDiffStats(r.Diff)
		key := fmt.Sprintf("%s/%s", r.Kind, r.Name)
		_, err = fmt.Fprintf(w, "%-40s %-8d %-8d %d\n", key, stats.Added, stats.Removed, stats.Total())
		if err != nil {
			return err
		}
	}
	return nil
}
