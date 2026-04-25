package output

import (
	"fmt"
	"io"
	"strings"
)

// WriteBaselineReport writes a human-readable comparison between the current
// diff and the previously captured baseline diff for a release.
func WriteBaselineReport(w io.Writer, release, namespace, currentDiff string, entry BaselineEntry) error {
	var sb strings.Builder

	sb.WriteString("Baseline Comparison\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n")
	fmt.Fprintf(&sb, "Release:   %s\n", release)
	fmt.Fprintf(&sb, "Namespace: %s\n", namespace)
	fmt.Fprintf(&sb, "Captured:  %s\n", entry.Captured.Format("2006-01-02 15:04:05 UTC"))
	sb.WriteString(strings.Repeat("─", 40) + "\n")

	baselineNorm := strings.TrimSpace(entry.Diff)
	currentNorm := strings.TrimSpace(currentDiff)

	if baselineNorm == currentNorm {
		sb.WriteString("Status: no change since baseline\n")
	} else if baselineNorm == "" && currentNorm != "" {
		sb.WriteString("Status: drift introduced since baseline\n")
		sb.WriteString("\nCurrent diff:\n")
		sb.WriteString(currentDiff)
	} else if baselineNorm != "" && currentNorm == "" {
		sb.WriteString("Status: drift resolved since baseline\n")
	} else {
		sb.WriteString("Status: drift changed since baseline\n")
		sb.WriteString("\nCurrent diff:\n")
		sb.WriteString(currentDiff)
	}

	sb.WriteString("\n")
	_, err := io.WriteString(w, sb.String())
	return err
}
