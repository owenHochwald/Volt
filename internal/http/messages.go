package http

import "time"

type LoadTestStartMsg struct {
	Config *JobConfig
}

type LoadTestStatsMsg struct {
	Stats    *LoadTestStats
	Progress float64 // 0.0 to 1.0
}

type LoadTestCompleteMsg struct {
	Stats    *LoadTestStats
	Duration time.Duration
}

type LoadTestErrorMsg struct {
	Error error
}
