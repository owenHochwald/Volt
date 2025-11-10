package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

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
		// TODO: add standard error handling logic
		fmt.Fprintf(os.Stderr, "Error highlighting response: %v", err)
	}
	return pretty.String()
}

func highlightJSON(content string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, content, "json", "terminal256", "monokai")

	if err != nil {
		// TODO: add standard error handling logic
		fmt.Fprintf(os.Stderr, "Error highlighting response: %v", err)
	}
	return buf.String()
}

func (m *ResponsePane) SetResponse(response *http.Response) {
	m.Response = response

	if m.Response != nil {

		// add more sophisticated parsing for json type
		formatted := formatJSON(m.Response.Body)
		content := highlightJSON(formatted)

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

	// add in response header bar for error routes
	if m.Response.Error != "" {
		m.viewport.SetContent(m.Response.Error)
		lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Background(lipgloss.Color("160")).Render(
			m.viewport.View(),
		)
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
