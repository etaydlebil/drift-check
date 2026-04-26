package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/example/drift-check/internal/compare"
)

// ChangeCountReport holds per-kind drift counts derived from a slice of DriftResults.
type ChangeCountReport struct {
	Total   int
	Drifted int
	ByKind  map[string]int
}

// NewChangeCountReport builds a ChangeCountReport from the provided results.
func NewChangeCountReport(results []compare.DriftResult) ChangeCountReport {
	r := ChangeCountReport{
		ByKind: make(map[string]int),
	}
	for _, res := range results {
		r.Total++
		if res.Diff != "" {
			r.Drifted++
			r.ByKind[res.Kind]++
		}
	}
	return r
}

// WriteChangeCountReport writes a human-readable change-count summary to w.
func WriteChangeCountReport(w io.Writer, results []compare.DriftResult) error {
	rep := NewChangeCountReport(results)

	_, err := fmt.Fprintf(w, "Change count report\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  Total resources : %d\n", rep.Total)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  Drifted         : %d\n", rep.Drifted)
	if err != nil {
		return err
	}

	if len(rep.ByKind) == 0 {
		_, err = fmt.Fprintf(w, "  No drift detected.\n")
		return err
	}

	_, err = fmt.Fprintf(w, "  By kind:\n")
	if err != nil {
		return err
	}

	kinds := make([]string, 0, len(rep.ByKind))
	for k := range rep.ByKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	for _, k := range kinds {
		padded := k + strings.Repeat(" ", max(0, 20-len(k)))
		_, err = fmt.Fprintf(w, "    %s %d\n", padded, rep.ByKind[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
