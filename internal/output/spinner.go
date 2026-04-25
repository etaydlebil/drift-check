package output

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// SpinnerWriter writes animated spinner output for long-running operations.
type SpinnerWriter struct {
	w       io.Writer
	mu      sync.Mutex
	frames  []string
	ticker  *time.Ticker
	done    chan struct{}
	label   string
	enabled bool
}

// NewSpinnerWriter creates a new SpinnerWriter. When enabled is false the
// spinner is a no-op (useful when stdout is not a TTY or verbosity is low).
func NewSpinnerWriter(w io.Writer, enabled bool) *SpinnerWriter {
	return &SpinnerWriter{
		w:       w,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		done:    make(chan struct{}),
		enabled: enabled,
	}
}

// Start begins spinning with the given label. It is safe to call multiple
// times; a running spinner is stopped before the new one starts.
func (s *SpinnerWriter) Start(label string) {
	if !s.enabled {
		return
	}
	s.Stop()
	s.mu.Lock()
	s.label = label
	s.done = make(chan struct{})
	s.ticker = time.NewTicker(80 * time.Millisecond)
	s.mu.Unlock()

	go func() {
		for i := 0; ; i++ {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.mu.Lock()
				fmt.Fprintf(s.w, "\r%s %s ", s.frames[i%len(s.frames)], s.label)
				s.mu.Unlock()
			}
		}
	}()
}

// Stop halts the spinner and clears the line.
func (s *SpinnerWriter) Stop() {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ticker == nil {
		return
	}
	s.ticker.Stop()
	close(s.done)
	s.ticker = nil
	fmt.Fprintf(s.w, "\r\033[K") // clear line
}

// Done stops the spinner and prints a completion message.
func (s *SpinnerWriter) Done(msg string) {
	s.Stop()
	if s.enabled {
		fmt.Fprintln(s.w, msg)
	}
}
