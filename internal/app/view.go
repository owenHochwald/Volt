package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/ui"
)

func (m Model) View() string {
	sidebarWidth := 20
	mainWidth := m.width - sidebarWidth - 4
	mainHeight := m.height - 2

	requestHeight := int(float64(mainHeight) / 2.2)
	responseHeight := int(float64(mainHeight) / 2.2)

	sidebar := ui.SidebarStyle.Width(sidebarWidth).
		Height(m.height - 2).
		Render(m.requests.View())

	request := ui.RequestStyle.Width(mainWidth - 10).
		Height(requestHeight).
		Render("Request editor")

	response := ui.ResponseStyle.Width(mainWidth - 10).
		Height(responseHeight).
		Render("Response editor")

	//main := ui.MainStyle.
	//	Width(mainWidth).
	//	Height(m.height - 2).
	//	Render(m.detailView())
	rightSide := lipgloss.JoinVertical(lipgloss.Right, request, response)
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightSide)
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
