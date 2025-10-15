package main

import (
	"fmt"

	"github.com/faizan2786/gobyexample/url"
)

func main() {
	url, _ := url.Parse("www.dummyurl.com")
	fmt.Println(url)
}
