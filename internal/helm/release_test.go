package helm_test

import (
	"testing"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/user/drift-check/internal/helm"
)

func makeRelease(manifest string, values map[string]interface{}) *release.Release {
	return &release.Release{
		Name:     "test-release",
		Manifest: manifest,
		Config:   values,
		Chart:    &chart.Chart{},
	}
}

func TestGetManifest_ReturnsManifestString(t *testing.T) {
	expected := "apiVersion: v1\nkind: Service\n"
	rel := makeRelease(expected, nil)

	got := helm.GetManifest(rel)
	if got != expected {
		t.Errorf("expected manifest %q, got %q", expected, got)
	}
}

func TestGetManifest_NilReleaseReturnsEmpty(t *testing.T) {
	got := helm.GetManifest(nil)
	if got != "" {
		t.Errorf("expected empty string for nil release, got %q", got)
	}
}

func TestGetValues_ReturnsConfigValues(t *testing.T) {
	values := map[string]interface{}{
		"replicaCount": 3,
		"image":        "nginx:latest",
	}
	rel := makeRelease("", values)

	got := helm.GetValues(rel)
	if len(got) != len(values) {
		t.Errorf("expected %d values, got %d", len(values), len(got))
	}
	if got["replicaCount"] != 3 {
		t.Errorf("expected replicaCount=3, got %v", got["replicaCount"])
	}
}

func TestGetValues_NilReleaseReturnsNil(t *testing.T) {
	got := helm.GetValues(nil)
	if got != nil {
		t.Errorf("expected nil for nil release, got %v", got)
	}
}
