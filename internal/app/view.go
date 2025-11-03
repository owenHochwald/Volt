package app

import (
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

	sidebar := ui.ApplyFocus(ui.SidebarStyle, m.focusedPanel == 0).Width(sidebarWidth).
		Height(mainHeight - 4).
		Render(m.requestsList.View())

	request := ui.ApplyFocus(ui.RequestStyle, m.focusedPanel == 1).Width(mainWidth - 10).
		Height(requestHeight).
		Render(m.requestView())

	response := ui.ApplyFocus(ui.ResponseStyle, m.focusedPanel == 2).Width(mainWidth - 10).
		Height(responseHeight).
		Render("Response editor")

	rightSide := lipgloss.JoinVertical(lipgloss.Right, request, response)
	bottomPanels := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightSide)
	return lipgloss.JoinVertical(lipgloss.Top, header, bottomPanels)
}

func (m Model) requestView() string {
	m.requestPane.SetFocused(m.focusedPanel == RequestPanel)

	return m.requestPane.View()
}
