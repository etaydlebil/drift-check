package output

import (
	"fmt"
	"io"
	"time"
)

// TrendDirection indicates whether drift is increasing, decreasing, or stable.
type TrendDirection int

const (
	TrendStable   TrendDirection = 0
	TrendIncreasing TrendDirection = 1
	TrendDecreasing TrendDirection = -1
)

func (t TrendDirection) String() string {
	switch t {
	case TrendIncreasing:
		return "increasing"
	case TrendDecreasing:
		return "decreasing"
	default:
		return "stable"
	}
}

// TrendPoint represents the drift count at a point in time.
type TrendPoint struct {
	Timestamp time.Time
	DriftCount int
}

// DriftTrend summarises the direction and magnitude of drift over a series of snapshots.
type DriftTrend struct {
	Points    []TrendPoint
	Direction TrendDirection
	Delta     int // change from first to last point
}

// ComputeDriftTrend derives a DriftTrend from an ordered slice of Snapshots.
func ComputeDriftTrend(snapshots []Snapshot) DriftTrend {
	if len(snapshots) == 0 {
		return DriftTrend{Direction: TrendStable}
	}

	points := make([]TrendPoint, 0, len(snapshots))
	for _, s := range snapshots {
		points = append(points, TrendPoint{
			Timestamp:  s.CapturedAt,
			DriftCount: s.DriftedCount(),
		})
	}

	first := points[0].DriftCount
	last := points[len(points)-1].DriftCount
	delta := last - first

	var dir TrendDirection
	switch {
	case delta > 0:
		dir = TrendIncreasing
	case delta < 0:
		dir = TrendDecreasing
	default:
		dir = TrendStable
	}

	return DriftTrend{Points: points, Direction: dir, Delta: delta}
}

// WriteDriftTrend writes a human-readable trend summary to w.
func WriteDriftTrend(w io.Writer, trend DriftTrend) error {
	_, err := fmt.Fprintf(w, "Drift trend: %s (delta: %+d) over %d snapshot(s)\n",
		trend.Direction, trend.Delta, len(trend.Points))
	if err != nil {
		return err
	}
	for _, p := range trend.Points {
		_, err = fmt.Fprintf(w, "  %s  drifted=%d\n", p.Timestamp.Format(time.RFC3339), p.DriftCount)
		if err != nil {
			return err
		}
	}
	return nil
}
