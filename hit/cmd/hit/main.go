package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
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

	// any args that comes AFTER the flags are "positional arguments" and can be
	// retrieved by arg[i] method after parsing the args by FlagSet

	// since the positional args are retrieved directly from the command line args (without a parser),
	// we need to set the usage message manually to include the positional args in the message
	flagSet.Usage = func() {

		// print the usage message
		fmt.Fprintf(
			flagSet.Output(),
			"usage: %s [options] url\noptions:\n",
			flagSet.Name(), // name of the flagset (defined earlier)
		)

		// print the default values for the flags
		flagSet.PrintDefaults()
	}

	// register all of the int type flags using our custom PositiveInt Value type (parser)
	flagSet.Var(asPositiveInt(&config.c), "c", "concurrency level")
	flagSet.Var(asPositiveInt(&config.n), "n", "number of requests to send")
	flagSet.Var(asPositiveInt(&config.rps), "rps", "requests per second")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// retrieve the 1st positional argument (i.e. url)
	// (Since the positional arguments don't have named flags (i.e. -flagname), their values are accessed directly by its position)
	config.url = flagSet.Arg(0) // returns empty string if there are no positional args provided

	// validate the provided argument values
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
