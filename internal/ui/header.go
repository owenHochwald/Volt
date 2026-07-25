package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Header struct {
	width int
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
	asciiArt := `██╗   ██╗ ██████╗ ██╗  ████████╗
██║   ██║██╔═══██╗██║  ╚══██╔══╝
██║   ██║██║   ██║██║     ██║
╚██╗ ██╔╝██║   ██║██║     ██║
 ╚████╔╝ ╚██████╔╝███████╗██║
  ╚═══╝   ╚═════╝ ╚══════╝╚═╝   `

	logo := HeaderLogoStyle.Render(asciiArt)

	help := HeaderHelpStyle.Render("⚡ v0.1 • [ ? ] Help • [ ctrl+c ] Quit")

	return lipgloss.JoinHorizontal(lipgloss.Left, logo, "\t", help)
}

func SetupHeader() *Header {
	return &Header{
		width: 80,
	}
}
