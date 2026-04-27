package output

import (
	"strings"

	"github.com/user/drift-check/internal/compare"
)

// NamespaceFilter restricts DriftResults to a set of allowed namespaces.
type NamespaceFilter struct {
	allowed map[string]struct{}
}

// NewNamespaceFilter creates a filter that passes results whose namespace is in
// the provided list. An empty list passes all results.
func NewNamespaceFilter(namespaces []string) *NamespaceFilter {
	allowed := make(map[string]struct{}, len(namespaces))
	for _, ns := range namespaces {
		allowed[strings.ToLower(ns)] = struct{}{}
	}
	return &NamespaceFilter{allowed: allowed}
}

// Apply returns only the results whose namespace matches the filter.
// If no namespaces were configured, all results are returned unchanged.
func (f *NamespaceFilter) Apply(results []compare.DriftResult) []compare.DriftResult {
	if len(f.allowed) == 0 {
		return results
	}
	out := make([]compare.DriftResult, 0, len(results))
	for _, r := range results {
		if _, ok := f.allowed[strings.ToLower(r.Namespace)]; ok {
			out = append(out, r)
		}
	}
	return out
}

// Passes reports whether a single DriftResult passes the filter.
func (f *NamespaceFilter) Passes(r compare.DriftResult) bool {
	if len(f.allowed) == 0 {
		return true
	}
	_, ok := f.allowed[strings.ToLower(r.Namespace)]
	return ok
}
