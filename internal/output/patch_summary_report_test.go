package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWritePatchSummaryReport_NoDrift(t *testing.T) {
	resources := []ResourcePatchSummary{
		{Namespace: "default", Kind: "Deployment", Name: "api", Summary: PatchSummary{}},
	}
	var buf bytes.Buffer
	if err := WritePatchSummaryReport(&buf, resources); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !containsStr(out, "(none)") {
		t.Errorf("expected '(none)' in output, got:\n%s", out)
	}
}

func TestWritePatchSummaryReport_WithDrift(t *testing.T) {
	resources := []ResourcePatchSummary{
		{
			Namespace: "production",
			Kind:      "Deployment",
			Name:      "web",
			Summary:   PatchSummary{Changed: 2, Added: 1},
		},
		{
			Namespace: "default",
			Kind:      "ConfigMap",
			Name:      "cfg",
			Summary:   PatchSummary{},
		},
	}
	var buf bytes.Buffer
	if err := WritePatchSummaryReport(&buf, resources); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !containsStr(out, "production") {
		t.Errorf("expected 'production' in output, got:\n%s", out)
	}
	if containsStr(out, "default") {
		t.Errorf("clean resource 'default/cfg' should be omitted, got:\n%s", out)
	}
	if !containsStr(out, "2 changed") {
		t.Errorf("expected change summary in output, got:\n%s", out)
	}
}

func TestWritePatchSummaryReport_HeaderAlwaysPresent(t *testing.T) {
	var buf bytes.Buffer
	if err := WritePatchSummaryReport(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"NAMESPACE", "KIND", "NAME", "CHANGES"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in header, got:\n%s", col, out)
		}
	}
}

func TestWritePatchSummaryReport_WriterError(t *testing.T) {
	w := &failWriter{}
	err := WritePatchSummaryReport(w, nil)
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}

type failWriter struct{}

func (f *failWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}
