package k8s

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func makeClient(t *testing.T, objs ...runtime.Object) *Client {
	t.Helper()
	scheme := runtime.NewScheme()
	dynClient := fake.NewSimpleDynamicClient(scheme, objs...)
	return &Client{dynamic: dynClient}
}

func TestGetLiveManifest_ReturnsMarshaledJSON(t *testing.T) {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "my-deploy",
				"namespace": "default",
				"uid":       "abc-123",
				"resourceVersion": "999",
			},
		},
	}
	obj.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})

	scheme := runtime.NewScheme()
	dynClient := fake.NewSimpleDynamicClient(scheme, obj)
	c := &Client{dynamic: dynClient}

	id := ResourceID{
		Group: "apps", Version: "v1", Resource: "deployments",
		Namespace: "default", Name: "my-deploy",
	}

	raw, err := c.GetLiveManifest(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("expected non-empty manifest")
	}

	var out map[string]interface{}
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	meta, _ := out["metadata"].(map[string]interface{})
	if _, exists := meta["uid"]; exists {
		t.Error("expected uid to be stripped from live manifest")
	}
	if _, exists := meta["resourceVersion"]; exists {
		t.Error("expected resourceVersion to be stripped from live manifest")
	}
}

func TestGetLiveManifest_NotFoundReturnsError(t *testing.T) {
	scheme := runtime.NewScheme()
	dynClient := fake.NewSimpleDynamicClient(scheme)
	c := &Client{dynamic: dynClient}

	id := ResourceID{
		Group: "apps", Version: "v1", Resource: "deployments",
		Namespace: "default", Name: "missing",
	}

	_, err := c.GetLiveManifest(context.Background(), id)
	if err == nil {
		t.Fatal("expected error for missing resource, got nil")
	}
}

// Ensure NewClient fails gracefully when no kubeconfig is available.
func TestNewClient_InvalidConfigReturnsError(t *testing.T) {
	_ = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	_, err := NewClient("/nonexistent/kubeconfig")
	if err == nil {
		t.Fatal("expected error with invalid kubeconfig path")
	}
}

// Compile-time check: fake.FakeDynamicClient satisfies dynamic.Interface.
var _ dynamic.Interface = (*fake.FakeDynamicClient)(nil)
