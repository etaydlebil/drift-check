package output

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// errWriter always returns an error on Write.
type errWriter struct{ err error }

func (e *errWriter) Write(_ []byte) (int, error) { return 0, e.err }

func TestMultiWriter_WritesToAllWriters(t *testing.T) {
	var a, b bytes.Buffer
	mw := NewMultiWriter(&a, &b)

	_, err := mw.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.String() != "hello" {
		t.Errorf("writer a: got %q, want %q", a.String(), "hello")
	}
	if b.String() != "hello" {
		t.Errorf("writer b: got %q, want %q", b.String(), "hello")
	}
}

func TestMultiWriter_ErrorPropagated(t *testing.T) {
	sentinel := errors.New("disk full")
	mw := NewMultiWriter(&errWriter{err: sentinel})

	_, err := mw.Write([]byte("data"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("got %v, want wrapped %v", err, sentinel)
	}
}

func TestMultiWriter_Add(t *testing.T) {
	mw := NewMultiWriter()
	if mw.Len() != 0 {
		t.Fatalf("expected 0 writers, got %d", mw.Len())
	}

	var buf bytes.Buffer
	mw.Add(&buf)
	if mw.Len() != 1 {
		t.Fatalf("expected 1 writer, got %d", mw.Len())
	}

	_, _ = mw.Write([]byte("added"))
	if buf.String() != "added" {
		t.Errorf("got %q, want %q", buf.String(), "added")
	}
}

func TestMultiWriter_ShortWrite(t *testing.T) {
	// io.Writer that only writes half the bytes
	half := &halfWriter{}
	mw := NewMultiWriter(half)

	_, err := mw.Write([]byte("abcd"))
	if !errors.Is(err, io.ErrShortWrite) {
		t.Errorf("expected ErrShortWrite, got %v", err)
	}
}

type halfWriter struct{ buf bytes.Buffer }

func (h *halfWriter) Write(p []byte) (int, error) {
	half := len(p) / 2
	h.buf.Write(p[:half])
	return half, nil
}
