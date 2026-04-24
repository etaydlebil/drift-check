package compare

import (
	"fmt"

	"github.com/yourusername/drift-check/internal/diff"
	"github.com/yourusername/drift-check/internal/helm"
	"github.com/yourusername/drift-check/internal/k8s"
)

// Result holds the outcome of a drift comparison for a release.
type Result struct {
	Release   string
	Namespace string
	HasDrift  bool
	Diff      string
}

// Runner orchestrates fetching Helm manifests and live k8s state,
// then comparing them for drift.
type Runner struct {
	helmClient *helm.Client
	k8sClient  *k8s.Client
}

// NewRunner constructs a Runner with the provided clients.
func NewRunner(helmClient *helm.Client, k8sClient *k8s.Client) *Runner {
	return &Runner{
		helmClient: helmClient,
		k8sClient:  k8sClient,
	}
}

// Run performs drift detection for the given release and namespace.
// It returns a Result describing whether drift was found.
func (r *Runner) Run(release, namespace string) (*Result, error) {
	helmManifest, err := r.helmClient.GetManifest(release, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetching helm manifest: %w", err)
	}

	liveManifest, err := r.k8sClient.GetLiveManifest(namespace, release)
	if err != nil {
		return nil, fmt.Errorf("fetching live manifest: %w", err)
	}

	diffOutput := diff.CompareManifests(helmManifest, liveManifest)

	return &Result{
		Release:   release,
		Namespace: namespace,
		HasDrift:  diffOutput != "",
		Diff:      diffOutput,
	}, nil
}
