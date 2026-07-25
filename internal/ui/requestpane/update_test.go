package requestpane

import (
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/owenHochwald/Volt/internal/storage"
	"github.com/owenHochwald/Volt/internal/ui"
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

func TestTabAndShiftTabOnlyNavigateFields(t *testing.T) {
	pane := newTestRequestPane(t)
	pane.SetFocused(true)

	updated, _ := pane.Update(keyPress(tea.KeyTab, "", 0))
	if got := updated.FocusManager.CurrentIndex(); got != int(FieldURL) {
		t.Fatalf("tab focus index = %d, want URL field", got)
	}

	updated, _ = updated.Update(keyPress(tea.KeyTab, "", tea.ModShift))
	if got := updated.FocusManager.CurrentIndex(); got != int(FieldMethodSelector) {
		t.Fatalf("shift+tab focus index = %d, want method selector", got)
	}
}

func TestArrowKeysDoNotChangeFocusedField(t *testing.T) {
	pane := newTestRequestPane(t)
	pane.SetFocused(true)
	pane.FocusManager.Next()
	pane.URLInput.SetValue("https://example.com")
	pane.URLInput.CursorEnd()

	updated, _ := pane.Update(keyPress(tea.KeyLeft, "", 0))

	if got := updated.FocusManager.CurrentIndex(); got != int(FieldURL) {
		t.Fatalf("left arrow changed focus to index %d", got)
	}
	if got, want := updated.URLInput.Position(), len("https://example.com")-1; got != want {
		t.Fatalf("left arrow cursor position = %d, want %d", got, want)
	}
}

func TestInvalidHeadersBlockSubmissionAndSurfaceError(t *testing.T) {
	pane := newTestRequestPane(t)
	pane.SetFocused(true)
	pane.URLInput.SetValue("https://example.com")
	pane.Headers.SetValue("missing delimiter")

	updated, cmd := pane.Update(keyPress(tea.KeyEnter, "", tea.ModCtrl))

	if updated.RequestInProgress {
		t.Fatal("invalid headers started a request")
	}
	if cmd == nil {
		t.Fatal("invalid headers did not produce a notification")
	}
	msg, ok := cmd().(ui.NotificationMsg)
	if !ok {
		t.Fatalf("message type = %T, want ui.NotificationMsg", cmd())
	}
	if !strings.Contains(strings.ToLower(msg.Notification.Text), "header") {
		t.Fatalf("notification = %q, want header error", msg.Notification.Text)
	}
}

func TestInvalidLoadTestConfigBlocksStartAndSurfacesError(t *testing.T) {
	pane := newTestRequestPane(t)
	pane.SetFocused(true)
	pane.toggleLoadTestMode()
	pane.URLInput.SetValue("https://example.com")
	pane.LoadTestConcurrency.SetValue("not-a-number")

	updated, cmd := pane.Update(keyPress(tea.KeyEnter, "", tea.ModCtrl))

	if updated.RequestInProgress {
		t.Fatal("invalid load-test config remained in progress")
	}
	if cmd == nil {
		t.Fatal("invalid load-test config did not produce a notification")
	}
	msg, ok := cmd().(ui.NotificationMsg)
	if !ok {
		t.Fatalf("message type = %T, want ui.NotificationMsg", cmd())
	}
	if !strings.Contains(strings.ToLower(msg.Notification.Text), "concurrency") {
		t.Fatalf("notification = %q, want concurrency error", msg.Notification.Text)
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
