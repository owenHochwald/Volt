package requestpane

import (
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/owenHochwald/volt/internal/storage"
	"github.com/owenHochwald/volt/internal/ui"
)

// SetupRequestPane creates and initializes a new RequestPane
func SetupRequestPane(db *storage.SQLiteStorage) RequestPane {
	methodSelector := ui.NewMethodSelector()

	// Use factories for text inputs
	urlInput := NewURLInput()
	nameInput := NewNameInput()

	// Create text areas
	headers := NewHeadersTextArea()
	body := NewBodyTextArea()

	submitButton := ui.NewSubmitButton()

	// Load test inputs using factory
	ltConcurrency := NewLoadTestInput("100", 5, 15)
	ltTotalReqs := NewLoadTestInput("10000", 10, 15)
	ltQPS := NewLoadTestInput("0 (unlimited)", 10, 15)
	ltTimeout := NewLoadTestInput("30s", 10, 15)

	// Initialize with normal mode
	normalMode := &NormalMode{}

	m := RequestPane{
		MethodSelector:      methodSelector,
		URLInput:            &urlInput,
		NameInput:           &nameInput,
		Headers:             &headers,
		Body:                &body,
		SubmitButton:        submitButton,
		Client:              http.InitClient(0, false),
		Stopwatch:           stopwatch.NewWithInterval(10 * time.Millisecond),
		Request:             http.NewDefaultRequest(),
		DB:                  db,
		LoadTestConcurrency: &ltConcurrency,
		LoadTestTotalReqs:   &ltTotalReqs,
		LoadTestQPS:         &ltQPS,
		LoadTestTimeout:     &ltTimeout,
		LoadTestMode:        false,
		currentMode:         normalMode,
	}

	m.FocusManager = normalMode.GetFocusManager(&m)

	return m
}
