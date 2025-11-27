package http

import (
	"time"

	"github.com/valyala/fasthttp"
)

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

func (f *FastClient) Do(req *Request) (*JobResult, error) {
	r := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(r)

	// set specifics of the request
	r.SetRequestURI(req.URL)
	r.Header.SetMethod(req.Method)
	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}
	if req.Body != "" {
		r.SetBodyString(req.Body)
	}

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	start := time.Now()
	err := f.client.DoTimeout(r, res, f.timeout)
	totalDuration := time.Since(start)

	if err != nil {
		return &JobResult{
			err:      err,
			duration: totalDuration,
		}, err
	}

	return &JobResult{
		err:           nil,
		statusCode:    res.StatusCode(),
		duration:      totalDuration,
		contentLength: int64(res.Header.ContentLength()),
	}, nil

}
