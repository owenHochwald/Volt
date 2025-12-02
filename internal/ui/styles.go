package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	focusColor   = lipgloss.Color("205")
	unfocusColor = lipgloss.Color("240")

	darkPurple = lipgloss.Color("#4C1D95")
	deepViolet = lipgloss.Color("#5B21B6")
	dimGray    = lipgloss.Color("240")
	keyStyle   = lipgloss.NewStyle().Foreground(darkPurple).Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("160")).
			Background(lipgloss.Color("52"))

	FocusedStyle   = lipgloss.NewStyle().Foreground(focusColor)
	UnfocusedStyle = lipgloss.NewStyle().Foreground(unfocusColor)

	InactiveTab = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("240")) // dimGray

	ActiveTab = lipgloss.NewStyle().
			Padding(0, 2).
			Background(lipgloss.Color("98")). // darkPurple
			Foreground(lipgloss.Color("255")).
			Bold(true)

	FocusedButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	UnfocusedButton = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	DocStyle = lipgloss.NewStyle().Margin(1, 2)

	// TODO: add parent base style
	HelpStyle = lipgloss.NewStyle().
			Margin(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Height(2)

	// Header ASCII art styles
	HeaderLogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(darkPurple)

	HeaderHelpStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(20)

	RequestStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	ResponseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	LoadTestBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("226")) // Yellow

	ModalOverlayStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("0")). // Black background
				Foreground(lipgloss.Color("255"))

	ModalContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(focusColor).
				Background(lipgloss.Color("235")). // Dark gray background
				Padding(1, 2).
				Bold(true)
)

func ApplyFocus(style lipgloss.Style, focus bool) lipgloss.Style {
	if focus {
		return style.
			BorderForeground(focusColor).
			Bold(true)
	}
	return style.BorderForeground(unfocusColor)
}
