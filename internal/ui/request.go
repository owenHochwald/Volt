package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
)

const (
	maxFocusIndex = 3

	focusMethod = iota - 1
	focusURL
	focusName
	focusSubmit
)

var (
	// HTTP Method Styles - Minimalistic, color-coded badges
	methodStyleBase = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)

	// Color-coded by HTTP method semantics (no borders for minimalism)
	getMethodStyle    = methodStyleBase.Foreground(lipgloss.Color("42"))  // Green
	postMethodStyle   = methodStyleBase.Foreground(lipgloss.Color("214")) // Orange
	putMethodStyle    = methodStyleBase.Foreground(lipgloss.Color("117")) // Blue
	patchMethodStyle  = methodStyleBase.Foreground(lipgloss.Color("141")) // Purple
	deleteMethodStyle = methodStyleBase.Foreground(lipgloss.Color("196")) // Red

	// Label styles - Very subtle
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Submit button - Minimal, just text with color change
	submitButtonFocused = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true)

	submitButtonBlurred = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))
)

type RequestPane struct {
	methods       []string
	currentMethod int
	panelFocused  bool

	focusComponentIndex int

	urlInput  textinput.Model
	nameInput textinput.Model

	request *http.Request

	// MISSING: Consider adding these fields for future features:
	// headersExpanded bool
	// headers         []HeaderPair
	// queryParams     []QueryParam
	// bodyExpanded    bool
	// bodyInput       textarea.Model
	// validationError error
}

func (m *RequestPane) syncRequest() {
	m.request.Method = m.methods[m.currentMethod]
	m.request.URL = m.urlInput.Value()
	m.request.Name = m.nameInput.Value()
}

func (m RequestPane) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RequestPane) SetFocused(focused bool) {
	m.panelFocused = focused
}

func (m *RequestPane) GetCurrentMethod() string {
	return m.methods[m.currentMethod]
}

func (m *RequestPane) blurCurrentComponent() {
	switch m.focusComponentIndex {
	case focusMethod:
		methodStyleBase.BorderForeground(unfocusColor)
	case focusURL:
		m.urlInput.Blur()
	case focusName:
		m.nameInput.Blur()
	default:
	}
}

func (m *RequestPane) focusCurrentComponent() {
	switch m.focusComponentIndex {
	case focusMethod:
		methodStyleBase.BorderForeground(focusColor)
	case focusURL:
		m.urlInput.Focus()
	case focusName:
		m.nameInput.Focus()
	default:
	}
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

			m.blurCurrentComponent()

			// focus component cycling
			if s == "up" {
				m.focusComponentIndex--
			} else {
				m.focusComponentIndex++
			}

			if m.focusComponentIndex > maxFocusIndex {
				m.focusComponentIndex = 0
			} else if m.focusComponentIndex < 0 {
				m.focusComponentIndex = maxFocusIndex
			}
			m.focusCurrentComponent()

		}

		switch m.focusComponentIndex {
		case focusMethod:
			switch msg.String() {
			case tea.KeyRight.String(), "l":
				m.currentMethod = (m.currentMethod + 1) % len(m.methods)
			case tea.KeyLeft.String(), "h":
				m.currentMethod = (m.currentMethod - 1 + len(m.methods)) % len(m.methods)
			}
		case focusURL:
			var cmd tea.Cmd
			m.urlInput, cmd = m.urlInput.Update(msg)
			return m, cmd
		case focusName:
			var cmd tea.Cmd
			m.nameInput, cmd = m.nameInput.Update(msg)
			return m, cmd
		case focusSubmit:
			switch msg.String() {
			case tea.KeyEnter.String():
				m.syncRequest()
				// TODO: Submit request and switch panels, display toast
				return m, nil
			}
		default:

			return m, nil

		}

	}

	m.syncRequest()
	return m, nil
}

func (m RequestPane) View() string {
	methodDisplay := m.methods[m.currentMethod]
	var methodStyle lipgloss.Style

	switch methodDisplay {
	case http.GET:
		methodStyle = getMethodStyle
	case http.POST:
		methodStyle = postMethodStyle
	case http.PUT:
		methodStyle = putMethodStyle
	case http.PATCH:
		methodStyle = patchMethodStyle
	case http.DELETE:
		methodStyle = deleteMethodStyle
	default:
		methodStyle = methodStyleBase
	}

	if m.focusComponentIndex == focusMethod {
		methodStyle = methodStyle.BorderForeground(focusColor)
	}
	methodRendered := methodStyle.Render(methodDisplay)
	primaryLine := lipgloss.JoinHorizontal(lipgloss.Left, methodRendered, " ", m.urlInput.View())

	nameLabel := labelStyle.Render("Name ")
	nameLine := lipgloss.JoinHorizontal(lipgloss.Left, nameLabel, m.nameInput.View())

	var button string
	if m.focusComponentIndex == focusSubmit {
		button = submitButtonFocused.Render("→ Send")
	} else {
		button = submitButtonBlurred.Render("→ Send")
	}

	helpText := HelpStyle.Render("tab/↑/↓: navigate • ←/→ or h/l: change method • enter: send • q: quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		primaryLine,
		nameLine,
		"",
		button,
		helpText,
	)
}

func SetupRequestPane() RequestPane {
	methods := []string{
		http.GET,
		http.POST,
		http.PUT,
		http.PATCH,
		http.DELETE,
	}

	m := RequestPane{
		methods:             methods,
		currentMethod:       0,
		panelFocused:        false,
		focusComponentIndex: focusMethod,
		request:             http.NewDefaultRequest(),
	}

	m.urlInput = textinput.New()
	m.urlInput.Placeholder = "http://localhost:"
	m.urlInput.CharLimit = 40
	m.urlInput.Width = 60

	m.nameInput = textinput.New()
	m.nameInput.Placeholder = "My new awesome request.."
	m.nameInput.CharLimit = 40
	m.nameInput.Width = 60

	return m
}
