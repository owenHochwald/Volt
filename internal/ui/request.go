package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var methodStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type RequestPane struct {
	methods       []string
	currentMethod int
	focused       bool
}

func (m RequestPane) Init() tea.Cmd {
	return nil
}

func (m *RequestPane) SetFocused(focused bool) {
	m.focused = focused
}

func (m RequestPane) GetCurrentMethod() string {
	return m.methods[m.currentMethod]
}

func (m RequestPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "right", "l":
			m.currentMethod = (m.currentMethod + 1) % len(m.methods)
		case "left", "h":
			m.currentMethod = (m.currentMethod - 1 + len(m.methods)) % len(m.methods)
		}
	}

	return m, nil
}

func (m RequestPane) View() string {
	methodDisplay := m.methods[m.currentMethod]

	style := methodStyle
	if m.focused {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	return docStyle.Render("HTTP Method: \n" + style.Render(methodDisplay))
}

func SetupRequestPane() RequestPane {
	methods := []string{
		http.GET,
		http.POST,
		http.DELETE,
		http.PUT,
	}

	return RequestPane{
		methods:       methods,
		currentMethod: 0,
		focused:       false,
	}
}
