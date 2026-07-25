package requestpane

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/owenHochwald/Volt/internal/storage"
	"github.com/owenHochwald/Volt/internal/ui/keybindings"
)

func TestPlainEnterDoesNotSubmitFromEditableField(t *testing.T) {
	pane := newTestRequestPane(t)
	pane.SetFocused(true)
	pane.FocusManager.Next()
	pane.URLInput.SetValue("https://example.com")

	updated, cmd := pane.Update(keyPress(tea.KeyEnter, "", 0))

	if updated.RequestInProgress {
		t.Fatal("plain enter started a request")
	}
	if cmd != nil {
		t.Fatal("plain enter produced a submit command")
	}
}

func TestModifiedEnterSubmitsFromAnyRequestField(t *testing.T) {
	tests := []struct {
		name string
		mod  tea.KeyMod
	}{
		{name: "control", mod: tea.ModCtrl},
		{name: "alt", mod: tea.ModAlt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pane := newTestRequestPane(t)
			pane.SetFocused(true)
			pane.FocusManager.Next()
			pane.URLInput.SetValue("https://example.com")

			updated, cmd := pane.Update(keyPress(tea.KeyEnter, "", tt.mod))

			if !updated.RequestInProgress {
				t.Fatal("modified enter did not start a request")
			}
			if cmd == nil {
				t.Fatal("modified enter did not produce a submit command")
			}
		})
	}
}

func TestPrintableCommandKeysRemainEditable(t *testing.T) {
	tests := []string{"?", "q", "h", "j", "k", "l"}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			pane := newTestRequestPane(t)
			pane.SetFocused(true)
			pane.FocusManager.Next()

			updated, _ := pane.Update(keyPress(rune(text[0]), text, 0))

			if got := updated.URLInput.Value(); got != text {
				t.Fatalf("URL value = %q, want %q", got, text)
			}
		})
	}
}

func TestSetFocusedBlursAndRestoresCurrentControl(t *testing.T) {
	pane := newTestRequestPane(t)
	pane.SetFocused(true)
	pane.FocusManager.Next()
	if !pane.URLInput.Focused() {
		t.Fatal("URL input is not focused before leaving panel")
	}

	pane.SetFocused(false)
	if pane.URLInput.Focused() {
		t.Fatal("URL input remained focused after leaving panel")
	}

	pane.SetFocused(true)
	if !pane.URLInput.Focused() {
		t.Fatal("URL input focus was not restored on panel entry")
	}
}

func newTestRequestPane(t *testing.T) RequestPane {
	t.Helper()

	db, err := storage.NewSQLiteStorage(filepath.Join(t.TempDir(), "volt.db"))
	if err != nil {
		t.Fatalf("create test storage: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close test storage: %v", err)
		}
	})
	return SetupRequestPane(db, keybindings.DefaultKeyMap())
}

func keyPress(code rune, text string, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{
		Code: code,
		Text: text,
		Mod:  mod,
	})
}
