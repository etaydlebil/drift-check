package output

import "fmt"

// ExitCode represents the process exit code returned by drift-check.
type ExitCode int

const (
	ExitOK    ExitCode = 0
	ExitDrift ExitCode = 1
	ExitError ExitCode = 2
)

func (e ExitCode) String() string {
	switch e {
	case ExitOK:
		return "OK"
	case ExitDrift:
		return "DRIFT"
	case ExitError:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(e))
	}
}

// Resolver maps a Summary to the appropriate ExitCode.
type Resolver struct{}

// NewResolver returns a new Resolver.
func NewResolver() *Resolver {
	return &Resolver{}
}

// Resolve returns the exit code that should be used for the given summary.
// Error takes precedence over drift.
func (r *Resolver) Resolve(s *Summary) ExitCode {
	if s == nil {
		return ExitError
	}
	if s.Err != nil {
		return ExitError
	}
	if s.HasDrift {
		return ExitDrift
	}
	return ExitOK
}
