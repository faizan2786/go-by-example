// This file defines a concurrent pipeline to send N requests to a server concurrently
// consists of 3 stages: a Producer, Throttler and a Dispatcher
// each stage returns a receive-only channel to deliver its output

package hit

import (
	"net/http"
	"sync"
	"time"
)

func runPipeline(n int, opts Options, req *http.Request) <-chan Result {

	requests := produce(n, req) // stage-1

	// throttle if RPS is given
	if opts.RPS > 0 {
		requests = throttle(opts.RPS, requests) // stage-2
	}

	return dispatch(opts, requests)
}

// produces [http.Request]s.
func produce(n int, req *http.Request) <-chan *http.Request {
	// step-1: make an output channel
	out := make(chan *http.Request)

	// step-2: spawn worker go routine(s)
	go func() {
		defer close(out) // this is IMP to unblock any receiver
		for range n {
			out <- req
		}
	}()

	// step-3: return output channel immediately
	return out
}

func throttle(rps int, in <-chan *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)

	if rps > 0 {
		interval := time.Second / time.Duration(rps) // time interval between each tick (i.e. request)
		go func() {
			defer close(out) // this is IMP to unblock any receiver
			t := time.NewTicker(interval)
			for req := range in {
				<-t.C      // wait until next tick
				out <- req // send it to output channel
			}
		}()
	}

	return out
}

func dispatch(opts Options, in <-chan *http.Request) <-chan Result {
	out := make(chan Result)

	var wg sync.WaitGroup
	for range opts.Concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// read the requests, invoke Send() and send result to out channel
			for req := range in {
				out <- opts.Send(req)
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
