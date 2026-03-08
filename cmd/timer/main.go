package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/BobbyGerace/workout-timer/internal/config"
	"github.com/BobbyGerace/workout-timer/internal/fifo"
	"github.com/BobbyGerace/workout-timer/internal/model"
	"github.com/BobbyGerace/workout-timer/internal/parser"
	"github.com/BobbyGerace/workout-timer/internal/stopwatch"
)

func main() {
	cfg := config.Default()

	m, err := buildInitialModel(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	fifo.Listen(cfg.FIFOPath, func(cmd string) {
		p.Send(model.FifoMsg{Command: cmd})
	})

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func buildInitialModel(cfg config.Config) (model.Model, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return model.New(), nil
	}

	if args[0] == "stopwatch" {
		return model.NewWithProgram(stopwatch.New()), nil
	}

	prog, err := parser.ParseSet("set "+strings.Join(args, " "), cfg.DefaultMode)
	if err != nil {
		return model.Model{}, fmt.Errorf("invalid arguments: %w", err)
	}
	return model.NewWithProgram(prog), nil
}
