package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen int

const (
	ScreenList Screen = iota
	ScreenDetail
)

type Item struct {
	title, desc string
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type Model struct {
	requests        list.Model
	selectedRequest *Item
	currentScreen   Screen
}

func InitialModel() Model {
	items := []list.Item{
		Item{title: "GET", desc: "Get a resource"},
		Item{title: "POST", desc: "Create a resource"},
		Item{title: "PUT", desc: "Update a resource"},
		Item{title: "DELETE", desc: "Delete a resource"},
		Item{title: "PATCH", desc: "Apply partial changes to a resource"},
	}

	m := Model{
		requests:        list.New(items, list.NewDefaultDelegate(), 0, 0),
		selectedRequest: nil,
		currentScreen:   ScreenList,
	}
	m.requests.Title = "HTTP Methods"

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
