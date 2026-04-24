package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/rest"
)

// ReleaseClient defines the interface for fetching Helm release information.
type ReleaseClient interface {
	GetRelease(name, namespace string) (*release.Release, error)
}

// Client wraps the Helm action configuration.
type Client struct {
	cfg *action.Configuration
}

// NewClient creates a new Helm client configured for the given namespace.
func NewClient(namespace string, restConfig *rest.Config) (*Client, error) {
	cfg := new(action.Configuration)
	getter := newRESTClientGetter(namespace, restConfig)
	if err := cfg.Init(getter, namespace, "secret", func(format string, v ...interface{}) {
		// suppress helm debug output
	}); err != nil {
		return nil, fmt.Errorf("initializing helm config: %w", err)
	}
	return &Client{cfg: cfg}, nil
}

// GetRelease retrieves a Helm release by name from the given namespace.
func (c *Client) GetRelease(name, namespace string) (*release.Release, error) {
	get := action.NewGet(c.cfg)
	rel, err := get.Run(name)
	if err != nil {
		return nil, fmt.Errorf("getting release %q in namespace %q: %w", name, namespace, err)
	}
	return rel, nil
}

// GetManifest returns the rendered manifest string from a Helm release.
func GetManifest(rel *release.Release) string {
	if rel == nil {
		return ""
	}
	return rel.Manifest
}

// GetValues returns the user-supplied values from a Helm release.
func GetValues(rel *release.Release) map[string]interface{} {
	if rel == nil {
		return nil
	}
	return rel.Config
}
