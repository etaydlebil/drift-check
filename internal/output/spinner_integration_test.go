package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/drift-check/internal/output"
)

// TestSpinner_IntegrationWithProgressWriter verifies that a SpinnerWriter and
// ProgressWriter can share the same underlying writer without data races.
func TestSpinner_IntegrationWithProgressWriter(t *testing.T) {
	var buf bytes.Buffer

	spinner := output.NewSpinnerWriter(&buf, true)
	progress := output.NewProgressWriter(&buf, output.VerbosityVerbose)

	spinner.Start("fetching helm release")
	time.Sleep(80 * time.Millisecond)

	progress.Write(output.NewEvent("helm", "release fetched", nil))

	spinner.Done("✓ helm release fetched")

	spinner.Start("fetching live manifest")
	time.Sleep(80 * time.Millisecond)

	spinner.Done("✓ live manifest fetched")

	out := buf.String()
	if !strings.Contains(out, "✓ helm release fetched") {
		t.Errorf("expected first done message, got %q", out)
	}
	if !strings.Contains(out, "✓ live manifest fetched") {
		t.Errorf("expected second done message, got %q", out)
	}
}

// TestSpinner_DisabledIntegration ensures disabled spinner does not interfere
// with other writers sharing the same buffer.
func TestSpinner_DisabledIntegration(t *testing.T) {
	var buf bytes.Buffer

	spinner := output.NewSpinnerWriter(&buf, false)
	progress := output.NewProgressWriter(&buf, output.VerbosityVerbose)

	spinner.Start("silent")
	progress.Write(output.NewEvent("k8s", "live manifest retrieved", nil))
	spinner.Done("done")

	out := buf.String()
	if strings.Contains(out, "silent") || strings.Contains(out, "done") {
		t.Errorf("disabled spinner should produce no output, got %q", out)
	}
	if !strings.Contains(out, "live manifest retrieved") {
		t.Errorf("expected progress event in output, got %q", out)
	}
}
