package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
)

type Panel int

const (
	SidebarPanel Panel = iota
	RequestPanel
	ResponsePanel
)

type HttpMethod struct {
	title, desc string
}

func (i HttpMethod) Title() string       { return i.title }
func (i HttpMethod) Description() string { return i.desc }
func (i HttpMethod) FilterValue() string { return i.title }

type RequestItem struct {
	title, desc string
}

func (i RequestItem) Title() string       { return i.title }
func (i RequestItem) Description() string { return i.desc }
func (i RequestItem) FilterValue() string { return i.title }

type Model struct {
	httpMethods  list.Model
	requestsList list.Model

	selectedRequest *HttpMethod

	// SQLite State
	savedRequests []http.Request

	focusedPanel Panel

	width, height int
}

func InitialModel() Model {
	// TODO: Use these for new request creation: Only show saved reqs in sidebar
	items := []list.Item{
		HttpMethod{title: "GET", desc: "Get a resource"},
		HttpMethod{title: "POST", desc: "Create a resource"},
		HttpMethod{title: "PUT", desc: "Update a resource"},
		HttpMethod{title: "DELETE", desc: "Delete a resource"},
		HttpMethod{title: "PATCH", desc: "Apply partial changes to a resource"},
	}

	mockRequestsList := []list.Item{
		RequestItem{title: "Get Users", desc: "test"},
		RequestItem{title: "Delete all Users", desc: "test"},
		RequestItem{title: "Update a Single User", desc: "test"},
	}

	m := Model{
		httpMethods:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		requestsList:    list.New(mockRequestsList, list.NewDefaultDelegate(), 0, 0),
		selectedRequest: nil,
		focusedPanel:    SidebarPanel,
	}
	m.httpMethods.Title = "HTTP Methods"
	m.requestsList.Title = "Saved Requests"
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
