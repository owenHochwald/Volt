package keybindings

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Matches checks if a key message matches a binding
func Matches(msg tea.KeyPressMsg, binding key.Binding) bool {
	if !binding.Enabled() {
		return false
	}
	if key.Matches(msg, binding) {
		return true
	}
	keystroke := msg.Keystroke()
	for _, registered := range binding.Keys() {
		if keystroke == registered {
			return true
		}
	}
	return false
}
