// This program defines a custom parser to parse the command line args
// (instead of using the "flag" package)

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

const usage = `Usage:
  -url
       HTTP server URL (required)
  -n
       Number of requests
  -c
       Concurrency level
  -rps
       Requests per second
`

type parseFunc func(string) error // define a value parser function type

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

func parseArgs(args []string, config *argConfig) error {

	// define a map of arg names to its parser function
	argMap := map[string]parseFunc{
		"url": strParser(&config.url),
		"n":   intParser(&config.n),
		"c":   intParser(&config.c),
		"rps": intParser(&config.rps),
	}

	for _, arg := range args {
		// here arg represents a command line argument in a form "name=value" as a single string
		name, value, _ := strings.Cut(arg, "=")
		name = strings.TrimPrefix(name, "-")

		parser, ok := argMap[name] // try access the parser from the map
		if !ok {
			return fmt.Errorf("flag provided but not defined: -%s", name)
		}
		// set the argument value using the parser function
		if err := parser(value); err != nil {
			return fmt.Errorf("invalid value %q for flag '-%s': %w", value, name, err)
		}
	}
	return nil
}

// this is a higher-order function that wraps a string parser function as a closure and returns it
// (the returned closure will then be able access the variable pointer v passed to its parent function)
func strParser(v *string) parseFunc {
	return func(s string) error {
		*v = s
		return nil
	}
}

// this is a higher-order function that wraps an int parser function as a closure and returns it
func intParser(v *int) parseFunc {
	return func(s string) error {
		var err error
		*v, err = strconv.Atoi(s)
		return err
	}
}
