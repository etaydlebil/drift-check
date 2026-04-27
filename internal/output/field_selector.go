package output

import (
	"strings"

	"github.com/snyk/drift-check/internal/compare"
)

// FieldSelector filters DriftResult slices to only those whose resource
// metadata fields (name, namespace, kind) match the provided selector.
type FieldSelector struct {
	name      string
	namespace string
	kind      string
}

// NewFieldSelector constructs a FieldSelector. Empty strings match any value.
func NewFieldSelector(name, namespace, kind string) *FieldSelector {
	return &FieldSelector{
		name:      strings.ToLower(name),
		namespace: strings.ToLower(namespace),
		kind:      strings.ToLower(kind),
	}
}

// Apply returns only those results that satisfy all non-empty selector fields.
func (f *FieldSelector) Apply(results []compare.DriftResult) []compare.DriftResult {
	if f.name == "" && f.namespace == "" && f.kind == "" {
		out := make([]compare.DriftResult, len(results))
		copy(out, results)
		return out
	}

	var out []compare.DriftResult
	for _, r := range results {
		if f.name != "" && strings.ToLower(r.Name) != f.name {
			continue
		}
		if f.namespace != "" && strings.ToLower(r.Namespace) != f.namespace {
			continue
		}
		if f.kind != "" && strings.ToLower(r.Kind) != f.kind {
			continue
		}
		out = append(out, r)
	}
	return out
}

// IsEmpty reports whether the selector has no constraints.
func (f *FieldSelector) IsEmpty() bool {
	return f.name == "" && f.namespace == "" && f.kind == ""
}
