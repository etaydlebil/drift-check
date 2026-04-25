package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/drift-check/internal/compare"
	"github.com/your-org/drift-check/internal/output"
)

func makeLabelResult(ns, name, kind string, added, removed map[string]string) compare.DriftResult {
	return compare.DriftResult{
		Namespace:     ns,
		Name:          name,
		Kind:          kind,
		AddedLabels:   added,
		RemovedLabels: removed,
	}
}

func TestWriteLabelReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteLabelReport(&buf, []compare.DriftResult{
		makeLabelResult("default", "my-deploy", "Deployment", nil, nil),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No label drift detected.") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteLabelReport_WithAddedLabels(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteLabelReport(&buf, []compare.DriftResult{
		makeLabelResult("default", "my-deploy", "Deployment",
			map[string]string{"team": "platform"},
			nil,
		),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ team=platform") {
		t.Errorf("expected added label in output, got: %s", out)
	}
}

func TestWriteLabelReport_WithRemovedLabels(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteLabelReport(&buf, []compare.DriftResult{
		makeLabelResult("kube-system", "coredns", "Deployment",
			nil,
			map[string]string{"app": "coredns"},
		),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- app=coredns") {
		t.Errorf("expected removed label in output, got: %s", out)
	}
}

func TestWriteLabelReport_HeaderAlwaysPresent(t *testing.T) {
	var buf bytes.Buffer
	_ = output.WriteLabelReport(&buf, nil)
	if !strings.Contains(buf.String(), "LABEL DRIFT REPORT") {
		t.Errorf("expected header in output, got: %s", buf.String())
	}
}

func TestWriteLabelReport_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	err := output.WriteLabelReport(&buf, []compare.DriftResult{
		makeLabelResult("default", "svc", "Service",
			map[string]string{"z-label": "last", "a-label": "first"},
			nil,
		),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	aIdx := strings.Index(out, "a-label")
	zIdx := strings.Index(out, "z-label")
	if aIdx > zIdx {
		t.Errorf("expected sorted labels: a-label before z-label")
	}
}

func TestWriteLabelReport_WriterError(t *testing.T) {
	err := output.WriteLabelReport(&errWriter{}, []compare.DriftResult{})
	if err == nil {
		t.Error("expected error from failing writer")
	}
}
