package pubsub

type Ring[T any] struct {
	data []T
	head int
	tail int
	size int
}

func newRing[T any](size int) *Ring[T] {
	return &Ring[T]{
		data: make([]T, size),
		size: size,
	}
}

func (r *Ring[T]) Push(item T) {
	r.data[r.head] = item
	r.head += 1
	r.head = r.head % r.size

	if r.head == r.tail {
		r.tail += 1
	}
}

func (r *Ring[T]) Iter(f func(T)) {
	if r.head == r.tail {
		return
	}

	tail := r.tail
	if r.head < r.tail {
		tail = r.head
	}

	for {
		f(r.data[tail])

		tail += 1
		tail = tail % r.size

		if tail == r.head {
			break
		}
	}
}
