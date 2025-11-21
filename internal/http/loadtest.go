package http

import (
	"net/http"
	"sync"
	"time"

	"github.com/influxdata/tdigest"
)

type LoadTestConfig struct {
	Request       *http.Request // base request to send
	Concurrency   int           // number of concurrent requests
	TotalRequests int           // total requests to send
	RateLimit     int           // max requests per second
	Timeout       time.Duration // time per request
}

// PercentileCalculator data structure for streaming percentile calculations
// using t-digest algorithm
type PercentileCalculator struct {
	// use github.com/influxdata/tdigest
	digest *tdigest.TDigest
}

// LoadTestStats holds aggregated stats about the load test
type LoadTestStats struct {
	StartTime         time.Time
	EndTime           time.Time
	TotalRequests     int
	CompletedRequests int
	FailedRequests    int

	// time tracking
	MinDuration   time.Duration
	MaxDuration   time.Duration
	TotalDuration time.Duration

	// streaming percentile calculations
	percentiles *PercentileCalculator

	// network stats
	BytesSent int64
	BytesRecv int64

	// error tracking
	Errors map[string]int64 // error code -> count

	// system metrics
	CPUUsage    float64 // percentage
	MemoryUsage uint64  // bytes

	mu sync.RWMutex // protects above fields from concurrent access
}

type LoadTestResult struct {
	Stats  *LoadTestStats
	IsDone bool
	Err    error
}

func NewLoadTestStats() *LoadTestStats {
	return &LoadTestStats{
		percentiles: &PercentileCalculator{
			digest: tdigest.NewWithCompression(100),
		},
		Errors: make(map[string]int64),
	}
}

func RunLoadTest(config *LoadTestConfig, updates chan<- LoadTestStats) error {
	// worker goroutines with shared work channel
	// collect stats with sync.RWMutex
	// cancel with context.Context

	// FUTURE: add rate limiting with time/rate
	return nil

}
