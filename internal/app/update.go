package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/ui"
	"github.com/owenHochwald/volt/internal/utils"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	var headerModel tea.Model
	headerModel, cmd = m.headerPane.Update(msg)
	m.headerPane = headerModel.(*ui.Header)

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
			if m.focusedPanel == utils.RequestPanel {
				m.focusedPanel = utils.SidebarPanel
				return m, nil
			}
		case tea.KeyEnter.String(), " ":
			if m.focusedPanel == utils.SidebarPanel {
				if item, ok := m.sidebarPane.SelectedItem(); ok {
					m.focusedPanel = utils.RequestPanel
					return m, ui.SetRequestPaneRequestCmd(item.Request)
				}
			}
		}
	case http.ResultMsg:
		m.requestPane.ResultMsgCleanup()
		m.responsePane.SetResponse(msg.Response)
		m.focusedPanel = utils.ResponsePanel
		return m, nil

	case ui.RequestSavedMsg:
		if msg.Err != nil {
			return m, nil
		}
		return m, tea.Batch(
			ui.LoadRequestsCmd(m.db),
		)
	case ui.RequestDeletedMsg:
		if msg.Err != nil {
			return m, nil
		}
		return m, tea.Batch(
			ui.LoadRequestsCmd(m.db),
		)
	case ui.RequestsLoadingMsg:
		if msg.Err != nil {
			return m, nil
		}
		var sidebarModel tea.Model
		sidebarModel, cmd = m.sidebarPane.Update(msg)
		m.sidebarPane = sidebarModel.(*ui.SidebarPane)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.sidebarPane.SetSize(m.width/2, (m.height-15)/2)
	}
	if m.focusedPanel == utils.SidebarPanel {
		var sidebarPaneModel tea.Model
		sidebarPaneModel, cmd = m.sidebarPane.Update(msg)
		m.sidebarPane = sidebarPaneModel.(*ui.SidebarPane)
		return m, cmd
	} else if m.focusedPanel == utils.RequestPanel {
		m.requestPane.SetFocused(true)
		var requestPaneModel tea.Model
		requestPaneModel, cmd = m.requestPane.Update(msg)
		m.requestPane = requestPaneModel.(ui.RequestPane)
		return m, cmd
	}
	return m, cmd
}
