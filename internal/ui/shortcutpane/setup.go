package shortcutpane

// SetupShortcutPane sets up the shortcut pane for use
func SetupShortcutPane() ShortcutPane {
	return ShortcutPane{
		activeTab: 0,
		height:    30,
		width:     40,
		Focused:   false,
		tabs:      GetShortcutTabs(),
	}

}
