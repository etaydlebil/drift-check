package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestProgressWriter_VerboseStart(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, true)
	pw.Write(ProgressEvent{Stage: StageFetchHelm, Done: false, Timestamp: fixedTime})
	if !strings.Contains(buf.String(), "…") {
		t.Errorf("expected in-progress marker, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), string(StageFetchHelm)) {
		t.Errorf("expected stage name in output, got: %s", buf.String())
	}
}

func TestProgressWriter_VerboseDone(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, true)
	pw.Write(ProgressEvent{Stage: StageCompare, Done: true, Timestamp: fixedTime})
	if !strings.Contains(buf.String(), "✓") {
		t.Errorf("expected done marker, got: %s", buf.String())
	}
}

func TestProgressWriter_NonVerboseSkipsNonError(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, false)
	pw.Write(ProgressEvent{Stage: StageFetchLive, Done: true, Timestamp: fixedTime})
	if buf.Len() != 0 {
		t.Errorf("expected no output in non-verbose mode, got: %s", buf.String())
	}
}

func TestProgressWriter_ErrorAlwaysWritten(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, false)
	pw.Write(ProgressEvent{Stage: StageReport, Err: errors.New("boom"), Timestamp: fixedTime})
	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("expected error message in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "✗") {
		t.Errorf("expected error marker, got: %s", buf.String())
	}
}

func TestProgressWriter_TimestampPresent(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, true)
	pw.Write(ProgressEvent{Stage: StageFetchHelm, Done: true, Timestamp: fixedTime})
	if !strings.Contains(buf.String(), "2024-01-15") {
		t.Errorf("expected timestamp in output, got: %s", buf.String())
	}
}
