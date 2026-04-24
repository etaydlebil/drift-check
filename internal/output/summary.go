package output

import "fmt"

// Summary holds the high-level result of a drift-check run.
type Summary struct {
	Release   string
	Namespace string
	HasDrift  bool
	Diff      string
	Err       error
}

// NewSummary constructs a Summary from the outputs of a drift run.
func NewSummary(release, namespace, diff string, err error) *Summary {
	return &Summary{
		Release:   release,
		Namespace: namespace,
		HasDrift:  diff != "",
		Diff:      diff,
		Err:       err,
	}
}

// StatusLine returns a single human-readable status line suitable for
// printing to a terminal.
func (s *Summary) StatusLine() string {
	if s.Err != nil {
		return fmt.Sprintf("[ERROR] %s/%s — %v", s.Namespace, s.Release, s.Err)
	}
	if s.HasDrift {
		return fmt.Sprintf("[DRIFT] %s/%s — configuration drift detected", s.Namespace, s.Release)
	}
	return fmt.Sprintf("[OK]    %s/%s — no drift detected", s.Namespace, s.Release)
}
