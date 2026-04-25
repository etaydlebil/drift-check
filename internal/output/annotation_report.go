package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// AnnotationDiff holds the annotation differences for a single resource.
type AnnotationDiff struct {
	Kind      string
	Name      string
	Namespace string
	Added     map[string]string
	Removed   map[string]string
	Changed   map[string][2]string // key -> [live, chart]
}

// WriteAnnotationReport writes a human-readable annotation drift report to w.
// It lists each resource that has annotation differences, grouped by kind.
func WriteAnnotationReport(w io.Writer, diffs []AnnotationDiff) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "No annotation drift detected.")
		return err
	}

	// Sort for deterministic output.
	sorted := make([]AnnotationDiff, len(diffs))
	copy(sorted, diffs)
	sort.Slice(sorted, func(i, j int) bool {
		ki := sorted[i].Kind + "/" + sorted[i].Namespace + "/" + sorted[i].Name
		kj := sorted[j].Kind + "/" + sorted[j].Namespace + "/" + sorted[j].Name
		return ki < kj
	})

	for _, d := range sorted {
		resource := fmt.Sprintf("%s %s/%s", d.Kind, d.Namespace, d.Name)
		if d.Namespace == "" {
			resource = fmt.Sprintf("%s %s", d.Kind, d.Name)
		}
		if _, err := fmt.Fprintf(w, "--- %s\n", resource); err != nil {
			return err
		}

		keys := sortedKeys(d.Added)
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "  + %s: %s\n", k, d.Added[k]); err != nil {
				return err
			}
		}

		keys = sortedKeys(d.Removed)
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "  - %s: %s\n", k, d.Removed[k]); err != nil {
				return err
			}
		}

		keys = sortedKeys(d.Changed)
		for _, k := range keys {
			pair := d.Changed[k]
			if _, err := fmt.Fprintf(w, "  ~ %s: %s -> %s\n", k, strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
