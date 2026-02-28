package model

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/BobbyGerace/workout-timer/internal/config"
	prog "github.com/BobbyGerace/workout-timer/internal/program"
)

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
	Input textinput.Model
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
	input := textinput.New()
	input.Placeholder = "set 1:30"
	input.CharLimit = 100

	return Model{
		config: config.Default(),
		prompt: Prompt{Input: input},
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}
