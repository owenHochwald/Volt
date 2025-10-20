package ui

import "github.com/charmbracelet/lipgloss"

var (
	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(20).
			BorderForeground(lipgloss.Color("62"))
	MainStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63"))
)
