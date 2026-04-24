package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps a dynamic Kubernetes client for fetching live resource state.
type Client struct {
	dynamic dynamic.Interface
}

// NewClient creates a new k8s Client using the provided kubeconfig path.
// If kubeConfigPath is empty the in-cluster config is attempted.
func NewClient(kubeConfigPath string) (*Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeConfigPath != "" {
		loadingRules.ExplicitPath = kubeConfigPath
	}
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("build rest config: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("create dynamic client: %w", err)
	}

	return &Client{dynamic: dynClient}, nil
}

// ResourceID uniquely identifies a Kubernetes resource.
type ResourceID struct {
	Group     string
	Version   string
	Resource  string
	Namespace string
	Name      string
}

// GetLiveManifest retrieves the live JSON representation of a resource,
// stripping managed-fields to reduce noise in comparisons.
func (c *Client) GetLiveManifest(ctx context.Context, id ResourceID) ([]byte, error) {
	gvr := schema.GroupVersionResource{
		Group:    id.Group,
		Version:  id.Version,
		Resource: id.Resource,
	}

	obj, err := c.dynamic.Resource(gvr).Namespace(id.Namespace).Get(ctx, id.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get resource %s/%s: %w", id.Namespace, id.Name, err)
	}

	// Strip fields that are not present in Helm-rendered manifests.
	unstructured := obj.Object
	delete(unstructured, "managedFields")
	if meta, ok := unstructured["metadata"].(map[string]interface{}); ok {
		delete(meta, "managedFields")
		delete(meta, "resourceVersion")
		delete(meta, "uid")
		delete(meta, "creationTimestamp")
		delete(meta, "generation")
	}

	out, err := obj.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal resource: %w", err)
	}
	return out, nil
}
