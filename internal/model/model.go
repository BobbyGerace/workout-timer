package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/BobbyGerace/workout-timer/internal/renderer"
)

var timerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

type Model struct {
	width  int
	height int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	rows := renderer.BigDigits("71:23")
	content := timerStyle.Render(strings.Join(rows, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
