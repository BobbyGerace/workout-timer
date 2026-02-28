package timer

import (
	"testing"
	"time"

	"github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

func TestNextSingleIntervalSingleRound(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second}, 1, types.ModeManual)
	timer.Start()
	timer.Next()
	if timer.State() != program.ProgramDone {
		t.Errorf("expected Done after single interval/round, got %v", timer.State())
	}
}

func TestNextMultipleIntervals(t *testing.T) {
	intervals := []time.Duration{10 * time.Second, 20 * time.Second, 30 * time.Second}
	timer := New(intervals, 1, types.ModeManual)
	timer.Start()

	timer.Next()
	if timer.CurrentInterval() != 1 {
		t.Errorf("expected interval 1, got %d", timer.CurrentInterval())
	}
	if timer.TimeDisplay() != 20*time.Second {
		t.Errorf("expected 20s, got %v", timer.TimeDisplay())
	}

	timer.Next()
	if timer.CurrentInterval() != 2 {
		t.Errorf("expected interval 2, got %d", timer.CurrentInterval())
	}

	timer.Next() // completes last interval of only round
	if timer.State() != program.ProgramDone {
		t.Errorf("expected Done, got %v", timer.State())
	}
}

func TestNextMultipleRounds(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second}, 3, types.ModeManual)
	timer.Start()

	timer.Next() // end of round 0
	if timer.State() == program.ProgramDone {
		t.Error("should not be Done after round 0 of 3")
	}
	if timer.CurrentRound() != 1 {
		t.Errorf("expected round 1, got %d", timer.CurrentRound())
	}

	timer.Next() // end of round 1
	if timer.CurrentRound() != 2 {
		t.Errorf("expected round 2, got %d", timer.CurrentRound())
	}

	timer.Next() // end of round 2 (last)
	if timer.State() != program.ProgramDone {
		t.Errorf("expected Done after 3 rounds, got %v", timer.State())
	}
}

func TestNextLoopsForeverWhenRoundsZero(t *testing.T) {
	timer := New([]time.Duration{10 * time.Second}, 0, types.ModeManual)
	timer.Start()

	for i := 0; i < 10; i++ {
		timer.Next()
		if timer.State() == program.ProgramDone {
			t.Errorf("loop-forever timer reached Done on iteration %d", i)
		}
	}
}

func TestAutoModeAdvancesAtZero(t *testing.T) {
	intervals := []time.Duration{5 * time.Second, 10 * time.Second}
	timer := New(intervals, 1, types.ModeAuto)
	timer.Start()

	timer.Tick(6 * time.Second) // past first interval
	if timer.CurrentInterval() != 1 {
		t.Errorf("expected auto-advance to interval 1, got %d", timer.CurrentInterval())
	}
	if timer.State() != program.ProgramRunning {
		t.Errorf("expected still Running, got %v", timer.State())
	}
}
