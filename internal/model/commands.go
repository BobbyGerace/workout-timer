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

	case "prompt":
		m, cmd := m.openPrompt()
		return m, cmd, nil

	case "set":
		p, err := parser.ParseSet(command)
		if err != nil {
			return m, nil, err
		}
		m.prog = p
		m.prompt.Open = false
		m.prompt.Error = ""
		m.prompt.Input.SetValue("")
		return m, nil, nil
	}

	return m, nil, fmt.Errorf("unknown command: %s", verb)
}
