package main

import (
	"fmt"
	"strings"
	"testing"
)

// define a struct with fake (testing) env variables for run method
// (Builder implements io.Writer/io.Reader interface and
// read/writes the input data into its internal buffer
// which we can use to inspect the run method's outputs)
type testEnv struct {
	stdout strings.Builder
	stderr strings.Builder
}

// a helper method to define test env and call the run() method
// it returns the created test env and any error returned by run() so that we can inspect them in our tests
func testRun(args ...string) (*testEnv, error) {
	testenv := testEnv{}

	// bind the testEnv's stdout & stderr to env struct's corresponding variables
	// so that we can extract the output and error messages after we call run()
	e := &env{
		stdout:   &testenv.stdout,
		stderr:   &testenv.stderr,
		args:     append([]string{"hit"}, args...), // run should get entire os.Args (including the program name)
		testMode: true,
	}

	err := run(e)

	return &testenv, err
}

// test run method for valid inputs
func TestRunValid(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name string
		args []string
		want string // message printed by the run method on stdout
	}{
		{
			name: "default_flags",
			args: []string{"http://www.myurl.com"},
			want: fmt.Sprintf("%s\nSending %d requests to %q at %d requests per second (concurrency=%d)\n", logo, 1000, "http://www.myurl.com", 100, 1),
		},
		{
			name: "custom_flags",
			args: []string{"-n=1500", "-c=4", "-rps=50", "http://www.myurl.com"},
			want: fmt.Sprintf("%s\nSending %d requests to %q at %d requests per second (concurrency=%d)\n", logo, 1500, "http://www.myurl.com", 50, 4),
		},
	}

	for _, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {

			testArgs := append([]string{}, tt.args...) // pass a copy of args (to prevent testRun from accidentally modifying the args)

			// call the test run method
			testEnv, err := testRun(testArgs...)

			// check for no errors
			if err != nil {
				t.Fatalf("got error: %v; want: <nil>\n", err)
			}

			// IMP: check that our tool didn't log any messages on stderr
			if testEnv.stderr.Len() != 0 {
				t.Errorf("got non-empty stderr; want len(stderr) = 0 bytes\n")
			}

			// check the message on stdout
			got := testEnv.stdout.String()
			if got != tt.want {
				t.Errorf("got: %s; want: %s", got, tt.want)
			}
		})
	}
}

// test run method for invalid input
func TestRunInvalid(t *testing.T) {

	t.Parallel()

	testArgs := []string{"-n=2000"}

	// call the test run method
	testEnv, err := testRun(testArgs...)

	// check for errors
	if err == nil {
		t.Fatal("got error: <nil>; want an error")
	}

	// IMP: check that our tool has logged error message on stderr
	if testEnv.stderr.Len() == 0 {
		t.Error("got len(stderr) = 0 bytes; want: non-zero bytes")
	}

	// check that there is no message on stdout (since our tool failed)
	if testEnv.stdout.Len() != 0 {
		t.Fatalf("got len(stdout) = %d bytes; want: 0 bytes\n", testEnv.stdout.Len())
	}
}
