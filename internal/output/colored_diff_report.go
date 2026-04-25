package output

import (
	"fmt"
	"io"
)

// WriteColoredDiff writes a syntax-highlighted unified diff to w.
// When color is disabled the output is identical to WriteColoredText.
// The highlighter is responsible for all ANSI escape handling so this
// function stays free of direct escape-sequence logic.
func WriteColoredDiff(w io.Writer, s *Summary, colorEnabled bool) error {
	c := NewColorizer(colorEnabled)
	h := NewDiffHighlighter(c)

	statusLine := s.StatusLine()
	var colored string
	switch {
	case s.Err != nil:
		colored = c.Error(statusLine)
	case s.HasDrift:
		colored = c.Warn(statusLine)
	default:
		colored = c.OK(statusLine)
	}

	if _, err := fmt.Fprintln(w, colored); err != nil {
		return err
	}

	if s.Diff == "" {
		return nil
	}

	if _, err := fmt.Fprintln(w, h.RenderAll(s.Diff)); err != nil {
		return err
	}

	return nil
}
