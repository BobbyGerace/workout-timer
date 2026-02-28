package timer

import (
	"testing"
	"time"

	"github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

func newManual(d time.Duration) *Timer {
	return New([]time.Duration{d}, 1, types.ModeManual)
}

func TestNewTimer(t *testing.T) {
	timer := newManual(10 * time.Second)
	if timer.TimeDisplay() != 10*time.Second {
		t.Errorf("expected 10s, got %v", timer.TimeDisplay())
	}
	if timer.State() != program.ProgramReady {
		t.Errorf("expected ProgramReady, got %v", timer.State())
	}
}

func TestTickOnlyAdvancesWhenRunning(t *testing.T) {
	timer := newManual(10 * time.Second)
	timer.Tick(2 * time.Second)
	if timer.TimeDisplay() != 10*time.Second {
		t.Errorf("ready timer should not advance, got %v", timer.TimeDisplay())
	}

	timer.Start()
	timer.Tick(2 * time.Second)
	if timer.TimeDisplay() != 8*time.Second {
		t.Errorf("expected 8s, got %v", timer.TimeDisplay())
	}
}

func TestPauseStopsTicking(t *testing.T) {
	timer := newManual(10 * time.Second)
	timer.Start()
	timer.Tick(3 * time.Second)
	timer.Pause()
	timer.Tick(3 * time.Second)
	if timer.TimeDisplay() != 7*time.Second {
		t.Errorf("paused timer should not advance, got %v", timer.TimeDisplay())
	}
}

func TestManualModeOverflow(t *testing.T) {
	timer := newManual(5 * time.Second)
	timer.Start()
	timer.Tick(7 * time.Second)

	if !timer.IsOverflow() {
		t.Error("expected overflow")
	}
	if timer.Overflow() != 2*time.Second {
		t.Errorf("expected 2s overflow, got %v", timer.Overflow())
	}
	// TimeDisplay always non-negative
	if timer.TimeDisplay() != 2*time.Second {
		t.Errorf("expected TimeDisplay 2s, got %v", timer.TimeDisplay())
	}
	// Manual mode does NOT transition to Done on overflow
	if timer.State() == program.ProgramDone {
		t.Error("manual mode should not auto-finish")
	}
}

func TestStateTransitions(t *testing.T) {
	timer := newManual(10 * time.Second)

	if timer.State() != program.ProgramReady {
		t.Errorf("expected Ready, got %v", timer.State())
	}
	timer.Start()
	if timer.State() != program.ProgramRunning {
		t.Errorf("expected Running, got %v", timer.State())
	}
	timer.Pause()
	if timer.State() != program.ProgramPaused {
		t.Errorf("expected Paused, got %v", timer.State())
	}
	timer.Start()
	if timer.State() != program.ProgramRunning {
		t.Errorf("expected Running after resume, got %v", timer.State())
	}
}

// Tests for Next() behavior are in next_test.go, added after implementation.
