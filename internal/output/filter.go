package output

import (
	"strings"
)

// DriftFilter controls which resources are included in output.
type DriftFilter struct {
	OnlyDrifted bool
	Namespaces  []string
	Kinds       []string
}

// NewDriftFilter returns a DriftFilter with default (pass-all) settings.
func NewDriftFilter() *DriftFilter {
	return &DriftFilter{}
}

// Match reports whether a resource identified by namespace, kind, and drift
// status should be included according to the filter rules.
func (f *DriftFilter) Match(namespace, kind string, hasDrift bool) bool {
	if f.OnlyDrifted && !hasDrift {
		return false
	}
	if len(f.Namespaces) > 0 && !containsIgnoreCase(f.Namespaces, namespace) {
		return false
	}
	if len(f.Kinds) > 0 && !containsIgnoreCase(f.Kinds, kind) {
		return false
	}
	return true
}

// WithOnlyDrifted sets the OnlyDrifted flag and returns the filter for chaining.
func (f *DriftFilter) WithOnlyDrifted(v bool) *DriftFilter {
	f.OnlyDrifted = v
	return f
}

// WithNamespaces restricts matches to the given namespaces.
func (f *DriftFilter) WithNamespaces(ns ...string) *DriftFilter {
	f.Namespaces = append(f.Namespaces, ns...)
	return f
}

// WithKinds restricts matches to the given resource kinds.
func (f *DriftFilter) WithKinds(kinds ...string) *DriftFilter {
	f.Kinds = append(f.Kinds, kinds...)
	return f
}

func containsIgnoreCase(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}
