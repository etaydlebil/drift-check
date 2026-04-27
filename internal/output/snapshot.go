package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/example/drift-check/internal/compare"
)

// Snapshot records drift results at a point in time for later comparison.
type Snapshot struct {
	CapturedAt time.Time                     `json:"captured_at"`
	Release    string                         `json:"release"`
	Namespace  string                         `json:"namespace"`
	Results    []compare.DriftResult          `json:"results"`
}

// NewSnapshot creates a Snapshot from the given drift results.
func NewSnapshot(release, namespace string, results []compare.DriftResult) *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Release:    release,
		Namespace:  namespace,
		Results:    results,
	}
}

// SaveSnapshot serialises the snapshot as JSON to the given file path.
func SaveSnapshot(path string, s *Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads and deserialises a snapshot from the given file path.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// DriftedCount returns the number of results that have non-empty diffs.
func (s *Snapshot) DriftedCount() int {
	count := 0
	for _, r := range s.Results {
		if r.Diff != "" {
			count++
		}
	}
	return count
}

// DriftedResults returns only the results that have non-empty diffs.
func (s *Snapshot) DriftedResults() []compare.DriftResult {
	var drifted []compare.DriftResult
	for _, r := range s.Results {
		if r.Diff != "" {
			drifted = append(drifted, r)
		}
	}
	return drifted
}
