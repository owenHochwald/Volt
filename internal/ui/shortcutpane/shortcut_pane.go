package shortcutpane

import tea "github.com/charmbracelet/bubbletea"

// ShortcutPane is the component responsible for displaying shortcuts
type ShortcutPane struct {
	activeTab     int
	height, width int

	Focused bool
	// Future: might need to refactor viewable area to be viewports
}

func (m ShortcutPane) Init() tea.Cmd {
	return nil
}

func (m *ShortcutPane) SetFocused(focused bool) {
	m.Focused = focused
}

func (m *ShortcutPane) SetHeight(height int) {
	m.height = height
}

func (m *ShortcutPane) SetWidth(width int) {
	m.width = width
}
