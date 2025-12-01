package url

import "testing"

var counter int = 0

func incr() {
	counter++
}

// if you run below tests in parallel, they will fail
// because they both modify the same global var counter

func TestGlobalAccess1(t *testing.T) {

	//t.Parallel()
	incr()
	if counter != 1 {
		t.Errorf("Expected counter = %d, got %d", 1, counter)
	}
}

func TestGlobalAccess2(t *testing.T) {

	//t.Parallel()
	incr()
	incr()
	if counter != 3 {
		t.Errorf("Expected counter = %d, got %d", 3, counter)
	}
}
