package ui

import "github.com/charmbracelet/lipgloss"

var (
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
