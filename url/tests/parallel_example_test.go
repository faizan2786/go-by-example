package url

import (
	"testing"
	"time"
)

func TestParallel1(t *testing.T) {
	t.Parallel()
	time.Sleep(5 * time.Second) // simulate a long running test
}

func TestParallel2(t *testing.T) {
	t.Parallel()
	time.Sleep(5 * time.Second) // simulate a long running test
}

func TestParallel3(t *testing.T) {
	t.Parallel()
	t.Run("SubTest1", func(t *testing.T) {
		t.Parallel()
		time.Sleep(5 * time.Second) // simulate a long running test
	})

	t.Run("SubTest2", func(t *testing.T) {
		t.Parallel()
		time.Sleep(5 * time.Second) // simulate a long running test
	})
}
