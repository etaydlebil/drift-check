package output

import (
	"strings"
)

// DiffHighlighter annotates unified-diff lines with semantic labels
// so that downstream renderers can apply colour or other formatting
// without re-parsing the diff themselves.
type DiffHighlighter struct {
	colorizer *Colorizer
}

// LineKind classifies a single line of a unified diff.
type LineKind int

const (
	LineContext LineKind = iota
	LineAdded
	LineRemoved
	LineHeader
)

// HighlightedLine pairs a raw diff line with its semantic kind.
type HighlightedLine struct {
	Text string
	Kind LineKind
}

// NewDiffHighlighter returns a DiffHighlighter backed by the given Colorizer.
func NewDiffHighlighter(c *Colorizer) *DiffHighlighter {
	return &DiffHighlighter{colorizer: c}
}

// Highlight parses a unified diff string and returns annotated lines.
func (h *DiffHighlighter) Highlight(diff string) []HighlightedLine {
	if diff == "" {
		return nil
	}
	raw := strings.Split(diff, "\n")
	out := make([]HighlightedLine, 0, len(raw))
	for _, line := range raw {
		out = append(out, HighlightedLine{
			Text: line,
			Kind: classifyLine(line),
		})
	}
	return out
}

// Render returns a coloured string for a single HighlightedLine.
func (h *DiffHighlighter) Render(hl HighlightedLine) string {
	switch hl.Kind {
	case LineAdded:
		return h.colorizer.OK(hl.Text)
	case LineRemoved:
		return h.colorizer.Error(hl.Text)
	case LineHeader:
		return h.colorizer.Warn(hl.Text)
	default:
		return hl.Text
	}
}

// RenderAll returns the full diff as a single coloured string.
func (h *DiffHighlighter) RenderAll(diff string) string {
	lines := h.Highlight(diff)
	if len(lines) == 0 {
		return ""
	}
	parts := make([]string, len(lines))
	for i, l := range lines {
		parts[i] = h.Render(l)
	}
	return strings.Join(parts, "\n")
}

func classifyLine(line string) LineKind {
	if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
		return LineAdded
	}
	if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--") {
		return LineRemoved
	}
	if strings.HasPrefix(line, "@@") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "++") {
		return LineHeader
	}
	return LineContext
}
