package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/BobbyGerace/workout-timer/internal/renderer"
)

var labelStyle = lipgloss.NewStyle().Faint(true)

var timerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15"))

var hintStyle = lipgloss.NewStyle().
	Faint(true)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9"))

var lowTimeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3"))

var overflowStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6"))

var pausedStyle = lipgloss.NewStyle().
	Faint(true)

func (m Model) View() string {
	promptLines := m.renderPrompt()
	promptHeight := len(promptLines)

	mainHeight := max(m.height-promptHeight, 0)

	var mainContent string
	switch m.AppState() {
	case Unconfigured:
		hint := hintStyle.Render("Press : to configure or ? for help")
		mainContent = lipgloss.Place(m.width, mainHeight, lipgloss.Center, lipgloss.Center, hint)
	default:
		content := m.renderTime(mainHeight)
		if m.AppState() == Paused {
			content += "\n\n" + pausedStyle.Render("PAUSED")
		}
		mainContent = lipgloss.Place(m.width, mainHeight, lipgloss.Center, lipgloss.Top, "\n"+content)
	}

	if promptHeight == 0 {
		return mainContent
	}
	return mainContent + "\n" + strings.Join(promptLines, "\n")
}

// bigDigitHeight is the fixed row count of the big-digit font.
const bigDigitHeight = 5

func (m Model) renderTime(availableHeight int) string {
	timeStr := formatTime(m.prog.TimeDisplay())
	rows := renderer.BigDigits(timeStr)

	style := timerStyle
	timeIsLow := m.prog.IsLowTime(time.Duration(m.config.LowTimeWarning) * time.Second)
	if timeIsLow {
		style = lowTimeStyle
	} else if m.prog.IsOverflow() {
		style = overflowStyle
	}

	result := style.Render(strings.Join(rows, "\n")) + "\n"

	// Budget remaining lines for labels (each costs 1 row + 1 blank separator).
	budgetLeft := availableHeight - bigDigitHeight - 2 // -2 for the leading and trailing  "\n" in Place

	intervalCur, intervalTotal := m.prog.IntervalProgress()
	if intervalTotal > 0 && budgetLeft >= 2 {
		result += "\n" + labelStyle.Render(fmt.Sprintf("Interval %d/%d", intervalCur, intervalTotal))
		budgetLeft -= 2
	}

	roundCur, roundTotal := m.prog.RoundProgress()
	if roundTotal > 0 && budgetLeft >= 2 {
		result += "\n" + labelStyle.Render(fmt.Sprintf("Round %d/%d", roundCur, roundTotal))
	}

	return result
}

func (m Model) renderPrompt() []string {
	if !m.prompt.Open {
		return nil
	}
	lines := []string{m.prompt.Input.View()}
	if m.prompt.Error != "" {
		lines = append(lines, errorStyle.Render(m.prompt.Error))
	}
	return lines
}

// formatTime formats a duration as M:SS.
func formatTime(d time.Duration) string {
	total := int(d.Seconds())
	minutes := total / 60
	seconds := total % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
