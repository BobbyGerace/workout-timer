package config

import "github.com/BobbyGerace/workout-timer/internal/types"

type Config struct {
	DefaultMode    types.Mode
	LowTimeWarning int               // seconds, default 30
	TimeIncrement  int               // seconds, default 30
	Beep           bool              // default true
	Keybindings    map[string]string // key → command string
	FIFOPath       string            // default /tmp/workout-timer.fifo
	LockPath       string            // default /tmp/workout-timer.lock
}

func Default() Config {
	return Config{
		DefaultMode:    types.ModeAuto,
		LowTimeWarning: 30,
		TimeIncrement:  30,
		Beep:           true,
		Keybindings: map[string]string{
			"enter": "next",
			"space": "pause",
			"p":     "pause",
			"+":     "add 30",
			"-":     "subtract 30",
			"b":     "back",
			"l":     "next",
			"?":     "help",
			":":     "prompt",
			"q":     "quit",
		},
		FIFOPath: "/tmp/workout-timer.fifo",
		LockPath: "/tmp/workout-timer.lock",
	}
}
