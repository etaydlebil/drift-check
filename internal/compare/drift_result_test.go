package compare_test

import (
	"testing"

	"github.com/your-org/drift-check/internal/compare"
)

func TestHasMetadataDrift_NoChanges(t *testing.T) {
	r := compare.DriftResult{}
	if r.HasMetadataDrift() {
		t.Error("expected no metadata drift on empty result")
	}
}

func TestHasMetadataDrift_AddedLabel(t *testing.T) {
	r := compare.DriftResult{
		AddedLabels: map[string]string{"env": "prod"},
	}
	if !r.HasMetadataDrift() {
		t.Error("expected metadata drift with added label")
	}
}

func TestHasMetadataDrift_RemovedLabel(t *testing.T) {
	r := compare.DriftResult{
		RemovedLabels: map[string]string{"old": "value"},
	}
	if !r.HasMetadataDrift() {
		t.Error("expected metadata drift with removed label")
	}
}

func TestHasMetadataDrift_AddedAnnotation(t *testing.T) {
	r := compare.DriftResult{
		AddedAnnotations: map[string]string{"note": "added"},
	}
	if !r.HasMetadataDrift() {
		t.Error("expected metadata drift with added annotation")
	}
}

func TestHasMetadataDrift_RemovedAnnotation(t *testing.T) {
	r := compare.DriftResult{
		RemovedAnnotations: map[string]string{"old-note": "gone"},
	}
	if !r.HasMetadataDrift() {
		t.Error("expected metadata drift with removed annotation")
	}
}

func TestDriftResult_DriftedFlagIndependentOfMetadata(t *testing.T) {
	r := compare.DriftResult{Drifted: true}
	// Drifted flag is about spec diff; metadata drift is separate
	if r.HasMetadataDrift() {
		t.Error("Drifted flag alone should not imply metadata drift")
	}
}
