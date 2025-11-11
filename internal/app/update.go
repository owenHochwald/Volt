package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
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
				return m, nil
			}
		case tea.KeyEnter.String(), " ":
			if m.focusedPanel == SidebarPanel {
				if _, ok := m.sidebarPane.SelectedItem(); ok {
					m.focusedPanel = RequestPanel
					// TODO: Load the selected request into requestPane
				}
			}
		}
	case http.ResultMsg:
		m.requestPane.ResultMsgCleanup()
		m.responsePane.SetResponse(msg.Response)
		m.focusedPanel = ResponsePanel
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.sidebarPane.SetSize(m.width/2, (m.height-15)/2)
	}
	if m.focusedPanel == SidebarPanel {
		cmd := m.sidebarPane.Update(msg)
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
