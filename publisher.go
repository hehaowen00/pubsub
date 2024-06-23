package pubsub

type Publisher[T any] struct {
	t *topic[T]
}

func (p *Publisher[T]) Publish(msg T) error {
	return p.t.publish(msg)
}
