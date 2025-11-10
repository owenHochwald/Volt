package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/utils"
)

type ResponsePane struct {
	Response      *http.Response
	height, width int

	viewport viewport.Model
}

func (m ResponsePane) Init() tea.Cmd {
	return nil
}

func (m *ResponsePane) SetResponse(response *http.Response) {
	m.Response = response
	if m.Response != nil {
		m.viewport.SetContent(m.Response.Body)
	}
}

func (m *ResponsePane) SetHeight(height int) {
	m.height = height
}

func (m *ResponsePane) SetWidth(width int) {
	m.width = width
}

func (m *ResponsePane) GetCurrentMethod() string {
	return ""
}

func (m ResponsePane) View() string {
	if m.Response == nil {
		return "Make a request to see the response here!"
	}
	return m.viewport.View() + utils.MapStatusCodeToColor(m.Response.StatusCode).Render(m.Response.Status)
}

func (m ResponsePane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.viewport.SetContent(m.Response.Body)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func SetupResponsePane() ResponsePane {
	return ResponsePane{
		viewport: viewport.New(20, 10),
		width:    20,
		height:   10,
	}
}
