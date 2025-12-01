package responsepane

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles Bubble Tea messages and state transitions
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
			m.activeTab = int(TabBody)
			m.updateViewportForActiveTab()
		case "2":
			m.activeTab = int(TabHeaders)
			m.updateViewportForActiveTab()
		case "3":
			m.activeTab = int(TabTiming)
			m.updateViewportForActiveTab()
		// Tab navigation
		case "h", tea.KeyLeft.String():
			maxTabs := m.getMaxTabs()
			m.activeTab = (m.activeTab - 1 + maxTabs) % maxTabs
			m.updateViewportForActiveTab()
		case "l", tea.KeyRight.String():
			maxTabs := m.getMaxTabs()
			m.activeTab = (m.activeTab + 1) % maxTabs
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

// copyToClipboard copies content to the system clipboard
func (m ResponsePane) copyToClipboard(content string) tea.Cmd {
	return func() tea.Msg {
		clipboard.WriteAll(content)
		return nil
	}
}

// updateViewportForActiveTab updates the viewport content based on the active tab
func (m *ResponsePane) updateViewportForActiveTab() {
	if m.isLoadTest {
		m.updateLoadTestTabContent()
		return
	}

	if m.Response == nil {
		return
	}

	var content string
	switch TabIndex(m.activeTab) {
	case TabBody:
		if m.Response.Error != "" {
			content = m.Response.Error
		} else {
			contentType := m.Response.ParseContentType()
			content = formatContentByType(m.Response.Body, contentType)
		}
	case TabHeaders:
		content = m.renderHeaders()
	case TabTiming:
		content = m.renderTiming()
	}
	m.viewport.SetContent(content)
}

// updateLoadTestTabContent updates the viewport content for load test tabs
func (m *ResponsePane) updateLoadTestTabContent() {
	if m.LoadTestStats == nil {
		m.viewport.SetContent("No data")
		return
	}

	var content string
	switch TabIndex(m.activeTab) {
	case TabLoadTestOverview:
		content = m.renderLoadTestOverview()
	case TabLoadTestLatency:
		content = m.renderLoadTestLatency()
	case TabLoadTestErrors:
		content = m.renderLoadTestErrors()
	}
	m.viewport.SetContent(content)
}
