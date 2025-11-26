package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestRunLoadTest(t *testing.T) {
	req := NewRequestWithParams("POST", "http://localhost:8080/pong")
	config := &LoadTestConfig{
		Request:       req,
		Concurrency:   5,
		TotalRequests: 20,
	}

	updates := make(chan LoadTestStats, 10)
	ctx := context.Background()

	go RunLoadTest(ctx, config, updates)

	for stats := range updates {
		fmt.Printf("Progress: %d/%d completed\n",
			stats.CompletedRequests, config.TotalRequests)
	}
}

func TestBasicLoadTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	config := &LoadTestConfig{
		Request: &Request{
			Method: "GET",
			URL:    server.URL,
		},
		Concurrency:   50,
		TotalRequests: 100,
		Timeout:       10 * time.Second,
	}

	updates := make(chan LoadTestStats, 10)

	go RunLoadTest(context.Background(), config, updates)

	var finalStats LoadTestStats
	for stats := range updates {
		finalStats = stats
		fmt.Printf("Progress: %d/%d completed\n",
			stats.CompletedRequests, config.TotalRequests)
	}

	if finalStats.CompletedRequests != config.TotalRequests {
		t.Errorf("Expected 50 completed, got %d", finalStats.CompletedRequests)
	}

	// check the final snapshot
	assert.Equal(t, config.TotalRequests, finalStats.CompletedRequests)

}

func TestExtremeLoadTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	config := &LoadTestConfig{
		Request: &Request{
			Method: "GET",
			URL:    server.URL,
		},
		Concurrency:   500,
		TotalRequests: 5000,
		Timeout:       10 * time.Second,
	}

	updates := make(chan LoadTestStats, 10)

	go RunLoadTest(context.Background(), config, updates)

	var finalStats LoadTestStats
	for stats := range updates {
		finalStats = stats
		fmt.Printf("Progress: %d/%d completed\n",
			stats.CompletedRequests, config.TotalRequests)
	}

	if finalStats.CompletedRequests != config.TotalRequests {
		t.Errorf("Expected 50 completed, got %d", finalStats.CompletedRequests)
	}

	// check the final snapshot
	assert.Equal(t, config.TotalRequests, finalStats.CompletedRequests)
}
