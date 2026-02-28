package model

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/BobbyGerace/workout-timer/internal/renderer"
	"github.com/BobbyGerace/workout-timer/internal/timer"
)

var timerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

type tickMsg time.Time

type Model struct {
	width    int
	height   int
	timer    timer.Timer
	lastTick time.Time
}

func New() Model {
	t := timer.New(10 * time.Second)
	t.Start()
	return Model{
		timer: t,
	}
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tick()
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

	case tickMsg:
		now := time.Time(msg)
		if !m.lastTick.IsZero() {
			elapsed := now.Sub(m.lastTick)
			m.timer.Tick(elapsed)
		}
		m.lastTick = now
		if !m.timer.Finished() {
			return m, tick()
		}
	}

	return m, nil
}

func (m Model) View() string {
	rows := renderer.BigDigits(formatTime(m.timer.Remaining()))
	content := timerStyle.Render(strings.Join(rows, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func formatTime(d time.Duration) string {
	total := int(d.Seconds())
	minutes := total / 60
	seconds := total % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
