package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// HistoryEntry records the result of a single drift-check run.
type HistoryEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Release     string    `json:"release"`
	Namespace   string    `json:"namespace"`
	Drifted     bool      `json:"drifted"`
	DriftCount  int       `json:"drift_count"`
	ErrorMsg    string    `json:"error,omitempty"`
}

// History holds an ordered list of past drift-check entries.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// NewHistory returns an empty History.
func NewHistory() *History {
	return &History{}
}

// Add appends a new entry, keeping the list sorted by timestamp ascending.
func (h *History) Add(e HistoryEntry) {
	h.Entries = append(h.Entries, e)
	sort.Slice(h.Entries, func(i, j int) bool {
		return h.Entries[i].Timestamp.Before(h.Entries[j].Timestamp)
	})
}

// Last returns the most recent entry, or nil if the history is empty.
func (h *History) Last() *HistoryEntry {
	if len(h.Entries) == 0 {
		return nil
	}
	e := h.Entries[len(h.Entries)-1]
	return &e
}

// SaveHistory serialises h to the given file path as JSON.
func SaveHistory(path string, h *History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadHistory deserialises a History from the given file path.
// If the file does not exist an empty History is returned.
func LoadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewHistory(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("unmarshal history: %w", err)
	}
	return &h, nil
}
