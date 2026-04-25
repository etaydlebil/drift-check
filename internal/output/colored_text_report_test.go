package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestWriteColoredText_NoDrift_NoColor(t *testing.T) {
	var buf bytes.Buffer
	s := Summary{Release: "myapp", Namespace: "default", HasDrift: false}
	if err := WriteColoredText(&buf, s, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK in output, got: %q", out)
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected release name in output, got: %q", out)
	}
}

func TestWriteColoredText_WithDrift_ContainsDiff(t *testing.T) {
	var buf bytes.Buffer
	s := Summary{
		Release:   "myapp",
		Namespace: "default",
		HasDrift:  true,
		Diff:      "+added: value\n-removed: value",
	}
	if err := WriteColoredText(&buf, s, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DRIFT") {
		t.Errorf("expected DRIFT in output, got: %q", out)
	}
	if !strings.Contains(out, "+added: value") {
		t.Errorf("expected diff line in output, got: %q", out)
	}
}

func TestWriteColoredText_WithError(t *testing.T) {
	var buf bytes.Buffer
	s := Summary{Err: errors.New("something broke")}
	if err := WriteColoredText(&buf, s, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %q", out)
	}
}

func TestWriteColoredText_ColorEnabled_ContainsEscapes(t *testing.T) {
	var buf bytes.Buffer
	s := Summary{Release: "myapp", Namespace: "default", HasDrift: false}
	if err := WriteColoredText(&buf, s, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("expected ANSI escape codes in colored output, got: %q", out)
	}
}

func TestWriteColoredText_DiffLineColors(t *testing.T) {
	var buf bytes.Buffer
	s := Summary{
		Release:   "myapp",
		Namespace: "default",
		HasDrift:  true,
		Diff:      "+added\n-removed\n@@ context @@\n unchanged",
	}
	if err := WriteColoredText(&buf, s, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// green for additions
	if !strings.Contains(out, "\x1b[32m+added") {
		t.Errorf("expected green addition line, got: %q", out)
	}
	// red for removals
	if !strings.Contains(out, "\x1b[31m-removed") {
		t.Errorf("expected red removal line, got: %q", out)
	}
	// cyan for hunk headers
	if !strings.Contains(out, "\x1b[36m@@ context @@") {
		t.Errorf("expected cyan hunk header, got: %q", out)
	}
}
