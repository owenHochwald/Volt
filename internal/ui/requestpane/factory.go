package requestpane

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

// TextInputConfig holds configuration for creating a textinput.Model
type TextInputConfig struct {
	Placeholder string
	CharLimit   int
	Width       int
	Value       string
	Suggestions []string
}

// NewConfiguredTextInput creates a textinput.Model with the given configuration
func NewConfiguredTextInput(config TextInputConfig) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = config.Placeholder
	ti.CharLimit = config.CharLimit
	ti.Width = config.Width

	if config.Value != "" {
		ti.SetValue(config.Value)
	}

	if len(config.Suggestions) > 0 {
		ti.SetSuggestions(config.Suggestions)
	}

	return ti
}

// NewURLInput creates a pre-configured URL input field
func NewURLInput() textinput.Model {
	// call repository method to get suggestions (all previous urls)!

	return NewConfiguredTextInput(TextInputConfig{
		Value:     "http://localhost:",
		CharLimit: 40,
		Width:     60,
		Suggestions: []string{
			"http://localhost:8080",
			"http://localhost:8081",
			"http://localhost:8082",
		},
	})
}

// NewNameInput creates a pre-configured name input field
func NewNameInput() textinput.Model {
	return NewConfiguredTextInput(TextInputConfig{
		Placeholder: "My new awesome request..",
		CharLimit:   40,
		Width:       60,
	})
}

// NewLoadTestInput creates a pre-configured load test input field
func NewLoadTestInput(placeholder string, charLimit, width int) textinput.Model {
	return NewConfiguredTextInput(TextInputConfig{
		Placeholder: placeholder,
		CharLimit:   charLimit,
		Width:       width,
	})
}

// NewHeadersTextArea creates a pre-configured headers textarea
func NewHeadersTextArea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Content-Type = multipart/form-data,\nAuthorization= Bearer ...,"
	return ta
}

// NewBodyTextArea creates a pre-configured body textarea
func NewBodyTextArea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "key = value,\nname = volt,\nversion=1.0"
	return ta
}
