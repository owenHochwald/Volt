package http

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/influxdata/tdigest"
)

type LoadTestConfig struct {
	Request       *Request      // base request to send
	Concurrency   int           // number of concurrent requests
	TotalRequests int           // total requests to send
	RateLimit     int           // max requests per second
	Timeout       time.Duration // time per request
}

// PercentileCalculator data structure for streaming percentile calculations
type PercentileCalculator struct {
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

// NewLoadTestStats creates a new LoadTestStats instance
func NewLoadTestStats(totalRequests int) *LoadTestStats {
	return &LoadTestStats{
		StartTime:     time.Now(),
		TotalRequests: totalRequests,
		MinDuration:   time.Duration(math.MaxInt64),
		percentiles: &PercentileCalculator{
			digest: tdigest.NewWithCompression(100),
		},
		Errors: make(map[string]int64),
	}
}

type Job struct {
	ID      int
	Request *Request
}

type Result struct {
	JobID      int
	StatusCode int
	Duration   time.Duration
	Error      *string
}

func (p *PercentileCalculator) Add(duration time.Duration) {
	p.digest.Add(float64(duration.Milliseconds()), 1)
}

func (p *PercentileCalculator) Percentile(percentile float64) time.Duration {
	ms := p.digest.Quantile(percentile / 100.0)
	return time.Duration(ms) * time.Millisecond
}

func (s *LoadTestStats) RecordResult(result Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CompletedRequests++
	if result.Error != nil {
		s.FailedRequests++
		statusString := strconv.Itoa(result.StatusCode)
		s.Errors[statusString]++
	}

	s.MinDuration = min(s.MinDuration, result.Duration)
	s.MaxDuration = max(s.MaxDuration, result.Duration)
	s.TotalDuration += result.Duration

	s.percentiles.Add(result.Duration)
}

func (s *LoadTestStats) GetSnapshot() LoadTestStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// make a deep copy of errors map
	errorsCopy := make(map[string]int64, len(s.Errors))
	for k, v := range s.Errors {
		errorsCopy[k] = v
	}

	return LoadTestStats{
		StartTime:         s.StartTime,
		EndTime:           s.EndTime,
		TotalRequests:     s.TotalRequests,
		CompletedRequests: s.CompletedRequests,
		FailedRequests:    s.FailedRequests,
		MinDuration:       s.MinDuration,
		MaxDuration:       s.MaxDuration,
		TotalDuration:     s.TotalDuration,
		percentiles:       s.percentiles, // can share, it's thread safe
		BytesSent:         s.BytesSent,
		BytesRecv:         s.BytesRecv,
		Errors:            errorsCopy,
		CPUUsage:          s.CPUUsage,
		MemoryUsage:       s.MemoryUsage,
	}
}

func worker(id int, client *Client, jobs chan Job, results chan Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		response := make(chan *Response, 1)
		client.Send(job.Request, response)
		responseObject := <-response
		result := Result{
			JobID:      job.ID,
			StatusCode: responseObject.StatusCode,
			Duration:   responseObject.Duration,
		}

		if responseObject.Error != "" {
			result.Error = &responseObject.Error
		}
		results <- result
	}
}

// RunLoadTest runs a load test with a config and sends out updates
func RunLoadTest(ctx context.Context, config *LoadTestConfig, updates chan<- LoadTestStats) error {
	//context := context.Background()

	var wg sync.WaitGroup

	// job channel with capacity for numWorkers
	numWorkers := config.Concurrency
	jobs := make(chan Job, numWorkers)
	results := make(chan Result, numWorkers*2)
	client := InitClient(30*time.Second, true)
	s := NewLoadTestStats(config.TotalRequests)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, client, jobs, results, &wg)
	}

	// send workers to work!
	go func() {
		for i := 0; i < config.TotalRequests; i++ {
			jobs <- Job{ID: i, Request: config.Request}
		}
		close(jobs)
	}()

	// collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	//send snapshots of results queue every 100 milliseconds
	updateTicker := time.NewTicker(time.Millisecond * 100)
	defer updateTicker.Stop()

	// brain of the load tester
	for {
		select {
		case result, ok := <-results:
			if !ok {
				// All workers finished - set end time and send final snapshot
				s.mu.Lock()
				s.EndTime = time.Now()
				s.mu.Unlock()
				updates <- s.GetSnapshot()
				close(updates)
				return nil
			}

			s.RecordResult(result)
			if result.Error != nil {
				fmt.Printf("Job %d failed: %v\n", result.JobID, result.Error)
			} else {
				fmt.Printf("Job %d completed successfully!\n", result.JobID)
			}
		case <-updateTicker.C:
			snapshot := s.GetSnapshot()
			updates <- snapshot
		case <-ctx.Done():
			s.mu.Lock()
			s.EndTime = time.Now()
			s.mu.Unlock()
			updates <- s.GetSnapshot()
			close(updates)
			return ctx.Err()
		}

	}
}
