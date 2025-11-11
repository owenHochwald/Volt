package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/storage"
	"github.com/owenHochwald/volt/internal/ui"
)

type Panel int

const (
	SidebarPanel Panel = iota
	RequestPanel
	ResponsePanel
)

type RequestItem struct {
	title, desc string
}

func (i RequestItem) Title() string       { return i.title }
func (i RequestItem) Description() string { return i.desc }
func (i RequestItem) FilterValue() string { return i.title }

type Model struct {
	db *storage.SQLiteStorage

	requestsList list.Model

	requestPane  ui.RequestPane
	responsePane ui.ResponsePane

	selectedRequest *RequestItem

	savedRequests []http.Request

	focusedPanel Panel

	width, height int
}

func SetupModel(db *storage.SQLiteStorage) Model {
	mockRequestsList := []list.Item{
		RequestItem{title: "Get Users", desc: "test"},
		RequestItem{title: "Delete a User", desc: "test"},
		RequestItem{title: "Update a User", desc: "test"},
	}

	m := Model{
		db:              db,
		requestsList:    list.New(mockRequestsList, list.NewDefaultDelegate(), 0, 0),
		requestPane:     ui.SetupRequestPane(),
		responsePane:    ui.SetupResponsePane(),
		selectedRequest: nil,
		focusedPanel:    SidebarPanel,
	}
	InitialSidebar(&m)
	return m
}

func InitialSidebar(m *Model) {
	m.requestsList.Title = fmt.Sprintf("Saved (%d)", len(m.requestsList.Items()))
	customKeys := newCustomReqKeys()
	m.requestsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.newItem,
			customKeys.delete,
		}
	}

	m.requestsList.SetShowHelp(true)
}

type customReqKeys struct {
	newItem key.Binding
	delete  key.Binding
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

func (m Model) Init() tea.Cmd {
	return nil
}
