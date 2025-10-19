package app

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.requests.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.requests, cmd = m.requests.Update(msg)
	return m, cmd
}
