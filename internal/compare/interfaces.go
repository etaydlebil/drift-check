package compare

// HelmManifestGetter abstracts fetching the desired state from a Helm release.
type HelmManifestGetter interface {
	GetManifest(release, namespace string) (string, error)
}

// LiveManifestGetter abstracts fetching the current live state from Kubernetes.
type LiveManifestGetter interface {
	GetLiveManifest(namespace, name string) (string, error)
}

// runnerFromInterfaces is an internal variant of Runner that accepts interfaces,
// enabling easy unit testing without real Helm / k8s connections.
type runnerFromInterfaces struct {
	helmClient HelmManifestGetter
	k8sClient  LiveManifestGetter
}

// NewRunnerFromInterfaces constructs a testable Runner from interface values.
func NewRunnerFromInterfaces(h HelmManifestGetter, k LiveManifestGetter) *runnerFromInterfaces {
	return &runnerFromInterfaces{helmClient: h, k8sClient: k}
}

// Run performs drift detection using the injected interface implementations.
func (r *runnerFromInterfaces) Run(release, namespace string) (*Result, error) {
	return runDrift(r.helmClient, r.k8sClient, release, namespace)
}
