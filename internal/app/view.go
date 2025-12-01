package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/ui"
	"github.com/owenHochwald/volt/internal/utils"
)

func (m Model) View() string {
	sidebarWidth := 20
	contentHeight := m.height - 5

	mainWidth := m.width - sidebarWidth - 4
	mainHeight := contentHeight - 2

	requestHeight := int(float64(mainHeight)/2.2) - 10
	responseHeight := int(float64(mainHeight)/2.2) - 2

	sidebar := m.sidebarView(mainHeight, sidebarWidth)

	// Conditional rendering for custom request pane border color
	var request string
	if m.requestPane.LoadTestMode {
		// Load test style already has yellow border, background, and bold - don't apply focus
		request = ui.LoadTestBorderStyle.Width(mainWidth - 10).
			Height(requestHeight).
			Render(m.requestView(requestHeight))
	} else {
		request = ui.ApplyFocus(ui.RequestStyle, m.focusedPanel == 1).Width(mainWidth - 10).
			Height(requestHeight).
			Render(m.requestView(requestHeight))
	}

	response := ui.ApplyFocus(ui.ResponseStyle, m.focusedPanel == 2).Width(mainWidth - 10).
		Height(responseHeight).
		Render(m.responseView(responseHeight, mainWidth-10))

	rightSide := lipgloss.JoinVertical(lipgloss.Right, request, response)
	bottomPanels := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rightSide)
	return lipgloss.JoinVertical(lipgloss.Top, m.headerView(m.width), bottomPanels)
}

func (m Model) headerView(width int) string {
	header := ui.HeaderStyle.Width(width).Render(m.headerPane.View())
	return header
}

func (m Model) sidebarView(height, width int) string {
	sidebar := ui.ApplyFocus(ui.SidebarStyle, m.focusedPanel == 0).Width(width).
		Height(height - 4).
		Render(m.sidebarPane.View())
	return sidebar
}

func (m Model) requestView(height int) string {
	m.requestPane.SetFocused(m.focusedPanel == utils.RequestPanel)
	//m.requestPane.LoadTestMode
	m.requestPane.SetHeight(height)

	return m.requestPane.View()
}

func (m Model) responseView(height, width int) string {
	m.responsePane.SetHeight(height)
	m.responsePane.SetWidth(width)

	return m.responsePane.View()
}
