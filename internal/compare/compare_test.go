package compare_test

import (
	"errors"
	"testing"

	"github.com/yourusername/drift-check/internal/compare"
)

// fakeHelmClient satisfies the interface used by Runner via duck-typing in tests.
type fakeHelmClient struct {
	manifest string
	err      error
}

func (f *fakeHelmClient) GetManifest(release, namespace string) (string, error) {
	return f.manifest, f.err
}

type fakeK8sClient struct {
	manifest string
	err      error
}

func (f *fakeK8sClient) GetLiveManifest(namespace, name string) (string, error) {
	return f.manifest, f.err
}

func TestRun_NoDrift(t *testing.T) {
	manifest := "apiVersion: v1\nkind: ConfigMap\n"
	runner := compare.NewRunnerFromInterfaces(
		&fakeHelmClient{manifest: manifest},
		&fakeK8sClient{manifest: manifest},
	)

	result, err := runner.Run("my-release", "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift {
		t.Errorf("expected no drift, got diff:\n%s", result.Diff)
	}
}

func TestRun_DetectsDrift(t *testing.T) {
	runner := compare.NewRunnerFromInterfaces(
		&fakeHelmClient{manifest: "replicas: 1\n"},
		&fakeK8sClient{manifest: "replicas: 3\n"},
	)

	result, err := runner.Run("my-release", "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift {
		t.Error("expected drift to be detected")
	}
}

func TestRun_HelmErrorPropagated(t *testing.T) {
	runner := compare.NewRunnerFromInterfaces(
		&fakeHelmClient{err: errors.New("helm error")},
		&fakeK8sClient{},
	)

	_, err := runner.Run("bad-release", "default")
	if err == nil {
		t.Fatal("expected error from helm client, got nil")
	}
}

func TestRun_K8sErrorPropagated(t *testing.T) {
	runner := compare.NewRunnerFromInterfaces(
		&fakeHelmClient{manifest: "some: manifest\n"},
		&fakeK8sClient{err: errors.New("k8s error")},
	)

	_, err := runner.Run("my-release", "default")
	if err == nil {
		t.Fatal("expected error from k8s client, got nil")
	}
}
