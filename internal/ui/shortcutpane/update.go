package shortcutpane

import (
	tea "github.com/charmbracelet/bubbletea"
)

// CloseHelpModalMsg signals the app to close the help modal
type CloseHelpModalMsg struct{}

func (m ShortcutPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Direct tab access
		case "1":
			m.activeTab = int(Global)
		case "2":
			m.activeTab = int(Sidebar)
		case "3":
			m.activeTab = int(Request)
		case "4":
			m.activeTab = int(Response)

		// Tab navigation (vim style)
		case "h", tea.KeyLeft.String(), tea.KeyShiftTab.String():
			m.activeTab = (m.activeTab - 1 + m.getMaxTabs()) % m.getMaxTabs()
		case "l", tea.KeyRight.String(), tea.KeyTab.String():
			m.activeTab = (m.activeTab + 1) % m.getMaxTabs()

		// Close modal
		case "q", "?", tea.KeyEscape.String():
			return m, func() tea.Msg {
				return CloseHelpModalMsg{}
			}
		}
	}

	return m, nil
}
