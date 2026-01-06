// This file defines the hit client that can send multiple requests
// to the server and provide its performance as a summary

package hit

import (
	"iter"
	"time"
)

// Result is performance metrics of a single [http.Request].
type Result struct {
	Status   int           // 200
	Bytes    int64         // Number of bytes received
	Duration time.Duration // Duration to complete a request
	Error    error
}

// Results is an iterator for a collection of [Result] values.
type Results iter.Seq[Result]

// Summary is the summary of [Result] values.
type Summary struct {
	Requests int           // Requests is the total number of requests made
	Errors   int           // Errors is the total number of failed requests
	Bytes    int64         // Bytes is the total number of bytes received
	Fastest  time.Duration // Fastest is the fastest request duration
	Slowest  time.Duration // Slowest is the slowest request duration
	Average  time.Duration // Average request duration - average response time for an individual request (i.e. Latency)
	Duration time.Duration // Duration is the total (clock) time taken by all the requests
	RPS      float64       // RPS is the number of requests served per second (i.e. Throughput)
	Success  float64       // Success is the ratio of successful requests
}

// Summarize returns a [Summary] of [Results].
func Summarize(results Results) Summary {
	var s Summary

	// handle nil results (because ranging over a nil iterator causes panic)
	if results == nil {
		return s // return a zero-value summary
	}

	var requestDurationSum time.Duration // sum of all request durations

	start := time.Now()
	for r := range results {
		s.Requests += 1
		s.Bytes += r.Bytes

		if r.Error != nil {
			s.Errors += 1
		}

		if s.Fastest == 0 || r.Duration < s.Fastest {
			s.Fastest = r.Duration
		}

		if r.Duration > s.Slowest {
			s.Slowest = r.Duration
		}

		requestDurationSum += r.Duration
	}

	s.Duration = time.Since(start)                     // total clock time
	s.RPS = float64(s.Requests) / s.Duration.Seconds() // throughput

	if s.Requests > 0 {
		s.Average = requestDurationSum / time.Duration(s.Requests) // latency
		s.Success = (float64(s.Requests-s.Errors) / float64(s.Requests)) * 100
	}

	return s
}
