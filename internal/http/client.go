package http

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	TIMEOUT = 10 * time.Second
)

type Client struct {
	Timeout   time.Duration
	RoundTrip bool

	Transport *http.Transport
	Client    *http.Client
}

//round trip logic
// Source - https://stackoverflow.com/questions/30526946/time-http-response-in-go
// Posted by Devatoria
// Retrieved 11/4/2025, License - CC-BY-SA 4.0

//func Get() int {
//	start := time.Now()
//	result, err := http.Get("http://www.google.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer result.Body.Close()
//	elapsed := time.Since(start).Seconds()
//	log.Println(elapsed)
//
//	return result.StatusCode
//}

// just response time logic:
// Source - https://stackoverflow.com/questions/30526946/time-http-response-in-go
// Posted by icza
// Retrieved 11/4/2025, License - CC-BY-SA 4.0

//conn, err := net.Dial("tcp", "google.com:80")
//if err != nil {
//panic(err)
//}
//defer conn.Close()
//conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
//
//start := time.Now()
//oneByte := make([]byte,1)
//_, err = conn.Read(oneByte)
//if err != nil {
//panic(err)
//}
//log.Println("First byte:", time.Since(start))
//
//_, err = ioutil.ReadAll(conn)
//if err != nil {
//panic(err)
//}
//log.Println("Everything:", time.Since(start))

func (c *Client) makeCustomRequest(req *Request) (*http.Response, error) {
	payload := strings.NewReader(req.Body)
	customReq, err := http.NewRequest(req.Method, req.URL, payload)
	if err != nil {
		return nil, err
	}

	for key, val := range req.Headers {
		customReq.Header.Set(key, val)
	}

	if c.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()
		customReq = customReq.WithContext(ctx)
	}

	res, err := c.Client.Do(customReq)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (c *Client) Send(req *Request) (*Response, error) {
	res, err := c.makeCustomRequest(req)
	if err != nil {
		return nil, err
	}

	// close body from earlier call
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Body:       string(body),
		Headers:    res.Header,
	}, nil
}

//Measure request duration
//Parse response headers and body
//Handle errors gracefully (network errors, timeouts, invalid URLs)
//Add proper context for cancellation
//Write unit tests with httptest
//
//Definition

func InitClient(timeout time.Duration, roundTrip bool) *Client {
	if timeout == 0 {
		timeout = TIMEOUT
	}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	return &Client{
		Timeout:   timeout,
		RoundTrip: roundTrip,
		Client:    client,
		Transport: tr,
	}

}
