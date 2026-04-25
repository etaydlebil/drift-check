package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/your-org/drift-check/internal/compare"
)

// WriteLabelReport writes a report of label drift between Helm-managed and
// live Kubernetes resources. Added labels appear with a '+' prefix, removed
// labels with a '-' prefix.
func WriteLabelReport(w io.Writer, results []compare.DriftResult) error {
	if _, err := fmt.Fprintln(w, "LABEL DRIFT REPORT"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, strings.Repeat("-", 40)); err != nil {
		return err
	}

	any := false
	for _, r := range results {
		if len(r.AddedLabels) == 0 && len(r.RemovedLabels) == 0 {
			continue
		}
		any = true
		header := fmt.Sprintf("%s/%s (%s)", r.Namespace, r.Name, r.Kind)
		if _, err := fmt.Fprintln(w, header); err != nil {
			return err
		}

		keys := sortedKeys(r.AddedLabels)
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "  + %s=%s\n", k, r.AddedLabels[k]); err != nil {
				return err
			}
		}

		keys = sortedKeys(r.RemovedLabels)
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "  - %s=%s\n", k, r.RemovedLabels[k]); err != nil {
				return err
			}
		}
	}

	if !any {
		if _, err := fmt.Fprintln(w, "No label drift detected."); err != nil {
			return err
		}
	}
	return nil
}

// sortedStringKeys returns sorted keys from a map[string]string.
func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
