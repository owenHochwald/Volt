package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/owenHochwald/volt/internal/http"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type SQLiteStorage struct {
	db *sql.DB
}

func serializeHeaders(headers map[string]string) (string, error) {
	return "", nil
}
func deserializeHeaders(jsonStr string) (map[string]string, error) {
	return nil, nil
}

func runMigrations(db *sql.DB) error {
	// Set the embedded filesystem for goose
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLiteStorage{db: db}, nil

}

func (db *SQLiteStorage) Close() error {
	return db.db.Close()
}

func (db *SQLiteStorage) Save(request http.Request) error {
	//TODO implement me
	panic("implement me")
}

func (db *SQLiteStorage) Load() ([]http.Request, error) {
	//TODO implement me
	panic("implement me")
}

func (db *SQLiteStorage) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}
