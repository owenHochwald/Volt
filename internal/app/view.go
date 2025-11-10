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

	requestHeight := int(float64(mainHeight)/2.2) - 9
	responseHeight := int(float64(mainHeight) / 2.2)

	header := ui.HeaderStyle.Width(m.width).
		Render("Volt - TUI HTTP Client - v0.1 [?] Help  [q] Quit")

	sidebar := m.sidebarView(mainHeight, sidebarWidth)

	request := ui.ApplyFocus(ui.RequestStyle, m.focusedPanel == 1).Width(mainWidth - 10).
		Height(requestHeight).
		Render(m.requestView(requestHeight))

	response := ui.ApplyFocus(ui.ResponseStyle, m.focusedPanel == 2).Width(mainWidth - 10).
		Height(responseHeight).
		Render(m.responseView(responseHeight, mainWidth-10))

	rightSide := lipgloss.JoinVertical(lipgloss.Right, request, response)
	bottomPanels := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightSide)
	return lipgloss.JoinVertical(lipgloss.Top, header, bottomPanels)
}

func (m Model) sidebarView(height, width int) string {
	sidebar := ui.ApplyFocus(ui.SidebarStyle, m.focusedPanel == 0).Width(width).
		Height(height - 4).
		Render(m.requestsList.View())
	return sidebar
}

func (m Model) requestView(height int) string {
	m.requestPane.SetFocused(m.focusedPanel == RequestPanel)
	m.requestPane.SetHeight(height)

	return m.requestPane.View()
}

func (m Model) responseView(height, width int) string {
	m.responsePane.SetFocused(m.focusedPanel == ResponsePanel)
	m.responsePane.SetHeight(height)
	m.responsePane.SetWidth(width)

	return m.responsePane.View()
}
