package output

import (
	"fmt"
	"strings"
)

// Truncator limits long diff output to a configurable number of lines,
// appending a summary footer when lines are omitted.
type Truncator struct {
	maxLines int
	enabled  bool
}

// NewTruncator returns a Truncator that clips output to maxLines.
// If maxLines <= 0 truncation is disabled.
func NewTruncator(maxLines int) *Truncator {
	return &Truncator{
		maxLines: maxLines,
		enabled:  maxLines > 0,
	}
}

// Truncate returns the (possibly clipped) text and a boolean indicating
// whether any lines were omitted.
func (t *Truncator) Truncate(text string) (string, bool) {
	if !t.enabled || text == "" {
		return text, false
	}

	lines := strings.Split(text, "\n")
	if len(lines) <= t.maxLines {
		return text, false
	}

	visible := lines[:t.maxLines]
	omitted := len(lines) - t.maxLines
	footer := fmt.Sprintf("... %d more line(s) omitted (use --max-diff-lines 0 to disable truncation)", omitted)
	return strings.Join(append(visible, footer), "\n"), true
}

// MaxLines returns the configured line limit (0 means unlimited).
func (t *Truncator) MaxLines() int {
	if !t.enabled {
		return 0
	}
	return t.maxLines
}
