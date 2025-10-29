package url

import "testing"

// benchmarking URL's String() function
func BenchmarkURLString(b *testing.B) {
	url := &URL{"http", "www.dummyurl.com", "mypage"}
	for b.Loop() {
		_ = url.String()
	}
}
