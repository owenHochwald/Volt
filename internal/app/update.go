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
			if m.focusedPanel == RequestPanel {
				m.focusedPanel = SidebarPanel
				m.selectedRequest = nil
				return m, nil
			}
		case "tab":
			m.focusedPanel = (m.focusedPanel + 1) % 3
		case "enter", " ":
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
	}
	return m, nil
}
