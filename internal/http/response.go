package http

import (
	"net/http"
	"time"
)

type Response struct {
	StatusCode int           `json:"status_code"`
	Status     string        `json:"status,omitempty"`
	Headers    http.Header   `json:"headers,omitempty"`
	Body       string        `json:"body,omitempty"`
	Duration   time.Duration `json:"duration,omitempty"`
	Error      string        `json:"error,omitempty"`
}
