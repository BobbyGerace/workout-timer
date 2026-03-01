package timer

import (
	"time"

	"github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

type TimerState int

const (
	TimerReady TimerState = iota
	TimerRunning
	TimerPaused
	TimerDone
)

type Timer struct {
	intervals       []time.Duration
	rounds          int // 0 = loop forever
	mode            types.Mode
	currentInterval int
	currentRound    int
	timeLeft        time.Duration // can be negative in manual mode
	state           TimerState
}

func New(intervals []time.Duration, rounds int, mode types.Mode) *Timer {
	return &Timer{
		intervals: intervals,
		rounds:    rounds,
		mode:      mode,
		timeLeft:  intervals[0],
		state:     TimerReady,
	}
}

func (t *Timer) Start() {
	if t.state == TimerReady {
		t.state = TimerRunning
	}
}

// For convenience, this will also start the timer if it is TimerReady,
// or restart it if it is TimerDone
func (t *Timer) TogglePause() {
	if t.state == TimerRunning {
		t.state = TimerPaused
	} else {
		t.state = TimerRunning
	}
}

func (t *Timer) Tick(elapsed time.Duration) {
	if t.state != TimerRunning {
		return
	}
	t.timeLeft -= elapsed
	if t.mode == types.ModeAuto && t.timeLeft <= 0 {
		t.Next()
	}
}

// Next advances to the next interval, or to the next round if
// the current interval is the last one. If the last interval of the last
// round is complete, it transitions to TimerDone. In manual mode this is
// called by the user pressing Enter; in auto mode Tick calls it automatically.
// rounds == 0 means loop forever and should never reach TimerDone.
func (t *Timer) Next() {
	// Always inc / reset the interval, even when done
	t.currentInterval = (t.currentInterval + 1) % len(t.intervals)

	// Reset the time
	t.timeLeft = t.intervals[t.currentInterval]

	// If the interval is zero now, it means the last round was completed
	if t.currentInterval == 0 {

		// If this was the final round, reset the rounds and transition to done
		if t.rounds > 0 && t.currentRound == t.rounds-1 {
			t.state = TimerDone
			t.currentRound = 0
		} else {
			// otherwise just increment
			t.currentRound++
		}
	}
}

// Back returns to the start of the previous interval, or the previous round
// if already at the first interval. Implemented in M10.
func (t *Timer) Back() {}

// TimeDisplay returns the duration to render — always non-negative.
func (t *Timer) TimeDisplay() time.Duration {
	if t.timeLeft < 0 {
		return -t.timeLeft
	}
	return t.timeLeft
}

func (t *Timer) IsOverflow() bool {
	return t.timeLeft < 0
}

func (t *Timer) IsLowTime(threshold time.Duration) bool {
	return t.timeLeft > 0 && t.timeLeft < threshold
}

func (t *Timer) State() program.ProgramState {
	switch t.state {
	case TimerRunning:
		return program.ProgramRunning
	case TimerPaused:
		return program.ProgramPaused
	case TimerDone:
		return program.ProgramDone
	default:
		return program.ProgramReady
	}
}

func (t *Timer) CurrentInterval() int { return t.currentInterval }
func (t *Timer) CurrentRound() int    { return t.currentRound }
func (t *Timer) TotalIntervals() int  { return len(t.intervals) }
func (t *Timer) TotalRounds() int     { return t.rounds }
