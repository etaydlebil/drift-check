package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yourusername/drift-check/internal/compare"
)

// WatchEvent represents a single drift-check result captured during a watch cycle.
type WatchEvent struct {
	// Timestamp is when the event was recorded.
	Timestamp time.Time
	// Results holds the drift results for this cycle.
	Results []compare.DriftResult
	// CycleIndex is the 1-based index of the watch cycle.
	CycleIndex int
}

// WatchReport summarises changes observed across multiple watch cycles.
type WatchReport struct {
	// Events holds all watch events in chronological order.
	Events []WatchEvent
	// Release is the Helm release being watched.
	Release string
	// Namespace is the Kubernetes namespace of the release.
	Namespace string
}

// NewWatchReport creates a WatchReport for the given release and namespace.
func NewWatchReport(release, namespace string) *WatchReport {
	return &WatchReport{
		Release:   release,
		Namespace: namespace,
	}
}

// AddEvent appends a new watch event to the report.
func (w *WatchReport) AddEvent(results []compare.DriftResult) {
	w.Events = append(w.Events, WatchEvent{
		Timestamp:  time.Now(),
		Results:    results,
		CycleIndex: len(w.Events) + 1,
	})
}

// TotalCycles returns the number of watch cycles recorded.
func (w *WatchReport) TotalCycles() int {
	return len(w.Events)
}

// DriftedCycles returns the number of cycles in which drift was detected.
func (w *WatchReport) DriftedCycles() int {
	count := 0
	for _, e := range w.Events {
		for _, r := range e.Results {
			if r.Diff != "" {
				count++
				break
			}
		}
	}
	return count
}

// WriteWatchReport writes a human-readable summary of all watch cycles to w.
// Each cycle is separated by a divider and shows the timestamp, cycle index,
// and a brief drift summary.
func WriteWatchReport(w io.Writer, report *WatchReport) error {
	if report == nil {
		return nil
	}

	divider := strings.Repeat("-", 60)

	_, err := fmt.Fprintf(w, "Watch Report — release: %s  namespace: %s\n",
		report.Release, report.Namespace)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Cycles: %d  Drifted: %d\n\n",
		report.TotalCycles(), report.DriftedCycles())
	if err != nil {
		return err
	}

	for _, event := range report.Events {
		if _, err = fmt.Fprintln(w, divider); err != nil {
			return err
		}

		ts := event.Timestamp.Format(time.RFC3339)
		if _, err = fmt.Fprintf(w, "Cycle #%d  %s\n", event.CycleIndex, ts); err != nil {
			return err
		}

		drifted := 0
		for _, r := range event.Results {
			if r.Diff != "" {
				drifted++
			}
		}

		status := "OK"
		if drifted > 0 {
			status = fmt.Sprintf("DRIFT (%d resource(s))", drifted)
		}

		if _, err = fmt.Fprintf(w, "Status: %s\n", status); err != nil {
			return err
		}

		for _, r := range event.Results {
			if r.Diff == "" {
				continue
			}
			if _, err = fmt.Fprintf(w, "  • %s/%s\n", r.Kind, r.Name); err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintln(w, divider)
	return err
}
