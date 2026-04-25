package output

import (
	"bytes"
	"testing"
)

func TestColorizer_DisabledReturnsPlain(t *testing.T) {
	c := NewColorizer(true)
	got := c.Sprint(ColorRed, "hello")
	if got != "hello" {
		t.Errorf("expected plain string, got %q", got)
	}
}

func TestColorizer_EnabledWrapsWithEscapes(t *testing.T) {
	c := NewColorizer(false)
	got := c.Sprint(ColorRed, "hello")
	want := "\x1b[31mhello\x1b[0m"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestColorizer_OKGreen(t *testing.T) {
	c := NewColorizer(false)
	got := c.OK("ok")
	want := "\x1b[32mok\x1b[0m"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestColorizer_WarnYellow(t *testing.T) {
	c := NewColorizer(false)
	got := c.Warn("warn")
	want := "\x1b[33mwarn\x1b[0m"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestColorizer_ErrorRed(t *testing.T) {
	c := NewColorizer(false)
	got := c.Error("err")
	want := "\x1b[31merr\x1b[0m"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestColorizer_InfoCyan(t *testing.T) {
	c := NewColorizer(false)
	got := c.Info("info")
	want := "\x1b[36minfo\x1b[0m"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestColorizer_FprintWritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	c := NewColorizer(true)
	c.Fprint(&buf, ColorGreen, "hello")
	if buf.String() != "hello" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestColorizer_FprintlnAddsNewline(t *testing.T) {
	var buf bytes.Buffer
	c := NewColorizer(true)
	c.Fprintln(&buf, ColorGreen, "hello")
	if buf.String() != "hello\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
