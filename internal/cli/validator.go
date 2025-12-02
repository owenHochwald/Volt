package cli

import (
	"errors"
	"strings"
)

// Validate checks BenchConfig for errors
func (c *BenchConfig) Validate() error {
	// URL is required
	if c.URL == "" {
		return errors.New("--url is required")
	}

	// URL must start with http:// or https://
	if !strings.HasPrefix(c.URL, "http://") && !strings.HasPrefix(c.URL, "https://") {
		return errors.New("URL must start with http:// or https://")
	}

	// Method must be valid
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true,
		"DELETE": true, "PATCH": true, "HEAD": true,
	}
	if !validMethods[strings.ToUpper(c.Method)] {
		return errors.New("invalid HTTP method")
	}

	// Concurrency must be positive
	if c.Concurrency <= 0 {
		return errors.New("concurrency must be > 0")
	}

	// Duration and TotalRequests are mutually exclusive
	if c.Duration == 0 && c.TotalRequests == 0 {
		return errors.New("must specify either -d (duration) or -n (total requests)")
	}
	if c.Duration > 0 && c.TotalRequests > 0 {
		return errors.New("-d and -n are mutually exclusive")
	}

	// Timeout must be positive
	if c.Timeout <= 0 {
		return errors.New("timeout must be > 0")
	}

	// Rate limit must be non-negative
	if c.RateLimit < 0 {
		return errors.New("rate limit must be >= 0")
	}

	return nil
}
