package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/ui"
)

func (m Model) View() string {
	sidebarWidth := 20
	contentHeight := m.height - 5

	mainWidth := m.width - sidebarWidth - 4
	mainHeight := contentHeight - 2

	requestHeight := int(float64(mainHeight) / 2.2)
	responseHeight := int(float64(mainHeight) / 2.2)

	header := ui.HeaderStyle.Width(m.width).
		Render("Volt - TUI HTTP Client - v0.1 [?] Help  [q] Quit")

	sidebar := ui.SidebarStyle.Width(sidebarWidth).
		Height(contentHeight - 2).
		Render(m.requests.View())

	request := ui.RequestStyle.Width(mainWidth - 10).
		Height(requestHeight).
		Render(m.requestView())

	response := ui.ResponseStyle.Width(mainWidth - 10).
		Height(responseHeight).
		Render("Response editor")

	rightSide := lipgloss.JoinVertical(lipgloss.Right, request, response)
	bottomPanels := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightSide)
	return lipgloss.JoinVertical(lipgloss.Top, header, bottomPanels)
}

func (m Model) requestView() string {

	if m.selectedRequest == nil {
		return "No item selected"
	}
	s := "REQUEST SECTION\n"
	s += fmt.Sprintf("HTTP Method: %s\n\n", m.selectedRequest.title)
	s += fmt.Sprintf("Description: %s\n\n", m.selectedRequest.desc)
	s += "---\n\n"
	s += "URL: [              ]\n"
	s += "Headers: [          ]\n"
	s += "Body: [             ]\n\n"
	s += "Press ESC to go back"

	return docStyle.Render(s)
}
