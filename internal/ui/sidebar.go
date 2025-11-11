package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/storage"
)

type RequestItem struct {
	title, desc string
	Request     *http.Request
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

	db *storage.SQLiteStorage
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

func (s *SidebarPane) SetRequests(items []list.Item) {
	s.requestsList = list.New(items, list.NewDefaultDelegate(), s.width, s.height)
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
