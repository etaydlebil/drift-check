package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSnapshot(driftCount int, at time.Time) Snapshot {
	results := make([]interface{}, driftCount)
	_ = results
	// Build a minimal Snapshot with the right DriftedCount via Results.
	s := Snapshot{CapturedAt: at}
	for i := 0; i < driftCount; i++ {
		s.Results = append(s.Results, makeDriftResult(true))
	}
	return s
}

func TestComputeDriftTrend_Empty(t *testing.T) {
	trend := ComputeDriftTrend(nil)
	assert.Equal(t, TrendStable, trend.Direction)
	assert.Equal(t, 0, trend.Delta)
	assert.Empty(t, trend.Points)
}

func TestComputeDriftTrend_Stable(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSnapshot(2, now.Add(-2*time.Hour)),
		makeSnapshot(2, now.Add(-time.Hour)),
		makeSnapshot(2, now),
	}
	trend := ComputeDriftTrend(snaps)
	assert.Equal(t, TrendStable, trend.Direction)
	assert.Equal(t, 0, trend.Delta)
	assert.Len(t, trend.Points, 3)
}

func TestComputeDriftTrend_Increasing(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSnapshot(1, now.Add(-time.Hour)),
		makeSnapshot(4, now),
	}
	trend := ComputeDriftTrend(snaps)
	assert.Equal(t, TrendIncreasing, trend.Direction)
	assert.Equal(t, 3, trend.Delta)
}

func TestComputeDriftTrend_Decreasing(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSnapshot(5, now.Add(-time.Hour)),
		makeSnapshot(2, now),
	}
	trend := ComputeDriftTrend(snaps)
	assert.Equal(t, TrendDecreasing, trend.Direction)
	assert.Equal(t, -3, trend.Delta)
}

func TestTrendDirection_String(t *testing.T) {
	assert.Equal(t, "stable", TrendStable.String())
	assert.Equal(t, "increasing", TrendIncreasing.String())
	assert.Equal(t, "decreasing", TrendDecreasing.String())
}

func TestWriteDriftTrend_ContainsDirection(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSnapshot(1, now.Add(-time.Hour)),
		makeSnapshot(3, now),
	}
	trend := ComputeDriftTrend(snaps)
	var buf bytes.Buffer
	err := WriteDriftTrend(&buf, trend)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "increasing")
	assert.Contains(t, out, "+2")
	assert.True(t, strings.Count(out, "drifted=") == 2)
}

func TestWriteDriftTrend_WriterError(t *testing.T) {
	trend := ComputeDriftTrend([]Snapshot{makeSnapshot(1, time.Now())})
	err := WriteDriftTrend(&errorWriter{}, trend)
	assert.Error(t, err)
}
