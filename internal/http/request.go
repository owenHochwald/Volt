package http

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
)

type Request struct {
	ID      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

func NewRequest() *Request {
	return &Request{
		Method: "",
	}
}

func NewRequestWithParams(method, url string) *Request {
	return &Request{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
	}
}
