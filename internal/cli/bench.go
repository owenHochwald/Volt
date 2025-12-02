package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/owenHochwald/volt/internal/http"
)

// RunBench executes the load test with given configuration
func RunBench(config *BenchConfig) error {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Build Request object
	req := &http.Request{
		Method:  config.Method,
		URL:     config.URL,
		Headers: config.Headers,
		Body:    config.Body,
	}

	// Calculate total requests if duration-based
	totalRequests := config.TotalRequests
	if totalRequests == 0 {
		// Estimate based on duration and concurrency
		// Assume each worker can handle ~1000 req/s for local tests
		// This is conservative for high-performance scenarios
		estimatedReqsPerWorkerPerSec := 1000
		totalRequests = config.Concurrency * int(config.Duration.Seconds()) * estimatedReqsPerWorkerPerSec

		// Ensure minimum of 1 request
		if totalRequests < 1 {
			totalRequests = 1
		}
	}

	// Build JobConfig
	jobConfig := &http.JobConfig{
		Request:       req,
		Concurrency:   config.Concurrency,
		TotalRequests: totalRequests,
		Timeout:       config.Timeout,
		QPS:           float64(config.RateLimit),
		StreamUpdates: false, // CLI only needs final result
	}

	// Create update channel - larger buffer to prevent blocking in requester
	updates := make(chan *http.LoadTestStats, 1000)

	// Handle Ctrl+C gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start load test in background
	go jobConfig.Run(updates)

	// Wait for completion or interruption
	// For CLI mode, we only care about the final stats (when channel closes)
	var finalStats *http.LoadTestStats

	for {
		select {
		case stats, ok := <-updates:
			if !ok {
				// Channel closed, test complete
				return FormatOutput(finalStats, config)
			}
			// Store latest stats, but don't process intermediate updates
			// This is much more efficient than the TUI which renders each update
			finalStats = stats

		case <-sigCh:
			// User interrupted
			fmt.Fprintln(os.Stderr, "\nTest interrupted by user")
			if finalStats != nil {
				return FormatOutput(finalStats, config)
			}
			return nil
		}
	}
}
