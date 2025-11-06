package http

import (
	"fmt"
	"net/http"
	"time"
)

type ResultMsg struct {
	Response *Response
}

func (r *ResultMsg) String() string {
	if r.Response.Error != "" {
		// TODO: add styles
		return r.Response.Error
	}
	// TODO: add styles
	return fmt.Sprintf("%s that took: %d ms", r.Response.StatusCode, r.Response.Duration)
}

type Response struct {
	StatusCode int           `json:"status_code"`
	Status     string        `json:"status,omitempty"`
	Headers    http.Header   `json:"headers,omitempty"`
	Body       string        `json:"body,omitempty"`
	Duration   time.Duration `json:"duration,omitempty"`
	Error      string        `json:"error,omitempty"`
}
