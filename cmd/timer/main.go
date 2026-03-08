package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/BobbyGerace/workout-timer/internal/config"
	"github.com/BobbyGerace/workout-timer/internal/fifo"
	"github.com/BobbyGerace/workout-timer/internal/model"
	"github.com/BobbyGerace/workout-timer/internal/parser"
	"github.com/BobbyGerace/workout-timer/internal/stopwatch"
)

func main() {
	cfg := config.Default()

	lock, err := acquireLock(cfg.LockPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer lock.Close()

	m, err := buildInitialModel(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	fifo.Listen(cfg.FIFOPath, func(cmd string) {
		p.Send(model.FifoMsg{Command: cmd})
	})

	// Forward SIGTERM so Bubbletea can restore the terminal before exit.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		<-sigs
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// acquireLock opens (or creates) the lock file and acquires an exclusive
// non-blocking flock. Returns an error if another instance is already running.
// The lock is held for as long as the returned file is open.
func acquireLock(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot open lock file: %w", err)
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		f.Close()
		if err == syscall.EWOULDBLOCK {
			return nil, fmt.Errorf("another instance is already running")
		}
		return nil, fmt.Errorf("cannot acquire lock: %w", err)
	}
	return f, nil
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
