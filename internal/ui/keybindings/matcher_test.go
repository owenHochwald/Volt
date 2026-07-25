package keybindings

import (
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func TestMatchesUsesStructuredV2KeyData(t *testing.T) {
	binding := key.NewBinding(key.WithKeys("alt+l"))
	tests := []struct {
		name string
		key  tea.KeyPressMsg
		want bool
	}{
		{
			name: "legacy alt event without text",
			key:  tea.KeyPressMsg(tea.Key{Code: 'l', Mod: tea.ModAlt}),
			want: true,
		},
		{
			name: "mac option event with associated text",
			key:  tea.KeyPressMsg(tea.Key{Code: 'l', BaseCode: 'l', Text: "¬", Mod: tea.ModAlt}),
			want: true,
		},
		{
			name: "printable text without alt modifier",
			key:  tea.KeyPressMsg(tea.Key{Code: 'l', Text: "l"}),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matches(tt.key, binding); got != tt.want {
				t.Fatalf("Matches(%q) = %t, want %t", tt.key.Keystroke(), got, tt.want)
			}
		})
	}
}
