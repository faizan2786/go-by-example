package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
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

func main() {

	config := argConfig{
		n:   1000,
		c:   1,
		rps: 100,
	}
	args := os.Args[1:] // get the command line args from the os package (ignore 1st arg - the program name)

	if err := parseArgs(args, &config); err != nil {
		os.Exit(1)
	}

	fmt.Printf("%s\nSending %d requests to %q at %d requests per second (concurrency=%d)\n", logo, config.n, config.url, config.rps, config.c)
}

// function to parse command line args and assigned them to a config variable (using the flag package)
func parseArgs(args []string, config *argConfig) error {

	flagSet := flag.NewFlagSet("hit", flag.ContinueOnError)

	// register all the flags (by its type) with its name, usage message and default value

	// register "-url" string type flag using the StringVar Value parser provided by the FlagSet type
	flagSet.StringVar(
		&config.url,
		"url",
		"",
		"http server `URL` to send requests to (required)") // text between backticks (`URL`) will be shown as a type of this flag (i.e. -url URL) in the usage message

	// register rest of the int type flags using our custom PositiveInt Value type (parser)
	flagSet.Var(asPositiveInt(&config.c), "c", "concurrency level")
	flagSet.Var(asPositiveInt(&config.n), "n", "number of requests to send")
	flagSet.Var(asPositiveInt(&config.rps), "rps", "requests per second")

	return flagSet.Parse(args)
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
