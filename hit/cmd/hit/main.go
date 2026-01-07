package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/faizan2786/gobyexample/hit"
)

const logo = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__  _\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/
  `

// define variables for the command line args
type argConfig struct {
	url string
	n   int
	c   int
	rps int
}

// define a struct to hold the configurable env parameters for the run method
type env struct {
	stdout   io.Writer
	stderr   io.Writer
	args     []string
	testMode bool // indicate if our program is running in test mode (useful for unit-tests)
}

func main() {

	args := os.Args // get the command line args from the os package (ignore 1st arg - the program name)

	env := &env{
		os.Stdout,
		os.Stderr,
		args,
		false,
	}

	if err := run(env); err != nil {
		fmt.Printf("Error while running the hit tool: %v\n", err)
		os.Exit(1)
	}
}

// run method receives our env parameters as argument and execute's the program logic
// (it is a substitute for main which can be tested via dependency injection
// i.e. using custom stdout and stderr such as string builder to capture messages during testing)
func run(e *env) error {

	config := argConfig{
		n: 1000,
		c: 1,
	}

	if err := parseArgs(e.args[1:], &config, e.stderr); err != nil {
		return err
	}

	fmt.Fprintf(e.stdout, "%s\nSending %d requests to %q (concurrency=%d)\n", logo, config.n, config.url, config.c)

	if e.testMode {
		return nil
	}

	// run the actual hit client
	err := runHit(config, e.stdout)

	return err
}

// run the HIT client with given args and print the requests summary
// (HIT client will send N requests to the server and measure its performance)
func runHit(config argConfig, stdout io.Writer) error {

	// define a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, config.url, http.NoBody)
	if err != nil {
		return fmt.Errorf("error while creating a new http request: %w", err)
	}

	opts := hit.Options{Concurrency: config.c, RPS: config.rps}

	// derive a signal notification context to catch os interrupt signals (e.g., SIGINT - generally caused by ctrl+c press)
	// this will cause the go runtime to catch interrupt signal and cancel the context (i.e. notify)
	// However, the go runtime will continue listening for the signal until the stop() function is called.
	// Hence, stop() must be called as soon as we finish the termination on first signal.

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	/*
		// call sendN and print the result
		results, err := hit.SendN(ctx, config.n, opts, req)
		if err != nil {
			return fmt.Errorf("error while sending requests: %w", err)
		}
		printResults(config.n, results, stdout)
	*/

	// call sendN and calculate the summary
	results, err := hit.SendN(ctx, config.n, opts, req)
	if err != nil {
		return fmt.Errorf("error while sending requests: %w", err)
	}
	summary := hit.Summarize(results)
	printSummary(summary, stdout)

	return ctx.Err() // returns an error if context was cancelled (for some reason) otherwise, returns nil
}

func printSummary(sum hit.Summary, stdout io.Writer) {
	fmt.Fprintf(stdout, `  
Summary:
    Success:  %.0f%%  
    RPS:      %.1f  
    Requests: %d
    Errors:   %d
    Bytes:    %d
    Duration: %s
    Fastest:  %s
    Slowest:  %s
    Average:  %s
`,
		sum.Success,
		math.Round(sum.RPS),
		sum.Requests,
		sum.Errors,
		sum.Bytes,
		sum.Duration.Round(time.Millisecond),
		sum.Fastest.Round(time.Millisecond),
		sum.Slowest.Round(time.Millisecond),
		sum.Average.Round(time.Millisecond),
	)
}

// prints results as they come with a progress bar
// it buffers the results and return an iterator for further consumption
func printResults(n int, results hit.Results, stdout io.Writer) hit.Results {

	res := make([]hit.Result, 0, n)
	curr := 0
	for r := range results {
		res = append(res, r)
		curr += 1
		printProgress(curr, n, stdout)
	}

	return hit.Results(slices.Values(res))
}

func printProgress(c int, n int, stdout io.Writer) {

	if c <= 0 || n <= 0 {
		return
	}

	const width int = 40
	filled := c * width / n

	// build string once and then print
	bar := "[" +
		strings.Repeat("=", filled) +
		strings.Repeat(" ", width-filled) +
		"] " +
		fmt.Sprintf("%d/%d", c, n) +
		fmt.Sprintf(" (%d%%)", c*100/n) // print progress in percentage

	// \r moves cursor to start of the line
	fmt.Fprintf(stdout, "\r%s", bar)
}

// function to parse command line args and assigned them to a config variable (using the flag package)
func parseArgs(args []string, config *argConfig, stderr io.Writer) error {

	flagSet := flag.NewFlagSet("hit", flag.ContinueOnError)
	flagSet.SetOutput(stderr) // set the destination for output messages (default is os's Stderr)

	// since the positional args are retrieved directly from the command line args (without a parser),
	// we need to set the usage message manually to include the positional args in the message
	flagSet.Usage = func() {

		fmt.Fprintf(
			flagSet.Output(), // returns the writer we set above
			"usage: %s [options] url\noptions:\n",
			flagSet.Name(),
		)

		// print the default values for the flags
		flagSet.PrintDefaults()
	}

	// register all of the int type flags using our custom PositiveInt Value type (parser)
	// the default values will be derived from the values already defined in the passed *config struct
	flagSet.Var(asPositiveInt(&config.c), "c", "concurrency level")
	flagSet.Var(asPositiveInt(&config.n), "n", "number of requests to send")
	flagSet.Var(asPositiveInt(&config.rps), "rps", "requests per second")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// any args that comes AFTER the flags are "positional arguments" and can be
	// retrieved by arg[i] method after parsing the args by FlagSet

	// retrieve the 1st positional argument (i.e. url)
	// (Since the positional arguments don't have named flags (i.e. -flagname), their values are accessed directly by its position)
	config.url = flagSet.Arg(0) // returns empty string if there are no positional args provided

	// validate any positional argument values
	if err := validateArgs(config); err != nil {
		// print the error message followed by the usage message
		fmt.Fprintln(flagSet.Output(), err)
		flagSet.Usage()
		return err
	}

	return nil
}

func validateArgs(config *argConfig) error {

	// parse the provided url (using net's url package)
	u, err := url.Parse(config.url)
	if err != nil {
		return fmt.Errorf("invalid value %q for url: %w", config.url, err)
	}

	if config.url == "" || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid value %q for url: requires a valid url with a scheme and host", config.url)
	}

	if config.c > config.n {
		return fmt.Errorf("value for flag -c(=%d) can not be greater than the value for flag -n(=%d)", config.c, config.n)
	}

	return nil
}

// define a positive int type that implements flag's Value interface
// (in order to force a custom type checking for an int flag)

type PositiveInt int

func (p *PositiveInt) String() string {
	return strconv.Itoa(int(*p))
}

func (p *PositiveInt) Set(s string) error {

	// parse the string to an int
	i, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return err
	}

	// return error if i is not a positive int
	if i <= 0 {
		return errors.New("value should be greater than 0")
	}

	*p = PositiveInt(i)
	return nil
}

// a helper function to wrap a pointer to int to a pointer to a PositiveInt
func asPositiveInt(i *int) *PositiveInt {
	return (*PositiveInt)(i) // the conversion works because both int and PositiveInt share same underlying type
}
