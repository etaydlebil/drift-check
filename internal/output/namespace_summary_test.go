package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/drift-check/internal/compare"
	"github.com/user/drift-check/internal/output"
)

func makeNSResults(pairs ...string) []compare.DriftResult {
	// pairs: namespace, diff (empty means no drift)
	var results []compare.DriftResult
	for i := 0; i+1 < len(pairs); i += 2 {
		results = append(results, compare.DriftResult{
			Namespace: pairs[i],
			Diff:      pairs[i+1],
		})
	}
	return results
}

func TestNewNamespaceSummary_Empty(t *testing.T) {
	s := output.NewNamespaceSummary(nil)
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
}

func TestNewNamespaceSummary_NoDrift(t *testing.T) {
	results := makeNSResults("default", "", "kube-system", "")
	s := output.NewNamespaceSummary(results)
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
}

func TestNewNamespaceSummary_WithDrift(t *testing.T) {
	results := makeNSResults(
		"default", "- old\n+ new",
		"default", "- a\n+ b",
		"staging", "- x\n+ y",
		"prod", "",
	)
	s := output.NewNamespaceSummary(results)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Counts["default"] != 2 {
		t.Errorf("expected default=2, got %d", s.Counts["default"])
	}
	if s.Counts["staging"] != 1 {
		t.Errorf("expected staging=1, got %d", s.Counts["staging"])
	}
}

func TestWriteNamespaceSummary_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	if err := output.WriteNamespaceSummary(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "NAMESPACE") {
		t.Error("expected header NAMESPACE")
	}
	if !strings.Contains(out, "(none)") {
		t.Error("expected (none) row")
	}
}

func TestWriteNamespaceSummary_WithDrift(t *testing.T) {
	results := makeNSResults(
		"alpha", "- old\n+ new",
		"beta", "- x\n+ y",
	)
	var buf bytes.Buffer
	if err := output.WriteNamespaceSummary(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Error("expected alpha in output")
	}
	if !strings.Contains(out, "beta") {
		t.Error("expected beta in output")
	}
	if !strings.Contains(out, "TOTAL") {
		t.Error("expected TOTAL row")
	}
}

func TestWriteNamespaceSummary_SortedOutput(t *testing.T) {
	results := makeNSResults(
		"zeta", "diff",
		"alpha", "diff",
		"mango", "diff",
	)
	var buf bytes.Buffer
	_ = output.WriteNamespaceSummary(&buf, results)
	out := buf.String()
	alphaIdx := strings.Index(out, "alpha")
	mangoIdx := strings.Index(out, "mango")
	zetaIdx := strings.Index(out, "zeta")
	if !(alphaIdx < mangoIdx && mangoIdx < zetaIdx) {
		t.Errorf("expected sorted output: alpha < mango < zeta, got indices %d %d %d", alphaIdx, mangoIdx, zetaIdx)
	}
}
