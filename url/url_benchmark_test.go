package url

import (
	"fmt"
	"strings"
	"testing"
)

// benchmarking URL's String() function
func BenchmarkURLStringShort(b *testing.B) {
	url := &URL{"https", "myurl.com", "myblog"}
	for b.Loop() {
		_ = url.String()
	}
}

func BenchmarkURLStringLong(b *testing.B) {

	scheme := "https"
	hostUnit := "h"
	pathUnit := "p"

	lengths := []int{10, 100, 1000, 10000}
	for _, n := range lengths {
		host := strings.Repeat(hostUnit, n)
		path := strings.Repeat(pathUnit, n)
		url := &URL{scheme, host, path}

		b.Run(fmt.Sprintf("Bench%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = url.String()
			}
		})
	}
}
