package hit

import (
	"fmt"
	"net/http"
	"time"
)

// Send sends an HTTP request and returns its performance metric as [Result].
func Send(_ *http.Client, _ *http.Request) Result {
	const roundTripTime = 100 * time.Millisecond

	// for now, simulate sending a request and return a successful result
	time.Sleep(roundTripTime)

	return Result{
		Status:   http.StatusOK,
		Bytes:    100,
		Duration: roundTripTime,
	}

}

// SendN sends N requests using [Send].
// It returns a [Results] iterator that
// pushes a [Result] for each [http.Request] sent.
func SendN(N int, opts Options, req *http.Request) (Results, error) {

	// fills opts with default values for unset/invalid options
	opts = withDefaults(opts)

	if N <= 0 {
		return nil, fmt.Errorf("n must be greater than 0: got %d", N)
	}

	results := runPipeline(N, opts, req)

	// define an iterator with a yield function that
	// reads a result from results channel and produces (i.e. yields) to the consumer
	iter := func(yield func(Result) bool) {
		for result := range results {
			if !yield(result) {
				return
			}
		}
	}

	// Note: we use iterator instead of a slice as
	// using a slice will need us to complete sending each request before we can return their results
	// whereas iterator produces each result lazily, allowing EARLY STOPPING if yield returns false
	// (i.e. if something goes wrong or consumer wants to stop receiving further values)
	// hence, saving further compute and memory allocations

	return iter, nil
}
