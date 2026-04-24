package output_test

import (
	"errors"
	"testing"

	"github.com/younited/drift-check/internal/output"
)

func TestResolve_NoDrift_ReturnsExitOK(t *testing.T) {
	r := output.NewResolver()
	s := &output.Summary{HasDrift: false, Err: nil}
	if got := r.Resolve(s); got != output.ExitOK {
		t.Fatalf("expected ExitOK, got %v", got)
	}
}

func TestResolve_WithDiff_ReturnsExitDrift(t *testing.T) {
	r := output.NewResolver()
	s := &output.Summary{HasDrift: true, Err: nil}
	if got := r.Resolve(s); got != output.ExitDrift {
		t.Fatalf("expected ExitDrift, got %v", got)
	}
}

func TestResolve_WithError_ReturnsExitError(t *testing.T) {
	r := output.NewResolver()
	s := &output.Summary{HasDrift: false, Err: errors.New("boom")}
	if got := r.Resolve(s); got != output.ExitError {
		t.Fatalf("expected ExitError, got %v", got)
	}
}

func TestResolve_ErrorTakesPrecedenceOverDiff(t *testing.T) {
	r := output.NewResolver()
	s := &output.Summary{HasDrift: true, Err: errors.New("boom")}
	if got := r.Resolve(s); got != output.ExitError {
		t.Fatalf("expected ExitError, got %v", got)
	}
}

func TestExitCode_String(t *testing.T) {
	cases := []struct {
		code output.ExitCode
		want string
	}{
		{output.ExitOK, "OK"},
		{output.ExitDrift, "DRIFT"},
		{output.ExitError, "ERROR"},
		{output.ExitCode(99), "UNKNOWN(99)"},
	}
	for _, tc := range cases {
		if got := tc.code.String(); got != tc.want {
			t.Errorf("ExitCode(%d).String() = %q, want %q", int(tc.code), got, tc.want)
		}
	}
}
