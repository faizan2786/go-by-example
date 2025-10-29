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
		name:   "full",
		rawURL: "https://myurl.com/myblog",
		want:   &URL{"https", "myurl.com", "myblog"},
		err:    nil,
	},
	{
		name:   "invalid",
		rawURL: "https//myurl.com",
		err:    errors.New("missing scheme"),
	},
	{
		name:   "opaque",
		rawURL: "data:text/json",
		want:   &URL{Scheme: "data"},
		err:    nil,
	},
	{
		name:   "no_scheme",
		rawURL: "://github.com",
		err:    errors.New("missing scheme"),
	},
	{
		name:   "no_host",
		rawURL: "https:///myblog",
		err:    errors.New("path without host"),
	},
	{
		name:   "no_path",
		rawURL: "https://myurl.com",
		want:   &URL{Scheme: "https", Host: "myurl.com"},
		err:    nil,
	},
	{
		name:   "host_with_trailing_slash",
		rawURL: "https://myurl.com/",
		want:   &URL{Scheme: "https", Host: "myurl.com"},
		err:    nil,
	},
}

func TestParse(t *testing.T) {

	for _, tt := range parseTestCases {

		// run each test case as a subtest
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test case: %q\n", tt.name)
			got, gotErr := Parse(tt.rawURL)

			// if error is not expected
			if tt.err == nil && gotErr != nil {
				t.Fatalf("Parse(%q) err = %v, want error: %v", tt.rawURL, gotErr, tt.err)
			}

			// if error is expected but content of the error we got is different
			if tt.err != nil && (gotErr == nil || tt.err.Error() != gotErr.Error()) {
				t.Fatalf("Parse(%q) err = %v, want error: %v", tt.rawURL, gotErr, tt.err)
			}

			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("Parse(%q) output mismatch (-want +got):\n%s", tt.rawURL, diff)
			}
		})
	}
}

func TestURLString(t *testing.T) {

	testCases := []struct {
		name string
		url  *URL
		want string
	}{
		{
			name: "valid",
			url:  &URL{"http", "www.dummyurl.com", "mypage"},
			want: "http://www.dummyurl.com/mypage",
		},
		{
			name: "empty",
			url:  new(URL), // create an empty URL (i.e. same as &URL{})
			want: "",
		},
		{
			name: "nil",
			url:  nil,
			want: "",
		},
		{
			name: "scheme_only",
			url:  &URL{Scheme: "https"},
			want: "https://",
		},
		{
			name: "no_path",
			url:  &URL{Scheme: "https", Host: "www.dummyurl.com"},
			want: "https://www.dummyurl.com",
		},
		{
			name: "no_host",
			url:  &URL{Scheme: "https", Path: "myblog"},
			want: "https:///myblog",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.url.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}

}
