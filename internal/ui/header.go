package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Header struct {
	width   int
	compact bool
}

func (h *Header) Init() tea.Cmd {
	return nil
}

func (h *Header) Update(msg tea.Msg) (*Header, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
	}
	return h, nil
}

func (h *Header) View() string {
	if h.compact {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			HeaderLogoStyle.Render("⚡ VOLT"),
			"  ",
			HeaderHelpStyle.Render("F1 help • alt+h/alt+l panels • ctrl+c quit"),
		)
	}

	asciiArt := `██╗   ██╗ ██████╗ ██╗  ████████╗
██║   ██║██╔═══██╗██║  ╚══██╔══╝
██║   ██║██║   ██║██║     ██║
╚██╗ ██╔╝██║   ██║██║     ██║
 ╚████╔╝ ╚██████╔╝███████╗██║
  ╚═══╝   ╚═════╝ ╚══════╝╚═╝   `

	logo := HeaderLogoStyle.Render(asciiArt)

	help := HeaderHelpStyle.Render("v0.1 • F1 help • alt+h/alt+l panels • ctrl+c quit")

	return lipgloss.JoinHorizontal(lipgloss.Left, logo, "\t", help)
}

func (h *Header) SetSize(width int, compact bool) {
	h.width = max(width, 0)
	h.compact = compact
}

func SetupHeader() *Header {
	return &Header{
		width: 80,
	}
}
