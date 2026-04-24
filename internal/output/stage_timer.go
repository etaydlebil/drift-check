package output

import (
	"sync"
	"time"
)

// StageTimer records the start and end times of named pipeline stages.
type StageTimer struct {
	mu     sync.Mutex
	starts map[Stage]time.Time
	ends   map[Stage]time.Time
}

// NewStageTimer returns an initialised StageTimer.
func NewStageTimer() *StageTimer {
	return &StageTimer{
		starts: make(map[Stage]time.Time),
		ends:   make(map[Stage]time.Time),
	}
}

// Start records the start time for stage s.
func (st *StageTimer) Start(s Stage) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.starts[s] = time.Now()
}

// End records the end time for stage s.
func (st *StageTimer) End(s Stage) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.ends[s] = time.Now()
}

// Elapsed returns the duration of stage s.
// Returns zero if the stage was never started or not yet ended.
func (st *StageTimer) Elapsed(s Stage) time.Duration {
	st.mu.Lock()
	defer st.mu.Unlock()
	start, ok := st.starts[s]
	if !ok {
		return 0
	}
	end, ok := st.ends[s]
	if !ok {
		return 0
	}
	return end.Sub(start)
}

// Stages returns all stages that have been started.
func (st *StageTimer) Stages() []Stage {
	st.mu.Lock()
	defer st.mu.Unlock()
	out := make([]Stage, 0, len(st.starts))
	for s := range st.starts {
		out = append(out, s)
	}
	return out
}
