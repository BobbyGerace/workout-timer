package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/BobbyGerace/workout-timer/internal/config"
	"github.com/BobbyGerace/workout-timer/internal/model"
	"github.com/BobbyGerace/workout-timer/internal/parser"
	"github.com/BobbyGerace/workout-timer/internal/stopwatch"
)

func main() {
	m, err := buildInitialModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func buildInitialModel() (model.Model, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return model.New(), nil
	}

	if args[0] == "stopwatch" {
		return model.NewWithProgram(stopwatch.New()), nil
	}

	cfg := config.Default()
	prog, err := parser.ParseSet("set "+strings.Join(args, " "), cfg.DefaultMode)
	if err != nil {
		return model.Model{}, fmt.Errorf("invalid arguments: %w", err)
	}
	return model.NewWithProgram(prog), nil
}
