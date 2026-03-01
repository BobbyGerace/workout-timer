package stopwatch

import (
	"time"

	"github.com/BobbyGerace/workout-timer/internal/program"
)

type StopwatchState int

const (
	StopwatchReady StopwatchState = iota
	StopwatchRunning
	StopwatchPaused
)

type Stopwatch struct {
	elapsed time.Duration
	laps    []time.Duration
	state   StopwatchState
}

func New() *Stopwatch {
	return &Stopwatch{state: StopwatchReady}
}

func (s *Stopwatch) Start() {
	if s.state == StopwatchReady {
		s.state = StopwatchRunning
	}
}

func (s *Stopwatch) TogglePause() {
	if s.state == StopwatchRunning {
		s.state = StopwatchPaused
	} else if s.state == StopwatchPaused {
		s.state = StopwatchRunning
	}
}

// Next records a lap, matching the Program interface. Equivalent to Lap().
func (s *Stopwatch) Next() {
	s.Lap()
}

func (s *Stopwatch) Tick(elapsed time.Duration) {
	if s.state == StopwatchRunning {
		s.elapsed += elapsed
	}
}

func (s *Stopwatch) Lap() {
	s.laps = append(s.laps, s.elapsed)
	s.elapsed = 0
}

func (s *Stopwatch) Laps() []time.Duration {
	return s.laps
}

func (s *Stopwatch) TimeDisplay() time.Duration {
	return s.elapsed
}

func (s *Stopwatch) IsOverflow() bool {
	return false
}

func (s *Stopwatch) IsLowTime(threshold time.Duration) bool {
	return false
}

func (s *Stopwatch) State() program.ProgramState {
	switch s.state {
	case StopwatchRunning:
		return program.ProgramRunning
	case StopwatchPaused:
		return program.ProgramPaused
	default:
		return program.ProgramReady
	}
}
