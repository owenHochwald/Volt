package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
)

type ResponsePane struct {
	response *http.Response
}

func (m ResponsePane) Init() tea.Cmd {
	return nil
}

func (m *ResponsePane) SetFocused(focused bool) {
}

func (m *ResponsePane) SetHeight(height int) {
}

func (m *ResponsePane) GetCurrentMethod() string {
	return ""
}

func (m ResponsePane) View() string {
	return ""
}

func (m ResponsePane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func SetupResponsePane() *ResponsePane {
	return &ResponsePane{}
}
