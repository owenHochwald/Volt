package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// BenchConfig holds parsed CLI flags for benchmarking
type BenchConfig struct {
	// Target configuration
	URL    string
	Method string
	Body   string

	// Headers (parsed from repeated -H flags)
	Headers map[string]string

	// Load parameters
	Concurrency   int
	Duration      time.Duration
	TotalRequests int

	// Performance tuning
	Timeout   time.Duration
	RateLimit int // requests per second, 0 = unlimited
	KeepAlive bool

	// Output options
	Quiet  bool
	JSON   bool
	Output string // file path for output
}

// headerFlags implements flag.Value for repeated -H flags
type headerFlags map[string]string

func (h headerFlags) String() string {
	return ""
}

func (h headerFlags) Set(value string) error {
	// Parse "Key: Value" format
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("header must be in format 'Key: Value'")
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	h[key] = val
	return nil
}

// ParseBenchFlags parses command-line flags for bench subcommand
func ParseBenchFlags(args []string) (*BenchConfig, error) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)

	config := &BenchConfig{
		Headers: make(map[string]string),
	}

	headers := make(headerFlags)

	// Target configuration
	fs.StringVar(&config.URL, "url", "", "Target URL (required)")
	fs.StringVar(&config.Method, "m", "GET", "HTTP method")
	fs.StringVar(&config.Body, "b", "", "Request body")
	fs.Var(headers, "H", "Custom header (repeatable, format: 'Key: Value')")

	// Load parameters
	fs.IntVar(&config.Concurrency, "c", 50, "Number of concurrent connections")
	fs.DurationVar(&config.Duration, "d", 10*time.Second, "Test duration (e.g., '30s', '5m')")
	fs.IntVar(&config.TotalRequests, "n", 0, "Total number of requests (0 = use duration)")

	// Performance tuning
	fs.DurationVar(&config.Timeout, "t", 30*time.Second, "Request timeout")
	fs.IntVar(&config.RateLimit, "rate", 0, "Rate limit (requests/sec, 0 = unlimited)")
	keepAlive := fs.Bool("keepalive", true, "Enable HTTP keep-alive")
	noKeepAlive := fs.Bool("no-keepalive", false, "Disable HTTP keep-alive")

	// Output options
	fs.BoolVar(&config.Quiet, "q", false, "Quiet mode (minimal output)")
	fs.BoolVar(&config.JSON, "json", false, "Output results as JSON")
	fs.StringVar(&config.Output, "o", "", "Write results to file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// If -n was provided, clear duration to avoid mutual exclusivity error
	if config.TotalRequests > 0 {
		config.Duration = 0
	}

	// Handle keep-alive flags
	if *noKeepAlive {
		config.KeepAlive = false
	} else {
		config.KeepAlive = *keepAlive
	}

	config.Headers = headers

	return config, nil
}
