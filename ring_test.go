package pubsub

import (
	"testing"
)

func TestRing(t *testing.T) {
	ring := newRing[int](16)
	for i := range 10 {
		ring.Push(i)
	}

	ring.Iter(func(i int) {
		t.Log(i)
	})

	t.Log(ring.head, ring.data)
}
