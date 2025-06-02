package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const statGap = 5
const spacer = ""

// View renders the current state of the Model as a string for display
func (m Model) View() string {
	if m.showResults {
		return m.renderResults()
	}

	var sections []string

	timer := m.renderTimer()
	sections = append(sections, timer)

	textDisplay := m.renderText()
	sections = append(sections, textDisplay)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// renderTimer formats the remaining time for display
func (m Model) renderTimer() string {
	remaining := m.game.GetRemainingTime()
	return timeStyle.Render(fmt.Sprintf("%d", remaining))
}

// renderText formats the text display with appropriate styles for typed, current, untyped characters
func (m Model) renderText() string {
	displayText := m.game.GetDisplayText()

	var rendered strings.Builder

	for i, char := range displayText {
		// Use helper to style character
		styledChar := m.styleChar(char, i)
		rendered.WriteString(styledChar)
	}

	// Format into lines
	content := rendered.String()
	lines := m.formatIntoLines(content)

	return textBoxStyle.Render(strings.Join(lines, "\n"))
}

// formatIntoLines formats the content into lines based on the game's display settings
func (m Model) formatIntoLines(plainContent string) []string {
	lines := m.game.DisplayLines

	maxLines := m.game.LinesPerView
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	var styledLines []string
	charIndex := 0

	for i, line := range lines {
		if i >= maxLines {
			break
		}

		if charIndex >= len(plainContent) {
			// If we run out of content, just render untyped
			styledLines = append(styledLines, mutedStyle.Render(line))
			continue
		}

		var styledLine strings.Builder

		for _, char := range line {
			if charIndex < len(plainContent) {
				// Style character using helper
				styledChar := m.styleChar(char, charIndex)
				styledLine.WriteString(styledChar)
				charIndex++
			} else {
				styledLine.WriteString(mutedStyle.Render(string(char)))
			}
		}

		styledLines = append(styledLines, styledLine.String())

		if charIndex < len(plainContent) && i < len(lines)-1 {
			charIndex++
		}
	}

	return styledLines
}

// styleChar determines the style of a character based on its position and error status
func (m Model) styleChar(char rune, index int) string {
	userPos := m.game.CurrentPos
	errorIndex := m.game.GlobalPos - (userPos - index)

	switch {
	case index < userPos:
		// Already typed
		if m.game.Errors != nil {
			if _, hasErr := m.game.Errors[errorIndex]; hasErr {
				return errorStyle.Render(string(char))
			}
		}
		return boldStyle.Render(string(char))
	case index == userPos:
		// Current character
		return cursorStyle.Render(string(char))
	default:
		// Not yet typed
		return mutedStyle.Render(string(char))
	}
}

// renderResults formats the final results of the typing test for display
func (m Model) renderResults() string {
	stats := m.finalStats

	accSection := lipgloss.JoinVertical(
		lipgloss.Right,
		mutedStyle.Render("acc"),
		boldStyle.Render(fmt.Sprintf("%.0f%%", stats.Accuracy)),
	)

	bestAccSection := lipgloss.JoinVertical(
		lipgloss.Center,
		mutedStyle.Render("best acc"),
		boldStyle.Render(fmt.Sprintf("%.0f%%", stats.BestAccuracy)),
	)

	wpmSection := lipgloss.JoinVertical(
		lipgloss.Right,
		mutedStyle.Render("wpm"),
		boldStyle.Render(fmt.Sprintf("%.0f", stats.WPM)),
	)

	bestWpmSection := lipgloss.JoinVertical(
		lipgloss.Center,
		mutedStyle.Render("best wpm"),
		boldStyle.Render(fmt.Sprintf("%.0f", stats.BestWPM)),
	)

	timeSection := lipgloss.JoinVertical(
		lipgloss.Right,
		mutedStyle.Render("time"),
		boldStyle.Render(fmt.Sprintf("%.0fs", stats.TimeElapsed.Seconds())),
	)

	languageSection := lipgloss.JoinVertical(
		lipgloss.Right,
		mutedStyle.Render("lang"),
		boldStyle.Render(m.language),
	)

	// Arrange stats horizontally
	statsRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		accSection,
		strings.Repeat(" ", statGap),
		bestAccSection,
		strings.Repeat(" ", statGap),
		wpmSection,
		strings.Repeat(" ", statGap),
		bestWpmSection,
		strings.Repeat(" ", statGap),
		timeSection,
		strings.Repeat(" ", statGap),
		languageSection,
	)

	instructions := mutedStyle.Align(lipgloss.Center).Render("Press Enter to restart • Esc to quit")

	// Results layout
	resultsContent := lipgloss.JoinVertical(
		lipgloss.Center,
		spacer,
		statsRow,
		spacer,
		instructions,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		resultsContainerStyle.Render(resultsContent),
	)
}
