package ui

import tea "github.com/charmbracelet/bubbletea"

type SubmitButton struct {
	focused bool
}

func NewSubmitButton() *SubmitButton {
	return &SubmitButton{focused: false}
}

func (s *SubmitButton) Focus() tea.Cmd {
	s.focused = true
	return nil
}

func (s *SubmitButton) Blur() {
	s.focused = false
}

func (s *SubmitButton) IsFocused() bool {
	return s.focused
}
