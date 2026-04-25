package output

import (
	"strings"
	"testing"
)

func TestTruncator_DisabledPassesThrough(t *testing.T) {
	tr := NewTruncator(0)
	input := strings.Repeat("line\n", 200)
	out, truncated := tr.Truncate(input)
	if truncated {
		t.Fatal("expected no truncation when disabled")
	}
	if out != input {
		t.Fatal("expected output to be unchanged")
	}
}

func TestTruncator_BelowLimitPassesThrough(t *testing.T) {
	tr := NewTruncator(10)
	input := "a\nb\nc"
	out, truncated := tr.Truncate(input)
	if truncated {
		t.Fatal("expected no truncation for short input")
	}
	if out != input {
		t.Fatalf("expected %q, got %q", input, out)
	}
}

func TestTruncator_ExceedsLimitIsTruncated(t *testing.T) {
	tr := NewTruncator(3)
	input := "line1\nline2\nline3\nline4\nline5"
	out, truncated := tr.Truncate(input)
	if !truncated {
		t.Fatal("expected truncation")
	}
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line3") {
		t.Error("expected visible lines to be present")
	}
	if strings.Contains(out, "line4") || strings.Contains(out, "line5") {
		t.Error("expected omitted lines to be absent")
	}
	if !strings.Contains(out, "2 more line(s) omitted") {
		t.Error("expected footer with omitted count")
	}
}

func TestTruncator_FooterMentionsFlag(t *testing.T) {
	tr := NewTruncator(1)
	input := "a\nb"
	out, _ := tr.Truncate(input)
	if !strings.Contains(out, "--max-diff-lines") {
		t.Error("expected footer to mention --max-diff-lines flag")
	}
}

func TestTruncator_EmptyInputNotTruncated(t *testing.T) {
	tr := NewTruncator(5)
	out, truncated := tr.Truncate("")
	if truncated {
		t.Fatal("empty input should not be truncated")
	}
	if out != "" {
		t.Fatal("expected empty output")
	}
}

func TestTruncator_MaxLinesReflectsConfig(t *testing.T) {
	if NewTruncator(42).MaxLines() != 42 {
		t.Error("MaxLines should return configured limit")
	}
	if NewTruncator(0).MaxLines() != 0 {
		t.Error("MaxLines should return 0 when disabled")
	}
}
