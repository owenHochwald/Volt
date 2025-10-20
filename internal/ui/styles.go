package ui

import "github.com/charmbracelet/lipgloss"

var (
	focusColor   = lipgloss.Color("205")
	unfocusColor = lipgloss.Color("240")

	// TODO: add parent base style
	HeaderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Height(2).
			BorderForeground(lipgloss.Color("61"))

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(20).
			BorderForeground(lipgloss.Color("62"))

	RequestStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("64"))

	ResponseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("65"))
)

func ApplyFocus(style lipgloss.Style, focus bool) lipgloss.Style {
	if focus {
		return style.
			BorderForeground(focusColor).
			Bold(true)
	}
	return style.BorderForeground(unfocusColor)
}
