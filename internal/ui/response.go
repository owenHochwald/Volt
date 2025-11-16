package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/utils"
)

// Added styles for the new tabs
var (
	InactiveTab = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("240"))
	ActiveTab   = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("76")).Foreground(lipgloss.Color("255"))
)

type ResponsePane struct {
	Response *http.Response
	height, width int

	viewport viewport.Model
	// activeTab tracks the currently selected tab
	activeTab int // 0=Body, 1=Headers, 2=Cookies, 3=Timing
}

func (m ResponsePane) Init() tea.Cmd {
	return nil
}

func formatJSON(content string) string {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(content), "", "    ")

	if err != nil {
		// return original content as fallback
		return content
	}
	return pretty.String()
}

func highlightContent(content, lexer string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, content, lexer, "terminal256", "monokai")

	if err != nil {
		// return original content as fallback
		return content
	}
	return buf.String()
}

func (m *ResponsePane) SetResponse(response *http.Response) {
	m.Response = response

	if m.Response != nil {

		if m.Response.Error != "" {
			m.viewport.SetContent(m.Response.Error)
			return
		}
		contentType := m.Response.ParseContentType()
		content := m.Response.Body

		switch {
		case strings.Contains(contentType, "application/json"):
			formatted := formatJSON(m.Response.Body)
			content = highlightContent(formatted, "json")

		case strings.Contains(contentType, "image/jpeg"):
			content = "Sorry, we don't support image/jpeg yet!"
		case strings.Contains(contentType, "text/html"):
			content = highlightContent(content, "html")
		case strings.Contains(contentType, "text/plain"):
			content = highlightContent(content, "plaintext")
		case strings.Contains(contentType, "application/xml"):
			content = highlightContent(content, "xml")
		case strings.Contains(contentType, "application/graphql"):
			content = "Sorry, we don't support graphql yet!"
		case strings.Contains(contentType, "multipart/form-data"):
			content = "Sorry, we don't support multipart/form-data yet!"
		default:
			content = fmt.Sprintf("Unhandled Content-Type: %s\n", contentType)
		}
		m.viewport.SetContent(content)
	}
}

func (m *ResponsePane) SetHeight(height int) {
	m.height = height
	m.viewport.Height = height
}

func (m *ResponsePane) SetWidth(width int) {
	m.width = width
	m.viewport.Width = width

}

func (m ResponsePane) renderHeaderBar() string {
	statusStyle := utils.MapStatusCodeToColor(m.Response.StatusCode)
	status := statusStyle.Render(m.Response.Status)
	duration := fmt.Sprintf(" %d ms", m.Response.Duration.Milliseconds())
	if m.Response.RoundTrip {
		duration += " (round trip)"
	} else {
		duration += " (direct)"
	}
	size := fmt.Sprintf(" %s", utils.FormatSize(len(m.Response.Body)))
	return lipgloss.JoinHorizontal(lipgloss.Left, " | ", status, " | ", duration, " | ", size)
}

func (m ResponsePane) View() string {
	if m.Response == nil {
		return "Make a request to see the response here!"
	}

	if m.Response.Error != "" {
		// Still show the error, but do it in the body panel
		statusBar := ErrorStyle.Render("ERROR")
		m.viewport.SetContent(m.Response.Error)
		// return lipgloss.JoinVertical(lipgloss.Left, statusBar, m.viewport.View())
	}

	statusBar := m.renderHeaderBar()
	// Render the new tab header
	tabHeader := m.renderTabs()

	// Set viewport height dynamically, leaving room for status and tabs
	m.viewport.Height = m.height - 2
	// Render the content for the active tab
	tabContent := m.renderActiveTabContent()

	return lipgloss.JoinVertical(lipgloss.Left, statusBar, tabHeader, tabContent)
}

func (m ResponsePane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Handle key presses for tab switching
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.activeTab = 0
		case "2":
			m.activeTab = 1
		case "3":
			m.activeTab = 2
		case "4":
			m.activeTab = 3
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func SetupResponsePane() ResponsePane {
	return ResponsePane{
		viewport:  viewport.New(20, 10),
		width:     20,
		height:    30,
		activeTab: 0, // Default to the "Body" tab
	}
}

// renderTabs draws the tab interface
func (m ResponsePane) renderTabs() string {
	tabs := []string{"[1] Body", "[2] Headers", "[3] Cookies", "[4] Timing"}
	renderedTabs := []string{}

	for i, tab := range tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, ActiveTab.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, InactiveTab.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...)
}

// renderActiveTabContent returns the content for the currently active tab
func (m ResponsePane) renderActiveTabContent() string {
	switch m.activeTab {
	case 0: // Body
		// The viewport content is already set
		return m.viewport.View()
	case 1: // Headers
		return m.renderHeaders()
	case 2: // Cookies
		return m.renderCookies()
	case 3: // Timing
		return m.renderTiming()
	default:
		return "Something went wrong."
	}
}

// renderHeaders creates a formatted string of all response headers
func (m ResponsePane) renderHeaders() string {
	// Placeholder: This function needs to be implemented
	return "Headers Content Goes Here\n\n(Not yet implemented)"
}

// renderCookies creates a formatted string of cookies
func (m ResponsePane) renderCookies() string {
	// Placeholder: This function needs to be implemented
	return "Cookies Content Goes Here\n\n(Not yet implemented)"
}

// renderTiming creates a formatted string of performance data
func (m ResponsePane) renderTiming() string {
	// Placeholder: This function needs to be implemented
	total := fmt.Sprintf("Total Duration: %s\n", m.Response.Duration)
	return total + "\nTiming Content Goes Here\n\n(Not yet implemented)"
}