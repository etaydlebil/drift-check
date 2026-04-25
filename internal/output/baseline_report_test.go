package output

import (
	"strings"
	"testing"
	"time"
)

func baselineEntry(diff string) BaselineEntry {
	return BaselineEntry{
		Release:   "my-app",
		Namespace: "default",
		Diff:      diff,
		Captured:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestWriteBaselineReport_NoChangeSinceBaseline(t *testing.T) {
	var buf strings.Builder
	diff := "+replica: 3\n"
	err := WriteBaselineReport(&buf, "my-app", "default", diff, baselineEntry(diff))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "no change since baseline") {
		t.Errorf("expected 'no change since baseline' in output, got:\n%s", out)
	}
}

func TestWriteBaselineReport_DriftIntroduced(t *testing.T) {
	var buf strings.Builder
	err := WriteBaselineReport(&buf, "my-app", "default", "+replica: 3\n", baselineEntry(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "drift introduced since baseline") {
		t.Errorf("expected 'drift introduced' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "+replica: 3") {
		t.Errorf("expected current diff in output, got:\n%s", out)
	}
}

func TestWriteBaselineReport_DriftResolved(t *testing.T) {
	var buf strings.Builder
	err := WriteBaselineReport(&buf, "my-app", "default", "", baselineEntry("+old: value\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "drift resolved since baseline") {
		t.Errorf("expected 'drift resolved' in output, got:\n%s", out)
	}
}

func TestWriteBaselineReport_DriftChanged(t *testing.T) {
	var buf strings.Builder
	err := WriteBaselineReport(&buf, "my-app", "default", "+new: val\n", baselineEntry("+old: val\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "drift changed since baseline") {
		t.Errorf("expected 'drift changed' in output, got:\n%s", out)
	}
}

func TestWriteBaselineReport_HeaderAlwaysPresent(t *testing.T) {
	var buf strings.Builder
	_ = WriteBaselineReport(&buf, "my-app", "default", "", baselineEntry(""))
	out := buf.String()
	for _, want := range []string{"Baseline Comparison", "Release:", "Namespace:", "Captured:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in header, got:\n%s", want, out)
		}
	}
}

func TestWriteBaselineReport_WriterError(t *testing.T) {
	err := WriteBaselineReport(&errorWriter{}, "app", "ns", "", baselineEntry(""))
	if err == nil {
		t.Error("expected error from failing writer")
	}
}
