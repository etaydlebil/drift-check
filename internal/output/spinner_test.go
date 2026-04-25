package output

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSpinner_DisabledProducesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerWriter(&buf, false)
	s.Start("loading")
	time.Sleep(200 * time.Millisecond)
	s.Stop()
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got %q", buf.String())
	}
}

func TestSpinner_DisabledDoneProducesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerWriter(&buf, false)
	s.Done("finished")
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got %q", buf.String())
	}
}

func TestSpinner_EnabledWritesFrames(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerWriter(&buf, true)
	s.Start("checking")
	time.Sleep(300 * time.Millisecond)
	s.Stop()
	out := buf.String()
	if !strings.Contains(out, "checking") {
		t.Errorf("expected label in output, got %q", out)
	}
}

func TestSpinner_DoneWritesMessage(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerWriter(&buf, true)
	s.Start("working")
	time.Sleep(100 * time.Millisecond)
	s.Done("✓ complete")
	out := buf.String()
	if !strings.Contains(out, "✓ complete") {
		t.Errorf("expected done message in output, got %q", out)
	}
}

func TestSpinner_StopIsIdempotent(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerWriter(&buf, true)
	s.Start("task")
	time.Sleep(50 * time.Millisecond)
	s.Stop()
	// second Stop should not panic
	s.Stop()
}

func TestSpinner_RestartClearsOldTicker(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinnerWriter(&buf, true)
	s.Start("first")
	time.Sleep(50 * time.Millisecond)
	s.Start("second") // should stop first cleanly
	time.Sleep(150 * time.Millisecond)
	s.Stop()
	out := buf.String()
	if !strings.Contains(out, "second") {
		t.Errorf("expected second label in output, got %q", out)
	}
}
