package pubsub

type Subscriber[T any] struct {
	id    int64
	recv  chan T
	topic *Topic[T]
}

func newSubscriber[T any](topic *Topic[T]) *Subscriber[T] {
	return &Subscriber[T]{
		recv:  make(chan T, 32),
		topic: topic,
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
