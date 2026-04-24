package output

import (
	"testing"
	"time"
)

func TestStageTimer_ElapsedAfterStartEnd(t *testing.T) {
	st := NewStageTimer()
	st.Start(StageFetchHelm)
	time.Sleep(5 * time.Millisecond)
	st.End(StageFetchHelm)

	elapsed := st.Elapsed(StageFetchHelm)
	if elapsed < 5*time.Millisecond {
		t.Errorf("expected elapsed >= 5ms, got %v", elapsed)
	}
}

func TestStageTimer_ElapsedWithoutEnd(t *testing.T) {
	st := NewStageTimer()
	st.Start(StageCompare)
	if d := st.Elapsed(StageCompare); d != 0 {
		t.Errorf("expected 0 when not ended, got %v", d)
	}
}

func TestStageTimer_ElapsedUnknownStage(t *testing.T) {
	st := NewStageTimer()
	if d := st.Elapsed(StageReport); d != 0 {
		t.Errorf("expected 0 for unknown stage, got %v", d)
	}
}

func TestStageTimer_StagesReturnsStarted(t *testing.T) {
	st := NewStageTimer()
	st.Start(StageFetchHelm)
	st.Start(StageFetchLive)

	stages := st.Stages()
	if len(stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(stages))
	}
}

func TestStageTimer_MultipleStagesIndependent(t *testing.T) {
	st := NewStageTimer()
	st.Start(StageFetchHelm)
	time.Sleep(2 * time.Millisecond)
	st.End(StageFetchHelm)

	st.Start(StageFetchLive)
	time.Sleep(4 * time.Millisecond)
	st.End(StageFetchLive)

	a := st.Elapsed(StageFetchHelm)
	b := st.Elapsed(StageFetchLive)
	if a >= b {
		t.Errorf("expected fetch-helm < fetch-live, got %v >= %v", a, b)
	}
}
