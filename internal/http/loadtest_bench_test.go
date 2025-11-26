package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkLoadTest_10_Workers(b *testing.B) {
	benchmarkLoadTest(b, 10, 1000)
}

func BenchmarkLoadTest_100_Workers(b *testing.B) {
	benchmarkLoadTest(b, 100, 1000)
}
func BenchmarkLoadTest_500_Workers(b *testing.B) {
	benchmarkLoadTest(b, 500, 1000)
}

func benchmarkLoadTest(b *testing.B, concurrency, totalRequests int) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	config := &LoadTestConfig{
		Request: &Request{
			Method: "GET",
			URL:    server.URL,
		},
		Concurrency:   concurrency,
		TotalRequests: totalRequests,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updates := make(chan LoadTestStats, 10)
		go RunLoadTest(context.Background(), config, updates)

		// Drain updates
		for range updates {
		}
	}
}
