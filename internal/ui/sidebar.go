package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/storage"
)

type RequestItem struct {
	title, desc string
	Request     *http.Request
}

func (i RequestItem) Title() string       { return i.title }
func (i RequestItem) Description() string { return i.desc }
func (i RequestItem) FilterValue() string { return i.title }

type SidebarPane struct {
	panelFocused  bool
	height, width int

	requestsList    list.Model
	selectedRequest *RequestItem

	db *storage.SQLiteStorage
}

func (s *SidebarPane) SetRequests(items []list.Item) {
	s.requestsList = list.New(items, list.NewDefaultDelegate(), s.width, s.height)
	s.requestsList.SetShowHelp(false)
}

func (s *SidebarPane) Init() tea.Cmd {
	return LoadRequestsCmd(s.db)
}

func (s *SidebarPane) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case RequestsLoadingMsg:
		if msg.Err != nil {
			s.SetRequests([]list.Item{})
			s.requestsList.Title = "Saved (0)"
			return nil

		}
		items := make([]list.Item, 0, len(msg.Requests))
		for _, req := range msg.Requests {
			items = append(items, RequestItem{
				title:   req.Name,
				desc:    req.URL[max(len(req.URL)-10, 0):],
				Request: &req,
			})
		}
		s.SetRequests(items)
		s.requestsList.Title = fmt.Sprintf("Saved (%d)", len(s.requestsList.Items()))
		return nil
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			// TODO: make delete request async
			item, ok := s.SelectedItem()
			if !ok || item.Request == nil || item.Request.ID == 0 {
				return nil
			}
			err := s.db.Delete(item.Request.ID)
			if err != nil {
				return nil
			}
			//s.SetRequests(s.requestsList.Items())
			return LoadRequestsCmd(s.db)
		}

	}

	s.requestsList, cmd = s.requestsList.Update(msg)

	// TODO: handle key presses for deleting requests
	return cmd
}

func (s *SidebarPane) View() string {
	helpText := HelpStyle.Render("n: new • d: delete •enter: send")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.requestsList.View(),
		lipgloss.NewStyle().Height(s.height-1).Render(""),
		helpText,
	)
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

func NewSidebar(db *storage.SQLiteStorage) *SidebarPane {
	loadingItems := []list.Item{
		RequestItem{
			title:   "Loading...",
			desc:    "Loading saved requests...",
			Request: nil,
		},
	}

	sidebar := &SidebarPane{
		panelFocused: false,
		height:       10,
		width:        10,
		db:           db,
		requestsList: list.New(loadingItems, list.NewDefaultDelegate(), 0, 0),
	}
	sidebar.requestsList.Title = "Saved (Loading...)"
	sidebar.requestsList.SetShowHelp(false)

	return sidebar
}
