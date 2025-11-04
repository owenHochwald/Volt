package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	focusColor   = lipgloss.Color("205")
	unfocusColor = lipgloss.Color("240")

	FocusedStyle   = lipgloss.NewStyle().Foreground(focusColor)
	UnfocusedStyle = lipgloss.NewStyle().Foreground(unfocusColor)

	FocusedButton   = FocusedStyle.Render("[ Submit ]")
	UnfocusedButton = fmt.Sprintf("[ %s ]", UnfocusedStyle.Render("Submit"))

	DocStyle = lipgloss.NewStyle().Margin(1, 2)

	// TODO: add parent base style
	HelpStyle = lipgloss.NewStyle().
			Margin(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Height(2)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(20)

	RequestStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	ResponseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())
)

func ApplyFocus(style lipgloss.Style, focus bool) lipgloss.Style {
	if focus {
		return style.
			BorderForeground(focusColor).
			Bold(true)
	}
	return style.BorderForeground(unfocusColor)
}
