package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Header struct {
	text        string
	progress    progress.Model
	percentFill int
}

func (h Header) Init() tea.Cmd {
	return nil
}

func (h Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return h, nil
}

func (h Header) View() string {
	return ""
}

func SetupHeader(text string) *Header {
	return &Header{text: text}
}
