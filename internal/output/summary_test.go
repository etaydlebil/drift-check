package output

import (
	"errors"
	"strings"
	"testing"
)

func TestNewSummary_NoDrift(t *testing.T) {
	s := NewSummary("my-app", "default", "", nil)
	if s.Drifted {
		t.Error("expected Drifted=false when diff is empty")
	}
	if s.Release != "my-app" {
		t.Errorf("unexpected release: %s", s.Release)
	}
	if s.Namespace != "default" {
		t.Errorf("unexpected namespace: %s", s.Namespace)
	}
}

func TestNewSummary_WithDrift(t *testing.T) {
	s := NewSummary("my-app", "default", "-old\n+new", nil)
	if !s.Drifted {
		t.Error("expected Drifted=true when diff is non-empty")
	}
}

func TestStatusLine_OK(t *testing.T) {
	s := NewSummary("svc", "prod", "", nil)
	line := s.StatusLine()
	if !strings.HasPrefix(line, "OK") {
		t.Errorf("expected OK prefix, got: %s", line)
	}
}

func TestStatusLine_Drift(t *testing.T) {
	s := NewSummary("svc", "prod", "-a\n+b\n-c\n+d", nil)
	line := s.StatusLine()
	if !strings.HasPrefix(line, "DRIFT") {
		t.Errorf("expected DRIFT prefix, got: %s", line)
	}
	if !strings.Contains(line, "changed_lines=4") {
		t.Errorf("expected changed_lines count, got: %s", line)
	}
}

func TestStatusLine_Error(t *testing.T) {
	s := NewSummary("svc", "prod", "", errors.New("helm timeout"))
	line := s.StatusLine()
	if !strings.HasPrefix(line, "ERROR") {
		t.Errorf("expected ERROR prefix, got: %s", line)
	}
	if !strings.Contains(line, "helm timeout") {
		t.Errorf("expected error message in status line, got: %s", line)
	}
}

func TestNewSummary_CheckedAtSet(t *testing.T) {
	s := NewSummary("svc", "prod", "", nil)
	if s.CheckedAt.IsZero() {
		t.Error("expected CheckedAt to be set")
	}
}
