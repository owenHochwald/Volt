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
	responsePane ui.ResponsePane
	headerPane   *ui.Header

	savedRequests []http.Request

	focusedPanel utils.Panel

	width, height int
}

func SetupModel(db *storage.SQLiteStorage) Model {
	m := Model{
		db:           db,
		sidebarPane:  ui.NewSidebar(db),
		requestPane:  ui.SetupRequestPane(db),
		responsePane: ui.SetupResponsePane(),
		focusedPanel: utils.SidebarPanel,
		headerPane:   ui.SetupHeader(),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return m.sidebarPane.Init()
}
