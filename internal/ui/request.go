package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/owenHochwald/volt/internal/http"
)

const (
	focusMethod = iota
	focusURL
	focusName
	focusSubmit
)

var methodStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type RequestPane struct {
	methods       []string
	currentMethod int
	panelFocused  bool

	// ISSUE #4: Confusing Focus Index System
	// focusComponentIndex is used for BOTH method selection (index 0) AND text inputs.
	// This creates confusion because:
	// - Index 0 = method selection
	// - Index 1+ should map to inputs, but inputs[0] is the first input!
	// The mapping is: focusComponentIndex 1 -> inputs[0], focusComponentIndex 2 -> inputs[1]
	// SOLUTION: Use constants for clarity:
	//   const (
	//       focusMethod = iota  // 0
	//       focusURL            // 1 -> inputs[0]
	//       focusName           // 2 -> inputs[1]
	//       focusSubmit         // 3
	//   )
	focusComponentIndex int

	// ISSUE #5: Input Order Confusion
	// inputs[0] = URL, inputs[1] = Name
	// But in SetupRequestPane, you initialize them in reverse visual order.
	// This works but is confusing. Consider naming them explicitly:
	//   urlInput  textinput.Model
	//   nameInput textinput.Model
	// Or at minimum, add comments documenting the indices.
	inputs []textinput.Model // [0] = URL input, [1] = Name input

	// ISSUE #6: Unused Request Field
	// This field is defined but NEVER populated with actual data!
	// The inputs change, but m.request is never updated.
	// SOLUTION: Add a syncRequest() method that's called after input changes:
	//   func (m *RequestPane) syncRequest() {
	//       m.request.Method = m.methods[m.currentMethod]
	//       m.request.URL = m.inputs[0].Value()
	//       m.request.Name = m.inputs[1].Value()
	//   }
	request http.Request

	// MISSING: Consider adding these fields for future features:
	// headersExpanded bool
	// headers         []HeaderPair
	// queryParams     []QueryParam
	// bodyExpanded    bool
	// bodyInput       textarea.Model
	// validationError error
}

func (m *RequestPane) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RequestPane) SetFocused(focused bool) {
	m.panelFocused = focused
}

func (m *RequestPane) GetCurrentMethod() string {
	return m.methods[m.currentMethod]
}

func (m *RequestPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.panelFocused {
			return m, nil
		}
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case tea.KeyTab.String(), tea.KeyUp.String(), tea.KeyDown.String():
			s := msg.String()
			// ISSUE #8: Submit button logic is in the wrong place
			// This checks for "enter" key, but the switch only matches tab/up/down!
			// This condition will NEVER be true because "enter" won't match this case.
			// SOLUTION: Move submit logic to its own case or handle after focus cycling.
			if s == "enter" && m.focusComponentIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// ISSUE #9: No Focus/Blur Management
			// When changing focus, you should call .Blur() on the previously focused input
			// and .Focus() on the newly focused input. Currently, all inputs stay in their
			// initial state (only inputs[0] is focused from SetupRequestPane).
			// SOLUTION: Add focus management:
			//   m.blurCurrentComponent()
			//   // ... change focusComponentIndex ...
			//   m.focusCurrentComponent()

			// focus component cycling
			if s == "up" {
				m.focusComponentIndex--
			} else {
				m.focusComponentIndex++
			}

			// ISSUE #10: Off-by-one in focus wrapping
			// You have 4 focusable components: method(0), url(1), name(2), submit(3)
			// But len(m.inputs) = 2, so max should be 3, not len(m.inputs)
			// Current logic: max = 2, but you need to support up to index 3 for submit button
			// SOLUTION: Define a constant maxFocusIndex = 3 (or len(inputs) + 1 for submit)
			if m.focusComponentIndex > len(m.inputs) {
				m.focusComponentIndex = 0
			} else if m.focusComponentIndex < 0 {
				m.focusComponentIndex = len(m.inputs)
			}

		}

		// ISSUE #11: Focus-based routing is incomplete
		// This switch only handles case 0 (method) and default (inputs).
		// What about the submit button? It's never explicitly handled.
		// When focusComponentIndex == 3 (submit), it falls to default and updates ALL inputs,
		// which is inefficient and incorrect.
		// SOLUTION: Add explicit cases:
		//   case 0: // method selection
		//   case 1: // URL input - update only inputs[0]
		//   case 2: // Name input - update only inputs[1]
		//   case 3: // Submit button - handle Enter key
		switch m.focusComponentIndex {
		// method selection
		case 0:
			switch msg.String() {
			case tea.KeyRight.String(), "l":
				m.currentMethod = (m.currentMethod + 1) % len(m.methods)
			case tea.KeyLeft.String(), "h":
				m.currentMethod = (m.currentMethod - 1 + len(m.methods)) % len(m.methods)
			}
		default:
			// ISSUE #12: Updates ALL inputs regardless of focus
			// Even if only inputs[1] is focused, this updates both inputs[0] and inputs[1].
			// This works but is inefficient. Better to update only the focused input:
			//   case 1: m.inputs[0], cmd = m.inputs[0].Update(msg)
			//   case 2: m.inputs[1], cmd = m.inputs[1].Update(msg)
			cmd := m.updateInputs(msg)
			return m, cmd

		}

	}

	// ISSUE #13: Missing request sync
	// After any update that changes input values, you should sync the request model:
	// m.syncRequest() - to keep m.request up to date with current input values

	return m, nil
}

// ISSUE #14: View should use pointer receiver for consistency
func (m *RequestPane) View() string {
	methodDisplay := m.methods[m.currentMethod]

	// ISSUE #15: Method style should change based on focus
	// Currently, the method border changes when the PANEL is focused,
	// but it should also change when the METHOD COMPONENT specifically is focused
	// (i.e., when focusComponentIndex == 0).
	// SOLUTION:
	//   style := methodStyle
	//   if m.panelFocused && m.focusComponentIndex == 0 {
	//       style = style.BorderForeground(focusColor)
	//   } else {
	//       style = style.BorderForeground(unfocusColor)
	//   }
	style := methodStyle
	if m.panelFocused {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	// ISSUE #16: Inputs don't show which one is focused
	// When rendering inputs, you should indicate which one is currently focused.
	// The textinput.Model handles some of this internally, but only if Focus/Blur are called.
	// Since you're not calling Focus/Blur in Update(), the inputs never change appearance.
	// SOLUTION: Ensure Focus/Blur are called when changing focusComponentIndex
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	// ISSUE #17: Poor Layout Structure
	// Current layout:
	//   HTTP Method:
	//   [GET]
	//   <URL input>
	//   <Name input>
	//
	// Better layout would be:
	//   [GET] https://example.com
	//   Request Name: <name input>
	//
	// This puts method and URL on the same line, which is more natural and saves space.
	// SOLUTION: Use lipgloss.JoinHorizontal to combine method and URL on one line:
	//   firstLine := lipgloss.JoinHorizontal(lipgloss.Left, style.Render(methodDisplay), " ", m.inputs[0].View())
	//   s := docStyle.Render(firstLine + "\n" + "Name: " + m.inputs[1].View())

	// ISSUE #18: No submit button rendered
	// You define focusedButton and blurredButton at the top, but never render them!
	// The View() should include a submit button at the bottom.
	// SOLUTION: Add submit button rendering:
	//   var button string
	//   if m.focusComponentIndex == 3 {
	//       button = focusedButton
	//   } else {
	//       button = blurredButton
	//   }
	//   s += "\n" + button

	// ISSUE #19: No visual hierarchy
	// Everything has the same visual weight. Consider:
	// - Making the method + URL more prominent (larger, bolder)
	// - Adding section separators
	// - Using different colors for different sections
	// - Adding labels: "URL:", "Name:", etc.

	//s := docStyle.Render("HTTP Method: \n")
	s := style.Render(methodDisplay) + "\n"
	s += b.String() + "\n"

	// ISSUE #20: No help text
	// Users won't know what keys to press. Add help text:
	//   help := helpStyle.Render("tab: next • ←/→: change method • enter: submit • esc: back")
	//   s += "\n" + help

	return s
}

func (m *RequestPane) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
func SetupRequestPane() RequestPane {
	// ISSUE #21: Method order should match HTTP semantics
	// Consider ordering: GET, POST, PUT, PATCH, DELETE
	// This follows the CRUD pattern more naturally
	methods := []string{
		http.GET,
		http.POST,
		http.DELETE,
		http.PUT,
	}

	m := RequestPane{
		methods:             methods,
		currentMethod:       0,
		panelFocused:        false,
		focusComponentIndex: 0, // ISSUE #22: Should this start at 0 (method) or 1 (URL)?
		// Currently starts at method selection, but URL input is focused in the loop below.
		// This is inconsistent! Either:
		// 1. Start focusComponentIndex at 0 and don't focus any input initially
		// 2. Start focusComponentIndex at 1 and focus inputs[0] (URL)
		inputs: make([]textinput.Model, 2),
	}

	// ISSUE #23: Request model not initialized
	// You should initialize the request model with default values:
	//   m.request = *http.NewDefaultRequest()
	// or at minimum:
	//   m.request.Method = methods[0]  // Set default to GET

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		//t.Cursor.Style = cursorStyle
		t.CharLimit = 100 // ISSUE #24: CharLimit should match validation in http.Request

		switch i {
		case 0:
			t.Placeholder = "URL"
			// ISSUE #25: Only URL input is focused on init
			// This is inconsistent with focusComponentIndex = 0 (method selection)
			// If focusComponentIndex starts at 0 (method), no input should be focused yet.
			// If focusComponentIndex starts at 1 (URL), then this is correct.
			t.Focus()
			//t.PromptStyle = focusedStyle
			//t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Request Name"
			t.CharLimit = 300 // ISSUE #26: This should be 40 per http.Request validation
			// See internal/http/request.go:70 - max name length is 40, not 300
		}

		m.inputs[i] = t
	}

	// ISSUE #27: Consider adding width configuration
	// Text inputs should have Width set based on expected terminal size
	// Example:
	//   m.inputs[0].Width = 60  // URL input
	//   m.inputs[1].Width = 40  // Name input

	return m
}

// ============================================================================
// SUMMARY OF ISSUES IN THIS FILE
// ============================================================================
//
// CRITICAL (Fix First):
// - Issue #4, #22, #25: Focus management is inconsistent and confusing
// - Issue #6: Request model is never populated with data
// - Issue #8: Submit button logic will never execute
// - Issue #9: Focus/Blur never called, so inputs don't visually change
// - Issue #18: Submit button never rendered
//
// HIGH PRIORITY (Fix Next):
// - Issue #3, #7, #14: Inconsistent receiver types (value vs pointer)
// - Issue #10: Off-by-one error in focus wrapping
// - Issue #11, #12: Inefficient input updating
// - Issue #26: CharLimit doesn't match validation rules
//
// MEDIUM PRIORITY (Improve UX):
// - Issue #1, #2: Style duplication and static rendering
// - Issue #15, #16: Focus indicators not clear enough
// - Issue #17: Poor layout structure
// - Issue #19: No visual hierarchy
// - Issue #20: Missing help text
//
// LOW PRIORITY (Nice to Have):
// - Issue #5: Input array vs named fields
// - Issue #13: Request sync not called
// - Issue #21: Method ordering
// - Issue #23, #24, #27: Better defaults and configuration
//
// See docs/REQUEST_PANE_GUIDE.md for detailed implementation guidance
// ============================================================================
