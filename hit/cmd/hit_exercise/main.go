package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type PositiveDuration time.Duration
type HttpMethod string

type config struct {
	flag    int           // flag
	timeout time.Duration // timeout duration
	method  string        // http method
}

func main() {

	args := os.Args[1:]
	c := config{}
	if err := parseArgs(&c, args); err != nil {
		os.Exit(1)
	}
	fmt.Println(c)
}

func parseArgs(c *config, args []string) error {

	fs := flag.NewFlagSet("hit_test", flag.ContinueOnError)
	fs.IntVar(&c.flag, "flag", 100, "an int flag")

	fs.Var(asPositiveDuration(&c.timeout), "timeout", "a positive timeout value (in seconds)")
	fs.Var(asHttpMethod(&c.method), "method", "http method to use: [GET, POST, PUT]")

	err := fs.Parse(args)

	// print the position args
	fmt.Println(fs.Args())

	return err
}

// construct positive duration from duration
func asPositiveDuration(d *time.Duration) *PositiveDuration {
	return (*PositiveDuration)(d)
}

func (pd *PositiveDuration) Set(val string) error {

	i, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return err
	}

	if i <= 0 {
		return fmt.Errorf("should be greater than 0")
	}

	*pd = PositiveDuration(i)
	return nil
}

func (pd *PositiveDuration) String() string {
	return strconv.Itoa(int(*pd))
}

func asHttpMethod(val *string) *HttpMethod {
	return (*HttpMethod)(val)
}

func (hm *HttpMethod) Set(val string) error {
	method := strings.ToUpper(val)
	if method != "GET" && method != "POST" && method != "PUT" {
		return fmt.Errorf("should be one of [GET, POST, PUT]")
	}
	*hm = HttpMethod(val)
	return nil
}

func (hm *HttpMethod) String() string {
	return string(*hm)
}
