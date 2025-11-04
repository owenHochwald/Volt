package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
)

var (
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
var methodStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type RequestPane struct {
	methods       []string
	currentMethod int
	panelFocused  bool

	// text input (name, url)
	focusComponentIndex int
	inputs              []textinput.Model

	request http.Request
}

func (m RequestPane) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RequestPane) SetFocused(focused bool) {
	m.panelFocused = focused
}

func (m RequestPane) GetCurrentMethod() string {
	return m.methods[m.currentMethod]
}

func (m RequestPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.panelFocused {
			return m, nil
		}
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case tea.KeyTab.String(), tea.KeyUp.String(), tea.KeyDown.String():
			s := msg.String()
			// submit button
			if s == "enter" && m.focusComponentIndex == len(m.inputs) {
				return m, tea.Quit
			}
			// focus component cycling
			if s == "up" {
				m.focusComponentIndex--
			} else {
				m.focusComponentIndex++
			}

			if m.focusComponentIndex > len(m.inputs) {
				m.focusComponentIndex = 0
			} else if m.focusComponentIndex < 0 {
				m.focusComponentIndex = len(m.inputs)
			}

		}
		switch m.focusComponentIndex {
		// method selection
		case 0:
			switch msg.String() {
			case tea.KeyRight.String(), "l":
				m.currentMethod = (m.currentMethod + 1) % len(m.methods)
			case tea.KeyLeft.String(), "h":
				m.currentMethod = (m.currentMethod - 1 + len(m.methods)) % len(m.methods)
			}
		default:
			cmd := m.updateInputs(msg)
			return m, cmd

		}

	}

	return m, nil
}

func (m RequestPane) View() string {
	methodDisplay := m.methods[m.currentMethod]

	style := methodStyle
	if m.panelFocused {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}
	s := docStyle.Render("HTTP Method: \n")
	s += style.Render(methodDisplay) + "\n"
	s += b.String() + "\n"

	return s
}

func (m *RequestPane) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
func SetupRequestPane() RequestPane {
	methods := []string{
		http.GET,
		http.POST,
		http.DELETE,
		http.PUT,
	}

	m := RequestPane{
		methods:             methods,
		currentMethod:       0,
		panelFocused:        false,
		focusComponentIndex: 0,
		inputs:              make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 100

		switch i {
		case 0:
			t.Placeholder = "URL"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Request Name"
			t.CharLimit = 300
		}

		m.inputs[i] = t
	}

	return m
}
