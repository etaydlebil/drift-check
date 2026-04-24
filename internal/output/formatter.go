package output

import (
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for drift results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Result holds the drift comparison result for a release.
type Result struct {
	Release   string
	Namespace string
	HasDrift  bool
	Diff      string
}

// Formatter writes drift results to an io.Writer in a specific format.
type Formatter struct {
	w      io.Writer
	format Format
}

// NewFormatter creates a new Formatter writing to w in the given format.
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// Write outputs the Result using the configured format.
func (f *Formatter) Write(r Result) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(r)
	default:
		return f.writeText(r)
	}
}

func (f *Formatter) writeText(r Result) error {
	if !r.HasDrift {
		_, err := fmt.Fprintf(f.w, "[OK] %s/%s: no drift detected\n", r.Namespace, r.Release)
		return err
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[DRIFT] %s/%s: drift detected\n", r.Namespace, r.Release))
	for _, line := range strings.Split(r.Diff, "\n") {
		if line != "" {
			sb.WriteString("  " + line + "\n")
		}
	}
	_, err := fmt.Fprint(f.w, sb.String())
	return err
}

func (f *Formatter) writeJSON(r Result) error {
	diffEscaped := strings.ReplaceAll(r.Diff, `"`, `\"`)
	diffEscaped = strings.ReplaceAll(diffEscaped, "\n", `\n`)
	hasDrift := "false"
	if r.HasDrift {
		hasDrift = "true"
	}
	json := fmt.Sprintf(
		`{"release":%q,"namespace":%q,"has_drift":%s,"diff":%q}`,
		r.Release, r.Namespace, hasDrift, r.Diff,
	)
	_, err := fmt.Fprintln(f.w, json)
	return err
}
