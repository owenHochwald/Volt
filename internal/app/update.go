package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.currentScreen == ScreenDetail {
				m.currentScreen = ScreenList
				m.selectedRequest = nil
				return m, nil
			}
		case "enter", " ":
			if m.currentScreen == ScreenList {
				if i, ok := m.requests.SelectedItem().(Item); ok {
					m.currentScreen = ScreenDetail
					m.selectedRequest = &i
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.requests.SetSize(msg.Width-h, msg.Height-v)
	}
	if m.currentScreen == ScreenList {
		var cmd tea.Cmd
		m.requests, cmd = m.requests.Update(msg)
		return m, cmd
	}
	return m, nil
}
