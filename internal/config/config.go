package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BobbyGerace/workout-timer/internal/types"
	"github.com/BurntSushi/toml"
)

type Config struct {
	DefaultMode    types.Mode
	Font           string
	LowTimeWarning int               // seconds, default 30
	Beep           bool              // default true
	Keybindings    map[string]string // key → command string
	FIFOPath       string            // default /tmp/workout-timer.fifo
	LockPath       string            // default /tmp/workout-timer.lock
}

// tomlFile mirrors the TOML file structure. Pointer fields let us detect
// which keys were explicitly set versus absent.
type tomlFile struct {
	DefaultMode    *string `toml:"default_mode"`
	Font           *string
	LowTimeWarning *int              `toml:"low_time_warning"`
	Beep           *bool             `toml:"beep"`
	FIFOPath       *string           `toml:"fifo_path"`
	LockPath       *string           `toml:"lock_path"`
	Keybindings    map[string]string `toml:"keybindings"`
}

// ConfigPath returns the default config file location.
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "workout-timer", "config.toml")
}

// Load reads the config file at path and merges it over the defaults.
// If the file does not exist, the defaults are returned with no error.
func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	var f tomlFile
	if _, err := toml.Decode(string(data), &f); err != nil {
		return cfg, fmt.Errorf("config parse error: %w", err)
	}

	if f.DefaultMode != nil {
		switch *f.DefaultMode {
		case "auto":
			cfg.DefaultMode = types.ModeAuto
		case "manual":
			cfg.DefaultMode = types.ModeManual
		default:
			return cfg, fmt.Errorf("unknown default_mode %q (must be \"auto\" or \"manual\")", *f.DefaultMode)
		}
	}
	if f.LowTimeWarning != nil {
		cfg.LowTimeWarning = *f.LowTimeWarning
	}
	if f.Beep != nil {
		cfg.Beep = *f.Beep
	}
	if f.Font != nil {
		cfg.Font = *f.Font
	}
	if f.FIFOPath != nil {
		cfg.FIFOPath = *f.FIFOPath
	}
	if f.LockPath != nil {
		cfg.LockPath = *f.LockPath
	}
	for k, v := range f.Keybindings {
		if v == "" {
			delete(cfg.Keybindings, k)
		} else {
			cfg.Keybindings[k] = v
		}
	}

	return cfg, nil
}

func Default() Config {
	return Config{
		DefaultMode:    types.ModeAuto,
		LowTimeWarning: 30,
		Beep:           true,
		Font:           "pixel",
		Keybindings: map[string]string{
			"space": "pause",
			"p":     "pause",
			"+":     "add 30",
			"-":     "subtract 30",
			"enter": "next",
			"n":     "next",
			"b":     "back",
			"l":     "next",
			"?":     "help",
			":":     "prompt",
			"q":     "quit",
			"1":     "set 1:00",
			"2":     "set 2:00",
			"3":     "set 3:00",
			"4":     "set 4:00",
			"5":     "set 5:00",
			"6":     "set 6:00",
			"7":     "set 7:00",
			"8":     "set 8:00",
			"9":     "set 9:00",
			"0":     "set 10:00",
		},
		FIFOPath: "/tmp/workout-timer.fifo",
		LockPath: "/tmp/workout-timer.lock",
	}
}
