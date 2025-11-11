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
	return nil
}

func (s *SidebarPane) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
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
	// TODO is there a way to make this async?
	requests, err := db.Load()
	var requestItems []list.Item

	if err != nil {
		requestItems = []list.Item{}
	} else {
		requestItems = make([]list.Item, len(requests))
		for _, req := range requests {
			requestItems = append(requestItems, RequestItem{
				title: req.Name,
				// TODO: add cool URL shading / fading styles
				desc:    req.URL[(len(req.URL) - 10):],
				Request: &req,
			})
		}
	}

	sidebar := &SidebarPane{
		panelFocused: false,
		height:       10,
		width:        10,
		db:           db,
		requestsList: list.New(requestItems, list.NewDefaultDelegate(), 0, 0),
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
