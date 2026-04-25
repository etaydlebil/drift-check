package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteAnnotationReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAnnotationReport(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No annotation drift") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestWriteAnnotationReport_WithAdded(t *testing.T) {
	diffs := []AnnotationDiff{
		{
			Kind:      "Deployment",
			Name:      "my-app",
			Namespace: "default",
			Added:     map[string]string{"team": "platform"},
		},
	}
	var buf bytes.Buffer
	if err := WriteAnnotationReport(&buf, diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Deployment default/my-app") {
		t.Errorf("expected resource header, got: %q", out)
	}
	if !strings.Contains(out, "+ team: platform") {
		t.Errorf("expected added annotation line, got: %q", out)
	}
}

func TestWriteAnnotationReport_WithRemoved(t *testing.T) {
	diffs := []AnnotationDiff{
		{
			Kind:    "Service",
			Name:    "svc",
			Removed: map[string]string{"deprecated": "true"},
		},
	}
	var buf bytes.Buffer
	if err := WriteAnnotationReport(&buf, diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- deprecated: true") {
		t.Errorf("expected removed annotation line, got: %q", out)
	}
}

func TestWriteAnnotationReport_WithChanged(t *testing.T) {
	diffs := []AnnotationDiff{
		{
			Kind:      "ConfigMap",
			Name:      "cfg",
			Namespace: "kube-system",
			Changed:   map[string][2]string{"version": {"v1", "v2"}},
		},
	}
	var buf bytes.Buffer
	if err := WriteAnnotationReport(&buf, diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ version: v1 -> v2") {
		t.Errorf("expected changed annotation line, got: %q", out)
	}
}

func TestWriteAnnotationReport_SortedOutput(t *testing.T) {
	diffs := []AnnotationDiff{
		{Kind: "Service", Name: "z-svc", Namespace: "default", Added: map[string]string{"a": "1"}},
		{Kind: "Deployment", Name: "a-app", Namespace: "default", Added: map[string]string{"b": "2"}},
	}
	var buf bytes.Buffer
	if err := WriteAnnotationReport(&buf, diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	idxDeployment := strings.Index(out, "Deployment")
	idxService := strings.Index(out, "Service")
	if idxDeployment > idxService {
		t.Errorf("expected Deployment before Service in sorted output")
	}
}

func TestWriteAnnotationReport_WriterError(t *testing.T) {
	diffs := []AnnotationDiff{
		{Kind: "Pod", Name: "p", Added: map[string]string{"k": "v"}},
	}
	if err := WriteAnnotationReport(&errWriter{}, diffs); err == nil {
		t.Error("expected error from failing writer")
	}
}

// errWriter is a shared test helper — reuse if already defined, else define locally.
type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}
