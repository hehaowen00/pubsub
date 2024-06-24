package pubsub

type ring[T any] struct {
	data []T
	head int
	tail int
	size int
}

func newRing[T any](size int) *ring[T] {
	return &ring[T]{
		data: make([]T, size),
		size: size,
	}
}

func (r *ring[T]) push(item T) {
	r.data[r.head] = item
	r.head += 1
	r.head = r.head % r.size

	if r.head == r.tail {
		r.tail += 1
	}
}

func (r *ring[T]) iter(f func(T)) {
	if r.head == r.tail {
		return
	}

	tail := r.tail
	if r.head < r.tail {
		tail = r.head
	}

	for {
		item := r.data[tail]
		f(item)

		tail += 1
		tail = tail % r.size

		if tail == r.head {
			break
		}
	}
}
