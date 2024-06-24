package pubsub

import (
	"testing"
)

func TestRing(t *testing.T) {
	ring := newRing[int](16)
	for i := range 10 {
		ring.push(i)
	}

	ring.iter(func(i int) {
		t.Log(i)
	})

	t.Log(ring.head, ring.data)
}
