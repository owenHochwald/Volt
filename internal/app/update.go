package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/ui"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyShiftTab:
			m.focusedPanel = (m.focusedPanel + 1) % 3
		default:
		}
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case tea.KeyEscape.String():
			if m.focusedPanel == RequestPanel {
				m.focusedPanel = SidebarPanel
				m.selectedRequest = nil
				return m, nil
			}
		case tea.KeyEnter.String(), " ":
			if m.focusedPanel == SidebarPanel {
				if i, ok := m.requestsList.SelectedItem().(RequestItem); ok {
					m.focusedPanel = RequestPanel
					m.selectedRequest = &i
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.httpMethods.SetSize(m.width/2, (m.height-15)/2)
		m.requestsList.SetSize(m.width/2, (m.height-15)/2)
	}
	if m.focusedPanel == SidebarPanel {
		var cmd tea.Cmd
		m.httpMethods, cmd = m.httpMethods.Update(msg)
		m.requestsList, cmd = m.requestsList.Update(msg)
		return m, cmd
	} else if m.focusedPanel == RequestPanel {
		// use request pane update
		m.requestPane.SetFocused(true)
		var requestPaneModel tea.Model
		requestPaneModel, cmd = m.requestPane.Update(msg)
		m.requestPane = requestPaneModel.(ui.RequestPane)
		return m, cmd
	}
	return m, cmd
}
