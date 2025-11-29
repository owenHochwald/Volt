package http

import (
	"time"

	"github.com/valyala/fasthttp"
)

type HeaderEntry struct {
	Key   []byte
	Value []byte
}

type FastRequest struct {
	Method  []byte
	URL     []byte
	Headers []HeaderEntry
	Body    []byte
}

type FastClient struct {
	client  *fasthttp.Client
	timeout time.Duration
}

func NewFastClient(timeout time.Duration, s *JobConfig) *FastClient {
	return &FastClient{
		timeout: timeout,
		client: &fasthttp.Client{
			MaxConnsPerHost:     max(500, s.Concurrency),
			MaxIdleConnDuration: 10 * time.Second,
			ReadTimeout:         s.Timeout,
			WriteTimeout:        s.Timeout,
			MaxConnDuration:     0,
		},
	}
}

func (f *FastClient) Do(
	fr *FastRequest,
	req *fasthttp.Request,
	res *fasthttp.Response,
) (status int, contentLen int64, err error) {
	req.Reset()
	res.Reset()

	// zero allocation setters
	req.Header.SetMethodBytes(fr.Method)
	req.SetRequestURIBytes(fr.URL)
	for _, entry := range fr.Headers {

		req.Header.SetBytesKV(entry.Key, entry.Value)
	}
	if fr.Body != nil {
		req.SetBodyRaw(fr.Body)
	}

	err = f.client.DoTimeout(req, res, f.timeout)
	if err != nil {
		return 0, 0, err
	}

	return res.StatusCode(), int64(res.Header.ContentLength()), nil
}
