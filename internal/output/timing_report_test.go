package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/drift-check/internal/output"
)

func makeTimerWithStages(t *testing.T) *output.StageTimer {
	t.Helper()
	timer := output.NewStageTimer()
	timer.Start("helm")
	time.Sleep(2 * time.Millisecond)
	timer.End("helm")
	timer.Start("k8s")
	time.Sleep(2 * time.Millisecond)
	timer.End("k8s")
	return timer
}

func TestWriteTimingReport_ContainsStageNames(t *testing.T) {
	var buf bytes.Buffer
	timer := makeTimerWithStages(t)
	if err := output.WriteTimingReport(&buf, timer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "helm") {
		t.Errorf("expected output to contain 'helm', got: %s", out)
	}
	if !strings.Contains(out, "k8s") {
		t.Errorf("expected output to contain 'k8s', got: %s", out)
	}
}

func TestWriteTimingReport_ContainsDurations(t *testing.T) {
	var buf bytes.Buffer
	timer := makeTimerWithStages(t)
	if err := output.WriteTimingReport(&buf, timer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ms") {
		t.Errorf("expected output to contain duration in ms, got: %s", out)
	}
}

func TestWriteTimingReport_EmptyTimer(t *testing.T) {
	var buf bytes.Buffer
	timer := output.NewStageTimer()
	if err := output.WriteTimingReport(&buf, timer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Timing") {
		t.Errorf("expected header even with no stages, got: %s", out)
	}
}

func TestWriteTimingReport_WriterError(t *testing.T) {
	fw := &failWriter{}
	timer := makeTimerWithStages(t)
	err := output.WriteTimingReport(fw, timer)
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}
