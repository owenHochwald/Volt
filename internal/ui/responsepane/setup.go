package responsepane

import "github.com/charmbracelet/bubbles/viewport"

// SetupResponsePane creates and initializes a new ResponsePane with default values
func SetupResponsePane() ResponsePane {
	return ResponsePane{
		viewport:   viewport.New(20, 10),
		width:      20,
		height:     30,
		activeTab:  int(TabBody), // Start on Body tab
		isLoadTest: false,
	}
}
