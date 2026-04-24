package output

import (
	"fmt"
	"io"
)

const (
	separator = "----------------------------------------"
)

// WriteText writes a human-readable report of the Summary to w.
func WriteText(w io.Writer, s Summary) error {
	fmt.Fprintln(w, separator)
	fmt.Fprintf(w, "Release   : %s\n", s.Release)
	fmt.Fprintf(w, "Namespace : %s\n", s.Namespace)
	fmt.Fprintf(w, "Checked   : %s\n", s.CheckedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(w, separator)

	if s.Error != nil {
		fmt.Fprintf(w, "Status    : ERROR\n")
		fmt.Fprintf(w, "Error     : %v\n", s.Error)
		fmt.Fprintln(w, separator)
		return nil
	}

	if s.Drifted {
		fmt.Fprintf(w, "Status    : DRIFT DETECTED\n")
		fmt.Fprintln(w, separator)
		fmt.Fprintln(w, s.Diff)
	} else {
		fmt.Fprintf(w, "Status    : OK (no drift)\n")
	}
	fmt.Fprintln(w, separator)
	return nil
}
