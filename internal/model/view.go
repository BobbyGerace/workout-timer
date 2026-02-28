package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/BobbyGerace/workout-timer/internal/renderer"
)

var timerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

var hintStyle = lipgloss.NewStyle().
	Faint(true)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9"))

func (m Model) View() string {
	promptLines := m.renderPrompt()
	promptHeight := len(promptLines)

	mainHeight := m.height - promptHeight
	if mainHeight < 0 {
		mainHeight = 0
	}

	var mainContent string
	switch m.AppState() {
	case Unconfigured:
		hint := hintStyle.Render("Press : to configure or ? for help")
		mainContent = lipgloss.Place(m.width, mainHeight, lipgloss.Center, lipgloss.Center, hint)
	default:
		timeStr := formatTime(m.prog.TimeDisplay())
		rows := renderer.BigDigits(timeStr)
		content := timerStyle.Render(strings.Join(rows, "\n"))
		mainContent = lipgloss.Place(m.width, mainHeight, lipgloss.Center, lipgloss.Center, content)
	}

	if promptHeight == 0 {
		return mainContent
	}
	return mainContent + "\n" + strings.Join(promptLines, "\n")
}

func (m Model) renderPrompt() []string {
	if !m.prompt.Open {
		return nil
	}
	lines := []string{": " + m.prompt.Input.View()}
	if m.prompt.Error != "" {
		lines = append(lines, errorStyle.Render(m.prompt.Error))
	}
	return lines
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
