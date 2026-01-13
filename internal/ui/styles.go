package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/stef16robbe/stamp/internal/adr"
)

var (
	// Colors
	Green   = lipgloss.Color("34")
	Yellow  = lipgloss.Color("214")
	Orange  = lipgloss.Color("208")
	Red     = lipgloss.Color("196")
	Purple  = lipgloss.Color("99")
	Gray    = lipgloss.Color("241")
	White   = lipgloss.Color("255")
	Cyan    = lipgloss.Color("87")
	Magenta = lipgloss.Color("213")

	// Feedback styles
	SuccessStyle = lipgloss.NewStyle().Foreground(Green)
	WarningStyle = lipgloss.NewStyle().Foreground(Yellow)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Red)
	MutedStyle   = lipgloss.NewStyle().Foreground(Gray)
	BoldStyle    = lipgloss.NewStyle().Bold(true)

	// Symbols
	SuccessIcon = SuccessStyle.Render("✓")
	WarningIcon = WarningStyle.Render("⚠")
	ErrorIcon   = ErrorStyle.Render("✗")
	ArrowIcon   = MutedStyle.Render("→")

	// Status badge styles (with background)
	StatusStyles = map[adr.Status]lipgloss.Style{
		adr.StatusDraft:      lipgloss.NewStyle().Background(Gray).Foreground(White).Padding(0, 1),
		adr.StatusProposed:   lipgloss.NewStyle().Background(Yellow).Foreground(lipgloss.Color("0")).Padding(0, 1),
		adr.StatusAccepted:   lipgloss.NewStyle().Background(Green).Foreground(White).Padding(0, 1),
		adr.StatusDeprecated: lipgloss.NewStyle().Background(Orange).Foreground(White).Padding(0, 1),
		adr.StatusSuperseded: lipgloss.NewStyle().Background(Purple).Foreground(White).Padding(0, 1),
		adr.StatusRejected:   lipgloss.NewStyle().Background(Red).Foreground(White).Padding(0, 1),
	}

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(Cyan).Padding(0, 1)
	TableCellStyle   = lipgloss.NewStyle().Padding(0, 1)
	TableBorder      = lipgloss.RoundedBorder()
)

// RenderStatus renders a status as a colored badge
func RenderStatus(status adr.Status) string {
	style, ok := StatusStyles[status]
	if !ok {
		style = lipgloss.NewStyle().Padding(0, 1)
	}
	return style.Render(string(status))
}

// RenderStatusTransition renders a status change with colored arrow
func RenderStatusTransition(from, to adr.Status) string {
	fromStyle := StatusStyles[from]
	toStyle := StatusStyles[to]
	arrow := lipgloss.NewStyle().Foreground(Magenta).Render(" → ")
	return fromStyle.Render(string(from)) + arrow + toStyle.Render(string(to))
}

// Success prints a success message with checkmark
func Success(message string) string {
	return SuccessIcon + " " + message
}

// Warning prints a warning message
func Warning(message string) string {
	return WarningIcon + " " + WarningStyle.Render(message)
}

// Error prints an error message
func Error(message string) string {
	return ErrorIcon + " " + ErrorStyle.Render(message)
}

// Muted renders text in gray
func Muted(text string) string {
	return MutedStyle.Render(text)
}

// Bold renders text in bold
func Bold(text string) string {
	return BoldStyle.Render(text)
}
