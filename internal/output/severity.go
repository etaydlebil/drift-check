package output

import "fmt"

// Severity represents the importance level of a drift finding.
type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
)

var severityNames = map[Severity]string{
	SeverityUnknown: "unknown",
	SeverityLow:     "low",
	SeverityMedium:  "medium",
	SeverityHigh:    "high",
}

var severityByName = map[string]Severity{
	"unknown": SeverityUnknown,
	"low":     SeverityLow,
	"medium":  SeverityMedium,
	"high":    SeverityHigh,
}

// String returns the lowercase name of the severity.
func (s Severity) String() string {
	if name, ok := severityNames[s]; ok {
		return name
	}
	return "unknown"
}

// ParseSeverity parses a severity level from a string.
// Returns an error if the string does not match a known severity.
func ParseSeverity(s string) (Severity, error) {
	if sev, ok := severityByName[s]; ok {
		return sev, nil
	}
	return SeverityUnknown, fmt.Errorf("unknown severity %q: must be one of %v", s, KnownSeverities())
}

// KnownSeverities returns all valid severity level strings.
func KnownSeverities() []string {
	return []string{
		SeverityLow.String(),
		SeverityMedium.String(),
		SeverityHigh.String(),
	}
}

// ClassifyDrift assigns a severity based on the number of differing lines.
func ClassifyDrift(diffLines int) Severity {
	switch {
	case diffLines == 0:
		return SeverityUnknown
	case diffLines <= 5:
		return SeverityLow
	case diffLines <= 20:
		return SeverityMedium
	default:
		return SeverityHigh
	}
}
