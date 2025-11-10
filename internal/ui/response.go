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

type ResponsePane struct {
	Response      *http.Response
	height, width int

	viewport viewport.Model
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
			fmt.Println("Sorry, we don't support image/jpeg yet!")
		case strings.Contains(contentType, "text/html"):
			content = highlightContent(content, "html")
		case strings.Contains(contentType, "text/plain"):
			content = highlightContent(content, "plaintext")
		case strings.Contains(contentType, "application/xml"):
			content = highlightContent(content, "xml")
		case strings.Contains(contentType, "application/graphql"):
			fmt.Println("Sorry, we don't support graphql yet!")
		case strings.Contains(contentType, "multipart/form-data"):
			fmt.Println("Sorry, we don't support multipart/form-data yet!")
		default:
			fmt.Printf("Unhandled Content-Type: %s\n", contentType)
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

func (m *ResponsePane) GetCurrentMethod() string {
	return ""
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
	size := fmt.Sprintf(" %d bytes", len([]byte(m.Response.Body)))

	return lipgloss.JoinHorizontal(lipgloss.Left, " | ", status, " | ", duration, " | ", size)
}

func (m ResponsePane) View() string {
	if m.Response == nil {
		return "Make a request to see the response here!"
	}

	if m.Response.Error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("160")).
			Background(lipgloss.Color("52"))
		statusBar := errorStyle.Render("ERROR")
		return lipgloss.JoinVertical(lipgloss.Left, statusBar, m.viewport.View())
	}

	statusBar := m.renderHeaderBar()
	return lipgloss.JoinVertical(lipgloss.Left, statusBar, m.viewport.View())
}

func (m ResponsePane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func SetupResponsePane() ResponsePane {
	return ResponsePane{
		viewport: viewport.New(20, 10),
		width:    20,
		height:   30,
	}
}
