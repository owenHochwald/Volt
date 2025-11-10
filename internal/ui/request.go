package ui

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/utils"
)

const (
	maxFocusIndex = 5

	focusMethod = iota - 1
	focusURL
	focusName
	focusHeaders
	focusBody
	focusSubmit
)

var (
	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
)

type RequestPane struct {
	client *http.Client

	stopwatch stopwatch.Model
	quitting  bool

	methods       []string
	currentMethod int
	panelFocused  bool
	methodSelector MethodSelector

	focusComponentIndex int

	urlInput  textinput.Model
	nameInput textinput.Model

	request *http.Request

	height int

	parseErrors []string

	headersExpanded bool
	headers         textarea.Model
	bodyExpanded    bool
	body            textarea.Model

	requestInProgress bool
}

func (m *RequestPane) syncRequest() {
	m.request.Method = m.methodSelector.Current()
	m.request.URL = m.urlInput.Value()
	m.request.Name = m.nameInput.Value()
	headerMap, headerErrors := utils.ParseKeyValuePairs(m.headers.Value())
	bodyMap, bodyErrors := utils.ParseKeyValuePairs(m.body.Value())
	jsonData, err := json.Marshal(bodyMap)
	if err != nil {
		// TODO: add standard error handling logic
		m.parseErrors = append(m.parseErrors, "JSON marshal error: "+err.Error())
		return
	}
	m.request.Headers = headerMap
	m.request.Body = string(jsonData)
	m.parseErrors = append(headerErrors, bodyErrors...)
}

func sendRequestCmd(client *http.Client, request *http.Request) tea.Cmd {
	return func() tea.Msg {
		res := make(chan *http.Response)
		go client.Send(request, res)

		responseObject := <-res

		return http.ResultMsg{
			Response: responseObject,
		}
	}
}

func (m RequestPane) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RequestPane) SetFocused(focused bool) {
	m.panelFocused = focused
}

func (m *RequestPane) SetHeight(height int) {
	m.height = height
}

func (m *RequestPane) GetCurrentMethod() string {
	return m.methodSelector.Current()
}

func (m *RequestPane) blurCurrentComponent() {
	switch m.focusComponentIndex {
	case focusMethod:
		methodStyleBase.BorderForeground(unfocusColor)
	case focusURL:
		m.urlInput.Blur()
	case focusName:
		m.nameInput.Blur()
	case focusHeaders:
		m.headers.Blur()
	case focusBody:
		m.body.Blur()
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
	case focusHeaders:
		m.headers.Focus()
	case focusBody:
		m.body.Focus()
	default:
	}
}

func (m *RequestPane) ResultMsgCleanup() {
	m.stopwatch.Stop()
	m.stopwatch = stopwatch.NewWithInterval(10 * time.Millisecond)
	m.requestInProgress = false
}

func (m RequestPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.stopwatch, cmd = m.stopwatch.Update(msg)

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
				m.methodSelector.Next()
			case tea.KeyLeft.String(), "h":
				m.methodSelector.Prev()
			}
		case focusURL:
			var cmd tea.Cmd
			m.urlInput, cmd = m.urlInput.Update(msg)
			return m, cmd
		case focusName:
			var cmd tea.Cmd
			m.nameInput, cmd = m.nameInput.Update(msg)
			return m, cmd
		case focusHeaders:
			var cmd tea.Cmd
			m.headers, cmd = m.headers.Update(msg)
			return m, cmd
		case focusBody:
			var cmd tea.Cmd
			m.headers, cmd = m.headers.Update(msg)
			return m, cmd
		case focusSubmit:
			switch msg.String() {
			case tea.KeyEnter.String():
				if m.requestInProgress {
					return m, nil
				}

				m.syncRequest()
				m.requestInProgress = true

				m.stopwatch.Reset()
				stopwatchCmd := m.stopwatch.Start()

				return m, tea.Batch(stopwatchCmd, sendRequestCmd(m.client, m.request))
			}
		default:
			return m, nil
		}
	}

	m.syncRequest()
	return m, cmd
}

func (m RequestPane) View() string {
	methodRendered := m.methodSelector.GetStyle().Render(m.methodSelector.Current())
	primaryLine := lipgloss.JoinHorizontal(lipgloss.Left, methodRendered, " ", m.urlInput.View())

	nameLabel := labelStyle.Render("Name ")
	nameLine := lipgloss.JoinHorizontal(lipgloss.Left, nameLabel, m.nameInput.View())

	headersLabel := labelStyle.Render("Headers ")
	headersLine := lipgloss.JoinHorizontal(lipgloss.Left, headersLabel, m.headers.View())

	bodyLabel := labelStyle.Render("Body    ")
	bodyLine := lipgloss.JoinHorizontal(lipgloss.Left, bodyLabel, m.body.View())

	var button string
	var stopwatchCount string
	if m.requestInProgress {
		button = FocusedButton.Render("Sending...")
		elapsed := m.stopwatch.Elapsed()
		milliseconds := elapsed.Milliseconds()
		seconds := float64(milliseconds) / 1000.0
		stopwatchCount = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render(fmt.Sprintf("%.3fs", seconds))
	} else if m.focusComponentIndex == focusSubmit {
		button = FocusedButton.Render("→ Send")
	} else {
		button = UnfocusedButton.Render("→ Send")
	}

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		primaryLine,
		nameLine,
		headersLine,
		bodyLine,
		"",
		button,
	)

	helpText := HelpStyle.Render("tab/↑/↓: navigate • ←/→ or h/l: change method • enter: send • q: quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		lipgloss.NewStyle().Height(m.height-10).Render(""),
		stopwatchCount,
		helpText,
	)

}

func SetupRequestPane() RequestPane {
	m := RequestPane{
		client:              http.InitClient(0, false),
		stopwatch:           stopwatch.NewWithInterval(10 * time.Millisecond),
		methods:             methods,
		currentMethod:       0,
		panelFocused:        false,
		focusComponentIndex: focusMethod,

		headers:         textarea.New(),
		headersExpanded: false,

		body:         textarea.New(),
		bodyExpanded: false,

		request: http.NewDefaultRequest(),
	}

	m.urlInput = textinput.New()
	//m.urlInput.Placeholder = "http://localhost:..."
	m.urlInput.SetValue("http://localhost:")
	m.urlInput.CharLimit = 40
	m.urlInput.Width = 60

	m.nameInput = textinput.New()
	m.nameInput.Placeholder = "My new awesome request.."
	m.nameInput.CharLimit = 40
	m.nameInput.Width = 60

	m.headers.Placeholder = "Content-Type = multipart/form-data,\nAuthorization= Bearer ...,"

	m.body.Placeholder = "key = value,\nname = volt,\nversion=1.0"
	return m
}
