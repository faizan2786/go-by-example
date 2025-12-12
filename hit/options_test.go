package hit

import "testing"

func TestDefaults(t *testing.T) {

	op := DefaultOptions()

	if op.Concurrency != 1 {
		t.Errorf("Concurrency = %d; want %d\f", op.Concurrency, 1)
	}

	if op.RPS != 0 {
		t.Errorf("RPS = %d; want %d\f", op.RPS, 0)
	}

	if op.Send == nil {
		t.Errorf("Send = <nil>; want a valid function of type %T\n", op.Send)
	}
}

func TestDefaultsForValidInputs(t *testing.T) {

	op := Options{RPS: 5, Concurrency: 2}
	op = withDefaults(op)

	if op.Concurrency != 2 {
		t.Errorf("Concurrency = %d; want %d\f", op.Concurrency, 2)
	}

	if op.RPS != 5 {
		t.Errorf("RPS = %d; want %d\f", op.RPS, 5)
	}
}

func TestDefaultsForInvalidInputs(t *testing.T) {

	op := Options{RPS: -4, Concurrency: 0}
	op = withDefaults(op)

	if op.Concurrency != 1 {
		t.Errorf("Concurrency = %d; want %d\f", op.Concurrency, 1)
	}

	if op.RPS != 0 {
		t.Errorf("RPS = %d; want %d\f", op.RPS, 0)
	}
}
