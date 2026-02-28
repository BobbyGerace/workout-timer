package model

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	prog "github.com/BobbyGerace/workout-timer/internal/program"
)

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tickMsg:
		return m.handleTick(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// ctrl+c always quits, regardless of state
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	if m.prompt.Open {
		switch msg.String() {
		case "esc":
			m.prompt.Open = false
			m.prompt.Error = ""
			m.prompt.Input.SetValue("")
		case "enter":
			var err error
			m, cmd, err = m.executeCommand(m.prompt.Input.Value())
			if err != nil {
				m.prompt.Error = err.Error()
				cmd = nil
			}
		default:
			m.prompt.Input, cmd = m.prompt.Input.Update(msg)
		}
		return m, cmd
	}

	if command, ok := m.config.Keybindings[msg.String()]; ok {
		m, cmd, _ = m.executeCommand(command)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleTick(msg tickMsg) (tea.Model, tea.Cmd) {
	now := time.Time(msg)
	if !m.lastTick.IsZero() && m.prog != nil && m.prog.State() == prog.ProgramRunning {
		elapsed := now.Sub(m.lastTick)
		m.prog.Tick(elapsed)
	}
	m.lastTick = now
	return m, tick()
}

// openPrompt focuses the textinput and returns the blink command.
func (m Model) openPrompt() (Model, tea.Cmd) {
	m.prompt.Open = true
	m.prompt.Error = ""
	m.prompt.Input.SetValue("")
	m.prompt.Input.Focus()
	return m, textinput.Blink
}
