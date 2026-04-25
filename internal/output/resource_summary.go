package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// ResourceSummary holds aggregated drift counts per Kubernetes resource kind.
type ResourceSummary struct {
	KindCounts map[string]kindStat
	Total      int
	Drifted    int
}

type kindStat struct {
	Total   int
	Drifted int
}

// NewResourceSummary builds a ResourceSummary from a slice of DriftResult.
func NewResourceSummary(results []DriftResult) ResourceSummary {
	counts := make(map[string]kindStat)
	total := 0
	drifted := 0

	for _, r := range results {
		stat := counts[r.Kind]
		stat.Total++
		if r.HasDrift() {
			stat.Drifted++
			drifted++
		}
		counts[r.Kind] = stat
		total++
	}

	return ResourceSummary{
		KindCounts: counts,
		Total:      total,
		Drifted:    drifted,
	}
}

// WriteResourceSummary writes a tabular breakdown of drift by resource kind.
func WriteResourceSummary(w io.Writer, rs ResourceSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, "KIND\tTOTAL\tDRIFTED"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(tw, "----\t-----\t-------"); err != nil {
		return err
	}

	kinds := make([]string, 0, len(rs.KindCounts))
	for k := range rs.KindCounts {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	for _, kind := range kinds {
		stat := rs.KindCounts[kind]
		if _, err := fmt.Fprintf(tw, "%s\t%d\t%d\n", kind, stat.Total, stat.Drifted); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(tw, "TOTAL\t%d\t%d\n", rs.Total, rs.Drifted); err != nil {
		return err
	}

	return tw.Flush()
}
