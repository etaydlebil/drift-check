package output

import (
	"fmt"
	"strings"
	"time"
)

// Summary holds the result of a drift-check run.
type Summary struct {
	Release   string
	Namespace string
	Drifted   bool
	Diff      string
	Error     error
	CheckedAt time.Time
}

// NewSummary constructs a Summary for the given release.
func NewSummary(release, namespace, diff string, err error) Summary {
	return Summary{
		Release:   release,
		Namespace: namespace,
		Drifted:   diff != "",
		Diff:      diff,
		Error:     err,
		CheckedAt: time.Now().UTC(),
	}
}

// StatusLine returns a single human-readable status string.
func (s Summary) StatusLine() string {
	if s.Error != nil {
		return fmt.Sprintf("ERROR release=%s namespace=%s error=%v", s.Release, s.Namespace, s.Error)
	}
	if s.Drifted {
		lines := len(strings.Split(strings.TrimSpace(s.Diff), "\n"))
		return fmt.Sprintf("DRIFT release=%s namespace=%s changed_lines=%d", s.Release, s.Namespace, lines)
	}
	return fmt.Sprintf("OK release=%s namespace=%s", s.Release, s.Namespace)
}
