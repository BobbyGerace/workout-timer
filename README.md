# workout-timer

A terminal-based workout timer designed to live in a tmux pane alongside your editor. Controlled via keyboard, an in-app command prompt, or commands piped from any external process.

```
 █████╗   ██╗  ██████╗   ██████╗
 ╚════██╗ ██║  ╚════██╗ ██╔════╝
  █████╔╝ ██║   █████╔╝ ███████╗
 ██╔═══╝  ██║  ██╔═══╝  ██╔══██║
 ███████╗ ██║  ███████╗ ╚██████╔╝
```

## Features

- Big-digit clock with graceful fallback to plain text at small terminal sizes
- **Auto mode** — advances to the next interval automatically
- **Manual mode** — beeps at zero, then counts up in cyan until you advance
- **Stopwatch mode** — counts up with lap recording
- Low-time warning (clock turns yellow, configurable threshold)
- Fully configurable keybindings
- External control via FIFO pipe (great for neovim integration)
- Single-instance enforcement via file lock

## Installation

Requires Go 1.21+.

```bash
git clone https://github.com/BobbyGerace/workout-timer
cd workout-timer
make install          # installs to ~/.local/bin/timer
```

To install elsewhere:

```bash
make install INSTALL_DIR=/usr/local/bin
```

## Usage

```bash
timer                           # launch idle, configure via prompt
timer 90                        # 90-second intervals, looping
timer 1:30                      # 1m30s intervals, looping
timer auto 1:30,60,4:00 x3      # 3 rounds of [1:30 → 60s → 4:00], auto-advance
timer manual 5:00 x10           # 10 × 5-minute rounds, manual advance
timer stopwatch                 # stopwatch mode
```

### Keybindings

| Key           | Action                                  |
| ------------- | --------------------------------------- |
| `space` / `p` | Start / pause / resume                  |
| `enter` / `n` | Next interval / record lap              |
| `b`           | Previous interval                       |
| `+`           | Add 30s to current interval             |
| `-`           | Subtract 30s (floors at 0:00)           |
| `1`–`9`       | Load 1–9 minute timer                   |
| `0`           | Load 10-minute timer                    |
| `:`           | Open command prompt                     |
| `?`           | Help overlay (scroll with j/k/↑↓/PgUp/PgDn) |
| `q`           | Quit                                    |

### Commands

Commands are entered via the `:` prompt or sent over the FIFO pipe.

**Set up a program:**
```
set <time>                    loop a single interval (e.g. set 90  or  set 1:30)
set auto|manual <time>        explicit mode
set <time> xN                 N rounds
set auto <t1>,<t2>,... xN     multiple intervals per round
stopwatch                     count up from zero; enter records a lap
```

**Playback:**
```
pause                         toggle pause
next                          advance interval / record lap
back                          previous interval
add <N>                       add N seconds
subtract <N>                  subtract N seconds
reset                         restart from beginning
clear                         return to idle
quit                          exit
```

## External Control (FIFO)

The timer listens on `/tmp/workout-timer.fifo` for commands. Any command valid in the prompt is valid over the pipe:

```bash
echo 'pause'              > /tmp/workout-timer.fifo
echo 'next'               > /tmp/workout-timer.fifo
echo 'set auto 60 x10'    > /tmp/workout-timer.fifo
```

This is the intended integration point for neovim, shell scripts, or any other tool.

## Configuration

Config file location: `~/.config/workout-timer/config.toml`

```toml
# Default mode for bare `set <time>` commands: "auto" or "manual"
default_mode = "auto"

# Clock display font: "pixel" (default) or "powerline" (requires Nerd Font)
font = "pixel"

# Seconds remaining at which the timer turns yellow
low_time_warning = 30

# Play a beep sound when an interval reaches zero
beep = true

# Uncomment to override pipe / lock paths
# fifo_path = "/tmp/workout-timer.fifo"
# lock_path  = "/tmp/workout-timer.lock"

[keybindings]
# Override any key binding. Set to "" to remove a default binding.
# Values are commands — the same syntax as the : prompt.
#
# Examples:
# "f" = "add 30"
# "+" = "add 60"       # change the increment
# "q" = ""             # disable q so you don't quit by accident
```

A full list of current bindings is always visible in the `?` help overlay.
