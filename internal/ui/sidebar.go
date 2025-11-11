package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type RequestItem struct {
	title, desc string
}

type customReqKeys struct {
	newItem key.Binding
	delete  key.Binding
}

func (i RequestItem) Title() string       { return i.title }
func (i RequestItem) Description() string { return i.desc }
func (i RequestItem) FilterValue() string { return i.title }

type SidebarPane struct {
	panelFocused  bool
	height, width int

	requestsList    list.Model
	selectedRequest *RequestItem
}

func newCustomReqKeys() customReqKeys {
	return customReqKeys{
		newItem: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new request"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete request"),
		),
	}
}

func (s *SidebarPane) Init() tea.Cmd {
	return nil
}

func (s *SidebarPane) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.requestsList, cmd = s.requestsList.Update(msg)
	return cmd
}

func (s *SidebarPane) View() string {
	return s.requestsList.View()
}

func (s *SidebarPane) SelectedItem() (RequestItem, bool) {
	if item := s.requestsList.SelectedItem(); item != nil {
		if reqItem, ok := item.(RequestItem); ok {
			return reqItem, true
		}
	}
	return RequestItem{}, false
}

func (s *SidebarPane) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.requestsList.SetSize(width, height)
}

func NewSidebar() *SidebarPane {
	mockRequestsList := []list.Item{
		RequestItem{title: "Get Users", desc: "test"},
		RequestItem{title: "Delete a User", desc: "test"},
		RequestItem{title: "Update a User", desc: "test"},
	}

	sidebar := &SidebarPane{
		panelFocused: false,
		height:       10,
		width:        10,
		requestsList: list.New(mockRequestsList, list.NewDefaultDelegate(), 0, 0),
	}

	sidebar.requestsList.Title = fmt.Sprintf("Saved (%d)", len(sidebar.requestsList.Items()))
	customKeys := newCustomReqKeys()
	sidebar.requestsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.newItem,
			customKeys.delete,
		}
	}

	sidebar.requestsList.SetShowHelp(true)

	return sidebar
}
