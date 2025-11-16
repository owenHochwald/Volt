package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Header struct {
	text      string
	progress  progress.Model
	percent   float64
	expanding bool
	width     int
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (h *Header) Init() tea.Cmd {
	return tickCmd()
}

func (h *Header) SetWidth(width int) {
	h.width = width
}

func (h *Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if h.expanding {
			h.percent += 0.01
			if h.percent > 1.0 {
				h.percent = 1.0
				h.expanding = false
			}
		} else {
			h.percent -= 0.01
			if h.percent < 0.0 {
				h.percent = 0.0
				h.expanding = true
			}
		}
		return h, tickCmd()
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.progress.Width = msg.Width - 4
	}

	return h, tickCmd()
}

func (h *Header) View() string {
	s := h.text + " - TUI HTTP Client - v0.1 [?] Help  [q] Quit"
	return s + "\n" + h.progress.ViewAs(h.percent)
}

func SetupHeader(text string) *Header {
	return &Header{
		text:  text,
		width: 80,
		progress: progress.New(
			progress.WithScaledGradient("#7C3AED", "#EC4899"),
			progress.WithoutPercentage()),
		expanding: true,
	}
}
