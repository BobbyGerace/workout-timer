package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	prog "github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/timer"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

// ParseSet parses a "set" command with the full grammar and returns a ready-to-use Program.
//
// Grammar:
//
//	set [auto|manual] <t1>[,<t2>,...] [xN]
//
// Examples:
//
//	set 90
//	set auto 1:30
//	set manual 60 x5
//	set auto 1:30,60,4:00 x3
//
// When no mode flag is given, defaultMode is used.
// When no round count is given, rounds defaults to 0 (loop forever).
func ParseSet(input string, defaultMode types.Mode) (prog.Program, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 || fields[0] != "set" {
		return nil, fmt.Errorf("usage: set [auto|manual] <duration>[,...] [xN]")
	}
	tokens := fields[1:] // drop "set"

	mode := defaultMode
	switch tokens[0] {
	case "manual":
		mode = types.ModeManual
		tokens = tokens[1:]
	case "auto":
		mode = types.ModeAuto
		tokens = tokens[1:]
	}

	// head should now be intervals
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Missing duration")
	}

	intervals, err := parseDurationList(tokens[0])
	if err != nil {
		return nil, err
	}

	rounds := 0
	if len(tokens) > 1 {
		rounds, err = parseRounds(tokens[1])
		if err != nil {
			return nil, err
		}
	}

	return timer.New(intervals, rounds, mode), nil
}

// ParseCommand validates a command string without executing it.
// Returns nil if the command is syntactically valid, or an error describing the problem.
// This is the canonical validator shared by the prompt, FIFO listener, and CLI.
func ParseCommand(input string, defaultMode types.Mode) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Errorf("empty command")
	}
	fields := strings.Fields(input)
	verb := fields[0]

	switch verb {
	case "quit", "q", "start", "next", "pause", "resume", "back",
		"reset", "clear", "status", "stopwatch":
		if len(fields) != 1 {
			return fmt.Errorf("%s takes no arguments", verb)
		}
		return nil

	case "add", "subtract":
		if len(fields) != 2 {
			return fmt.Errorf("%s requires a duration (e.g. 30 or 1:30)", verb)
		}
		d, err := ParseDuration(fields[1])
		if err != nil || d <= 0 {
			return fmt.Errorf("invalid duration: %q", fields[1])
		}
		return nil

	case "set":
		_, err := ParseSet(input, defaultMode)
		return err

	default:
		return fmt.Errorf("unknown command: %q", verb)
	}
}

// parseDurationList splits a comma-separated duration string and parses each segment.
func parseDurationList(s string) ([]time.Duration, error) {
	parts := strings.Split(s, ",")
	durations := make([]time.Duration, 0, len(parts))
	for _, p := range parts {
		d, err := ParseDuration(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		durations = append(durations, d)
	}
	return durations, nil
}

// ParseDuration converts a duration string into a time.Duration.
// Accepts plain seconds ("90") or m:ss format ("1:30").
// Returns an error for invalid input (non-numeric, bad m:ss format, seconds >= 60).
func ParseDuration(s string) (time.Duration, error) {
	fields := strings.Split(s, ":")
	if len(fields) == 2 {
		minutes, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, fmt.Errorf("Invalid minutes: %s", fields[0])
		} else if minutes < 0 || minutes > 59 {
			return 0, fmt.Errorf("Minutes out of range: %s", fields[0])
		}

		seconds, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, fmt.Errorf("Invalid seconds: %s", fields[1])
		} else if seconds < 0 || seconds > 59 {
			return 0, fmt.Errorf("Seconds out of range: %s", fields[1])
		}

		return time.Duration(minutes*60+seconds) * time.Second, nil
	} else if len(fields) == 1 {
		seconds, err := strconv.Atoi(fields[0])
		if err != nil || seconds < 0 {
			return 0, fmt.Errorf("Invalid duration: %s", fields[0])
		}

		return time.Duration(seconds) * time.Second, nil
	} else {
		return 0, fmt.Errorf("Invalid duration syntax")
	}
}

// parseRounds parses an "xN" token and returns N.
// Returns an error if the format is invalid or N < 1.
func parseRounds(s string) (int, error) {
	if len(s) < 2 || s[0] != 'x' {
		return 0, fmt.Errorf("invalid round count %q: expected format xN (e.g. x3)", s)
	}
	n, err := strconv.Atoi(s[1:])
	if err != nil || n < 1 {
		return 0, fmt.Errorf("invalid round count %q: N must be a positive integer", s)
	}
	return n, nil
}

// ensure *timer.Timer satisfies prog.Program at compile time
var _ prog.Program = (*timer.Timer)(nil)
