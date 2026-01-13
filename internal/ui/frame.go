package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Frame wraps content in a styled border with a title
func Frame(title string, content string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true)

	frameStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Gray).
		Padding(1, 2)

	header := titleStyle.Render(fmt.Sprintf("─ %s ", title))

	framedContent := frameStyle.Render(content)

	// Insert title into top border
	lines := []rune(framedContent)
	if len(lines) > 3 {
		// Find position after the corner
		headerRunes := []rune(header)
		insertPos := 1 // After the corner character
		for i, r := range headerRunes {
			if insertPos+i < len(lines) {
				lines[insertPos+i] = r
			}
		}
	}

	return string(lines)
}

// SimpleFrame wraps content in a simple border
func SimpleFrame(content string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Gray).
		Padding(0, 1).
		Render(content)
}
