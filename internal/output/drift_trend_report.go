package output

import (
	"fmt"
	"io"
	"strings"
)

const trendBarWidth = 20

// WriteDriftTrendReport writes a visual bar-chart style trend report to w.
// Each snapshot is represented as a row with a proportional bar.
func WriteDriftTrendReport(w io.Writer, trend DriftTrend) error {
	if _, err := fmt.Fprintln(w, "=== Drift Trend Report ==="); err != nil {
		return err
	}

	if len(trend.Points) == 0 {
		_, err := fmt.Fprintln(w, "No snapshot data available.")
		return err
	}

	// Find max for scaling.
	maxCount := 1
	for _, p := range trend.Points {
		if p.DriftCount > maxCount {
			maxCount = p.DriftCount
		}
	}

	for _, p := range trend.Points {
		barLen := (p.DriftCount * trendBarWidth) / maxCount
		bar := strings.Repeat("█", barLen)
		line := fmt.Sprintf("  %s  [%-*s] %d\n",
			p.Timestamp.Format("2006-01-02 15:04"),
			trendBarWidth, bar,
			p.DriftCount,
		)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "\nOverall: %s (delta: %+d)\n", trend.Direction, trend.Delta)
	return err
}
