package timer

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	timer := New(10 * time.Second)
	if timer.Remaining() != 10*time.Second {
		t.Errorf("expected 10s remaining, got %v", timer.Remaining())
	}
	if timer.State() != TimerPaused {
		t.Errorf("expected TimerPaused state, got %v", timer.State())
	}
}

func TestTickOnlyAdvancesWhenRunning(t *testing.T) {
	timer := New(10 * time.Second)
	// Paused — tick should have no effect
	timer.Tick(2 * time.Second)
	if timer.Remaining() != 10*time.Second {
		t.Errorf("paused timer should not advance, got %v", timer.Remaining())
	}

	timer.Start()
	timer.Tick(2 * time.Second)
	if timer.Remaining() != 8*time.Second {
		t.Errorf("expected 8s remaining, got %v", timer.Remaining())
	}
}

func TestTickClampsAtZero(t *testing.T) {
	timer := New(5 * time.Second)
	timer.Start()
	// Tick more than remaining
	timer.Tick(10 * time.Second)
	if timer.Remaining() != 0 {
		t.Errorf("expected 0 remaining, got %v", timer.Remaining())
	}
}

func TestFinishedStateOnZero(t *testing.T) {
	timer := New(5 * time.Second)
	timer.Start()
	timer.Tick(5 * time.Second)
	if !timer.Finished() {
		t.Error("expected timer to be finished")
	}
	if timer.State() != TimerFinished {
		t.Errorf("expected TimerFinished state, got %v", timer.State())
	}
}

func TestTickDoesNothingWhenFinished(t *testing.T) {
	timer := New(5 * time.Second)
	timer.Start()
	timer.Tick(5 * time.Second)
	// Another tick after finished — remaining stays at 0
	timer.Tick(2 * time.Second)
	if timer.Remaining() != 0 {
		t.Errorf("finished timer should stay at 0, got %v", timer.Remaining())
	}
}

func TestPauseStopsTicking(t *testing.T) {
	timer := New(10 * time.Second)
	timer.Start()
	timer.Tick(3 * time.Second)
	timer.Pause()
	timer.Tick(3 * time.Second)
	if timer.Remaining() != 7*time.Second {
		t.Errorf("paused timer should not advance, got %v", timer.Remaining())
	}
}
