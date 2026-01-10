package hit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Send sends an HTTP request and returns its performance metric as [Result].
func Send(client *http.Client, req *http.Request) Result {
	var (
		bytes  int64
		status int
	)

	// send the request to the server
	start := time.Now()
	res, err := client.Do(req)

	// read the response
	if err == nil {
		defer res.Body.Close()
		status = res.StatusCode
		// we just need to know number of bytes in the response
		// so stream the response efficiently (vi io.copy) and discard its content
		bytes, err = io.Copy(io.Discard, res.Body)
	}

	return Result{
		Status:   status,
		Bytes:    bytes,
		Duration: time.Since(start),
		Error:    err,
	}
}

// SendN sends N requests using [Send].
// It returns a [Results] iterator that
// pushes a [Result] for each [http.Request] sent.
func SendN(ctx context.Context, N int, opts Options, req *http.Request) (Results, error) {

	// fills opts with default values for unset/invalid options
	opts = withDefaults(opts)

	if N <= 0 {
		return nil, fmt.Errorf("n must be greater than 0: got %d", N)
	}

	// create a new child context from the received context
	// this new context will enable us trigger the cancellation in the pipeline even when the parent context is alive
	// e.g. when the iterator stops early (when the consumer wants to consume only part of the results)
	ctx, cancel := context.WithCancel(ctx)

	results := runPipeline(ctx, N, opts, req)

	// define an iterator with a yield function that
	// reads a result from results channel and produces (i.e. yields) to the consumer
	iter := func(yield func(Result) bool) {
		defer cancel() // cancel the derived context right before returning - in turn, cause the pipeline to stop
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
