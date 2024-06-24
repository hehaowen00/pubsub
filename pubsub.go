package pubsub

import (
	"errors"
	"sync"
)

var (
	topics = map[string]any{}
	rw     = sync.RWMutex{}

	ErrTopicDoesNotExist = errors.New("topic does not exist")
	ErrTopicExists       = errors.New("topic exists")
	ErrTypeMismatch      = errors.New("type mismatch")
)

func NewTopic[T any](name string, history ...ICache[T]) error {
	return newTopic[T](name)
}

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

func NewPublisher[T any](name string) (*Publisher[T], error) {
	rw.RLock()
	t, ok := topics[name]
	rw.RUnlock()

	if !ok {
		// t = newTopic[T](name)
		return nil, ErrTopicDoesNotExist
	}

	v, ok := t.(*topic[T])
	if !ok {
		return nil, ErrTypeMismatch
	}

	p := &Publisher[T]{
		t: v,
	}

	return p, nil
}

func NewSubscriber[T any](name string) (*Subscriber[T], error) {
	rw.RLock()
	t, ok := topics[name]
	rw.RUnlock()

	if !ok {
		// t = newTopic[T](name)
		return nil, ErrTopicDoesNotExist
	}

	v, ok := t.(*topic[T])
	if !ok {
		return nil, ErrTypeMismatch
	}

	sub := v.subscribe()

	return sub, nil
}
