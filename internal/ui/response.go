package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/utils"
)

var (
	inactiveTab = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("240"))
	activeTab   = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("76")).Foreground(lipgloss.Color("255"))
)

type ResponsePane struct {
	Response      *http.Response
	LoadTestStats *http.LoadTestStats
	isLoadTest    bool
	height, width int

	viewport  viewport.Model
	activeTab int
}

func (m ResponsePane) Init() tea.Cmd {
	return nil
}

func formatJSON(content string) string {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(content), "", "    ")

	if err != nil {
		return content
	}
	return pretty.String()
}

func highlightContent(content, lexer string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, content, lexer, "terminal256", "monokai")

	if err != nil {
		return content
	}
	return buf.String()
}

func (m *ResponsePane) SetResponse(response *http.Response) {
	m.Response = response
	m.isLoadTest = false

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

func (m *ResponsePane) SetLoadTestStats(stats *http.LoadTestStats) {
	m.LoadTestStats = stats
	m.isLoadTest = true
	m.activeTab = 0 // Reset to Overview tab
	m.updateViewportForActiveTab()
}

func (m *ResponsePane) ClearLoadTestStats() {
	m.LoadTestStats = nil
	m.isLoadTest = false
}

func (m *ResponsePane) SetHeight(height int) {
	m.height = height
	m.viewport.Height = height
}

func (m *ResponsePane) SetWidth(width int) {
	m.width = width
	m.viewport.Width = width
}

func (m *ResponsePane) updateViewportForActiveTab() {
	if m.isLoadTest {
		m.updateLoadTestTabContent()
		return
	}

	if m.Response == nil {
		return
	}

	var content string
	switch m.activeTab {
	case 0: // Body - re-render the response body
		if m.Response.Error != "" {
			content = m.Response.Error
		} else {
			contentType := m.Response.ParseContentType()
			content = m.Response.Body

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
		}
	case 1: // Headers
		content = m.renderHeaders()
	case 2: // Timing
		content = m.renderTiming()
	}
	m.viewport.SetContent(content)
}

func (m *ResponsePane) updateLoadTestTabContent() {
	if m.LoadTestStats == nil {
		m.viewport.SetContent("No data")
		return
	}

	var content string
	switch m.activeTab {
	case 0:
		content = m.renderLoadTestOverview()
	case 1:
		content = m.renderLoadTestLatency()
	case 2:
		content = m.renderLoadTestErrors()
	}
	m.viewport.SetContent(content)
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
	if m.isLoadTest {
		return m.renderLoadTestView()
	}

	if m.Response == nil {
		return "Make a request to see the response here!"
	}

	var statusBar string

	if m.Response.Error != "" {
		statusBar = ErrorStyle.Render("ERROR")
		m.viewport.SetContent(m.Response.Error)
	} else {
		statusBar = m.renderHeaderBar()
	}

	tabHeader := m.renderTabs()

	m.viewport.Height = m.height - 2
	tabContent := m.renderActiveTabContent()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		"\n",
		tabHeader,
		tabContent,
	)
}

func (m ResponsePane) renderLoadTestView() string {
	if m.LoadTestStats == nil {
		return "No load test data"
	}

	var b strings.Builder

	// Status line
	status := "Load Test "
	if m.LoadTestStats.EndTime.IsZero() {
		status += "In Progress..."
	} else {
		status += "Complete"
	}
	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		Render(status)
	b.WriteString(statusBar)
	b.WriteString("\n")

	// Tab header (reuse existing tab rendering for load test tabs)
	tabs := []string{"[1] Overview", "[2] Latency", "[3] Errors"}
	renderedTabs := []string{}
	for i, tab := range tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTab.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, inactiveTab.Render(tab))
		}
	}
	tabHeader := lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...)
	b.WriteString(tabHeader)
	b.WriteString("\n")

	// Tab content
	b.WriteString(m.viewport.View())

	return b.String()
}

func (m *ResponsePane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Direct tab access
		case "1":
			m.activeTab = 0
			m.updateViewportForActiveTab()
		case "2":
			m.activeTab = 1
			m.updateViewportForActiveTab()
		case "3":
			m.activeTab = 2
			m.updateViewportForActiveTab()
		// Tab navigation
		case "k", tea.KeyLeft.String():
			m.activeTab = (m.activeTab - 1 + 3) % 3
			m.updateViewportForActiveTab()
		case "j", tea.KeyRight.String():
			m.activeTab = (m.activeTab + 1) % 3
			m.updateViewportForActiveTab()
		// Copy handling
		case "y", "Y":
			if m.Response != nil && !m.isLoadTest {
				return m, m.copyToClipboard(m.Response.Body)
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ResponsePane) copyToClipboard(content string) tea.Cmd {
	return func() tea.Msg {
		clipboard.WriteAll(content)
		return nil
	}
}

func SetupResponsePane() ResponsePane {
	return ResponsePane{
		viewport:  viewport.New(20, 10),
		width:     20,
		height:    30,
		activeTab: 0,
	}
}

func (m ResponsePane) renderTabs() string {
	tabs := []string{"[1] Body", "[2] Headers", "[3] Timing"}
	renderedTabs := []string{}

	for i, tab := range tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTab.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, inactiveTab.Render(tab))
		}
	}
	renderedTabs = append(renderedTabs, "\n")

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderedTabs...,
	)
}

func (m ResponsePane) renderActiveTabContent() string {
	switch m.activeTab {
	case 0:
		return m.viewport.View()
	case 1:
		return m.renderHeaders()
	case 2:
		return m.renderTiming()
	default:
		return "Something went wrong."
	}
}

func (m ResponsePane) renderLoadTestOverview() string {
	stats := m.LoadTestStats
	if stats == nil {
		return "No data"
	}

	var b strings.Builder

	successCount := stats.CompletedRequests - stats.FailedRequests
	successRate := 0.0
	if stats.CompletedRequests > 0 {
		successRate = float64(successCount) / float64(stats.CompletedRequests) * 100
	}

	b.WriteString(fmt.Sprintf("Requests:       %d / %d\n", stats.CompletedRequests, stats.TotalRequests))
	b.WriteString(fmt.Sprintf("Success:        %d (%.1f%%)\n", successCount, successRate))
	b.WriteString(fmt.Sprintf("Failed:         %d (%.1f%%)\n", stats.FailedRequests, 100-successRate))
	b.WriteString("\n")

	// Throughput
	elapsed := time.Since(stats.StartTime)
	if !stats.EndTime.IsZero() {
		elapsed = stats.EndTime.Sub(stats.StartTime)
	}
	throughput := 0.0
	if elapsed.Seconds() > 0 {
		throughput = float64(stats.CompletedRequests) / elapsed.Seconds()
	}
	b.WriteString(fmt.Sprintf("Throughput:     %.1f req/s\n", throughput))
	b.WriteString(fmt.Sprintf("Duration:       %s\n", elapsed.Round(time.Millisecond)))

	return b.String()
}

func (m ResponsePane) renderLoadTestLatency() string {
	stats := m.LoadTestStats
	if stats == nil || stats.Percentiles == nil {
		return "No latency data"
	}

	var b strings.Builder
	b.WriteString("Latency Distribution:\n\n")
	b.WriteString(fmt.Sprintf("  Min:    %s\n", stats.MinDuration.Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("  p50:    %s\n", stats.Percentiles.Percentile(50).Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("  p90:    %s\n", stats.Percentiles.Percentile(90).Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("  p95:    %s\n", stats.Percentiles.Percentile(95).Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("  p99:    %s\n", stats.Percentiles.Percentile(99).Round(time.Millisecond)))
	b.WriteString(fmt.Sprintf("  Max:    %s\n", stats.MaxDuration.Round(time.Millisecond)))

	return b.String()
}

func (m ResponsePane) renderLoadTestErrors() string {
	stats := m.LoadTestStats
	if stats == nil {
		return "No error data"
	}

	var b strings.Builder

	if len(stats.Errors) == 0 {
		b.WriteString("No errors encountered!\n")
		return b.String()
	}

	b.WriteString("Error Breakdown:\n\n")
	for code, count := range stats.Errors {
		b.WriteString(fmt.Sprintf("  HTTP %s: %d occurrences\n", code, count))
	}

	return b.String()
}

// renderHeaders renders the response headers to the headers viewport
func (m ResponsePane) renderHeaders() string {
	if m.Response == nil || m.Response.Headers == nil {
		return "No headers available"
	}

	var b strings.Builder
	b.WriteString("Response Headers:\n")
	b.WriteString(strings.Repeat("-", len("Response Headers:")+2) + "\n\n")

	keys := make([]string, 0, len(m.Response.Headers))
	for k := range m.Response.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := m.Response.Headers[key]
		for _, value := range values {
			b.WriteString(keyStyle.Render(key))
			b.WriteString(": ")
			b.WriteString(value)
			b.WriteString("\n")
		}
	}
	return b.String()

}

// renderTiming renders the response timing data to the timing viewport
func (m ResponsePane) renderTiming() string {
	if m.Response == nil {
		return "No timing data available"
	}

	var b strings.Builder
	b.WriteString("Request Timing\n")
	b.WriteString(strings.Repeat("─", 60) + "\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(darkPurple).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	// Total duration
	b.WriteString(labelStyle.Render("Total Duration"))
	b.WriteString(": ")
	b.WriteString(valueStyle.Render(m.Response.Duration.String()))
	b.WriteString("\n\n")

	// Format in milliseconds for readability
	ms := m.Response.Duration.Milliseconds()
	b.WriteString(labelStyle.Render("Milliseconds"))
	b.WriteString(": ")
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d ms", ms)))
	b.WriteString("\n\n")

	// Connection type
	b.WriteString(labelStyle.Render("Connection Type"))
	b.WriteString(": ")
	if m.Response.RoundTrip {
		b.WriteString(valueStyle.Render("Round Trip (new connection)"))
	} else {
		b.WriteString(valueStyle.Render("Direct (keep-alive)"))
	}
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Faint(true).Render(
		"Note: Detailed timing breakdown (DNS, TLS, TTFB) coming in a future release!",
	))

	return b.String()
}
