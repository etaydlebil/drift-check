package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// BaselineEntry holds a snapshot of drift state for a single release.
type BaselineEntry struct {
	Release   string    `json:"release"`
	Namespace string    `json:"namespace"`
	Diff      string    `json:"diff"`
	Captured  time.Time `json:"captured"`
}

// Baseline is a collection of baseline entries keyed by "namespace/release".
type Baseline struct {
	Entries map[string]BaselineEntry `json:"entries"`
}

// NewBaseline returns an empty Baseline.
func NewBaseline() *Baseline {
	return &Baseline{Entries: make(map[string]BaselineEntry)}
}

// Set stores or replaces the entry for the given release.
func (b *Baseline) Set(release, namespace, diff string) {
	key := fmt.Sprintf("%s/%s", namespace, release)
	b.Entries[key] = BaselineEntry{
		Release:   release,
		Namespace: namespace,
		Diff:      diff,
		Captured:  time.Now().UTC(),
	}
}

// Get retrieves the baseline entry for a release, returning false if absent.
func (b *Baseline) Get(release, namespace string) (BaselineEntry, bool) {
	key := fmt.Sprintf("%s/%s", namespace, release)
	e, ok := b.Entries[key]
	return e, ok
}

// SaveBaseline writes the baseline to the given file path as JSON.
func SaveBaseline(path string, b *Baseline) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// LoadBaseline reads a baseline from the given file path.
func LoadBaseline(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: open %s: %w", path, err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode %s: %w", path, err)
	}
	if b.Entries == nil {
		b.Entries = make(map[string]BaselineEntry)
	}
	return &b, nil
}
