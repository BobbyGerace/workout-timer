package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/BobbyGerace/workout-timer/internal/parser"
)

// executeCommand dispatches a command string, returning the updated model,
// any tea.Cmd to run, and an error suitable for display in the prompt.
// Errors from keybinding dispatch are silently ignored by the caller.
func (m Model) executeCommand(command string) (Model, tea.Cmd, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return m, nil, nil
	}

	verb := strings.Fields(command)[0]

	switch verb {
	case "quit", "q":
		return m, tea.Quit, nil

	case "start":
		if m.prog != nil {
			m.prog.Start()
		}
		return m, nil, nil

	case "next":
		if m.prog != nil {
			m.prog.Next()
		}
		return m, nil, nil

	case "pause", "resume":
		if m.prog != nil {
			m.prog.TogglePause()
		}
		return m, nil, nil

	case "back":
		if m.prog != nil {
			m.prog.Back()
		}
		return m, nil, nil

	case "add", "subtract":
		if m.prog == nil {
			return m, nil, nil
		}
		parts := strings.Fields(command)
		if len(parts) != 2 {
			return m, nil, fmt.Errorf("%s requires a duration (e.g. 30 or 1:30)", verb)
		}
		d, err := parser.ParseDuration(parts[1])
		if err != nil || d <= 0 {
			return m, nil, fmt.Errorf("invalid duration: %q", parts[1])
		}
		if verb == "add" {
			m.prog.Add(d)
		} else {
			m.prog.Subtract(d)
		}
		return m, nil, nil

	case "reset":
		if m.prog != nil {
			m.prog.Reset()
			m.completionMsg = ""
		}
		return m, nil, nil

	case "clear":
		m.prog = nil
		m.completionMsg = ""
		return m, nil, nil

	case "prompt":
		m, cmd := m.openPrompt()
		return m, cmd, nil

	case "set":
		p, err := parser.ParseSet(command, m.config.DefaultMode)
		if err != nil {
			return m, nil, err
		}
		m.prog = p
		m.completionMsg = ""
		return m, nil, nil
	}

	return m, nil, fmt.Errorf("unknown command: %s", verb)
}
