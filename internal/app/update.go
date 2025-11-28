package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/ui"
	"github.com/owenHochwald/volt/internal/utils"
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
		return m, ui.LoadRequestsCmd(m.db)

	case ui.RequestDeletedMsg:
		if msg.Err != nil {
			return m, nil
		}
		return m, ui.LoadRequestsCmd(m.db)

	case ui.RequestsLoadingMsg:
		if msg.Err != nil {
			return m, nil
		}
		var sidebarModel tea.Model
		sidebarModel, cmd = m.sidebarPane.Update(msg)
		m.sidebarPane = sidebarModel.(*ui.SidebarPane)
		return m, cmd

	case http.LoadTestStartMsg:
		updates := make(chan *http.LoadTestStats, 100)
		m.loadTestUpdates = updates

		// start load test in background
		go func() {
			msg.Config.Init()
			msg.Config.Run(updates)
		}()

		m.responsePane.ClearLoadTestStats()
		return m, ui.WaitForLoadTestUpdatesCmd(updates, msg.Config.TotalRequests)

	case http.LoadTestStatsMsg:
		m.responsePane.SetLoadTestStats(msg.Stats)

		if m.loadTestUpdates != nil {
			return m, ui.WaitForLoadTestUpdatesCmd(m.loadTestUpdates, msg.Stats.TotalRequests)
		}
		return m, nil

	case http.LoadTestCompleteMsg:
		// final update
		m.loadTestUpdates = nil
		if msg.Stats != nil {
			m.responsePane.SetLoadTestStats(msg.Stats)
		}
		m.requestPane.ExitLoadTestMode()
		m.focusedPanel = utils.ResponsePanel // Switch focus to results
		return m, nil

	case http.LoadTestErrorMsg:
		m.loadTestUpdates = nil
		m.requestPane.ExitLoadTestMode()
		// TODO: Display error message
		return m, nil

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
	} else if m.focusedPanel == utils.ResponsePanel {
		var responsePaneModel tea.Model
		responsePaneModel, cmd = m.responsePane.Update(msg)
		m.responsePane = responsePaneModel.(*ui.ResponsePane)
		return m, cmd
	}

	return m, cmd
}
