package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfig  string
	namespace   string
	helmRelease string
	outputFmt   string
)

var rootCmd = &cobra.Command{
	Use:   "drift-check",
	Short: "Detect configuration drift between running Kubernetes workloads and their source Helm charts",
	Long: `drift-check compares the live state of Kubernetes resources against
the desired state defined in Helm chart templates, reporting any differences
that indicate configuration drift.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if helmRelease == "" {
			return fmt.Errorf("--release flag is required")
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Checking drift for release %q in namespace %q\n", helmRelease, namespace)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file (defaults to in-cluster or KUBECONFIG env)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to inspect")
	rootCmd.Flags().StringVarP(&helmRelease, "release", "r", "", "Helm release name to compare against (required)")
	rootCmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "output format: text|json")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
