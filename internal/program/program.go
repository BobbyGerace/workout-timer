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
}
