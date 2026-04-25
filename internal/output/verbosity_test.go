package output

import (
	"testing"
)

func TestVerbosity_String(t *testing.T) {
	cases := []struct {
		v    Verbosity
		want string
	}{
		{VerbositySilent, "silent"},
		{VerbosityNormal, "normal"},
		{VerbosityVerbose, "verbose"},
		{VerbosityDebug, "debug"},
		{Verbosity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.v.String(); got != tc.want {
			t.Errorf("Verbosity(%d).String() = %q, want %q", tc.v, got, tc.want)
		}
	}
}

func TestParseVerbosity_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  Verbosity
	}{
		{"silent", VerbositySilent},
		{"normal", VerbosityNormal},
		{"verbose", VerbosityVerbose},
		{"debug", VerbosityDebug},
	}
	for _, tc := range cases {
		got, err := ParseVerbosity(tc.input)
		if err != nil {
			t.Fatalf("ParseVerbosity(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseVerbosity(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseVerbosity_Invalid(t *testing.T) {
	_, err := ParseVerbosity("loud")
	if err == nil {
		t.Error("expected error for unknown verbosity, got nil")
	}
}

func TestKnownVerbosities_ContainsAll(t *testing.T) {
	known := KnownVerbosities()
	expected := []string{"silent", "normal", "verbose", "debug"}
	if len(known) != len(expected) {
		t.Fatalf("KnownVerbosities() len = %d, want %d", len(known), len(expected))
	}
	set := make(map[string]bool, len(known))
	for _, k := range known {
		set[k] = true
	}
	for _, e := range expected {
		if !set[e] {
			t.Errorf("KnownVerbosities() missing %q", e)
		}
	}
}
