package app

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/owenHochwald/Volt/internal/http"
	"github.com/owenHochwald/Volt/internal/ui"
	"github.com/owenHochwald/Volt/internal/ui/keybindings"
	"github.com/owenHochwald/Volt/internal/ui/shortcutpane"
	"github.com/owenHochwald/Volt/internal/utils"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if keybindings.Matches(msg, m.keys.ForceQuit) {
			return m, tea.Quit
		}

		if m.showHelpModal {
			if keybindings.Matches(msg, m.keys.GlobalHelp) {
				m.closeHelp()
				return m, nil
			}
			m.shortcutPane, cmd = m.shortcutPane.Update(msg)
			return m, cmd
		}

		if keybindings.Matches(msg, m.keys.GlobalHelp) {
			m.openHelp(keybindings.ContextGlobal)
			return m, nil
		}

		isEditing := m.focusedPanel == utils.RequestPanel && m.requestPane.IsEditing()
		if keybindings.Matches(msg, m.keys.ContextHelp) && !isEditing {
			m.openHelp(m.focusedContext())
			return m, nil
		}
		if keybindings.Matches(msg, m.keys.PreviousPanel) {
			m.movePanel(-1)
			return m, nil
		}
		if keybindings.Matches(msg, m.keys.NextPanel) {
			m.movePanel(1)
			return m, nil
		}
		if keybindings.Matches(msg, m.keys.Quit) &&
			(m.focusedPanel == utils.SidebarPanel || m.focusedPanel == utils.ResponsePanel) {
			return m, tea.Quit
		}
		if keybindings.Matches(msg, m.keys.EscapePanel) {
			if m.focusedPanel != utils.SidebarPanel {
				m.setFocusedPanel(utils.SidebarPanel)
				return m, nil
			}
		}
		if keybindings.Matches(msg, m.keys.LoadRequest) {
			if m.focusedPanel == utils.SidebarPanel {
				if item, ok := m.sidebarPane.SelectedItem(); ok {
					m.setFocusedPanel(utils.RequestPanel)
					return m, ui.SetRequestPaneRequestCmd(item.Request)
				}
			}
		}
	case shortcutpane.CloseHelpModalMsg:
		m.closeHelp()
		return m, nil

	case http.ResultMsg:
		m.requestPane.ResultMsgCleanup()
		m.responsePane.SetResponse(msg.Response)
		m.setFocusedPanel(utils.ResponsePanel)
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
		m.sidebarPane, cmd = m.sidebarPane.Update(msg)
		return m, cmd

	case http.LoadTestStartMsg:
		updates := make(chan *http.LoadTestStats, 100)
		m.loadTestUpdates = updates

		// start load test in background
		go func() {
			msg.Config.Run(context.Background(), updates)
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
		m.setFocusedPanel(utils.ResponsePanel)
		return m, nil

	case http.LoadTestErrorMsg:
		m.loadTestUpdates = nil
		m.requestPane.ExitLoadTestMode()
		// TODO: Display error message
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.applyLayout(calculateLayout(m.width, m.height))
	}

	// Existing panel update routing (only when help modal is closed)
	if !m.showHelpModal {
		if m.focusedPanel == utils.SidebarPanel {
			m.sidebarPane, cmd = m.sidebarPane.Update(msg)
			return m, cmd
		} else if m.focusedPanel == utils.RequestPanel {
			m.requestPane, cmd = m.requestPane.Update(msg)
			return m, cmd
		} else if m.focusedPanel == utils.ResponsePanel {
			m.responsePane, cmd = m.responsePane.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
}
