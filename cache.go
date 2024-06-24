package pubsub

type ICache[T any] interface {
	Push(item T)
	Iter(func(T))
}

type EmptyCache[T any] struct {
}

func NewEmptyCache[T any]() ICache[T] {
	return &EmptyCache[T]{}
}

func (h *EmptyCache[T]) Push(item T) {
}

func (h *EmptyCache[T]) Iter(f func(T)) {
}

type RingCache[T any] struct {
	internal *ring[T]
}

func NewRingCache[T any](bufferSize ...int) ICache[T] {
	size := 16
	if len(bufferSize) == 1 {
		size = bufferSize[0]
	}

	ring := newRing[T](size)

	return &RingCache[T]{
		internal: ring,
	}
}

func (r *RingCache[T]) Push(item T) {
	r.internal.push(item)
}

func (r *RingCache[T]) Iter(f func(T)) {
	r.internal.iter(f)
}
