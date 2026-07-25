package shortcutpane

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/key"
	"github.com/charmbracelet/x/ansi"
	"github.com/owenHochwald/Volt/internal/ui/keybindings"
)

func TestHelpTabLabelsFollowConfiguredBinding(t *testing.T) {
	keys := keybindings.DefaultKeyMap()
	keys.DirectHelpTab = key.NewBinding(key.WithKeys("a", "b", "c", "d"))
	pane := SetupShortcutPane(keys)

	rendered := pane.renderTabs()
	for _, expected := range []string{"[a] Global", "[b] Sidebar", "[c] Request", "[d] Response"} {
		if !strings.Contains(rendered, expected) {
			t.Errorf("help tab bar does not contain %q: %s", expected, rendered)
		}
	}
}

func TestHelpDoesNotWrapLongRegistryKeys(t *testing.T) {
	pane := SetupShortcutPane(keybindings.DefaultKeyMap())
	pane.SetWidth(72)
	pane.SetHeight(30)
	pane.SetContext(keybindings.ContextRequest)

	rendered := ansi.Strip(pane.View())
	if !strings.Contains(rendered, "ctrl+j") {
		t.Fatalf("long submit binding wrapped unexpectedly:\n%s", rendered)
	}
}
