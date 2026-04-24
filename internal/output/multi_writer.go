package output

import (
	"fmt"
	"io"
)

// MultiWriter fans out writes to multiple io.Writer targets.
type MultiWriter struct {
	writers []io.Writer
}

// NewMultiWriter creates a MultiWriter that writes to all provided writers.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Write implements io.Writer. It writes p to every underlying writer.
// If any writer returns an error the write stops and the error is returned.
func (m *MultiWriter) Write(p []byte) (int, error) {
	for _, w := range m.writers {
		n, err := w.Write(p)
		if err != nil {
			return n, fmt.Errorf("multi-writer: %w", err)
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

// Add appends a new writer to the fan-out list.
func (m *MultiWriter) Add(w io.Writer) {
	m.writers = append(m.writers, w)
}

// Len returns the number of underlying writers.
func (m *MultiWriter) Len() int {
	return len(m.writers)
}
