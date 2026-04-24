package diff

import (
	"fmt"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

// Result holds the outcome of a diff comparison.
type Result struct {
	HasDrift bool
	Details  string
}

// CompareManifests compares a live Kubernetes manifest string against the
// expected manifest rendered by Helm and returns a Result describing any drift.
func CompareManifests(live, expected string) (Result, error) {
	normLive := normalise(live)
	normExpected := normalise(expected)

	if normLive == normExpected {
		return Result{HasDrift: false}, nil
	}

	diff, err := unified(normExpected, normLive)
	if err != nil {
		return Result{}, fmt.Errorf("generating diff: %w", err)
	}

	return Result{HasDrift: true, Details: diff}, nil
}

// normalise strips trailing whitespace from each line and ensures the string
// ends with a single newline so that cosmetic differences are ignored.
func normalise(s string) string {
	lines := strings.Split(s, "\n")
	trimmed := make([]string, 0, len(lines))
	for _, l := range lines {
		trimmed = append(trimmed, strings.TrimRight(l, " \t"))
	}
	return strings.TrimRight(strings.Join(trimmed, "\n"), "\n") + "\n"
}

// unified produces a unified diff string between two texts.
func unified(expected, live string) (string, error) {
	ud := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(live),
		FromFile: "expected (helm)",
		ToFile:   "live (cluster)",
		Context:  3,
	}
	return difflib.GetUnifiedDiffString(ud)
}
