package http

import (
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/influxdata/tdigest"
)

// largest size of the buffer for the result channel
const maxResult = 1_000_000

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
	Percentiles *PercentileCalculator

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

type JobResult struct {
	err           error
	statusCode    int
	duration      time.Duration
	contentLength int64
}

type JobConfig struct {
	// Request config
	Request *Request // base request to send

	// Load parameters
	Concurrency   int           // number of concurrent requests
	TotalRequests int           // total requests to send
	RateLimit     int           // max requests per second
	Timeout       time.Duration // time per request
	QPS           float64       // rate limit for queries per second

	// Internal state
	results chan *JobResult
	stopCh  chan struct{}
	start   time.Time
	once    sync.Once
	stats   *LoadTestStats
}

// Init creates the singleton JobConfig internal representation
func (s *JobConfig) Init() {
	s.once.Do(func() {
		s.results = make(chan *JobResult, min(s.Concurrency*1000, maxResult))
		s.stopCh = make(chan struct{}, s.Concurrency)
	})
}

// NewLoadTestStats creates a new LoadTestStats instance
func NewLoadTestStats(totalRequests int) *LoadTestStats {
	return &LoadTestStats{
		StartTime:     time.Now(),
		TotalRequests: totalRequests,
		MinDuration:   time.Duration(math.MaxInt64),
		Percentiles: &PercentileCalculator{
			digest: tdigest.NewWithCompression(100),
		},
		Errors: make(map[string]int64),
	}
}

// Run makes all the requests and streams results
func (s *JobConfig) Run(updates chan<- *LoadTestStats) {
	s.Init()
	s.start = time.Now()
	s.stats = NewLoadTestStats(s.TotalRequests)

	// collect results in the background
	go s.collectResults(updates)

	s.runWorkers()
	s.finish()
}

func (s *JobConfig) Stop() {
	for range s.Concurrency {
		s.stopCh <- struct{}{}
	}
}

func (s *JobConfig) Finish() {
	close(s.results)
	//TODO: finish implementation
	// total := time.Now().Sub(s.start)

	// await the tdigest report being completed
	// <-s
	// finalize with the total time
}

func (s *JobConfig) makeRequest(c *FastClient) {
	// res contains metrics data
	res, err := c.Do(s.Request)
	if err != nil {
		s.results <- &JobResult{err: err}
		return
	}
	s.results <- res
}

func (p *PercentileCalculator) Add(duration time.Duration) {
	p.digest.Add(float64(duration.Milliseconds()), 1)
}

func (p *PercentileCalculator) Percentile(percentile float64) time.Duration {
	ms := p.digest.Quantile(percentile / 100.0)
	return time.Duration(ms) * time.Millisecond
}

func (s *LoadTestStats) RecordResult(result JobResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CompletedRequests++
	if result.err != nil {
		s.FailedRequests++
		if result.statusCode > 0 {
			statusString := strconv.Itoa(result.statusCode)
			s.Errors[statusString]++
		}
	}

	s.MinDuration = min(s.MinDuration, result.duration)
	s.MaxDuration = max(s.MaxDuration, result.duration)
	s.TotalDuration += result.duration

	s.Percentiles.Add(result.duration)
}

func (s *LoadTestStats) GetSnapshot() LoadTestStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

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
		Percentiles:       s.Percentiles, // can share, it's thread safe
		BytesSent:         s.BytesSent,
		BytesRecv:         s.BytesRecv,
		Errors:            errorsCopy,
		CPUUsage:          s.CPUUsage,
		MemoryUsage:       s.MemoryUsage,
	}
}

func (s *JobConfig) runWorker(client *FastClient, wg *sync.WaitGroup) {
	defer wg.Done()

	var ticker *time.Ticker
	var throttle <-chan time.Time
	if s.QPS > 0 {
		ticker = time.NewTicker(time.Duration(1e6/s.QPS) * time.Microsecond)
		defer ticker.Stop()
		throttle = ticker.C
	}

	requestsPerWorker := s.TotalRequests / s.Concurrency

	for i := 0; i < requestsPerWorker; i++ {
		if s.QPS > 0 {
			<-throttle
		}
		select {
		case <-s.stopCh:
			// halt the load test if a stop result is detected
			return
		default:
			s.makeRequest(client)
		}
	}
}

func (s *JobConfig) collectResults(updates chan<- *LoadTestStats) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case result, ok := <-s.results:
			if !ok {
				snapshot := s.stats.GetSnapshot()
				updates <- &snapshot
				close(updates)
				return
			}
			// actually update the t-digest
			s.stats.RecordResult(*result)

		case <-ticker.C:
			snapshot := s.stats.GetSnapshot()
			updates <- &snapshot
		}

	}
}

func (s *JobConfig) runWorkers() {
	var wg sync.WaitGroup

	// Use optimized FastClient constructor
	client := NewFastClient(s.Timeout, s)

	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		go s.runWorker(client, &wg)
	}

	wg.Wait()
}

func (s *JobConfig) finish() {
	close(s.results)
}
