package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/storage"
	"github.com/owenHochwald/volt/internal/ui"
	"github.com/owenHochwald/volt/internal/utils"
)

type Model struct {
	db *storage.SQLiteStorage

	sidebarPane  *ui.SidebarPane
	requestPane  ui.RequestPane
	responsePane *ui.ResponsePane

	savedRequests []http.Request

	focusedPanel utils.Panel

	width, height int
}

func SetupModel(db *storage.SQLiteStorage) Model {
	responsePane := ui.SetupResponsePane()
	m := Model{
		db:           db,
		sidebarPane:  ui.NewSidebar(db),
		requestPane:  ui.SetupRequestPane(db),
		responsePane: &responsePane,
		focusedPanel: utils.SidebarPanel,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.sidebarPane.Init())

	return tea.Batch(cmds...)
}
