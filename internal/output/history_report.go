package output

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteHistoryReport writes a human-readable table of drift history entries
// to w. When there are no entries a short message is printed instead.
func WriteHistoryReport(w io.Writer, h *History) error {
	if h == nil || len(h.Entries) == 0 {
		_, err := fmt.Fprintln(w, "no history recorded")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, err := fmt.Fprintln(tw, "TIMESTAMP\tRELEASE\tNAMESPACE\tDRIFTED\tCHANGES\tERROR")
	if err != nil {
		return err
	}

	for _, e := range h.Entries {
		drifted := "no"
		if e.Drifted {
			drifted = "yes"
		}
		errStr := "-"
		if e.ErrorMsg != "" {
			errStr = e.ErrorMsg
		}
		_, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Release,
			e.Namespace,
			drifted,
			e.DriftCount,
			errStr,
		)
		if err != nil {
			return err
		}
	}
	return tw.Flush()
}
