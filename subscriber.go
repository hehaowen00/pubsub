package pubsub

import "sync"

type Subscriber[T any] struct {
	id int64

	out  []T
	recv chan struct{}
	mu   sync.Mutex

	topic *topic[T]
}

func newSubscriber[T any](t *topic[T]) *Subscriber[T] {
	return &Subscriber[T]{
		recv:  make(chan struct{}, 1),
		topic: t,
	}
}

func (s *Subscriber[T]) send(msg T) {
	s.mu.Lock()
	s.out = append(s.out, msg)
	s.mu.Unlock()

	if len(s.recv) == 0 {
		s.recv <- struct{}{}
	}
}

func (s *Subscriber[T]) Recv() <-chan struct{} {
	return s.recv
}

func (s *Subscriber[T]) Read() []T {
	data := s.out
	s.out = nil
	return data
}

func (s *Subscriber[T]) Unsubscribe() {
	s.topic.unsubscribe(s)
}
