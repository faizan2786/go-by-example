// This file defines a concurrent pipeline to send N requests to a server concurrently
// consists of 3 stages: a Producer, Throttler and a Dispatcher
// each stage returns a receive-only channel to deliver its output

package hit

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func runPipeline(ctx context.Context, n int, opts Options, req *http.Request) <-chan Result {

	requests := produce(ctx, n, req) // stage-1

	// throttle if RPS is given
	if opts.RPS > 0 {
		requests = throttle(ctx, opts.RPS, requests) // stage-2
	}

	return dispatch(ctx, opts, requests)
}

// produces [http.Request]s.
func produce(ctx context.Context, n int, req *http.Request) <-chan *http.Request {
	// step-1: make an output channel
	out := make(chan *http.Request)

	// step-2: spawn worker go routine(s) that writes to the output channel
	go func() {
		defer close(out) // it is IMP to close the channel before we return (to unblock any receiver)

		for range n {
			// send or return
			select {
			case out <- req.Clone(ctx): // clone the request with the passed context and send to output channel
			case <-ctx.Done():
				return // exit goroutine if the context is cancelled (i.e. Done)
			}

			// Note that, above way of putting send (out <- ...) on case statement will allow runtime to choose ctx.Done()
			// when send is blocked and context is cancelled.
			// This way go routine can terminate gracefully on context termination without getting stuck on send
			// Putting send under a default: case will not achieve the same effect and will block if send is blocked.
		}
	}()

	// step-3: return output channel immediately
	return out

	// Note:
	// when producer context is cancelled it closes the out channel and returns
	// which in turn causes downstream consumers of the pipeline (i.e. throttler, dispatcher, sendN function, etc.)
	// to terminate their for...range loop and exit as well.

	// However, this will not cancel the ongoing operations in the consumers (e.g., if they happened to be blocked on send).
	// Hence, to terminate each of the component of the pipeline gracefully on parent's context cancel,
	// we must pass the context to each component and use select...case in each component
}

func throttle(ctx context.Context, rps int, in <-chan *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)

	if rps > 0 {
		interval := time.Second / time.Duration(rps) // time interval between each tick (i.e. request)
		go func() {
			defer close(out)
			t := time.NewTicker(interval)
			for req := range in {
				select {
				case <-t.C: // wait until next tick
					// send or return
					select {
					case out <- req:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return out
}

func dispatch(ctx context.Context, opts Options, in <-chan *http.Request) <-chan Result {
	out := make(chan Result)

	var wg sync.WaitGroup
	for range opts.Concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// read the requests, invoke Send() and send result to out channel
			for req := range in {
				// send or return
				select {
				case out <- opts.Send(req):
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	// spawn a go routine to close the output channel when all workers are done
	// (i.e. when the input channel closes)
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
