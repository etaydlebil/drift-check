package output

import (
	"fmt"
	"io"
	"time"
)

// Stage represents a named step in the drift-check pipeline.
type Stage string

const (
	StageFetchHelm  Stage = "fetch-helm"
	StageFetchLive  Stage = "fetch-live"
	StageCompare    Stage = "compare"
	StageReport     Stage = "report"
)

// ProgressEvent carries information about a pipeline stage transition.
type ProgressEvent struct {
	Stage     Stage
	Done      bool
	Err       error
	Timestamp time.Time
}

// ProgressWriter writes human-readable progress lines to an io.Writer.
type ProgressWriter struct {
	w       io.Writer
	verbose bool
}

// NewProgressWriter returns a ProgressWriter that writes to w.
// When verbose is false only error events are emitted.
func NewProgressWriter(w io.Writer, verbose bool) *ProgressWriter {
	return &ProgressWriter{w: w, verbose: verbose}
}

// Write emits a progress line for the given event.
func (p *ProgressWriter) Write(ev ProgressEvent) {
	if ev.Err != nil {
		fmt.Fprintf(p.w, "[%s] ✗ %s: %v\n", ev.Timestamp.Format(time.RFC3339), ev.Stage, ev.Err)
		return
	}
	if !p.verbose {
		return
	}
	status := "…"
	if ev.Done {
		status = "✓"
	}
	fmt.Fprintf(p.w, "[%s] %s %s\n", ev.Timestamp.Format(time.RFC3339), status, ev.Stage)
}

// NewEvent is a convenience constructor that stamps the current time onto a
// ProgressEvent, reducing boilerplate at call sites.
func NewEvent(stage Stage, done bool, err error) ProgressEvent {
	return ProgressEvent{
		Stage:     stage,
		Done:      done,
		Err:       err,
		Timestamp: time.Now(),
	}
}
