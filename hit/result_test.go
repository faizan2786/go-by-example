package hit

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

func TestSummarize(t *testing.T) {

	// define a collection of results
	results := []Result{
		{
			Bytes:    100,
			Duration: 100 * time.Millisecond,
		},
		{
			Bytes:    500,
			Duration: 300 * time.Millisecond,
		},
		{
			Duration: 500 * time.Millisecond,
			Error:    fmt.Errorf("server can't be reached."),
		},
	}

	res := Results(slices.Values(results)) // convert iter.Seq[Result] returned by slice.Values to Results
	s := Summarize(res)

	want := Summary{
		Bytes:    600,
		Fastest:  100 * time.Millisecond,
		Slowest:  500 * time.Millisecond,
		Average:  300 * time.Millisecond,
		Requests: 3,
		Errors:   1,
		Success:  66.667,
	}

	if s.Bytes != want.Bytes {
		t.Errorf("Bytes: got = %d, want = %d\n", s.Bytes, want.Bytes)
	}

	if s.Fastest != want.Fastest {
		t.Errorf("Fastest: got = %v, want = %v\n", s.Fastest, want.Fastest)
	}

	if s.Slowest != want.Slowest {
		t.Errorf("Slowest: got = %v, want = %v\n", s.Slowest, want.Slowest)
	}

	if s.Average != want.Average {
		t.Errorf("Average: got = %v, want = %v\n", s.Average, want.Average)
	}

	if s.Errors != want.Errors {
		t.Errorf("Errors: got = %d, want = %d\n", s.Errors, want.Errors)
	}

	if s.Requests != want.Requests {
		t.Errorf("Requests: got = %d, want = %d\n", s.Requests, want.Requests)
	}

	// compare success rate up to 2 decimals
	if fmt.Sprintf("%.2f", s.Success) != fmt.Sprintf("%.2f", want.Success) {
		t.Errorf("Success rate: got = %.2f, want = %.2f\n", s.Success, want.Success)
	}
}

// Test that Summarize doesn't panic when receiving a nil Results
func TestSummarizeNilResults(t *testing.T) {

	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("should not panic: %v\n", err)
		}
	}()

	_ = Summarize(nil)
}
