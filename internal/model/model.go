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

type Prompt struct {
	Input textinput.Model
	Error string
	Open  bool
}

type tickMsg time.Time

// FifoMsg carries a raw command string received from the FIFO pipe.
type FifoMsg struct{ Command string }

type Model struct {
	width, height int
	prog          prog.Program // nil when Unconfigured
	lastTick      time.Time
	prompt        Prompt
	showHelp      bool          // (M19)
	config        config.Config // (M18)
	completionMsg string
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

func New(cfg config.Config) Model {
	return NewWithProgram(nil, cfg)
}

func NewWithProgram(p prog.Program, cfg config.Config) Model {
	input := textinput.New()
	input.Placeholder = "set 1:30"
	input.CharLimit = 100

	return Model{
		config:  cfg,
		prompt:  Prompt{Input: input},
		prog:    p,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}
