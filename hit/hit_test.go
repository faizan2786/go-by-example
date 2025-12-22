package hit

import (
	"net/http"
	"testing"
	"time"
)

func TestSend(t *testing.T) {

	got := Send(nil, nil)

	want := Result{
		Status:   http.StatusOK,
		Bytes:    100,
		Duration: 100 * time.Millisecond,
	}

	if got != want {
		t.Fatalf("Send() = %+v; want = %+v\n", got, want)
	}
}

// Test function for Send()
func getSendFunc() SendFunc {
	return func(_ *http.Request) Result {
		const roundTripTime = 10 * time.Millisecond
		time.Sleep(roundTripTime)

		return Result{
			Status:   http.StatusOK,
			Bytes:    50,
			Duration: roundTripTime,
		}
	}
}

func getHttpRequest() *http.Request {
	req, err := http.NewRequest(http.MethodGet, "www.myurl.com", http.NoBody)
	if err != nil {
		panic("Failed to get a new http GET request for testing")
	}
	return req
}

func TestSendN(t *testing.T) {

	N := 100
	opts := Options{}
	opts.Send = getSendFunc()
	req := getHttpRequest()

	res, err := SendN(N, opts, req)
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
