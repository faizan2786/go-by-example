package url

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// test cases for Parse function
var parseTestCases = []struct {
	name   string // name of the test case
	rawURL string
	want   *URL
	err    error
}{
	{
		name:   "URL without path",
		rawURL: "https://myurl.com",
		want:   &URL{Scheme: "https", Host: "myurl.com"},
		err:    nil,
	},

	{
		name:   "Full URL",
		rawURL: "https://myurl.com/myblog",
		want:   &URL{"https", "myurl.com", "myblog"},
		err:    nil,
	},
	{
		name:   "Invalid URL",
		rawURL: "https//myurl.com",
		err:    errors.New("missing ':' in the provided url string"),
	},
	{
		name:   "Opaque URL",
		rawURL: "data:text/json",
		want:   &URL{Scheme: "data"},
		err:    nil,
	},
}

func TestURLString(t *testing.T) {

	url := URL{Scheme: "http", Host: "www.dummyurl.com", Path: "mypage"}
	got := url.String()
	want := "http://www.dummyurl.com/mypage"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestParse(t *testing.T) {

	for _, tt := range parseTestCases {

		// run each test case as a subtest
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test case: %q\n", tt.name)
			got, gotErr := Parse(tt.rawURL)

			// if error is not expected
			if tt.err == nil && gotErr != nil {
				t.Fatalf("Parse(%q) err = %q, want %q", tt.rawURL, gotErr, tt.err)
			}

			// if error is expected but content of the error we got is different
			if tt.err != nil && (gotErr == nil || tt.err.Error() != gotErr.Error()) {
				t.Fatalf("Parse(%q) err = %q, want %q", tt.rawURL, gotErr, tt.err)
			}

			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("Parse(%q) output mismatch (-want +got):\n%s", tt.rawURL, diff)
			}
		})
	}
}
