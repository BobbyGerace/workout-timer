package parser

import (
	"fmt"
	"strings"
	"strconv"
	"time"

	prog "github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/timer"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

// ParseSet parses a "set <duration>" command and returns a ready-to-use Program.
// Supported formats for M4: "set 90" (seconds) or "set 1:30" (m:ss).
func ParseSet(input string) (prog.Program, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 || fields[0] != "set" {
		return nil, fmt.Errorf("usage: set <duration>  (e.g. set 90 or set 1:30)")
	}

	d, err := parseDuration(fields[1])
	if err != nil {
		return nil, err
	}

	return timer.New([]time.Duration{d}, 0, types.ModeAuto), nil
}

// parseDuration converts a duration string into a time.Duration.
// Accepts plain seconds ("90") or m:ss format ("1:30").
// Returns an error for invalid input (non-numeric, bad m:ss format, seconds >= 60).
func parseDuration(s string) (time.Duration, error) {
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

		return time.Duration(minutes * 60 + seconds) * time.Second, nil
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
