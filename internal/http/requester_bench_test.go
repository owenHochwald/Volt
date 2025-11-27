package http

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

// setupServer starts a high-performance fasthttp server on a random port
func setupServer(latency time.Duration) (string, func()) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	handler := func(ctx *fasthttp.RequestCtx) {
		if latency > 0 {
			time.Sleep(latency)
		}
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	server := &fasthttp.Server{
		Handler:               handler,
		NoDefaultServerHeader: true, // Reduce overhead
		DisableKeepalive:      false,
	}

	go func() {
		if err := server.Serve(ln); err != nil {
		}
	}()

	return "http://" + ln.Addr().String(), func() {
		ln.Close()
	}
}

// BenchmarkEngine_Throughput measures maximum request throughput at different concurrency levels.
//
// Usage:
//
//	go test -bench=BenchmarkEngine_Throughput -benchmem ./internal/http
//
// To run a soak test (e.g., 10 seconds):
//
//	go test -bench=BenchmarkEngine_Throughput -benchtime=10s ./internal/http
func BenchmarkEngine_Throughput(b *testing.B) {
	// Start a server with 0 latency to measure pure engine overhead
	serverURL, stop := setupServer(0)
	defer stop()

	req := NewDefaultRequest()
	req.URL = serverURL
	req.Method = "GET"

	concurrencyLevels := []int{10, 50, 100, 500, 1000}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			batchSize := 10000
			if b.N < batchSize {
				batchSize = b.N
			}

			b.ReportAllocs()
			b.ResetTimer()

			var totalRequestsProcessed int64

			// We loop until we have processed b.N requests
			for totalRequestsProcessed < int64(b.N) {
				currentBatch := batchSize
				remaining := int64(b.N) - totalRequestsProcessed
				if remaining < int64(batchSize) {
					currentBatch = int(remaining)
				}

				updates := make(chan *LoadTestStats, 100)

				// Drain stats asynchronously
				go func() {
					for range updates {
					}
				}()

				config := &JobConfig{
					Request:       req,
					Concurrency:   concurrency,
					TotalRequests: currentBatch,
					Timeout:       5 * time.Second,
				}

				config.Run(updates)

				totalRequestsProcessed += int64(currentBatch)
			}

			b.StopTimer()

			elapsed := b.Elapsed().Seconds()
			rps := float64(b.N) / elapsed
			b.ReportMetric(rps, "req/s")
		})
	}
}

// BenchmarkEngine_Latency handles a scenario where the server is slow (simulating real world).
// This tests how efficient the engine is at waiting.
func BenchmarkEngine_Latency(b *testing.B) {
	serverURL, stop := setupServer(50 * time.Millisecond)
	defer stop()

	req := NewDefaultRequest()
	req.URL = serverURL
	req.Method = "GET"

	concurrency := 200

	b.ResetTimer()
	b.ReportAllocs()

	var totalRequestsProcessed int64
	batchSize := 1000

	for totalRequestsProcessed < int64(b.N) {
		currentBatch := batchSize
		remaining := int64(b.N) - totalRequestsProcessed
		if remaining < int64(batchSize) {
			currentBatch = int(remaining)
		}

		updates := make(chan *LoadTestStats, 100)
		go func() {
			for range updates {
			}
		}()

		config := &JobConfig{
			Request:       req,
			Concurrency:   concurrency,
			TotalRequests: currentBatch,
			Timeout:       10 * time.Second, // Increased timeout for latency
		}

		config.Run(updates)
		totalRequestsProcessed += int64(currentBatch)
	}

	elapsed := b.Elapsed().Seconds()
	rps := float64(b.N) / elapsed
	b.ReportMetric(rps, "req/s")
}
