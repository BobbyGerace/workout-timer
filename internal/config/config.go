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
			"space": "pause",
			"p":     "pause",
			"+":     "add 30",
			"-":     "subtract 30",
			"enter": "next",
			"n":     "next",
			"b":     "back",
			"l":     "next",
			"?":     "help",
			":":     "prompt",
			"q":     "quit",
			"1":		 "set 1:00",
			"2":		 "set 2:00",
			"3":		 "set 3:00",
			"4":		 "set 4:00",
			"5":		 "set 5:00",
			"6":		 "set 6:00",
			"7":		 "set 7:00",
			"8":		 "set 8:00",
			"9":		 "set 9:00",
			"0":		 "set 10:00",
		},
		FIFOPath: "/tmp/workout-timer.fifo",
		LockPath: "/tmp/workout-timer.lock",
	}
}
