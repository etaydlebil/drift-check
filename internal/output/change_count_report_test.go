package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/drift-check/internal/compare"
)

func makeResults(specs []struct {
	kind string
	diff string
}) []compare.DriftResult {
	out := make([]compare.DriftResult, 0, len(specs))
	for _, s := range specs {
		out = append(out, compare.DriftResult{Kind: s.kind, Diff: s.diff})
	}
	return out
}

func TestNewChangeCountReport_Empty(t *testing.T) {
	rep := NewChangeCountReport(nil)
	if rep.Total != 0 || rep.Drifted != 0 || len(rep.ByKind) != 0 {
		t.Errorf("expected zero report, got %+v", rep)
	}
}

func TestNewChangeCountReport_NoDrift(t *testing.T) {
	results := makeResults([]struct{ kind, diff string }{
		{"Deployment", ""},
		{"Service", ""},
	})
	rep := NewChangeCountReport(results)
	if rep.Total != 2 {
		t.Errorf("expected Total=2, got %d", rep.Total)
	}
	if rep.Drifted != 0 {
		t.Errorf("expected Drifted=0, got %d", rep.Drifted)
	}
	if len(rep.ByKind) != 0 {
		t.Errorf("expected empty ByKind, got %v", rep.ByKind)
	}
}

func TestNewChangeCountReport_WithDrift(t *testing.T) {
	results := makeResults([]struct{ kind, diff string }{
		{"Deployment", "- old\n+ new"},
		{"Deployment", "- a\n+ b"},
		{"ConfigMap", "- x\n+ y"},
		{"Service", ""},
	})
	rep := NewChangeCountReport(results)
	if rep.Total != 4 {
		t.Errorf("expected Total=4, got %d", rep.Total)
	}
	if rep.Drifted != 3 {
		t.Errorf("expected Drifted=3, got %d", rep.Drifted)
	}
	if rep.ByKind["Deployment"] != 2 {
		t.Errorf("expected Deployment=2, got %d", rep.ByKind["Deployment"])
	}
	if rep.ByKind["ConfigMap"] != 1 {
		t.Errorf("expected ConfigMap=1, got %d", rep.ByKind["ConfigMap"])
	}
}

func TestWriteChangeCountReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	err := WriteChangeCountReport(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}

func TestWriteChangeCountReport_WithDrift(t *testing.T) {
	results := makeResults([]struct{ kind, diff string }{
		{"Deployment", "- old\n+ new"},
		{"ConfigMap", "- x\n+ y"},
	})
	var buf bytes.Buffer
	err := WriteChangeCountReport(&buf, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Deployment", "ConfigMap", "Drifted", "2"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteChangeCountReport_SortedKinds(t *testing.T) {
	results := makeResults([]struct{ kind, diff string }{
		{"Service", "diff"},
		{"ConfigMap", "diff"},
		{"Deployment", "diff"},
	})
	var buf bytes.Buffer
	_ = WriteChangeCountReport(&buf, results)
	out := buf.String()
	cmIdx := strings.Index(out, "ConfigMap")
	depIdx := strings.Index(out, "Deployment")
	svcIdx := strings.Index(out, "Service")
	if !(cmIdx < depIdx && depIdx < svcIdx) {
		t.Errorf("expected alphabetical kind order in output:\n%s", out)
	}
}
