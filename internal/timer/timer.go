package timer

import "time"

type TimerState int

const (
	TimerPaused TimerState = iota
	TimerRunning
	TimerFinished
)

type Timer struct {
	duration  time.Duration
	remaining time.Duration
	state     TimerState
}

func New(duration time.Duration) Timer {
	return Timer{
		duration:  duration,
		remaining: duration,
		state:     TimerPaused,
	}
}

func (t *Timer) Start() {
	if t.state == TimerPaused {
		t.state = TimerRunning
	}
}

func (t *Timer) Pause() {
	if t.state == TimerRunning {
		t.state = TimerPaused
	}
}

func (t *Timer) Remaining() time.Duration {
	return t.remaining
}

func (t *Timer) State() TimerState {
	return t.state
}

func (t *Timer) Finished() bool {
	return t.state == TimerFinished
}

func (t *Timer) Tick(elapsed time.Duration) {
	if t.state == TimerPaused || t.state == TimerFinished {
		return
	}

	t.remaining = max(t.remaining - elapsed, 0)

	if (t.remaining == 0) {
		t.state = TimerFinished
	}
}
