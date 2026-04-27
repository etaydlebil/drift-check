package output

import (
	"strings"

	"github.com/your-org/drift-check/internal/compare"
)

// ResourceKindFilter restricts drift results to a specific set of Kubernetes resource kinds.
type ResourceKindFilter struct {
	allowedKinds map[string]struct{}
}

// NewResourceKindFilter creates a filter that passes only results whose Kind is in the
// provided list. An empty list means all kinds are accepted (no filtering).
func NewResourceKindFilter(kinds []string) *ResourceKindFilter {
	allowed := make(map[string]struct{}, len(kinds))
	for _, k := range kinds {
		allowed[strings.ToLower(strings.TrimSpace(k))] = struct{}{}
	}
	return &ResourceKindFilter{allowedKinds: allowed}
}

// Apply returns the subset of results whose Kind matches the filter.
// If the filter was created with no kinds, all results are returned unchanged.
func (f *ResourceKindFilter) Apply(results []compare.DriftResult) []compare.DriftResult {
	if len(f.allowedKinds) == 0 {
		return results
	}
	out := make([]compare.DriftResult, 0, len(results))
	for _, r := range results {
		if _, ok := f.allowedKinds[strings.ToLower(r.Kind)]; ok {
			out = append(out, r)
		}
	}
	return out
}

// Kinds returns the sorted list of configured kind filters (lower-cased).
func (f *ResourceKindFilter) Kinds() []string {
	if len(f.allowedKinds) == 0 {
		return nil
	}
	list := make([]string, 0, len(f.allowedKinds))
	for k := range f.allowedKinds {
		list = append(list, k)
	}
	return sortedStringKeys(f.allowedKinds)
}
