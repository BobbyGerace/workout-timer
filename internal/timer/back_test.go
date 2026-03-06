package timer

import (
	"testing"
	"time"

	"github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

// --- Back() ---

func TestBackFromMiddleInterval(t *testing.T) {
	intervals := []time.Duration{10 * time.Second, 20 * time.Second, 30 * time.Second}
	timer := New(intervals, 1, types.ModeManual)
	timer.Start()
	timer.Next() // → interval 1
	timer.Next() // → interval 2
	timer.Back() // → interval 1
	if timer.CurrentInterval() != 1 {
		t.Errorf("expected interval 1, got %d", timer.CurrentInterval())
	}
	if timer.TimeDisplay() != 20*time.Second {
		t.Errorf("expected 20s, got %v", timer.TimeDisplay())
	}
}

func TestBackFromFirstIntervalWrapsRound(t *testing.T) {
	intervals := []time.Duration{10 * time.Second, 20 * time.Second}
	timer := New(intervals, 3, types.ModeManual)
	timer.Start()
	timer.Next() // → interval 1, round 0
	timer.Next() // → interval 0, round 1
	timer.Back() // → interval 1, round 0
	if timer.CurrentRound() != 0 {
		t.Errorf("expected round 0, got %d", timer.CurrentRound())
	}
	if timer.CurrentInterval() != 1 {
		t.Errorf("expected interval 1, got %d", timer.CurrentInterval())
	}
	if timer.TimeDisplay() != 20*time.Second {
		t.Errorf("expected 20s, got %v", timer.TimeDisplay())
	}
}

func TestBackAtStartResetsTime(t *testing.T) {
	// At (interval=0, round=0), Back doesn't move position but does reset time.
	timer := New([]time.Duration{30 * time.Second}, 3, types.ModeManual)
	timer.Start()
	timer.Tick(10 * time.Second) // 20s left
	timer.Back()
	if timer.CurrentInterval() != 0 || timer.CurrentRound() != 0 {
		t.Error("position should not change at start")
	}
	if timer.TimeDisplay() != 30*time.Second {
		t.Errorf("expected time reset to 30s, got %v", timer.TimeDisplay())
	}
}

func TestBackNoOpWhenReady(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second, 20 * time.Second}, 2, types.ModeManual)
	timer.Back()
	if timer.State() != program.ProgramReady {
		t.Error("Back in Ready state should be a no-op")
	}
}

func TestBackNoOpWhenDone(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second}, 1, types.ModeManual)
	timer.Start()
	timer.Next() // → Done
	if timer.State() != program.ProgramDone {
		t.Fatal("expected Done state")
	}
	timer.Back()
	if timer.State() != program.ProgramDone {
		t.Error("Back in Done state should be a no-op")
	}
}

func TestBackPreservesRunningState(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second, 20 * time.Second}, 1, types.ModeManual)
	timer.Start()
	timer.Next()
	timer.Back()
	if timer.State() != program.ProgramRunning {
		t.Errorf("expected Running after Back, got %v", timer.State())
	}
}

func TestBackPreservesPausedState(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second, 20 * time.Second}, 1, types.ModeManual)
	timer.Start()
	timer.Next()
	timer.TogglePause()
	timer.Back()
	if timer.State() != program.ProgramPaused {
		t.Errorf("expected Paused after Back, got %v", timer.State())
	}
}

// --- Add() ---

func TestAdd(t *testing.T) {
	timer := newManual(30 * time.Second)
	timer.Start()
	timer.Tick(10 * time.Second) // 20s left
	timer.Add(15 * time.Second)
	if timer.TimeDisplay() != 35*time.Second {
		t.Errorf("expected 35s, got %v", timer.TimeDisplay())
	}
}

func TestAddNoOpDuringOverflow(t *testing.T) {
	// Two intervals so the first doesn't auto-advance when it hits zero.
	timer := New([]time.Duration{10 * time.Second, 20 * time.Second}, 1, types.ModeManual)
	timer.Start()
	timer.Tick(15 * time.Second) // 5s overflow on first interval
	timer.Add(30 * time.Second)
	if !timer.IsOverflow() {
		t.Error("expected to still be in overflow")
	}
	if timer.TimeDisplay() != 5*time.Second {
		t.Errorf("expected overflow display 5s unchanged, got %v", timer.TimeDisplay())
	}
}

// --- Subtract() ---

func TestSubtract(t *testing.T) {
	timer := newManual(30 * time.Second)
	timer.Start()
	timer.Tick(10 * time.Second) // 20s left
	timer.Subtract(15 * time.Second)
	if timer.TimeDisplay() != 5*time.Second {
		t.Errorf("expected 5s, got %v", timer.TimeDisplay())
	}
}

func TestSubtractClampsAtZero(t *testing.T) {
	timer := newManual(30 * time.Second)
	timer.Start()
	timer.Tick(10 * time.Second) // 20s left
	timer.Subtract(50 * time.Second)
	if timer.TimeDisplay() != 0 {
		t.Errorf("expected 0, got %v", timer.TimeDisplay())
	}
	if timer.IsOverflow() {
		t.Error("subtract should not cause overflow")
	}
}

func TestSubtractNoOpDuringOverflow(t *testing.T) {
	// Two intervals so the first doesn't auto-advance when it hits zero.
	timer := New([]time.Duration{10 * time.Second, 20 * time.Second}, 1, types.ModeManual)
	timer.Start()
	timer.Tick(15 * time.Second) // 5s overflow on first interval
	timer.Subtract(30 * time.Second)
	if !timer.IsOverflow() {
		t.Error("expected to still be in overflow")
	}
	if timer.TimeDisplay() != 5*time.Second {
		t.Errorf("expected overflow display 5s unchanged, got %v", timer.TimeDisplay())
	}
}
