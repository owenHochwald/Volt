package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	progress progress.Model
	width    int
}

func (h *Header) Init() tea.Cmd {
	return nil
}

func (h *Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width - 4
		h.progress.Width = h.width
	}
	return h, nil
}

func (h *Header) View() string {
	voltStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		PaddingRight(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	title := voltStyle.Render("VOLT ⚡")
	help := helpStyle.Render("│ v0.1 • [?] Help • [q] Quit")
	headerText := lipgloss.JoinHorizontal(lipgloss.Left, title, help)

	progressBar := h.progress.ViewAs(1.0)

	return lipgloss.JoinVertical(lipgloss.Left, headerText, progressBar)
}

func SetupHeader() *Header {
	return &Header{
		width: 80,
		progress: progress.New(
			progress.WithScaledGradient("#4C1D95", "#7C3AED"), // Dark purple → lighter purple
			progress.WithoutPercentage()),
	}
}
