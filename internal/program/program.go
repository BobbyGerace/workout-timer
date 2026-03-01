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
	Tick(elapsed time.Duration)
	Start()
	TogglePause()
	Next()
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
