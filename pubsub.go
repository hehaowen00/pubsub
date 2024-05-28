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

func Close[T any](name string) {
	rw.Lock()

	t, ok := topics[name]
	if !ok {
		return
	}

	v, ok := t.(*topic[T])
	if !ok {
		return
	}

	v.close()

	delete(topics, name)

	rw.Unlock()
}

func Publish[T any](name string, msg T) error {
	rw.RLock()
	t, ok := topics[name]
	rw.RUnlock()

	if !ok {
		t = newTopic[T](name)
	}

	v, ok := t.(*topic[T])
	if !ok {
		return ErrTypeMismatch
	}

	v.publish(msg)

	return nil
}

func Subscribe[T any](name string) (*Subscriber[T], error) {
	rw.RLock()
	t, ok := topics[name]
	rw.RUnlock()

	if !ok {
		t = newTopic[T](name)
	}

	v, ok := t.(*topic[T])
	if !ok {
		return nil, ErrTypeMismatch
	}

	sub := v.subscribe()

	return sub, nil
}
