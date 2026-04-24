package output

import "fmt"

// Format represents a supported output format.
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

// ParseFormat converts a string to a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON, FormatTable:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of text, json, table", s)
	}
}

// String implements fmt.Stringer.
func (f Format) String() string {
	return string(f)
}

// KnownFormats returns all supported format values as strings.
func KnownFormats() []string {
	return []string{
		FormatText.String(),
		FormatJSON.String(),
		FormatTable.String(),
	}
}
