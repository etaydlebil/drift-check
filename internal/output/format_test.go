package output

import (
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"text", FormatText},
		{"json", FormatJSON},
		{"table", FormatTable},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestFormat_String(t *testing.T) {
	if FormatJSON.String() != "json" {
		t.Errorf("expected 'json', got %s", FormatJSON.String())
	}
}

func TestKnownFormats_ContainsAll(t *testing.T) {
	formats := KnownFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 formats, got %d", len(formats))
	}
	seen := map[string]bool{}
	for _, f := range formats {
		seen[f] = true
	}
	for _, expected := range []string{"text", "json", "table"} {
		if !seen[expected] {
			t.Errorf("missing format %q in KnownFormats", expected)
		}
	}
}
