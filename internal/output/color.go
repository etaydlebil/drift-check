package output

import (
	"fmt"
	"io"
)

// Color represents an ANSI terminal color code.
type Color int

const (
	ColorReset  Color = 0
	ColorRed    Color = 31
	ColorGreen  Color = 32
	ColorYellow Color = 33
	ColorCyan   Color = 36
)

// Colorizer writes colored output to a writer when color is enabled.
type Colorizer struct {
	enabled bool
}

// NewColorizer returns a Colorizer. Color is enabled when the provided
// writer is considered a TTY-capable destination and noColor is false.
func NewColorizer(noColor bool) *Colorizer {
	return &Colorizer{enabled: !noColor}
}

// Sprint returns the string wrapped in ANSI escape codes when color is
// enabled, otherwise it returns the string unchanged.
func (c *Colorizer) Sprint(color Color, s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[%dm", color, s, ColorReset)
}

// Fprint writes a colored string to w.
func (c *Colorizer) Fprint(w io.Writer, color Color, s string) (int, error) {
	return fmt.Fprint(w, c.Sprint(color, s))
}

// Fprintln writes a colored string followed by a newline to w.
func (c *Colorizer) Fprintln(w io.Writer, color Color, s string) (int, error) {
	return fmt.Fprintln(w, c.Sprint(color, s))
}

// OK returns s colored green.
func (c *Colorizer) OK(s string) string { return c.Sprint(ColorGreen, s) }

// Warn returns s colored yellow.
func (c *Colorizer) Warn(s string) string { return c.Sprint(ColorYellow, s) }

// Error returns s colored red.
func (c *Colorizer) Error(s string) string { return c.Sprint(ColorRed, s) }

// Info returns s colored cyan.
func (c *Colorizer) Info(s string) string { return c.Sprint(ColorCyan, s) }
