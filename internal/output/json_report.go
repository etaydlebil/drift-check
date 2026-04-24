package output

import (
	"encoding/json"
	"io"
	"time"
)

// jsonReport is the JSON-serialisable representation of a Summary.
type jsonReport struct {
	Release   string `json:"release"`
	Namespace string `json:"namespace"`
	Drifted   bool   `json:"drifted"`
	Diff      string `json:"diff,omitempty"`
	Error     string `json:"error,omitempty"`
	CheckedAt string `json:"checked_at"`
}

// WriteJSON serialises the Summary as a JSON object to w.
func WriteJSON(w io.Writer, s Summary) error {
	rep := jsonReport{
		Release:   s.Release,
		Namespace: s.Namespace,
		Drifted:   s.Drifted,
		Diff:      s.Diff,
		CheckedAt: s.CheckedAt.Format(time.RFC3339),
	}
	if s.Error != nil {
		rep.Error = s.Error.Error()
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
