package storage

import (
	"github.com/owenHochwald/volt/internal/http"
)

type Storage interface {
	Save(requests *http.Request) error
	Load() ([]http.Request, error)
	Delete(id int64) error
	GetAllURLs() ([]string, error)
}
