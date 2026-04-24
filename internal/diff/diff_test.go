package diff

import (
	"strings"
	"testing"
)

const baseManifest = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 2
`

func TestCompareManifests_NoDrift(t *testing.T) {
	result, err := CompareManifests(baseManifest, baseManifest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift {
		t.Errorf("expected no drift but got drift: %s", result.Details)
	}
}

func TestCompareManifests_DetectsDrift(t *testing.T) {
	live := strings.Replace(baseManifest, "replicas: 2", "replicas: 5", 1)

	result, err := CompareManifests(live, baseManifest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift {
		t.Error("expected drift but none was detected")
	}
	if !strings.Contains(result.Details, "replicas") {
		t.Errorf("diff details missing expected field; got: %s", result.Details)
	}
}

func TestCompareManifests_TrailingWhitespacIgnored(t *testing.T) {
	liveWithSpaces := strings.ReplaceAll(baseManifest, "\n", "   \n")

	result, err := CompareManifests(liveWithSpaces, baseManifest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift {
		t.Errorf("trailing whitespace should be ignored but drift was reported: %s", result.Details)
	}
}

func TestCompareManifests_EmptyInputs(t *testing.T) {
	result, err := CompareManifests("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift {
		t.Error("two empty strings should produce no drift")
	}
}
