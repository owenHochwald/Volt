package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/storage"
)

type RequestsLoadingMsg struct {
	Requests []http.Request
	Err      error
}

type RequestSavedMsg struct {
	Request *http.Request
	Err     error
}

func SaveRequestCmd(db *storage.SQLiteStorage, request *http.Request) tea.Cmd {
	return func() tea.Msg {
		err := db.Save(request)
		return RequestSavedMsg{
			Request: request,
			Err:     err,
		}
	}
}

func LoadRequestsCmd(db *storage.SQLiteStorage) tea.Cmd {
	return func() tea.Msg {
		requests, err := db.Load()
		return RequestsLoadingMsg{
			Requests: requests,
			Err:      err,
		}
	}
}
