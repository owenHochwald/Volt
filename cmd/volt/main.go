package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/owenHochwald/volt/internal/app"
	"github.com/owenHochwald/volt/internal/storage"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(homeDir, ".volt", "volt.db")
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		fmt.Printf("Error connecting to database: %v", err)
		return
	}
	defer store.Close()

	p := tea.NewProgram(app.SetupModel(store), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
