package http

import "net/http"

type Response struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
	Body       string      `json:"body,omitempty"`
	Duration   int64       `json:"duration,omitempty"`
	Error      string      `json:"error,omitempty"`
}
