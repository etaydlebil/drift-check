package compare

import (
	"fmt"

	"github.com/yourusername/drift-check/internal/diff"
)

// runDrift is the shared implementation used by both Runner and runnerFromInterfaces.
func runDrift(helmClient HelmManifestGetter, k8sClient LiveManifestGetter, release, namespace string) (*Result, error) {
	helmManifest, err := helmClient.GetManifest(release, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetching helm manifest: %w", err)
	}

	liveManifest, err := k8sClient.GetLiveManifest(namespace, release)
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

// Ensure Runner also uses the shared implementation.
func (r *Runner) Run(release, namespace string) (*Result, error) {
	return runDrift(r.helmClient, r.k8sClient, release, namespace)
}
