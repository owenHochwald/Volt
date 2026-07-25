package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/owenHochwald/Volt/internal/http"
	"github.com/owenHochwald/Volt/internal/storage"
	"github.com/owenHochwald/Volt/internal/ui"
	"github.com/owenHochwald/Volt/internal/ui/keybindings"
	"github.com/owenHochwald/Volt/internal/ui/requestpane"
	"github.com/owenHochwald/Volt/internal/ui/responsepane"
	"github.com/owenHochwald/Volt/internal/ui/shortcutpane"
	"github.com/owenHochwald/Volt/internal/utils"
)

type Model struct {
	db   *storage.SQLiteStorage
	keys keybindings.KeyMap

	sidebarPane  *ui.SidebarPane
	requestPane  requestpane.RequestPane
	responsePane *responsepane.ResponsePane
	headerPane   *ui.Header
	shortcutPane shortcutpane.ShortcutPane

	savedRequests []http.Request

	focusedPanel utils.Panel

	width, height int

	loadTestUpdates <-chan *http.LoadTestStats
	showHelpModal   bool
}

func SetupModel(db *storage.SQLiteStorage) Model {
	keys := keybindings.DefaultKeyMap()
	responsePane := responsepane.SetupResponsePane(keys)
	shortcutPane := shortcutpane.SetupShortcutPane(keys)

	m := Model{
		db:            db,
		keys:          keys,
		sidebarPane:   ui.NewSidebar(db, keys),
		requestPane:   requestpane.SetupRequestPane(db, keys),
		responsePane:  &responsePane,
		shortcutPane:  shortcutPane,
		focusedPanel:  utils.SidebarPanel,
		headerPane:    ui.SetupHeader(),
		showHelpModal: false,
		width:         80,
		height:        24,
	}
	m.setFocusedPanel(utils.SidebarPanel)
	m.applyLayout(calculateLayout(m.width, m.height))
	return m
}

func (m Model) Init() tea.Cmd {
	return m.sidebarPane.Init()
}
