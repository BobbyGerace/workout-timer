package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/BobbyGerace/workout-timer/internal/renderer"
	"github.com/BobbyGerace/workout-timer/internal/stopwatch"
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

var completionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")).
	Bold(true)

func (m Model) View() string {
	promptLines := m.renderPrompt()
	promptHeight := len(promptLines)

	mainHeight := max(m.height-promptHeight, 0)

	var mainContent string
	if m.AppState() == Unconfigured {
		hint := hintStyle.Render("Press : to configure or ? for help")
		mainContent = lipgloss.Place(m.width, mainHeight, lipgloss.Center, lipgloss.Center, hint)
	} else {
		switch displayTier(mainHeight) {
		case 1:
			mainContent = m.renderTier1(mainHeight)
		case 2:
			mainContent = m.renderTier2(mainHeight)
		default:
			mainContent = m.renderTier3(mainHeight)
		}
	}

	if promptHeight == 0 {
		return mainContent
	}
	return mainContent + "\n" + strings.Join(promptLines, "\n")
}

// displayTier returns 1, 2, or 3 based on available height.
// Thresholds are fixed regardless of program mode or content.
//
//	T1 (≥11): big digits + labels
//	T2 (≥4):  plain text clock + labels, left-aligned
//	T3 (<4):  single line, left-aligned
func displayTier(height int) int {
	if height >= 11 {
		return 1
	} else if height >= 4 {
		return 2
	}
	return 3
}

// timeStyle returns the lipgloss style for the clock based on timer state.
func (m Model) timeStyle() lipgloss.Style {
	if m.prog.IsLowTime(time.Duration(m.config.LowTimeWarning) * time.Second) {
		return lowTimeStyle
	} else if m.prog.IsOverflow() {
		return overflowStyle
	}
	return timerStyle
}

// styledTime returns the plain-text clock string with color applied.
func (m Model) styledTime() string {
	return m.timeStyle().Render(formatTime(m.prog.TimeDisplay()))
}

// compactStateLabel returns a state indicator for T2/T3, or "" when running.
func (m Model) compactStateLabel() string {
	switch m.AppState() {
	case Done:
		return completionStyle.Render(m.completionMsg)
	case Paused:
		return pausedStyle.Render("PAUSED")
	case Ready:
		return hintStyle.Render("READY")
	}
	return ""
}

// bigDigitHeight is the fixed row count of the big-digit font.
const bigDigitHeight = 5

func (m Model) renderTier1(mainHeight int) string {
	timeStr := formatTime(m.prog.TimeDisplay())
	rows := renderer.BigDigits(timeStr)
	content := m.timeStyle().Render(strings.Join(rows, "\n")) + "\n"

	// Budget remaining lines for labels (each costs 1 row + 1 blank separator).
	budgetLeft := mainHeight - bigDigitHeight - 2 // -2 for the leading and trailing "\n" in Place

	intervalCur, intervalTotal := m.prog.IntervalProgress()
	if intervalTotal > 0 && budgetLeft >= 2 {
		content += "\n" + labelStyle.Render(fmt.Sprintf("Interval %d/%d", intervalCur, intervalTotal))
		budgetLeft -= 2
	}

	roundCur, roundTotal := m.prog.RoundProgress()
	if roundTotal > 0 && budgetLeft >= 2 {
		content += "\n" + labelStyle.Render(fmt.Sprintf("Round %d/%d", roundCur, roundTotal))
		budgetLeft -= 2
	}

	if sw, ok := m.prog.(*stopwatch.Stopwatch); ok {
		laps := sw.Laps()
		lapNum := len(laps) + 1
		if budgetLeft >= 2 {
			content += "\n" + labelStyle.Render(fmt.Sprintf("Lap %02d", lapNum))
			budgetLeft -= 2
		}
		content += renderLapList(laps, budgetLeft)
	}

	switch m.AppState() {
	case Done:
		content += "\n\n" + completionStyle.Render(m.completionMsg)
	case Paused:
		content += "\n\n" + pausedStyle.Render("PAUSED")
	case Ready:
		content += "\n\n" + hintStyle.Render("Press space to start")
	}

	return lipgloss.Place(m.width, mainHeight, lipgloss.Center, lipgloss.Top, "\n"+content)
}

func (m Model) renderTier2(mainHeight int) string {
	lines := []string{m.styledTime()}

	intervalCur, intervalTotal := m.prog.IntervalProgress()
	if intervalTotal > 0 {
		lines = append(lines, labelStyle.Render(fmt.Sprintf("Interval %d/%d", intervalCur, intervalTotal)))
	}

	roundCur, roundTotal := m.prog.RoundProgress()
	if roundTotal > 0 {
		lines = append(lines, labelStyle.Render(fmt.Sprintf("Round %d/%d", roundCur, roundTotal)))
	}

	if sw, ok := m.prog.(*stopwatch.Stopwatch); ok {
		lapNum := len(sw.Laps()) + 1
		lines = append(lines, labelStyle.Render(fmt.Sprintf("Lap %02d", lapNum)))
	}

	if label := m.compactStateLabel(); label != "" {
		lines = append(lines, label)
	}

	content := strings.Join(lines, "\n")
	return lipgloss.Place(m.width, mainHeight, lipgloss.Left, lipgloss.Top, content)
}

// renderTier3 renders everything on a single left-aligned line.
// Includes: time, interval progress (if > 1 interval), round progress
// (if > 1 round), lap count (stopwatch), and state indicator (PAUSED/READY).
func (m Model) renderTier3(mainHeight int) string {
	elements := []string{m.styledTime()}

	intervalCur, intervalTotal := m.prog.IntervalProgress()
	if intervalTotal > 0 {
		elements = append(elements, labelStyle.Render(fmt.Sprintf("Interval %d/%d", intervalCur, intervalTotal)))
	}

	roundCur, roundTotal := m.prog.RoundProgress()
	if roundTotal > 0 {
		elements = append(elements, labelStyle.Render(fmt.Sprintf("Round %d/%d", roundCur, roundTotal)))
	}

	if sw, ok := m.prog.(*stopwatch.Stopwatch); ok {
		lapNum := len(sw.Laps()) + 1
		elements = append(elements, labelStyle.Render(fmt.Sprintf("Lap %02d", lapNum)))
	}

	if label := m.compactStateLabel(); label != "" {
		elements = append(elements, label)
	}

	content := strings.Join(elements, "  ")
	return lipgloss.Place(m.width, mainHeight, lipgloss.Left, lipgloss.Top, content)
}

// renderLapList formats completed lap splits for display.
// Shows the most recent laps that fit within budget lines,
// prepending a "..." line if older laps are hidden.
// Each entry is formatted as "  01  0:45" with a leading newline.
// Returns an empty string if there are no laps or no budget.
func renderLapList(laps []time.Duration, budget int) string {
	output := ""

	if len(laps) == 0 || budget <= 0 {
		return ""
	}

	startAt := max(0, len(laps)-budget)
	if startAt > 0 {
		startAt = max(0, len(laps)-(budget-1)) // reserve 1 line for "..."
		output += fmt.Sprintf("\n  ...")
	}

	for i := range laps[startAt:] {
		j := startAt + i
		split := laps[j]
		if j > 0 {
			split = laps[j] - laps[j-1]
		}

		mm := math.Floor(split.Minutes())
		ss := math.Floor(split.Seconds() - mm*60)

		output += fmt.Sprintf("\n  %02d  %02d:%02d", j+1, int(mm), int(ss))
	}

	return output
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
