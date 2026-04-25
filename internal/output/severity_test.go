package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/drift-check/internal/compare"
)

func TestParseSeverity_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  Severity
	}{
		{"low", SeverityLow},
		{"medium", SeverityMedium},
		{"high", SeverityHigh},
		{"unknown", SeverityUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseSeverity(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseSeverity_Invalid(t *testing.T) {
	_, err := ParseSeverity("critical")
	if err == nil {
		t.Fatal("expected error for unknown severity")
	}
}

func TestSeverity_String(t *testing.T) {
	if SeverityHigh.String() != "high" {
		t.Errorf("expected 'high', got %q", SeverityHigh.String())
	}
}

func TestKnownSeverities_ContainsAll(t *testing.T) {
	known := KnownSeverities()
	for _, s := range []string{"low", "medium", "high"} {
		found := false
		for _, k := range known {
			if k == s {
				found = true
			}
		}
		if !found {
			t.Errorf("KnownSeverities missing %q", s)
		}
	}
}

func TestClassifyDrift_Boundaries(t *testing.T) {
	cases := []struct {
		lines int
		want  Severity
	}{
		{0, SeverityUnknown},
		{1, SeverityLow},
		{5, SeverityLow},
		{6, SeverityMedium},
		{20, SeverityMedium},
		{21, SeverityHigh},
	}
	for _, tc := range cases {
		got := ClassifyDrift(tc.lines)
		if got != tc.want {
			t.Errorf("ClassifyDrift(%d) = %v, want %v", tc.lines, got, tc.want)
		}
	}
}

func TestWriteSeverityReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	results := []compare.DriftResult{
		{Kind: "Deployment", Name: "api", Namespace: "default", Diff: ""},
	}
	if err := WriteSeverityReport(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected.") {
		t.Errorf("expected no-drift message, got:\n%s", buf.String())
	}
}

func TestWriteSeverityReport_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	diff := strings.Repeat("+added line\n", 10)
	results := []compare.DriftResult{
		{Kind: "Deployment", Name: "api", Namespace: "default", Diff: diff},
	}
	if err := WriteSeverityReport(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[MEDIUM]") {
		t.Errorf("expected [MEDIUM] severity tag, got:\n%s", out)
	}
	if !strings.Contains(out, "Deployment/api") {
		t.Errorf("expected resource name in output, got:\n%s", out)
	}
}
