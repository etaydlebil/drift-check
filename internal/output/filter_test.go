package output_test

import (
	"testing"

	"github.com/yourusername/drift-check/internal/output"
)

func TestFilter_DefaultPassesAll(t *testing.T) {
	f := output.NewDriftFilter()
	if !f.Match("default", "Deployment", false) {
		t.Error("default filter should pass non-drifted resources")
	}
	if !f.Match("default", "Deployment", true) {
		t.Error("default filter should pass drifted resources")
	}
}

func TestFilter_OnlyDrifted(t *testing.T) {
	f := output.NewDriftFilter().WithOnlyDrifted(true)
	if f.Match("default", "Deployment", false) {
		t.Error("OnlyDrifted filter should reject non-drifted")
	}
	if !f.Match("default", "Deployment", true) {
		t.Error("OnlyDrifted filter should pass drifted")
	}
}

func TestFilter_NamespaceRestriction(t *testing.T) {
	f := output.NewDriftFilter().WithNamespaces("production")
	if f.Match("staging", "Deployment", true) {
		t.Error("should reject resource in wrong namespace")
	}
	if !f.Match("production", "Deployment", true) {
		t.Error("should pass resource in allowed namespace")
	}
}

func TestFilter_NamespaceCaseInsensitive(t *testing.T) {
	f := output.NewDriftFilter().WithNamespaces("Production")
	if !f.Match("production", "Deployment", true) {
		t.Error("namespace match should be case-insensitive")
	}
}

func TestFilter_KindRestriction(t *testing.T) {
	f := output.NewDriftFilter().WithKinds("Deployment")
	if f.Match("default", "ConfigMap", true) {
		t.Error("should reject resource of wrong kind")
	}
	if !f.Match("default", "Deployment", true) {
		t.Error("should pass resource of allowed kind")
	}
}

func TestFilter_KindCaseInsensitive(t *testing.T) {
	f := output.NewDriftFilter().WithKinds("deployment")
	if !f.Match("default", "Deployment", true) {
		t.Error("kind match should be case-insensitive")
	}
}

func TestFilter_CombinedRules(t *testing.T) {
	f := output.NewDriftFilter().
		WithOnlyDrifted(true).
		WithNamespaces("default").
		WithKinds("Deployment")

	if f.Match("default", "Deployment", false) {
		t.Error("should fail: not drifted")
	}
	if f.Match("staging", "Deployment", true) {
		t.Error("should fail: wrong namespace")
	}
	if f.Match("default", "ConfigMap", true) {
		t.Error("should fail: wrong kind")
	}
	if !f.Match("default", "Deployment", true) {
		t.Error("should pass all combined rules")
	}
}

func TestFilter_MultipleNamespaces(t *testing.T) {
	f := output.NewDriftFilter().WithNamespaces("default", "production")
	if !f.Match("default", "Pod", true) {
		t.Error("should match first namespace")
	}
	if !f.Match("production", "Pod", true) {
		t.Error("should match second namespace")
	}
	if f.Match("staging", "Pod", true) {
		t.Error("should not match unlisted namespace")
	}
}
