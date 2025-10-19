package app

func (m Model) View() string {
	return docStyle.Render(m.requests.View())
}
