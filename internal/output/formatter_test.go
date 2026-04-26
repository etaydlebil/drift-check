package output_test

import (
	"strings"
	"testing"

	"github.com/your-org/drift-check/internal/output"
)

func TestFormatter_TextNoDrift(t *testing.T) {
	var buf strings.Builder
	f := output.NewFormatter(&buf, output.FormatText)
	err := f.Write(output.Result{
		Release:   "my-app",
		Namespace: "default",
		HasDrift:  false,
		Diff:      "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[OK]") {
		t.Errorf("expected [OK] in output, got: %q", got)
	}
	if !strings.Contains(got, "no drift detected") {
		t.Errorf("expected 'no drift detected' in output, got: %q", got)
	}
}

func TestFormatter_TextWithDrift(t *testing.T) {
	var buf strings.Builder
	f := output.NewFormatter(&buf, output.FormatText)
	err := f.Write(output.Result{
		Release:   "my-app",
		Namespace: "staging",
		HasDrift:  true,
		Diff:      "+replicas: 3\n-replicas: 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[DRIFT]") {
		t.Errorf("expected [DRIFT] in output, got: %q", got)
	}
	if !strings.Contains(got, "+replicas: 3") {
		t.Errorf("expected diff content in output, got: %q", got)
	}
}

func TestFormatter_JSONNoDrift(t *testing.T) {
	var buf strings.Builder
	f := output.NewFormatter(&buf, output.FormatJSON)
	err := f.Write(output.Result{
		Release:   "api",
		Namespace: "prod",
		HasDrift:  false,
		Diff:      "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, `"has_drift":false`) {
		t.Errorf("expected has_drift:false in JSON output, got: %q", got)
	}
	if !strings.Contains(got, `"release":"api"`) {
		t.Errorf("expected release field in JSON output, got: %q", got)
	}
}

func TestFormatter_JSONWithDrift(t *testing.T) {
	var buf strings.Builder
	f := output.NewFormatter(&buf, output.FormatJSON)
	err := f.Write(output.Result{
		Release:   "worker",
		Namespace: "prod",
		HasDrift:  true,
		Diff:      "+image: v2\n-image: v1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, `"has_drift":true`) {
		t.Errorf("expected has_drift:true in JSON output, got: %q", got)
	}
}

func TestFormatter_TextIncludesReleaseAndNamespace(t *testing.T) {
	var buf strings.Builder
	f := output.NewFormatter(&buf, output.FormatText)
	err := f.Write(output.Result{
		Release:   "frontend",
		Namespace: "qa",
		HasDrift:  false,
		Diff:      "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "frontend") {
		t.Errorf("expected release name in output, got: %q", got)
	}
	if !strings.Contains(got, "qa") {
		t.Errorf("expected namespace in output, got: %q", got)
	}
}
