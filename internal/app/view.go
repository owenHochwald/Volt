package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/ui"
)

func (m Model) View() string {
	sidebarWidth := 20
	mainWidth := m.width - sidebarWidth - 4

	sidebar := ui.SidebarStyle.Width(sidebarWidth).
		Height(m.height - 2).
		Render(m.requests.View())

	main := ui.MainStyle.
		Width(mainWidth).
		Height(m.height - 2).
		Render(m.detailView())
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
}

func (m Model) detailView() string {
	if m.selectedRequest == nil {
		return "No item selected"
	}

	s := fmt.Sprintf("HTTP Method: %s\n\n", m.selectedRequest.title)
	s += fmt.Sprintf("Description: %s\n\n", m.selectedRequest.desc)
	s += "---\n\n"
	s += "URL: [              ]\n"
	s += "Headers: [          ]\n"
	s += "Body: [             ]\n\n"
	s += "Press ESC to go back"

	return docStyle.Render(s)
}
