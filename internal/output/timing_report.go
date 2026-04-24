package output

import (
	"fmt"
	"io"
	"sort"
)

// allStages defines the canonical display order for stage timings.
var allStages = []Stage{
	StageFetchHelm,
	StageFetchLive,
	StageCompare,
	StageReport,
}

// WriteTimingReport writes a human-readable timing summary to w using the
// durations recorded in t. Stages that were never started are omitted.
func WriteTimingReport(w io.Writer, t *StageTimer) error {
	type row struct {
		stage Stage
		ms    float64
	}

	var rows []row
	for _, s := range allStages {
		d := t.Elapsed(s)
		if d == 0 {
			continue
		}
		rows = append(rows, row{stage: s, ms: float64(d.Microseconds()) / 1000.0})
	}

	// Include any stages not in the canonical list (future-proofing).
	known := make(map[Stage]bool)
	for _, s := range allStages {
		known[s] = true
	}
	extra := t.Stages()
	sort.Slice(extra, func(i, j int) bool { return extra[i] < extra[j] })
	for _, s := range extra {
		if known[s] {
			continue
		}
		d := t.Elapsed(s)
		if d == 0 {
			continue
		}
		rows = append(rows, row{stage: s, ms: float64(d.Microseconds()) / 1000.0})
	}

	if len(rows) == 0 {
		_, err := fmt.Fprintln(w, "no timing data available")
		return err
	}

	fmt.Fprintln(w, "Stage timings:")
	for _, r := range rows {
		_, err := fmt.Fprintf(w, "  %-14s %.2f ms\n", r.stage, r.ms)
		if err != nil {
			return err
		}
	}
	return nil
}
