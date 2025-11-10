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
	m.focused = true
}

func (m *MethodSelector) Blur() {
	m.focused = false

}

func (m *MethodSelector) Current() string {
	return m.methods[m.currentMethod]
}

func (m *MethodSelector) Next() {
	m.currentMethod = (m.currentMethod + 1) % len(m.methods)
}

func (m *MethodSelector) Prev() {
	m.currentMethod--
	if m.currentMethod < 0 {
		m.currentMethod = len(m.methods) - 1
	}
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

	if m.focused {
		methodStyle = methodStyle.BorderForeground(focusColor)
	} else {
		methodStyle = methodStyle.BorderForeground(unfocusColor)
	}

	return methodStyle
}

func NewMethodSelector(methods []string) *MethodSelector {
	return &MethodSelector{
		methods:       methods,
		currentMethod: 0,
		focused:       false,
	}
}
