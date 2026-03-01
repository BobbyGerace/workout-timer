# Workout Timer — Roadmap

Each milestone produces a runnable binary with something new to manually verify.
No milestone breaks the previous one — the app always launches cleanly.

---

## Project Structure

```
workout-timer/
  cmd/timer/          # main.go — CLI entrypoint
  internal/
    model/            # Bubbletea Model, Update, View (wires everything together)
    timer/            # Pure countdown/stopwatch logic, no TUI dependency
    parser/           # Shared command grammar (CLI args, prompt input, FIFO)
    renderer/         # Big-digit font + lipgloss layout
    audio/            # Beep abstraction
    fifo/             # FIFO listener goroutine
    config/           # TOML config loading
  specs/
```

`timer/` and `parser/` are kept free of Bubbletea so they can be unit tested
independently. Everything else lives in `model/` and calls into them.

### Architecture: Keybindings Dispatch Commands

All user-triggerable operations are named commands (e.g. `pause`, `next`,
`back`, `add 30`). Keybindings are a config-time map of `key → command string`.
The model's `Update` function dispatches the same command regardless of whether
it arrived via a keypress, the `:` prompt, or the FIFO pipe. This means:

- `internal/parser` is the single source of truth for command semantics
- Keybinding overrides in the config file are just `"<key>" = "<command>"`
- Users can bind any key to any valid command; `"<key>" = ""` unsets a key
- No action-specific code paths exist in `Update` — only command dispatch

### App States

| State          | Description                                                    |
| -------------- | -------------------------------------------------------------- |
| `Unconfigured` | No program loaded; shows hint text, no timer displayed         |
| `Ready`        | Program loaded but not yet started; timer shows starting value |
| `Running`      | Timer is ticking                                               |
| `Paused`       | Timer is frozen                                                |
| `Done`         | Program completed; displays a completion message               |

---

## Milestone 1 — Scaffold ✓

**Delivers:** A Bubbletea app that launches, shows "Workout Timer" placeholder
text, and exits cleanly.

**Verify:**

- `go run ./cmd/timer` launches without error
- `q` and `Ctrl+C` both quit

**Notes:**

- `go mod init`, establish directory layout above
- Wire up `tea.Program` with a minimal Model/Update/View
- Add `lipgloss` and `bubbletea` as dependencies
- `git init`, create `.gitignore` (Go template), and set up GitHub repo
  (user will assist with auth/remote setup)

---

## Milestone 2 — Big Digit Renderer ✓

**Delivers:** The app displays a hardcoded time (`1:23`) in the custom block
character font, centered in the terminal.

**Verify:**

- Block digits render correctly for all characters (0–9, `:`)
- Display is horizontally and vertically centered
- Resize the terminal — centering adjusts

**Notes:**

- Encode the font from `specs/character-font-reference.txt` as a map of `rune`
  → `[5]string` (5 rows, each padded to exactly 6 chars; colon has empty 5th row)
- Renderer takes a string like `"1:23"` and returns a `[]string` of composed rows
- Use `lipgloss` for centering; `tea.WindowSizeMsg` updates stored width/height
  in the model

---

## Milestone 3 — Live Countdown ✓

**Delivers:** A hardcoded 10-second countdown that ticks down and stops at
`0:00`.

**Verify:**

- Timer counts down in real time
- Stops cleanly at `0:00` (no wrap-around)

**Notes:**

- Introduce `internal/timer` package with a pure `Timer` struct: duration,
  remaining, running bool
- Bubbletea tick: emit a `tickMsg` every 100ms via `tea.Tick`; `Update` advances
  the timer by elapsed wall time (not fixed 100ms) for accuracy
- Model holds the timer struct and calls `View` to render it
- Include unit tests covering: tick accuracy, stop-at-zero, pause/resume
  behavior, overflow (count-up in a later milestone)

---

## Milestone 3.5 — Full Data Model ✓

**Delivers:** All core structs and enums defined across packages. No new
user-visible behavior — this is a foundation milestone that prevents
architecture surprises in later milestones.

**Verify:**

- `go build ./...` passes with all new types in place
- All existing tests still pass

### `internal/types`

Shared primitives imported by any package that needs them, breaking the
`model` ↔ `config` import cycle:

```go
type Mode int

const (
    ModeAuto Mode = iota
    ModeManual
)
```

### `internal/timer`

`Timer` is fully self-contained — it owns interval/round progression and
its own state machine. The existing `TimerState` enum is replaced by a
richer state that lives inside `Timer`. `timeLeft` is allowed to go
negative in manual mode; `overflow` is derived as `max(0, -timeLeft)` and
is not stored.

```go
type TimerState int

const (
    TimerReady    TimerState = iota
    TimerRunning
    TimerPaused
    TimerDone
)

type Timer struct {
    intervals       []time.Duration
    rounds          int           // 0 = loop forever
    mode            types.Mode
    currentInterval int
    currentRound    int
    timeLeft        time.Duration // can be negative in manual mode
    state           TimerState
}
```

Key behaviors:

- `Tick(elapsed)` advances `timeLeft`; in auto mode, hitting zero advances
  to the next interval (or transitions to `TimerDone`); in manual mode,
  `timeLeft` continues past zero into negative territory
- `Next()` manually advances to the next interval (used in manual mode and
  by the `next` command)
- `Back()` returns to the previous interval
- `Overflow() time.Duration` returns `max(0, -timeLeft)` — derived, not stored

### `internal/stopwatch`

```go
type StopwatchState int

const (
    StopwatchReady   StopwatchState = iota
    StopwatchRunning
    StopwatchPaused
)

type Stopwatch struct {
    elapsed  time.Duration
    laps     []time.Duration
    state    StopwatchState
}
```

Key behaviors:

- `Tick(elapsed)` advances `elapsed` when running
- `Lap()` appends current `elapsed` to `laps` and resets `elapsed` to zero
- `Laps()` returns the lap slice

### `internal/program`

`Program` is a Go interface implemented by both `Timer` and `Stopwatch`.
It is the only thing `Model` interacts with — `Model` never reaches inside.

```go
type ProgramState int

const (
    ProgramReady   ProgramState = iota
    ProgramRunning
    ProgramPaused
    ProgramDone    // not applicable to Stopwatch
)

type Program interface {
    Tick(elapsed time.Duration)
    Start()
    Pause()
    State() ProgramState
    // TimeDisplay returns the duration to render (always non-negative)
    TimeDisplay() time.Duration
    // IsOverflow reports whether we are past zero in manual mode
    IsOverflow() bool
}
```

`*Timer` and `*Stopwatch` both implement `Program`.

### `internal/model`

`AppState` is **not stored** — it is derived from `program`:

```go
func (m Model) AppState() AppState {
    if m.program == nil {
        return Unconfigured
    }
    switch m.program.State() {
    case program.ProgramReady:   return Ready
    case program.ProgramRunning: return Running
    case program.ProgramPaused:  return Paused
    case program.ProgramDone:    return Done
    }
}
```

`Toast` and `Prompt` are extracted as sub-structs:

```go
type Toast struct {
    Message string
    Expiry  time.Time
}

type Prompt struct {
    Input textinput.Model
    Error string
    Open  bool
}

type Model struct {
    width, height int
    program       program.Program  // nil when Unconfigured
    lastTick      time.Time
    toast         Toast
    prompt        Prompt
    showHelp      bool            // (M19)
    config        config.Config   // (M18)
}
```

### `internal/config`

```go
type Config struct {
    DefaultMode    types.Mode
    LowTimeWarning int               // seconds, default 30
    TimeIncrement  int               // seconds, default 30
    Beep           bool              // default true
    Keybindings    map[string]string // key → command string
    FIFOPath       string            // default /tmp/workout-timer.fifo
    LockPath       string            // default /tmp/workout-timer.lock
}

func Default() Config // returns a Config with all defaults populated
```

**Notes:**

- `promptInput textinput.Model` requires `github.com/charmbracelet/bubbles/textinput` — add the dependency in this milestone.
- Fields for future milestones can be stubbed with zero values; they don't
  need to be wired up yet.

---

## Milestone 4 — Command Prompt ✓

**Delivers:** `:` opens an inline command prompt. `set <seconds>` or
`set <m:ss>` configures the timer, transitioning to `Ready` state.

**Verify:**

- Launch shows `Unconfigured` screen with hint: `Press : to configure or ? for help`
- `:` opens a prompt at the bottom of the screen
- `set 90` loads a 90-second timer and displays `1:30` in big digits (`Ready` state)
- `set 1:30` does the same
- `Esc` cancels the prompt without changing anything
- Invalid input shows an inline error message; prompt stays open

**Notes:**

- Use the `bubbles/textinput` component
- Introduce `internal/parser` package; begin with `ParseSet` supporting a single
  duration only
- On submit, parser result is applied to model; on error, error string stored
  and rendered below the prompt
- `Ready` state: timer is displayed at its full starting value, not yet ticking
- Include unit tests for `parser.ParseSet` covering valid inputs (`90`, `1:30`,
  `0`), invalid inputs, and edge cases

---

## Milestone 5 — Start / Pause

**Delivers:** `Enter` starts a loaded timer. `Space`/`p` toggles pause.

**Verify:**

- `set 90` → `Enter` starts the countdown
- `Space` and `p` pause and unpause
- Paused state is visually indicated (dim timer or "PAUSED" label)

**Notes:**

- `Enter` transitions `Ready → Running`; tick only fires when `Running`
- `Space`/`p` toggle between `Running` and `Paused`
- Both keys dispatch commands (`next` / `pause`) through the command dispatch
  layer established in M4 — no direct action handling

---

## Milestone 6 — Low-Time Color Warning

**Delivers:** Timer text turns yellow when remaining time is below the warning
threshold (default: 30s).

**Verify:**

- Timer is white above 30s, turns yellow at 29s and below
- Works correctly after a pause/unpause

**Notes:**

- `lipgloss` color applied in renderer based on remaining seconds
- Hardcode threshold for now; will be config-driven in M18

---

## Milestone 7 — Full `set` Grammar + Intervals

**Delivers:** Full `set` command syntax — multiple intervals, round counts,
auto/manual mode flag. Interval counter (`3/10`) is displayed.

**Verify:**

- `set auto 90` — single interval, auto mode
- `set manual 60 x5` — 5 rounds, manual mode
- `set auto 1:30,60,4:00 x3` — multi-interval program
- Interval counter (e.g. `2/3`) appears below the timer
- Round counter (e.g. `Round 1/3`) appears if rounds > 1

**Notes:**

- Extend `internal/parser` to handle the full grammar; extend unit tests to
  cover multi-interval and round syntax
- Model gains `Program` struct: slice of durations, round count, mode, current
  interval index, current round
- Labels are rendered only if there is space (check against terminal height)

---

## Milestone 8 — Auto-Advance

**Delivers:** In auto mode, the timer advances to the next interval when it
reaches zero. After the last interval of the last round, the program ends.

**Verify:**

- Multi-interval program cycles through automatically
- Interval and round counters update correctly
- End of program transitions to `Done` state and displays an enthusiastic
  completion message (e.g., "Done!", "Finished!", "Great work!") below the timer

**Notes:**

- `timer.Timer` reaching zero triggers an `intervalDoneMsg` in `Update`
- Advance logic: increment interval index; if past end, increment round; if
  past last round, transition to `Done`
- Completion message: randomize from a small list; displayed until the user
  presses a key or issues a new `set`/`stopwatch` command

---

## Milestone 9 — Manual Mode + Audio

**Delivers:** In manual mode, timer beeps at zero and begins counting up in
cyan. `Enter` advances to the next interval.

**Verify:**

- At zero: audible beep, timer flips to cyan and counts upward
- `Enter` advances to next interval at any point (before or after zero)
- Count-up is visually distinct from normal countdown

**Notes:**

- Audio: start with `exec.Command("afplay", ...)` on macOS, `aplay` on Linux,
  as a pluggable call in `internal/audio` — abstract the interface so it can be
  swapped later
- `timer.Timer` gains an `Overflow` state: once it hits zero it starts counting
  up
- Renderer checks `timer.Mode == ModeManual && timer.Overflow` to apply cyan
- `Enter` dispatches the `next` command through the command layer — no new
  keybinding-specific code paths

---

## Milestone 10 — Back + Add/Subtract

**Delivers:** `b` goes back to the previous interval. `+` adds 30s. `-`
subtracts 30s (floors at 0:00).

**Verify:**

- `b` from interval 3/5 goes to interval 2/5, resets that interval's time
- `+` and `-` adjust the running timer visibly
- `-` does not go below 0:00

**Notes:**

- Back navigates the `Program` struct: decrement interval; wrap to previous
  round if needed
- `add`/`subtract` are dispatched as commands (keybindings call `add 30` and
  `subtract 30`)

---

## Milestone 11 — Toast Notifications

**Delivers:** A transient notification area at the bottom of the screen for
brief auto-dismissing messages.

**Verify:**

- `b` at the first interval of the first round shows "Already at the first interval"
- Notification disappears automatically after ~2 seconds
- A new notification immediately replaces any existing one

**Notes:**

- Store `toastMsg string` and `toastExpiry time.Time` in the model
- On each tick, clear toast if `time.Now().After(toastExpiry)`
- Rendered as a single dim line at the bottom; lowest display priority —
  dropped first if vertical space is tight
- Retrofit `b` at first interval to use this instead of silently doing nothing

---

## Milestone 12 — Remaining Prompt Commands

**Delivers:** All playback commands work from the `:` prompt: `pause`, `next`,
`back`, `add <N>`, `subtract <N>`, `reset`, `clear`, `status`, `quit`/`q`.

**Verify:**

- Each command from the spec table works as described
- `status` prints current config + mode + progress (temporary bottom line or
  overlay)
- `clear` returns to `Unconfigured` state
- `reset` restarts the current program from the beginning

**Notes:**

- Extend `internal/parser` with `ParseCommand` covering all commands; add unit
  tests
- Parser is now complete and shared by the FIFO listener in M16

---

## Milestone 13 — Stopwatch Mode

**Delivers:** `stopwatch` command starts a count-up timer. `Enter`/`l` records
a lap. Lap history displays below the timer if space permits.

**Verify:**

- `stopwatch` from `Unconfigured` starts counting up from 0:00
- `Enter` records a lap split; lap count increments
- With enough terminal height, lap history (split times) shows below
- With limited height, only current lap time + count is shown

**Notes:**

- `timer.Timer` gains `ModeStopwatch`; tick increments instead of decrements
- Lap list stored in model; renderer conditionally includes it based on height

---

## Milestone 14 — CLI Arguments

**Delivers:** Command-line arguments configure the timer on launch, using the
same grammar as the `set` command.

**Verify:**

- `go run ./cmd/timer 90` launches in `Ready` state with a 90s timer
- `go run ./cmd/timer auto 1:30,60 x3` launches with a full program
- `go run ./cmd/timer stopwatch` launches directly into stopwatch mode
- `go run ./cmd/timer` (no args) launches `Unconfigured`

**Notes:**

- `cmd/timer/main.go` calls `parser.ParseSet(os.Args[1:])` and passes the
  result to the initial model
- Reuses the existing parser; no new grammar to define

---

## Milestone 15 — Display Priority + Resize

**Delivers:** Display degrades gracefully when the terminal is small, following
the priority order from the spec.

**Verify:**

- Very tall terminal: time + interval counter + round counter + labels all visible
- Medium terminal: time + interval counter + round counter
- Small terminal: time + interval counter only
- Very small terminal: time only
- Resize in real time — display updates immediately

**Notes:**

- All layout decisions are made in `View` based on stored `width` and `height`
- Each display element has a minimum height cost; conditionally include
  bottom-up based on remaining space
- "Downsize" for very small terminals: consider a half-height fallback font or
  plain `MM:SS` text if the big digits don't fit

---

## Milestone 16 — FIFO Pipe

**Delivers:** The timer listens on `/tmp/workout-timer.fifo` for commands. Any
command valid in the prompt is valid over the pipe.

**Verify:**

- From a second terminal: `echo 'pause' > /tmp/workout-timer.fifo` toggles pause
- `echo 'next' > /tmp/workout-timer.fifo` advances the interval
- `echo 'set auto 60 x5' > /tmp/workout-timer.fifo` reconfigures a running timer

**Notes:**

- `internal/fifo` runs a goroutine that opens the FIFO for reading in a loop
  and sends each line as a `fifoCommandMsg` into the Bubbletea program via
  `program.Send()`
- FIFO is created on startup if it doesn't exist; path is configurable

---

## Milestone 17 — Process Management

**Delivers:** Only one instance can run at a time. A second launch exits with a
clear error message.

**Verify:**

- Launch the timer; open a second terminal and run `timer` again — it exits with
  an error
- Kill the first instance; the second launch now succeeds

**Notes:**

- Acquire an exclusive file lock on `/tmp/workout-timer.lock` at startup using
  `syscall.Flock` (or `golang.org/x/sys/unix`)
- Release lock on clean exit; also release on `SIGTERM`/`SIGINT`
- FIFO is cleaned up on exit (optional, per spec)

---

## Milestone 18 — Config File

**Delivers:** `~/.config/workout-timer/config.toml` is read on startup and
applied to defaults.

**Verify:**

- Set `low_time_warning = 10` in config — yellow triggers at 10s
- Set `time_increment = 60` — `+` key adds 60s
- Set `default_mode = "manual"` — bare `set <time>` uses manual mode
- Add a `[keybindings]` entry — custom key works, and setting a key to `""`
  unsets it

**Notes:**

- Use `github.com/BurntSushi/toml` for parsing
- `internal/config` defines a `Config` struct with all fields and their defaults
- Keybindings are a `[keybindings]` TOML table: e.g. `"b" = "back"`,
  `"space" = "pause"`. Set `"<key>" = ""` to unset a binding entirely.
- The config loader populates the same `map[string]string` that `Update` uses
  for command dispatch — no separate action enum needed
- Config is loaded once at startup and passed into the initial model

---

## Milestone 19 — Help Overlay

**Delivers:** `?` opens a full-screen help overlay showing all keybindings and
commands. A live status line shows the current timer state while the overlay
is open.

**Verify:**

- `?` shows a readable help screen
- A status line at the top of the overlay shows current timer state in plain
  text (e.g., `▶ 1:23 remaining — Interval 2/3`) and updates live
- Keybinding overrides from config are reflected in the help text
- Any keypress dismisses it

**Notes:**

- Help content is generated from the same keybinding map used by `Update`, so
  it always stays in sync with overrides
- Status line reuses model's existing timer state; rendered as one extra line
  before the keybinding table
- Rendered with `lipgloss` box/padding; no new packages needed

---

## Milestone 20 — Polish Pass

**Delivers:** Final edge-case handling and UX cleanup before the app is
considered complete.

**Items:**

- Confirm `b` behavior at first interval/first round (toast added in M11,
  but verify the exact message and no-op behavior is correct)
- `status` command output formatting
- Smooth handling of rapid keypresses (debounce or queue)
- Test FIFO behavior when the pipe has no reader / when commands arrive fast
- Verify audio fallback path on Linux (`aplay`/`paplay` detection)
- Clean up any lipgloss layout edge cases at unusual terminal sizes
