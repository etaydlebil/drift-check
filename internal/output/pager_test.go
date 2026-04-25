package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPager_WritePassesThrough(t *testing.T) {
	var buf bytes.Buffer
	p := NewPager(&buf, false)

	_, err := p.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestPager_WriteAll_DisabledWritesDirect(t *testing.T) {
	var buf bytes.Buffer
	p := NewPager(&buf, false)

	content := strings.Repeat("line\n", 200)
	if err := p.WriteAll(content, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != content {
		t.Error("content mismatch: expected full content in buffer")
	}
}

func TestPager_WriteAll_BelowThresholdWritesDirect(t *testing.T) {
	var buf bytes.Buffer
	p := NewPager(&buf, true)
	p.pagerCmd = "nonexistent-pager-xyz" // ensure pager is unavailable

	content := "line1\nline2\nline3\n"
	if err := p.WriteAll(content, 100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != content {
		t.Errorf("expected %q, got %q", content, buf.String())
	}
}

func TestPager_WriteAll_PagerUnavailableFallsBack(t *testing.T) {
	var buf bytes.Buffer
	p := NewPager(&buf, true)
	p.pagerCmd = "nonexistent-pager-xyz"

	content := strings.Repeat("line\n", 50)
	if err := p.WriteAll(content, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != content {
		t.Error("fallback write failed: buffer does not match content")
	}
}

func TestNewPager_DefaultCmd(t *testing.T) {
	t.Setenv("PAGER", "")
	p := NewPager(&bytes.Buffer{}, true)
	if p.pagerCmd != "less" {
		t.Errorf("expected default pager \"less\", got %q", p.pagerCmd)
	}
}

func TestNewPager_EnvOverride(t *testing.T) {
	t.Setenv("PAGER", "more")
	p := NewPager(&bytes.Buffer{}, true)
	if p.pagerCmd != "more" {
		t.Errorf("expected pager \"more\", got %q", p.pagerCmd)
	}
}
