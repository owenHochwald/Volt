package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
)

type MethodSelector struct {
	methods       []string
	currentMethod int
	focused       bool
}

func (m *MethodSelector) Focus() {
	methodStyleBase.BorderForeground(focusColor)
	m.focused = true
}

func (m *MethodSelector) Blur() {
	methodStyleBase.BorderForeground(unfocusColor)
	m.focused = false

}

func (m *MethodSelector) Current() string {
	return m.methods[m.currentMethod]
}

func (m *MethodSelector) Next() {
	m.currentMethod = (m.currentMethod + 1) % len(m.methods)
}

func (m *MethodSelector) Prev() {
	m.currentMethod = (m.currentMethod - 1) % len(m.methods)
}

func (m *MethodSelector) GetStyle() lipgloss.Style {
	currMethod := m.Current()
	var methodStyle lipgloss.Style

	switch currMethod {
	case http.GET:
		methodStyle = getMethodStyle
	case http.POST:
		methodStyle = postMethodStyle
	case http.PUT:
		methodStyle = putMethodStyle
	case http.PATCH:
		methodStyle = patchMethodStyle
	case http.DELETE:
		methodStyle = deleteMethodStyle
	default:
		methodStyle = methodStyleBase
	}

	if m.focused == focusMethod {
		methodStyle = methodStyle.BorderForeground(focusColor)
	} else {
		methodStyle = methodStyle.BorderForeground(unfocusColor)
	}

	return methodStyle
}
