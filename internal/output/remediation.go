package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/your-org/drift-check/internal/compare"
)

// RemediationHint holds a suggested kubectl command to reconcile drift.
type RemediationHint struct {
	Resource  string
	Namespace string
	Kind      string
	Command   string
}

// NewRemediationHints builds a slice of hints from drift results.
func NewRemediationHints(results []compare.DriftResult) []RemediationHint {
	var hints []RemediationHint
	for _, r := range results {
		if r.Diff == "" {
			continue
		}
		hint := RemediationHint{
			Resource:  r.Name,
			Namespace: r.Namespace,
			Kind:      r.Kind,
			Command:   buildCommand(r),
		}
		hints = append(hints, hint)
	}
	return hints
}

func buildCommand(r compare.DriftResult) string {
	var sb strings.Builder
	if r.Namespace != "" {
		sb.WriteString(fmt.Sprintf(
			"helm upgrade --reuse-values --namespace %s <release> <chart>",
			r.Namespace,
		))
	} else {
		sb.WriteString("helm upgrade --reuse-values <release> <chart>")
	}
	return sb.String()
}

// WriteRemediationReport writes suggested remediation commands to w.
func WriteRemediationReport(w io.Writer, results []compare.DriftResult) error {
	hints := NewRemediationHints(results)
	if len(hints) == 0 {
		_, err := fmt.Fprintln(w, "No remediation required — no drift detected.")
		return err
	}
	if _, err := fmt.Fprintln(w, "Remediation suggestions:"); err != nil {
		return err
	}
	for _, h := range hints {
		line := fmt.Sprintf("  [%s/%s] %s", h.Kind, h.Resource, h.Command)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
