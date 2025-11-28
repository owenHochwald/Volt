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

	// Load test mode
	loadTestMode         bool
	loadTestFocusManager *FocusManager
	loadTestConcurrency  *textinput.Model
	loadTestTotalReqs    *textinput.Model
	loadTestQPS          *textinput.Model
	loadTestTimeout      *textinput.Model
}

func (m *RequestPane) syncRequest() {
	if m.request.Headers == nil {
		m.request.Headers = make(map[string]string)
	}

	m.request.Method = m.methodSelector.Current()
	m.request.URL = m.urlInput.Value()
	m.request.Name = m.nameInput.Value()
	headerMap, headerErrors := utils.ParseKeyValuePairs(m.headers.Value())
	bodyMap, bodyErrors := utils.ParseKeyValuePairs(m.body.Value())
	jsonData, err := json.Marshal(bodyMap)
	if err != nil {
		m.parseErrors = append(m.parseErrors, "JSON marshal error: "+err.Error())
		m.request.Headers = headerMap
		m.request.Body = "{}" // Set to valid empty JSON
		m.parseErrors = append(m.parseErrors, headerErrors...)
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

func (m *RequestPane) buildJobConfig() (*http.JobConfig, error) {
	var parseErrors []string

	concurrency := 100
	if m.loadTestConcurrency.Value() != "" {
		n, err := fmt.Sscanf(m.loadTestConcurrency.Value(), "%d", &concurrency)
		if err != nil || n != 1 || concurrency <= 0 {
			parseErrors = append(parseErrors, "Invalid concurrency (must be positive integer)")
			concurrency = 100
		}
	}

	totalRequests := 10000
	if m.loadTestTotalReqs.Value() != "" {
		n, err := fmt.Sscanf(m.loadTestTotalReqs.Value(), "%d", &totalRequests)
		if err != nil || n != 1 || totalRequests <= 0 {
			parseErrors = append(parseErrors, "Invalid total requests (must be positive integer)")
			totalRequests = 10000
		}
	}

	qps := 0.0
	if m.loadTestQPS.Value() != "" {
		n, err := fmt.Sscanf(m.loadTestQPS.Value(), "%f", &qps)
		if err != nil || n != 1 || qps < 0 {
			parseErrors = append(parseErrors, "Invalid QPS (must be non-negative number)")
			qps = 0.0
		}
	}

	timeout := 30 * time.Second
	if m.loadTestTimeout.Value() != "" {
		parsedTimeout, err := time.ParseDuration(m.loadTestTimeout.Value())
		if err != nil {
			parseErrors = append(parseErrors, "Invalid timeout format (use 30s, 1m, etc.)")
			timeout = 30 * time.Second
		} else if parsedTimeout <= 0 {
			parseErrors = append(parseErrors, "Timeout must be positive")
			timeout = 30 * time.Second
		} else {
			timeout = parsedTimeout
		}
	}

	m.parseErrors = append(m.parseErrors, parseErrors...)

	if err := m.request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	return &http.JobConfig{
		Request:       m.request,
		Concurrency:   concurrency,
		TotalRequests: totalRequests,
		QPS:           qps,
		Timeout:       timeout,
	}, nil
}

func (m *RequestPane) ExitLoadTestMode() {
	m.loadTestMode = false
	m.requestInProgress = false
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
	case SetRequestPaneRequestMsg:
		m.reinitRequesetPane(msg.Request)
		return m, nil
	case tea.KeyMsg:
		if !m.panelFocused {
			return m, nil
		}
		switch msg.String() {
		case "alt+l":
			m.loadTestMode = !m.loadTestMode
			if m.loadTestMode {
				unifiedComponents := []Focusable{
					m.methodSelector,
					m.urlInput,
					m.nameInput,
					m.headers,
					m.body,
					m.loadTestConcurrency,
					m.loadTestTotalReqs,
					m.loadTestQPS,
					m.loadTestTimeout,
					m.submitButton,
				}
				m.focusManager = NewFocusManager(unifiedComponents)
			} else {
				components := []Focusable{
					m.methodSelector,
					m.urlInput,
					m.nameInput,
					m.headers,
					m.body,
					m.submitButton,
				}
				m.focusManager = NewFocusManager(components)
			}
			return m, nil
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

		// Load test mode - unified navigation through all fields
		if m.loadTestMode {
			var cmds []tea.Cmd

			switch m.focusManager.CurrentIndex() {
			case 0:
				// Method selector
				switch msg.String() {
				case tea.KeyRight.String(), "l":
					m.methodSelector.Next()
				case tea.KeyLeft.String(), "h":
					m.methodSelector.Prev()
				}
			case 1:
				// URL input
				*m.urlInput, cmd = m.urlInput.Update(msg)
				cmds = append(cmds, cmd)
			case 2:
				// Name input
				*m.nameInput, cmd = m.nameInput.Update(msg)
				cmds = append(cmds, cmd)
			case 3:
				// Headers
				*m.headers, cmd = m.headers.Update(msg)
				cmds = append(cmds, cmd)
			case 4:
				// Body
				*m.body, cmd = m.body.Update(msg)
				cmds = append(cmds, cmd)
			case 5:
				// Load test concurrency
				*m.loadTestConcurrency, cmd = m.loadTestConcurrency.Update(msg)
				cmds = append(cmds, cmd)
			case 6:
				// Load test total requests
				*m.loadTestTotalReqs, cmd = m.loadTestTotalReqs.Update(msg)
				cmds = append(cmds, cmd)
			case 7:
				// Load test QPS
				*m.loadTestQPS, cmd = m.loadTestQPS.Update(msg)
				cmds = append(cmds, cmd)
			case 8:
				// Load test timeout
				*m.loadTestTimeout, cmd = m.loadTestTimeout.Update(msg)
				cmds = append(cmds, cmd)
			case 9:
				// Submit button
				switch msg.String() {
				case tea.KeyEnter.String():
					if m.requestInProgress {
						return m, nil
					}

					m.syncRequest()
					m.requestInProgress = true

					config, err := m.buildJobConfig()
					if err != nil {
						m.parseErrors = append(m.parseErrors, "Load test config error: "+err.Error())
						m.requestInProgress = false
						return m, nil
					}
					return m, StartLoadTestCmd(config)
				}
			}

			return m, tea.Batch(cmds...)
		}

		// Normal mode - handle regular fields
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
		if m.loadTestMode {
			button = FocusedButton.Render("Running Load Test...")
		} else {
			button = FocusedButton.Render("Sending...")
			elapsed := m.stopwatch.Elapsed()
			milliseconds := elapsed.Milliseconds()
			seconds := float64(milliseconds) / 1000.0
			stopwatchCount = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Render(fmt.Sprintf("%.3fs", seconds))
		}
	} else if m.submitButton.IsFocused() {
		button = FocusedButton.Render("→ Send")
	} else {
		button = UnfocusedButton.Render("→ Send")
	}

	var mainContent string
	if m.loadTestMode {
		// Show load test configuration
		ltConcurrencyLabel := labelStyle.Render("Concurrency:    ")
		ltConcurrencyLine := lipgloss.JoinHorizontal(lipgloss.Left,
			ltConcurrencyLabel, m.loadTestConcurrency.View())

		ltTotalLabel := labelStyle.Render("Total Requests: ")
		ltTotalLine := lipgloss.JoinHorizontal(lipgloss.Left,
			ltTotalLabel, m.loadTestTotalReqs.View())

		ltQPSLabel := labelStyle.Render("QPS (limit):    ")
		ltQPSLine := lipgloss.JoinHorizontal(lipgloss.Left,
			ltQPSLabel, m.loadTestQPS.View())

		ltTimeoutLabel := labelStyle.Render("Timeout:        ")
		ltTimeoutLine := lipgloss.JoinHorizontal(lipgloss.Left,
			ltTimeoutLabel, m.loadTestTimeout.View())

		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			primaryLine,
			nameLine,
			headersLine,
			bodyLine,
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true).Render("Load Test Configuration:"),
			ltConcurrencyLine,
			ltTotalLine,
			ltQPSLine,
			ltTimeoutLine,
			"",
			button,
		)
	} else {
		// Normal view
		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			primaryLine,
			nameLine,
			headersLine,
			bodyLine,
			"",
			button,
		)
	}

	var helpText string
	if m.loadTestMode {
		helpText = HelpStyle.Render("alt+l: exit load test mode • tab/↑/↓: navigate • enter: start load test")
	} else {
		helpText = HelpStyle.Render("alt+l: load test mode • tab/↑/↓: navigate • ←/→ or h/l: change method • enter: send • ctrl+s: save")
	}

	finalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		lipgloss.NewStyle().Height(m.height-10).Render(""),
		stopwatchCount,
		helpText,
	)

	// Apply yellow border if in load test mode
	if m.loadTestMode {
		borderStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("226")). // Yellow
			Padding(0, 1)
		return borderStyle.Render(finalContent)
	}

	return finalContent
}

func (m *RequestPane) reinitRequesetPane(request *http.Request) {
	m.request = request

	m.methodSelector.SetCurrentIndex(request.Method)
	m.urlInput.SetValue(request.URL)
	m.nameInput.SetValue(request.Name)
	m.headers.SetValue(utils.ParseMapToString(request.Headers))
	m.body.SetValue(request.Body[1 : len(request.Body)-1])
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

	ltConcurrency := textinput.New()
	ltConcurrency.Placeholder = "100"
	ltConcurrency.CharLimit = 5
	ltConcurrency.Width = 15
	// Don't call Focus() here - let focus manager handle it

	ltTotalReqs := textinput.New()
	ltTotalReqs.Placeholder = "10000"
	ltTotalReqs.CharLimit = 10
	ltTotalReqs.Width = 15

	ltQPS := textinput.New()
	ltQPS.Placeholder = "0 (unlimited)"
	ltQPS.CharLimit = 10
	ltQPS.Width = 15

	ltTimeout := textinput.New()
	ltTimeout.Placeholder = "30s"
	ltTimeout.CharLimit = 10
	ltTimeout.Width = 15

	components := []Focusable{
		methodSelector,
		&urlInput,
		&nameInput,
		&headers,
		&body,
		submitButton,
	}
	focusManager := NewFocusManager(components)

	// Create focus manager for load test inputs
	loadTestComponents := []Focusable{
		&ltConcurrency,
		&ltTotalReqs,
		&ltQPS,
		&ltTimeout,
	}
	loadTestFocusManager := NewFocusManager(loadTestComponents)

	m := RequestPane{
		methodSelector:       methodSelector,
		urlInput:             &urlInput,
		nameInput:            &nameInput,
		headers:              &headers,
		body:                 &body,
		focusManager:         focusManager,
		loadTestFocusManager: loadTestFocusManager, // NOW INITIALIZED
		client:               http.InitClient(0, false),
		stopwatch:            stopwatch.NewWithInterval(10 * time.Millisecond),
		panelFocused:         false,
		submitButton:         submitButton,
		headersExpanded:      false,
		bodyExpanded:         false,
		request:              http.NewDefaultRequest(),
		db:                   db,
		loadTestMode:         false,
		loadTestConcurrency:  &ltConcurrency,
		loadTestTotalReqs:    &ltTotalReqs,
		loadTestQPS:          &ltQPS,
		loadTestTimeout:      &ltTimeout,
	}

	return m

}
