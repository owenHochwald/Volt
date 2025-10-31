package http

type Response struct {
	StatusCode int
	Status     string
	Headers    map[string]string
	Body       string
	Duration   int64
	Error      error
}
