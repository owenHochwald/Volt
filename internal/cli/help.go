package cli

import "fmt"

// PrintHelp displays usage information
func PrintHelp() {
	fmt.Println(`Volt - Terminal HTTP Client and Load Tester

USAGE:
  volt             Launch interactive TUI
  volt bench       Run CLI load test

BENCH FLAGS:
  -url <string>     Target URL (required)
  -c <int>          Number of concurrent connections (default: 50)
  -d <duration>     Test duration, e.g. "30s", "5m" (default: 10s)
  -n <int>          Total number of requests (mutually exclusive with -d)
  -m <string>       HTTP method (default: GET)
  -H <string>       Custom header, repeatable (format: "Key: Value")
  -b <string>       Request body
  -t <duration>     Request timeout (default: 30s)
  -rate <int>       Rate limit (requests/sec, 0 = unlimited)
  -keepalive        Enable HTTP keep-alive (default: true)
  -no-keepalive     Disable HTTP keep-alive
  -q                Quiet mode (minimal output)
  -json             Output results as JSON
  -o <file>         Write results to file

EXAMPLES:
  # Basic throughput test
  volt bench -url http://localhost:8080 -c 100 -d 30s

  # POST request with custom headers
  volt bench -url http://localhost:8080/api -m POST \
    -b '{"test":true}' -H "Content-Type: application/json"

  # JSON output to file for CI/CD
  volt bench -url http://localhost:8080 -c 50 -d 60s -json -o results.json

  # Rate-limited testing
  volt bench -url http://localhost:8080 -c 10 -d 30s -rate 1000

  # Quiet mode (just final stats)
  volt bench -url http://localhost:8080 -c 100 -n 10000 -q`)
}
