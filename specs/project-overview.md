# Workout Timer — Spec

A terminal-based workout timer designed for use in a tmux pane alongside neovim. Controlled via keyboard, command prompt, or external commands over a FIFO pipe.

## Display

- Timer displays in large block characters (e.g., using braille or box-drawing characters)
- Characters downsize gracefully if the terminal window is small
- Display priority when space is constrained (highest to lowest):
  1. Current time
  2. Interval counter (e.g., `3/10`)
  3. Round counter (e.g., `Round 2/3`)
  4. Labels / mode indicator
- Timer text changes color to yellow when time is low (default: <30s, configurable)
- When in manual mode and counting up after zero, the count-up time displays in cyan

## Modes

### Interval Mode

User configures one or more intervals with optional round counts. Two sub-modes:

- **Auto:** Timer automatically advances to the next interval when it reaches zero.
- **Manual:** Timer beeps at zero, then counts up (in cyan) until the user manually advances. The count-up represents elapsed rest and is expected behavior, not an error state.

### Stopwatch Mode

Timer counts up from zero indefinitely until paused. Supports laps — each lap records a split time. Lap history displays below the timer if space permits; otherwise only the current lap time and lap count are shown.

## Keybindings

All keybindings are configurable via the config file.

| Key           | Action                                                    |
| ------------- | --------------------------------------------------------- |
| `Enter`       | Advance to next interval / lap (stopwatch mode)           |
| `Space` / `p` | Start Timer / Pause / unpause (toggle)                    |
| `+`           | Add time to current timer (default: 30s, configurable)    |
| `-`           | Subtract time from current timer (floors at 0:00)         |
| `b`           | Go back to previous interval                              |
| `l`           | Alias for `Enter` (convenient for laps in stopwatch mode) |
| `?`           | Show help overlay (keybindings and commands)              |
| `:`           | Open command prompt                                       |

## Commands

All commands are entered via the `:` command prompt or received over the FIFO pipe.

### Timer Configuration

```
set <seconds|m:ss>                   # Loop a single interval forever (uses default mode)
set auto <seconds|m:ss>              # Same, with auto-advance
set manual <seconds|m:ss>            # Same, with manual advance
set <time> x<N>                      # N rounds of a single interval
set auto <t1>,<t2>,<t3> x<N>         # N rounds of multiple intervals
stopwatch                            # Start counting up from zero
```

Examples:

```
set 60                               # 60s intervals, default mode, looping
set auto 60                          # 60s intervals, auto-advance, looping
set manual 60 x10                    # 60s intervals, manual advance, 10 rounds
set auto 1:30,60,4:00 x3             # 3 rounds of [1:30 → 60s → 4:00]
```

### Playback Control

| Command        | Description                                                |
| -------------- | ---------------------------------------------------------- |
| `pause`        | Toggle pause                                               |
| `next`         | Advance to next interval (records a lap in stopwatch mode) |
| `back`         | Return to previous interval                                |
| `add <N>`      | Add N seconds to current timer                             |
| `subtract <N>` | Subtract N seconds (floors at 0:00)                        |
| `reset`        | Restart from the beginning of the current program          |
| `clear`        | Remove the current program and return to idle state        |
| `status`       | Display current configuration, mode, and progress          |
| `quit` / `q`   | Exit the program                                           |

## Audio

- Beep sound when any interval reaches zero
- Configurable (on/off, sound type) via config file

## CLI Usage

The command line interface mirrors the `set` / `stopwatch` grammar so there's one syntax to learn:

```bash
timer 90                             # Launch with 90s intervals, default mode
timer auto 1:30,60 x3               # Launch with a full program
timer stopwatch                      # Launch directly into stopwatch mode
timer                                # Launch idle, configure via command prompt
```

When launched idle, the screen displays a hint: `Press ? for help or : to configure`.

## External Control (FIFO Pipe)

The timer listens on a named pipe at `/tmp/workout-timer.fifo` for commands. Any command from the command prompt is valid over the pipe. This enables integration with neovim or any other process.

```bash
echo 'next' > /tmp/workout-timer.fifo
echo 'pause' > /tmp/workout-timer.fifo
echo 'set auto 60 x10' > /tmp/workout-timer.fifo
```

### Process Management

- On startup, the timer acquires a file lock at `/tmp/workout-timer.lock` to prevent multiple instances. If the lock is held, it exits with an error.
- The FIFO is created if it doesn't exist and reused if it does.
- On exit, the lock is released. The FIFO is optionally cleaned up.

## Configuration

Settings are read from a `.toml` config file (e.g., `~/.config/workout-timer/config.toml`).

Configurable values include:

- Default mode (`auto` or `manual`)
- Low-time warning threshold (default: 30s)
- Time increment for `+` key (default: 30s)
- Beep on/off and sound type
- Keybinding overrides
- FIFO and lock file paths
