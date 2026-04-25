package output

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// ResourcePatchSummary pairs a resource identifier with its PatchSummary.
type ResourcePatchSummary struct {
	Namespace string
	Kind      string
	Name      string
	Summary   PatchSummary
}

// WritePatchSummaryReport writes a tabular report of per-resource patch
// summaries to w. Resources with no changes are omitted unless all are clean.
func WritePatchSummaryReport(w io.Writer, resources []ResourcePatchSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	_, err := fmt.Fprintln(tw, "NAMESPACE\tKIND\tNAME\tCHANGES")
	if err != nil {
		return err
	}

	drifted := make([]ResourcePatchSummary, 0, len(resources))
	for _, r := range resources {
		if r.Summary.Added > 0 || r.Summary.Removed > 0 || r.Summary.Changed > 0 {
			drifted = append(drifted, r)
		}
	}

	if len(drifted) == 0 {
		_, err = fmt.Fprintln(tw, "(none)\t\t\t")
		if err != nil {
			return err
		}
		return tw.Flush()
	}

	for _, r := range drifted {
		_, err = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			r.Namespace, r.Kind, r.Name, r.Summary.String())
		if err != nil {
			return err
		}
	}

	return tw.Flush()
}
