package output

import "fmt"

// ExitCode represents the exit code to be returned by the CLI.
type ExitCode int

const (
	// ExitOK indicates no drift was detected.
	ExitOK ExitCode = 0
	// ExitDrift indicates drift was detected.
	ExitDrift ExitCode = 1
	// ExitError indicates an unexpected error occurred.
	ExitError ExitCode = 2
)

// Resolver determines the appropriate exit code based on drift results.
type Resolver struct{}

// NewResolver creates a new ExitCode Resolver.
func NewResolver() *Resolver {
	return &Resolver{}
}

// Resolve returns the correct ExitCode for the given drift output and error.
// - If err is non-nil, ExitError is returned.
// - If diff is non-empty, ExitDrift is returned.
// - Otherwise ExitOK is returned.
func (r *Resolver) Resolve(diff string, err error) ExitCode {
	if err != nil {
		return ExitError
	}
	if diff != "" {
		return ExitDrift
	}
	return ExitOK
}

// String returns a human-readable label for the exit code.
func (e ExitCode) String() string {
	switch e {
	case ExitOK:
		return "ok"
	case ExitDrift:
		return "drift"
	case ExitError:
		return "error"
	default:
		return fmt.Sprintf("unknown(%d)", int(e))
	}
}
