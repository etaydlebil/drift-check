package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CommonFlags holds flag values shared across sub-commands.
type CommonFlags struct {
	Release      string
	Namespace    string
	Format       string
	Verbose      bool
	NoColor      bool
	MaxDiffLines int
}

// BindCommonFlags registers the standard flag set on cmd and populates f.
func BindCommonFlags(cmd *cobra.Command, f *CommonFlags) {
	cmd.Flags().StringVarP(&f.Release, "release", "r", "", "Helm release name (required)")
	cmd.Flags().StringVarP(&f.Namespace, "namespace", "n", "default", "Kubernetes namespace")
	cmd.Flags().StringVar(&f.Format, "format", "text", fmt.Sprintf("Output format: text, colored-text, json, table"))
	cmd.Flags().BoolVarP(&f.Verbose, "verbose", "v", false, "Enable verbose progress output")
	cmd.Flags().BoolVar(&f.NoColor, "no-color", false, "Disable ANSI color codes in output")
	cmd.Flags().IntVar(&f.MaxDiffLines, "max-diff-lines", 50, "Maximum diff lines to display per resource (0 = unlimited)")
}

// Validate returns an error if required flags are missing or invalid.
func (f *CommonFlags) Validate() error {
	if f.Release == "" {
		return fmt.Errorf("--release flag is required")
	}
	if f.MaxDiffLines < 0 {
		return fmt.Errorf("--max-diff-lines must be >= 0, got %d", f.MaxDiffLines)
	}
	return nil
}
