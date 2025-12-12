package hit

import "net/http"

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
			return Send(http.DefaultClient, req)
		}
	}

	return op
}
