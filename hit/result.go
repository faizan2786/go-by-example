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
	RPS      float64       // RPS is the number of requests sent per second
	Duration time.Duration // Duration is the total time taken by requests
	Fastest  time.Duration // Fastest is the fastest request duration
	Slowest  time.Duration // Slowest is the slowest request duration
	Success  float64       // Success is the ratio of successful requests
}

// Summarize returns a [Summary] of [Results].
func Summarize(results Results) Summary {
	var s Summary

	// handle nil results (because ranging over a nil iterator causes panic)
	if results == nil {
		return s // return a zero-value summary
	}

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
	}

	s.Duration = time.Since(start)
	s.RPS = float64(s.Requests) / s.Duration.Seconds()

	if s.Requests > 0 {
		s.Success = (float64(s.Requests-s.Errors) / float64(s.Requests)) * 100
	}

	return s
}
