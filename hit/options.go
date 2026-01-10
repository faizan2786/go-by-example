package hit

import (
	"net/http"
	"time"
)

type SendFunc func(*http.Request) Result

// Options defines options for sending requests.
// Uses default values for unset options.
type Options struct {

	// number of concurrent requests to send
	// Default: 1
	Concurrency int

	// number of requests to send per second
	// Default: 0 (no rate limiting)
	RPS int

	// a request processing function
	// Default: uses [Send].
	Send SendFunc
}

// returns [Options] with defaults.
func DefaultOptions() Options {
	options := Options{}
	return withDefaults(options)
}

func withDefaults(op Options) Options {
	if op.Concurrency <= 0 {
		op.Concurrency = 1
	}

	if op.RPS < 0 {
		op.RPS = 0
	}

	if op.Send == nil {

		// a closure that wraps the hit.Send function with a default http client
		op.Send = func(req *http.Request) Result {

			// configure the http client to maintain a TCP connection pool
			// (so that each dispatcher goroutine can establish a TCP connection only once
			// and reuse it for subsequent requests)
			// Also, disable http redirects using the check redirect option
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: op.Concurrency,
				},
				CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Timeout: 10 * time.Second, // timeout per request
			}

			return Send(client, req)
		}
	}

	return op
}
