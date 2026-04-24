package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// WriteTable writes a human-readable table report of the drift summary to w.
func WriteTable(w io.Writer, s Summary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "FIELD\tHELM\tLIVE")
	fmt.Fprintln(tw, strings.Repeat("-", 6)+"\t"+strings.Repeat("-", 4)+"\t"+strings.Repeat("-", 4))

	if s.Drift == "" {
		fmt.Fprintln(tw, "(no drift detected)\t\t")
		return tw.Flush()
	}

	for _, line := range strings.Split(s.Drift, "\n") {
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "-"):
			fmt.Fprintf(tw, "%s\t%s\t%s\n", strings.TrimPrefix(line, "-"), strings.TrimPrefix(line, "-"), "")
		case strings.HasPrefix(line, "+"):
			fmt.Fprintf(tw, "%s\t%s\t%s\n", strings.TrimPrefix(line, "+"), "", strings.TrimPrefix(line, "+"))
		default:
			fmt.Fprintf(tw, "%s\t\t\n", line)
		}
	}

	if s.Error != "" {
		fmt.Fprintf(tw, "\nError:\t%s\t\n", s.Error)
	}

	return tw.Flush()
}
