package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/drift-check/internal/compare"
)

// NamespaceSummary holds drift counts grouped by namespace.
type NamespaceSummary struct {
	Counts map[string]int // namespace -> drifted resource count
	Total  int
}

// NewNamespaceSummary builds a NamespaceSummary from a slice of DriftResults.
func NewNamespaceSummary(results []compare.DriftResult) NamespaceSummary {
	counts := make(map[string]int)
	total := 0
	for _, r := range results {
		if r.Diff != "" {
			counts[r.Namespace]++
			total++
		}
	}
	return NamespaceSummary{Counts: counts, Total: total}
}

// WriteNamespaceSummary writes a tabular namespace-level drift summary to w.
func WriteNamespaceSummary(w io.Writer, results []compare.DriftResult) error {
	summary := NewNamespaceSummary(results)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAMESPACE\tDRIFTED RESOURCES")

	if summary.Total == 0 {
		fmt.Fprintln(tw, "(none)\t0")
		return tw.Flush()
	}

	namespaces := make([]string, 0, len(summary.Counts))
	for ns := range summary.Counts {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	for _, ns := range namespaces {
		fmt.Fprintf(tw, "%s\t%d\n", ns, summary.Counts[ns])
	}
	fmt.Fprintf(tw, "TOTAL\t%d\n", summary.Total)
	return tw.Flush()
}
