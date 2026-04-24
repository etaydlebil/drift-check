package output

import (
	"errors"
	"testing"
)

func TestResolve_NoDrift_ReturnsExitOK(t *testing.T) {
	r := NewResolver()
	code := r.Resolve("", nil)
	if code != ExitOK {
		t.Errorf("expected ExitOK (0), got %d", code)
	}
}

func TestResolve_WithDiff_ReturnsExitDrift(t *testing.T) {
	r := NewResolver()
	code := r.Resolve("--- a\n+++ b\n", nil)
	if code != ExitDrift {
		t.Errorf("expected ExitDrift (1), got %d", code)
	}
}

func TestResolve_WithError_ReturnsExitError(t *testing.T) {
	r := NewResolver()
	code := r.Resolve("", errors.New("something went wrong"))
	if code != ExitError {
		t.Errorf("expected ExitError (2), got %d", code)
	}
}

func TestResolve_ErrorTakesPrecedenceOverDiff(t *testing.T) {
	r := NewResolver()
	code := r.Resolve("--- a\n+++ b\n", errors.New("helm error"))
	if code != ExitError {
		t.Errorf("expected ExitError (2) when both diff and error present, got %d", code)
	}
}

func TestExitCode_String(t *testing.T) {
	tests := []struct {
		code     ExitCode
		wantStr string
	}{
		{ExitOK, "ok"},
		{ExitDrift, "drift"},
		{ExitError, "error"},
		{ExitCode(99), "unknown(99)"},
	}
	for _, tt := range tests {
		t.Run(tt.wantStr, func(t *testing.T) {
			if got := tt.code.String(); got != tt.wantStr {
				t.Errorf("ExitCode(%d).String() = %q, want %q", tt.code, got, tt.wantStr)
			}
		})
	}
}
