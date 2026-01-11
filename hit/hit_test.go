package hit

import (
	"context"
	"net/http"
	"testing"
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
		t.Fatalf("want: %d, got: %d\n", http.StatusAccepted, res.Status)
	}
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
		t.Fatalf("SendN() returned %d results; want %d\n results", gotN, N)
	}
}
