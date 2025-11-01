package http

type Response struct {
	StatusCode int               `json:"status_code"`
	Status     string            `json:"status,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
	Duration   int64             `json:"duration,omitempty"`
	Error      string            `json:"error,omitempty"`
}
