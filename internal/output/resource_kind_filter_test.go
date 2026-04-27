package output_test

import (
	"testing"

	"github.com/your-org/drift-check/internal/compare"
	"github.com/your-org/drift-check/internal/output"
)

func makeKindResults() []compare.DriftResult {
	return []compare.DriftResult{
		{Kind: "Deployment", Name: "api", Diff: "some diff"},
		{Kind: "Service", Name: "api-svc", Diff: ""},
		{Kind: "ConfigMap", Name: "cfg", Diff: "changed"},
		{Kind: "deployment", Name: "worker", Diff: "drift"},
	}
}

func TestResourceKindFilter_EmptyKindsPassesAll(t *testing.T) {
	f := output.NewResourceKindFilter(nil)
	results := makeKindResults()
	got := f.Apply(results)
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestResourceKindFilter_SingleKind(t *testing.T) {
	f := output.NewResourceKindFilter([]string{"Deployment"})
	got := f.Apply(makeKindResults())
	if len(got) != 2 {
		t.Fatalf("expected 2 deployments, got %d", len(got))
	}
	for _, r := range got {
		if r.Kind != "Deployment" && r.Kind != "deployment" {
			t.Errorf("unexpected kind %q in results", r.Kind)
		}
	}
}

func TestResourceKindFilter_MultipleKinds(t *testing.T) {
	f := output.NewResourceKindFilter([]string{"Service", "ConfigMap"})
	got := f.Apply(makeKindResults())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestResourceKindFilter_CaseInsensitive(t *testing.T) {
	f := output.NewResourceKindFilter([]string{"configmap"})
	got := f.Apply(makeKindResults())
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Name != "cfg" {
		t.Errorf("unexpected result name %q", got[0].Name)
	}
}

func TestResourceKindFilter_NoMatchReturnsEmpty(t *testing.T) {
	f := output.NewResourceKindFilter([]string{"StatefulSet"})
	got := f.Apply(makeKindResults())
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestResourceKindFilter_KindsReturnNilWhenEmpty(t *testing.T) {
	f := output.NewResourceKindFilter([]string{})
	if f.Kinds() != nil {
		t.Error("expected nil kinds for empty filter")
	}
}

func TestResourceKindFilter_KindsReturnConfigured(t *testing.T) {
	f := output.NewResourceKindFilter([]string{"Deployment", "Service"})
	kinds := f.Kinds()
	if len(kinds) != 2 {
		t.Fatalf("expected 2 kinds, got %d", len(kinds))
	}
}
