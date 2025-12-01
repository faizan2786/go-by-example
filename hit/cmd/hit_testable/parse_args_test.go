package main

import (
	"strings"
	"testing"
)

func TestParseArgsValid(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name string
		args []string
		want argConfig
	}{
		{
			name: "only_url",
			args: []string{"http://www.myurl.com"},
			want: argConfig{n: 0, c: 0, rps: 0, url: "http://www.myurl.com"},
		},
		{
			name: "custom_flags",
			args: []string{"-n=1500", "-c=4", "-rps=50", "http://www.myurl.com"},
			want: argConfig{n: 1500, c: 4, rps: 50, url: "http://www.myurl.com"},
		},
	}

	for _, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {

			var gotConfig argConfig
			var testStdErr strings.Builder

			err := parseArgs(tt.args, &gotConfig, &testStdErr)

			// check for no errors
			if err != nil {
				t.Fatalf("got error = %v; want: <nil>\n", err)
			}

			// IMP: check that our tool didn't log any messages on stderr
			if testStdErr.Len() != 0 {
				t.Fatalf("got len(stderr) = %d bytes; want: 0 bytes\n", testStdErr.Len())
			}

			// check the config parameters
			if gotConfig != tt.want {
				t.Errorf("got: %+v; want: %+v\n", gotConfig, tt.want)
			}
		})
	}
}

func TestParseArgsInvalid(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "invalid_url",
			args: []string{"www.myurl.com"},
		},
		{
			name: "invalid_flag",
			args: []string{"-n=1500", "-c=TWO", "http://www.myurl.com"},
		},
		{
			name: "non_positive_flag",
			args: []string{"-n=100", "-c=0", "-rps=50", "http://www.myurl.com"},
		},
	}

	for _, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {

			var gotConfig argConfig
			var testStdErr strings.Builder

			err := parseArgs(tt.args, &gotConfig, &testStdErr)

			// check for errors
			if err == nil {
				t.Fatal("got error = <nil>; want an error")
			}

			// IMP: check that our tool log the error messages on stderr
			if testStdErr.Len() == 0 {
				t.Error("got: len(stderr) = 0 bytes; want non-zero bytes")
			}
		})
	}
}
