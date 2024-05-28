package pubsub

type Subscriber[T any] struct {
	id    int64
	recv  chan T
	topic *topic[T]
}

func newSubscriber[T any](t *topic[T]) *Subscriber[T] {
	return &Subscriber[T]{
		recv:  make(chan T, 32),
		topic: t,
	}
}

func (s *Subscriber[T]) send(msg T) {
	s.recv <- msg
}

func (s *Subscriber[T]) Recv() <-chan T {
	return s.recv
}

func (s *Subscriber[T]) Unsubscribe() {
	s.topic.unsubscribe(s)
}
