package compare

// DriftResult holds the outcome of comparing a single Kubernetes resource
// between its Helm-rendered manifest and the live cluster state.
type DriftResult struct {
	// Namespace is the Kubernetes namespace of the resource.
	Namespace string
	// Name is the resource name.
	Name string
	// Kind is the Kubernetes kind (e.g. Deployment, Service).
	Kind string
	// Diff contains the unified diff string when drift is detected.
	Diff string
	// Drifted is true when the live resource differs from the Helm manifest.
	Drifted bool
	// Error holds any error encountered while fetching or comparing the resource.
	Error error
	// AddedLabels are labels present on the live resource but absent in the
	// Helm manifest.
	AddedLabels map[string]string
	// RemovedLabels are labels present in the Helm manifest but absent on the
	// live resource.
	RemovedLabels map[string]string
	// AddedAnnotations are annotations present on the live resource but absent
	// in the Helm manifest.
	AddedAnnotations map[string]string
	// RemovedAnnotations are annotations present in the Helm manifest but
	// absent on the live resource.
	RemovedAnnotations map[string]string
}

// HasMetadataDrift returns true if any label or annotation drift was detected.
func (d DriftResult) HasMetadataDrift() bool {
	return len(d.AddedLabels) > 0 ||
		len(d.RemovedLabels) > 0 ||
		len(d.AddedAnnotations) > 0 ||
		len(d.RemovedAnnotations) > 0
}
