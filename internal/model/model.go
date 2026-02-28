package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/BobbyGerace/workout-timer/internal/config"
	prog "github.com/BobbyGerace/workout-timer/internal/program"
	"github.com/BobbyGerace/workout-timer/internal/renderer"
	"github.com/BobbyGerace/workout-timer/internal/timer"
	"github.com/BobbyGerace/workout-timer/internal/types"
)

var timerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

type AppState int

const (
	Unconfigured AppState = iota
	Ready
	Running
	Paused
	Done
)

type Toast struct {
	Message string
	Expiry  time.Time
}

type Prompt struct {
	// Input textinput.Model  // added in M4
	Error string
	Open  bool
}

type tickMsg time.Time

type Model struct {
	width, height int
	prog          prog.Program // nil when Unconfigured
	lastTick      time.Time
	toast         Toast
	prompt        Prompt
	showHelp      bool          // (M19)
	config        config.Config // (M18)
}

func (m Model) AppState() AppState {
	if m.prog == nil {
		return Unconfigured
	}
	switch m.prog.State() {
	case prog.ProgramReady:
		return Ready
	case prog.ProgramRunning:
		return Running
	case prog.ProgramPaused:
		return Paused
	case prog.ProgramDone:
		return Done
	}
	return Unconfigured
}

func New() Model {
	t := timer.New([]time.Duration{10 * time.Second}, 1, types.ModeManual)
	t.Start()
	return Model{
		prog:   t,
		config: config.Default(),
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
		if !m.lastTick.IsZero() && m.prog != nil {
			elapsed := now.Sub(m.lastTick)
			m.prog.Tick(elapsed)
		}
		m.lastTick = now
		if m.prog == nil || m.prog.State() != prog.ProgramDone {
			return m, tick()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var timeStr string
	if m.prog != nil {
		timeStr = formatTime(m.prog.TimeDisplay())
	} else {
		timeStr = "0:00"
	}
	rows := renderer.BigDigits(timeStr)
	content := timerStyle.Render(strings.Join(rows, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// formatTime formats a duration as M:SS.
// Uses ceiling so that a countdown displays "10" for the full second before
// dropping to "9". This is correct for countdown display.
// TODO: count-up display (stopwatch mode and manual-mode overflow) should use
// floor instead — "0" until a full second has elapsed, then "1", etc.
// formatTime will need a parameter or a sibling function when those are wired up.
func formatTime(d time.Duration) string {
	total := int(math.Ceil(d.Seconds()))
	minutes := total / 60
	seconds := total % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
