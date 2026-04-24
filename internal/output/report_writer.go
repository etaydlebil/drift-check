package output

import (
	"fmt"
	"io"
)

// Format selects the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// SupportedFormats lists all valid output formats.
var SupportedFormats = []Format{FormatText, FormatJSON}

// ReportWriter writes a drift report in the requested format to w.
type ReportWriter struct {
	format Format
	w      io.Writer
}

// NewReportWriter creates a ReportWriter that emits reports in the given format.
func NewReportWriter(format Format, w io.Writer) (*ReportWriter, error) {
	switch format {
	case FormatText, FormatJSON:
		// valid
	default:
		return nil, fmt.Errorf("unsupported format %q: choose \"text\" or \"json\"", format)
	}
	return &ReportWriter{format: format, w: w}, nil
}

// Write emits the summary to the underlying writer using the chosen format.
func (r *ReportWriter) Write(s *Summary) error {
	switch r.format {
	case FormatJSON:
		return WriteJSON(r.w, s)
	default:
		return WriteText(r.w, s)
	}
}

// Format returns the configured output format.
func (r *ReportWriter) Format() Format {
	return r.format
}

// ParseFormat converts a raw string into a Format, returning an error if the
// value is not recognised. This is useful when reading format flags from CLI
// arguments or configuration files.
func ParseFormat(s string) (Format, error) {
	f := Format(s)
	for _, supported := range SupportedFormats {
		if f == supported {
			return f, nil
		}
	}
	return "", fmt.Errorf("unsupported format %q: choose \"text\" or \"json\"", s)
}
