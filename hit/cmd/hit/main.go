package main

import (
	"flag"
	"fmt"
	"os"
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

	flagSet.StringVar(
		&config.url,
		"url",
		"",
		"http server `URL` to send requests to (required)") // text between backticks (`URL`) will be shown as a type of this flag (i.e. -url URL) in the usage message
	flagSet.IntVar(&config.c, "c", 1, "concurrency level")
	flagSet.IntVar(&config.n, "n", 1000, "number of requests to send")
	flagSet.IntVar(&config.rps, "rps", 100, "requests per second")

	return flagSet.Parse(args)
}
