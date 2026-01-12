package hit

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

// define a test http request
func getTestHttpRequest() *http.Request {
	// define a test http request
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	if err != nil {
		panic("Failed to get a new http GET request")
	}
	return req
}

// test Send function (using a fake round tripper)

// define a function type that satisfies the http RoundTripper interface
type roundTripperFunc func(*http.Request) (*http.Response, error)

// implement the Roundtrip method of the interface
func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestSend(t *testing.T) {

	// define a fake round tripper function
	fakeRoundTripper := func(_ *http.Request) (*http.Response, error) {
		// return a test status code
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(strings.NewReader("fake response.")),

			// NoOpCloser wraps the argument and makes it a "Closer" by implementing Close() method with no operation
			// This is necessary because Body expects a ReadCloser and strings.NewReader returns just a Reader
		}, nil
	}

	// create a client with a roundTripperFunc type
	// (convert fakeRoundTripper to the roundTripperFunc type so that
	// client can use the fake round tripper as its Transport)
	client := &http.Client{
		Transport: roundTripperFunc(fakeRoundTripper), // converts the fakeRoundTripper function to the roundTripperFunc type
	}

	req := getTestHttpRequest()

	res := Send(client, req)

	if res.Status != http.StatusAccepted {
		t.Errorf("got: status = %d, want: %d\n", res.Status, http.StatusAccepted)
	}
	if res.Bytes != int64(len("fake response.")) {
		t.Errorf("got: %d bytes, want: %d\n", res.Bytes, len("fake response."))
	}
}

// test send function with time duration in the response.
// we use time.synctest package which runs tests in a "bubble".
// Specifically, it mocks time by allowing time to advance by an exact required amount
// (instead of waiting for said time, synctest allows time calls inside the test to jump by said amount)
// hence, the test is repeatable and fast
func TestSendWithDuration(t *testing.T) {

	// a round tripper that uses time.Sleep to simulate a long running request
	fakeRoundTripper := func(_ *http.Request) (*http.Response, error) {
		time.Sleep(30 * time.Second)

		res := &http.Response{
			StatusCode: http.StatusInternalServerError,
		}
		return res, nil
	}
	client := &http.Client{
		Transport: roundTripperFunc(fakeRoundTripper),
	}

	// run the test in a "bubble" using synctest
	synctest.Test(t, func(t *testing.T) {

		res := Send(client, getTestHttpRequest()) // call finishes "instantly" - NOT after 30s

		if res.Duration != (30 * time.Second) {
			t.Errorf("got: %v duration, want: %v\n", res.Duration, 30*time.Second)
		}
	})
}

// test SendN (using a fake send function)
func TestSendN(t *testing.T) {
	const N int = 50

	opts := Options{}

	// define a fake Send() function
	opts.Send = func(_ *http.Request) Result {
		const roundTripTime = 10 * time.Millisecond
		time.Sleep(roundTripTime)

		return Result{
			Status:   http.StatusOK,
			Bytes:    50,
			Duration: roundTripTime,
		}
	}

	// get the test http request
	req := getTestHttpRequest()

	res, err := SendN(context.Background(), N, opts, req)
	if err != nil {
		t.Fatalf("SendN() = %v; want no error\n", err)
	}

	// count number of results received from the iterator
	gotN := 0
	for range res {
		gotN += 1
	}

	if gotN != N {
		t.Errorf("SendN() returned %d results; want %d\n results", gotN, N)
	}
}
