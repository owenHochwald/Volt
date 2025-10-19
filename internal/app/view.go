package app

import "fmt"

func (m Model) View() string {
	switch m.currentScreen {
	case ScreenList:
		return docStyle.Render(m.requests.View())
	case ScreenDetail:
		return m.detailView()
	default:
		return docStyle.Render(m.requests.View())
	}
}

func (m Model) detailView() string {
	if m.selectedRequest == nil {
		return "No item selected"
	}

	s := fmt.Sprintf("HTTP Method: %s\n\n", m.selectedRequest.title)
	s += fmt.Sprintf("Description: %s\n\n", m.selectedRequest.desc)
	s += "---\n\n"
	s += "URL: [              ]\n"
	s += "Headers: [          ]\n"
	s += "Body: [             ]\n\n"
	s += "Press ESC to go back"

	return docStyle.Render(s)
}
