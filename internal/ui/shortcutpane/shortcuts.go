package shortcutpane

// Shortcut represents a keyboard shortcut and its description
type Shortcut struct {
	Key         string
	Description string
}

// ShortcutTab represents a category of shortcuts
type ShortcutTab struct {
	Name      string
	Shortcuts []Shortcut
}

// GetShortcutTabs returns all shortcut tabs for the help modal
func GetShortcutTabs() []ShortcutTab {
	return []ShortcutTab{
		{
			Name: "Global",
			Shortcuts: []Shortcut{
				{"?", "Show this help"},
				{"Shift+Tab", "Cycle panels"},
				{"q, Ctrl+C", "Quit"},
			},
		},
		{
			Name: "Sidebar",
			Shortcuts: []Shortcut{
				{"Enter/Space", "Select request"},
				{"d", "Delete request"},
				{"/", "Filter requests"},
				{"j/k", "Navigate up/down"},
			},
		},
		{
			Name: "Request",
			Shortcuts: []Shortcut{
				{"Enter", "Send request"},
				{"Tab", "Next field"},
				{"h/l", "Change method"},
				{"Ctrl+S", "Save request"},
				{"Alt+L", "Toggle load test"},
			},
		},
		{
			Name: "Response",
			Shortcuts: []Shortcut{
				{"1,2,3", "Jump to tab"},
				{"h/l", "Navigate tabs"},
				{"y/Y", "Copy response"},
				{"j/k", "Scroll"},
			},
		},
	}
}
