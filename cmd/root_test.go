package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd_NoReleaseFlagReturnsError(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Reset flags to ensure clean state
	helmRelease = ""

	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --release is not provided, got nil")
	}
}

func TestRootCmd_WithRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"--release", "my-app", "--namespace", "staging"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected output, got empty string")
	}

	expected := `Checking drift for release "my-app" in namespace "staging"`
	if !bytes.Contains([]byte(output), []byte(expected)) {
		t.Errorf("output %q does not contain expected string %q", output, expected)
	}
}

func TestRootCmd_DefaultNamespace(t *testing.T) {
	// Reset to defaults
	namespace = "default"
	helmRelease = ""

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"--release", "nginx"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if namespace != "default" {
		t.Errorf("expected default namespace, got %q", namespace)
	}
}
