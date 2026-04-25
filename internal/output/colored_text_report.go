package output

import (
	"fmt"
	"io"
	"strings"
)

// WriteColoredText writes a human-readable drift report with ANSI color
// support to w. Pass noColor=true to disable escape codes (e.g. when
// redirecting to a file).
func WriteColoredText(w io.Writer, s Summary, noColor bool) error {
	c := NewColorizer(noColor)

	var statusLine string
	switch {
	case s.Err != nil:
		statusLine = c.Error(fmt.Sprintf("ERROR  %s", s.Err))
	case s.HasDrift:
		statusLine = c.Warn(fmt.Sprintf("DRIFT  release=%s namespace=%s", s.Release, s.Namespace))
	default:
		statusLine = c.OK(fmt.Sprintf("OK     release=%s namespace=%s", s.Release, s.Namespace))
	}

	if _, err := fmt.Fprintln(w, statusLine); err != nil {
		return err
	}

	if s.HasDrift && s.Diff != "" {
		lines := strings.Split(s.Diff, "\n")
		for _, line := range lines {
			var colored string
			switch {
			case strings.HasPrefix(line, "+"):
				colored = c.Sprint(ColorGreen, line)
			case strings.HasPrefix(line, "-"):
				colored = c.Sprint(ColorRed, line)
			case strings.HasPrefix(line, "@@"):
				colored = c.Sprint(ColorCyan, line)
			default:
				colored = line
			}
			if _, err := fmt.Fprintln(w, colored); err != nil {
				return err
			}
		}
	}

	return nil
}
