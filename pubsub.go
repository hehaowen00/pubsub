package pubsub

import (
	"errors"
	"sync"
)

var (
	topics = map[string]any{}
	rw     = sync.RWMutex{}

	ErrTypeMismatch = errors.New("type mismatch")
)

func Close[T any](topic string) {
	rw.Lock()

	t, ok := topics[topic]
	if !ok {
		return
	}

	v, ok := t.(*Topic[T])
	if !ok {
		return
	}

	v.close()

	delete(topics, topic)

	rw.Unlock()
}

func Publish[T any](topic string, msg T) error {
	rw.RLock()
	t, ok := topics[topic]
	rw.RUnlock()

	if !ok {
		return nil
	}

	v, ok := t.(*Topic[T])
	if !ok {
		return ErrTypeMismatch
	}

	v.publish(msg)

	return nil
}

func Subscribe[T any](topic string) (*Subscriber[T], error) {
	rw.RLock()
	t, ok := topics[topic]
	rw.RUnlock()

	if !ok {
		t = newTopic[T](topic)
	}

	v, ok := t.(*Topic[T])
	if !ok {
		return nil, ErrTypeMismatch
	}

	sub := v.subscribe()

	return sub, nil
}
