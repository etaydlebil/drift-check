package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/drift-check/internal/compare"
)

// WriteSeverityReport writes a severity-annotated summary of drift results to w.
// Each resource with drift is classified as low / medium / high based on the
// number of changed lines in its diff.
func WriteSeverityReport(w io.Writer, results []compare.DriftResult) error {
	_, err := fmt.Fprintln(w, "SEVERITY REPORT")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, strings.Repeat("-", 40))
	if err != nil {
		return err
	}

	if len(results) == 0 {
		_, err = fmt.Fprintln(w, "No resources checked.")
		return err
	}

	driftFound := false
	for _, r := range results {
		if r.Diff == "" {
			continue
		}
		driftFound = true
		lines := countDiffLines(r.Diff)
		sev := ClassifyDrift(lines)
		_, err = fmt.Fprintf(w, "%-8s %s/%s (%s)\n",
			"["+strings.ToUpper(sev.String())+"]",
			r.Kind, r.Name, r.Namespace,
		)
		if err != nil {
			return err
		}
	}

	if !driftFound {
		_, err = fmt.Fprintln(w, "No drift detected.")
	}
	return err
}

// countDiffLines counts lines that start with '+' or '-' (excluding the
// unified-diff header lines that start with '---' / '+++').
func countDiffLines(diff string) int {
	count := 0
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
			count++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--") {
			count++
		}
	}
	return count
}
