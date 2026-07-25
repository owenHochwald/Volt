package app

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/owenHochwald/Volt/internal/storage"
	"github.com/owenHochwald/Volt/internal/ui/keybindings"
	"github.com/owenHochwald/Volt/internal/utils"
)

func TestQuestionMarkOpensFocusedPanelHelp(t *testing.T) {
	model := newTestModel(t)

	updated, _ := model.Update(appKeyPress('?', "?", 0))
	got := updated.(Model)

	if !got.showHelpModal {
		t.Fatal("question mark did not open help")
	}
	if context := got.shortcutPane.ActiveContext(); context != keybindings.ContextSidebar {
		t.Fatalf("help context = %s, want %s", context, keybindings.ContextSidebar)
	}
}

func TestQuestionMarkAndQRemainEditableInRequestFields(t *testing.T) {
	for _, text := range []string{"?", "q"} {
		t.Run(text, func(t *testing.T) {
			model := newTestModel(t)
			model.setFocusedPanel(utils.RequestPanel)
			model.requestPane.FocusManager.Next()

			updated, _ := model.Update(appKeyPress(rune(text[0]), text, 0))
			got := updated.(Model)

			if got.showHelpModal {
				t.Fatalf("%q opened help while editing", text)
			}
			if value := got.requestPane.URLInput.Value(); value != text {
				t.Fatalf("URL value = %q, want %q", value, text)
			}
		})
	}
}

func TestF1AlwaysOpensGlobalHelp(t *testing.T) {
	model := newTestModel(t)
	model.setFocusedPanel(utils.RequestPanel)
	model.requestPane.FocusManager.Next()

	updated, _ := model.Update(appKeyPress(tea.KeyF1, "", 0))
	got := updated.(Model)

	if !got.showHelpModal {
		t.Fatal("F1 did not open help while editing")
	}
	if context := got.shortcutPane.ActiveContext(); context != keybindings.ContextGlobal {
		t.Fatalf("help context = %s, want %s", context, keybindings.ContextGlobal)
	}
}

func TestPanelNavigationBlursRequestControl(t *testing.T) {
	model := newTestModel(t)
	model.setFocusedPanel(utils.RequestPanel)
	model.requestPane.FocusManager.Next()
	if !model.requestPane.URLInput.Focused() {
		t.Fatal("URL input is not focused before panel navigation")
	}

	updated, _ := model.Update(appKeyPress('l', "", tea.ModAlt))
	got := updated.(Model)

	if got.focusedPanel != utils.ResponsePanel {
		t.Fatalf("focused panel = %d, want response", got.focusedPanel)
	}
	if got.requestPane.URLInput.Focused() {
		t.Fatal("URL input remained focused after leaving request panel")
	}
}

func newTestModel(t *testing.T) Model {
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
	return SetupModel(db)
}

func appKeyPress(code rune, text string, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{
		Code: code,
		Text: text,
		Mod:  mod,
	})
}
