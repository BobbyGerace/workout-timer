package program

import "time"

type ProgramState int

const (
	ProgramReady ProgramState = iota
	ProgramRunning
	ProgramPaused
	ProgramDone
)

type Program interface {
	Tick(elapsed time.Duration) bool
	Start()
	TogglePause()
	Next()
	// Back returns to the start of the previous interval (or previous round).
	// No-op at the first interval of the first round.
	Back()
	// Add increases remaining time by d. No-op when in overflow.
	Add(d time.Duration)
	// Subtract decreases remaining time by d, floored at 0. No-op when in overflow.
	Subtract(d time.Duration)
	State() ProgramState
	// TimeDisplay returns the duration to render (always non-negative)
	TimeDisplay() time.Duration
	// IsOverflow reports whether we are past zero in manual mode
	IsOverflow() bool
	// IsLowTime determines if the timer (only in countdown mode) is less than the threshhold
	IsLowTime(threshold time.Duration) bool
	// IntervalProgress returns (current, total) interval for display.
	// Returns (0, 0) if not applicable (e.g. single interval or stopwatch).
	IntervalProgress() (current, total int)
	// RoundProgress returns (current, total) round for display.
	// Returns (0, 0) if looping forever or not applicable.
	RoundProgress() (current, total int)
}
