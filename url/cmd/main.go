package main

import (
	"fmt"

	"github.com/faizan2786/gobyexample/url"
)

func main() {
	url, _ := url.Parse("http://www.dummyurl.com")
	fmt.Println(url)
}
