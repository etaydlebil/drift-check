package output

import (
	"testing"

	"github.com/snyk/drift-check/internal/compare"
)

func makeFieldResults() []compare.DriftResult {
	return []compare.DriftResult{
		{Name: "frontend", Namespace: "default", Kind: "Deployment", Diff: "some diff"},
		{Name: "backend", Namespace: "default", Kind: "Deployment", Diff: ""},
		{Name: "redis", Namespace: "cache", Kind: "StatefulSet", Diff: "other diff"},
		{Name: "nginx", Namespace: "ingress", Kind: "DaemonSet", Diff: ""},
	}
}

func TestFieldSelector_EmptyPassesAll(t *testing.T) {
	fs := NewFieldSelector("", "", "")
	results := makeFieldResults()
	out := fs.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestFieldSelector_FilterByName(t *testing.T) {
	fs := NewFieldSelector("frontend", "", "")
	out := fs.Apply(makeFieldResults())
	if len(out) != 1 || out[0].Name != "frontend" {
		t.Fatalf("expected 1 result named frontend, got %v", out)
	}
}

func TestFieldSelector_FilterByNamespace(t *testing.T) {
	fs := NewFieldSelector("", "default", "")
	out := fs.Apply(makeFieldResults())
	if len(out) != 2 {
		t.Fatalf("expected 2 results in default namespace, got %d", len(out))
	}
}

func TestFieldSelector_FilterByKind(t *testing.T) {
	fs := NewFieldSelector("", "", "Deployment")
	out := fs.Apply(makeFieldResults())
	if len(out) != 2 {
		t.Fatalf("expected 2 Deployment results, got %d", len(out))
	}
}

func TestFieldSelector_CaseInsensitive(t *testing.T) {
	fs := NewFieldSelector("", "", "deployment")
	out := fs.Apply(makeFieldResults())
	if len(out) != 2 {
		t.Fatalf("expected 2 results with case-insensitive kind match, got %d", len(out))
	}
}

func TestFieldSelector_CombinedFilters(t *testing.T) {
	fs := NewFieldSelector("redis", "cache", "StatefulSet")
	out := fs.Apply(makeFieldResults())
	if len(out) != 1 || out[0].Name != "redis" {
		t.Fatalf("expected 1 result for redis/cache/StatefulSet, got %v", out)
	}
}

func TestFieldSelector_NoMatch(t *testing.T) {
	fs := NewFieldSelector("nonexistent", "", "")
	out := fs.Apply(makeFieldResults())
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestFieldSelector_IsEmpty(t *testing.T) {
	if !NewFieldSelector("", "", "").IsEmpty() {
		t.Fatal("expected IsEmpty true for zero-value selector")
	}
	if NewFieldSelector("x", "", "").IsEmpty() {
		t.Fatal("expected IsEmpty false when name is set")
	}
}
