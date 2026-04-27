package output_test

import (
	"testing"

	"github.com/user/drift-check/internal/compare"
	"github.com/user/drift-check/internal/output"
)

func makeNSFilterResults() []compare.DriftResult {
	return []compare.DriftResult{
		{Namespace: "default", Diff: "diff1"},
		{Namespace: "staging", Diff: ""},
		{Namespace: "prod", Diff: "diff2"},
	}
}

func TestNamespaceFilter_EmptyPassesAll(t *testing.T) {
	f := output.NewNamespaceFilter(nil)
	results := makeNSFilterResults()
	out := f.Apply(results)
	if len(out) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(out))
	}
}

func TestNamespaceFilter_SingleNamespace(t *testing.T) {
	f := output.NewNamespaceFilter([]string{"default"})
	out := f.Apply(makeNSFilterResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Namespace != "default" {
		t.Errorf("expected default, got %s", out[0].Namespace)
	}
}

func TestNamespaceFilter_MultipleNamespaces(t *testing.T) {
	f := output.NewNamespaceFilter([]string{"default", "prod"})
	out := f.Apply(makeNSFilterResults())
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestNamespaceFilter_CaseInsensitive(t *testing.T) {
	f := output.NewNamespaceFilter([]string{"DEFAULT"})
	out := f.Apply(makeNSFilterResults())
	if len(out) != 1 {
		t.Errorf("expected 1 result, got %d", len(out))
	}
}

func TestNamespaceFilter_NoMatch(t *testing.T) {
	f := output.NewNamespaceFilter([]string{"unknown"})
	out := f.Apply(makeNSFilterResults())
	if len(out) != 0 {
		t.Errorf("expected 0 results, got %d", len(out))
	}
}

func TestNamespaceFilter_Passes(t *testing.T) {
	f := output.NewNamespaceFilter([]string{"staging"})
	r := compare.DriftResult{Namespace: "staging"}
	if !f.Passes(r) {
		t.Error("expected Passes to return true for staging")
	}
	r2 := compare.DriftResult{Namespace: "prod"}
	if f.Passes(r2) {
		t.Error("expected Passes to return false for prod")
	}
}

func TestNamespaceFilter_PassesEmptyFilterAlwaysTrue(t *testing.T) {
	f := output.NewNamespaceFilter(nil)
	if !f.Passes(compare.DriftResult{Namespace: "anything"}) {
		t.Error("empty filter should always pass")
	}
}
