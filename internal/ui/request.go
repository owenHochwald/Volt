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
	"github.com/owenHochwald/volt/internal/storage"
	"github.com/owenHochwald/volt/internal/utils"
)

var (
	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
)

type RequestPane struct {
	client *http.Client

	stopwatch stopwatch.Model
	quitting  bool

	panelFocused bool

	focusManager *FocusManager

	methodSelector *MethodSelector
	urlInput       *textinput.Model
	nameInput      *textinput.Model
	headers        *textarea.Model
	body           *textarea.Model
	submitButton   *SubmitButton

	request *http.Request

	height int

	parseErrors []string

	headersExpanded bool

	requestInProgress bool
	// queryParams     []QueryParam
	bodyExpanded bool
	// validationError error

	db *storage.SQLiteStorage
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
		case tea.KeyCtrlS.String(), tea.KeyShiftDown.String():
			m.syncRequest()
			return m, SaveRequestCmd(m.db, m.request)
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case tea.KeyTab.String(), tea.KeyDown.String():
			m.focusManager.Next()
		case tea.KeyUp.String():
			m.focusManager.Prev()
		}

		switch m.focusManager.CurrentIndex() {
		case 0:
			switch msg.String() {
			case tea.KeyRight.String(), "l":
				m.methodSelector.Next()
			case tea.KeyLeft.String(), "h":
				m.methodSelector.Prev()
			}
		case 1:
			var cmd tea.Cmd
			*m.urlInput, cmd = m.urlInput.Update(msg)
			return m, cmd
		case 2:
			var cmd tea.Cmd
			*m.nameInput, cmd = m.nameInput.Update(msg)
			return m, cmd
		case 3:
			var cmd tea.Cmd
			*m.headers, cmd = m.headers.Update(msg)
			return m, cmd
		case 4:
			var cmd tea.Cmd
			*m.body, cmd = m.body.Update(msg)
			return m, cmd
		case 5:
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
	} else if m.submitButton.IsFocused() {
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

	helpText := HelpStyle.Render("tab/↑/↓: navigate • ←/→ or h/l: change method • enter: send • ctrl+s: save request")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		lipgloss.NewStyle().Height(m.height-10).Render(""),
		stopwatchCount,
		helpText,
	)

}

func SetupRequestPane(db *storage.SQLiteStorage) RequestPane {
	methodSelector := NewMethodSelector()

	urlInput := textinput.New()
	urlInput.SetValue("http://localhost:")
	urlInput.CharLimit = 40
	urlInput.Width = 60

	nameInput := textinput.New()
	nameInput.Placeholder = "My new awesome request.."
	nameInput.CharLimit = 40
	nameInput.Width = 60

	headers := textarea.New()
	headers.Placeholder = "Content-Type = multipart/form-data,\nAuthorization= Bearer ...,"

	body := textarea.New()
	body.Placeholder = "key = value,\nname = volt,\nversion=1.0"

	submitButton := NewSubmitButton()

	components := []Focusable{
		methodSelector,
		&urlInput,
		&nameInput,
		&headers,
		&body,
		submitButton,
	}
	focusManager := NewFocusManager(components)

	m := RequestPane{
		methodSelector:  methodSelector,
		urlInput:        &urlInput,
		nameInput:       &nameInput,
		headers:         &headers,
		body:            &body,
		focusManager:    focusManager,
		client:          http.InitClient(0, false),
		stopwatch:       stopwatch.NewWithInterval(10 * time.Millisecond),
		panelFocused:    false,
		submitButton:    submitButton,
		headersExpanded: false,
		bodyExpanded:    false,
		request:         http.NewDefaultRequest(),
		db:              db,
	}

	return m

}
