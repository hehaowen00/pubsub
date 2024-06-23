package pubsub

import (
	"context"
	"runtime"
	"slices"
	"sync"
)

type topic[T any] struct {
	id      int64
	history *Ring[T]

	recv        chan T
	subscribers []*Subscriber[T]

	ctx    context.Context
	cancel context.CancelFunc

	subEvents   chan *Subscriber[T]
	unsubEvents chan *Subscriber[T]

	once sync.Once
}

func newTopic[T any](name string) *topic[T] {
	rw.Lock()

	v, ok := topics[name]
	if !ok {
		ctx, cancel := context.WithCancel(context.Background())
		t := &topic[T]{
			history: newRing[T](16),
			recv:    make(chan T),

			ctx:    ctx,
			cancel: cancel,

			subEvents:   make(chan *Subscriber[T]),
			unsubEvents: make(chan *Subscriber[T]),
		}

		go t.run()
		topics[name] = t

		rw.Unlock()

		return t
	}

	rw.Unlock()

	t, ok := v.(*topic[T])
	if !ok {
		panic("invalid type")
	}

	return t
}

func (t *topic[T]) run() {
outer:
	for {
		select {
		case <-t.ctx.Done():
			break outer
		case s, ok := <-t.subEvents:
			if !ok {
				break outer
			}

			t.id += 1
			s.id = t.id

			t.history.Iter(func(i T) {
				s.recv <- i
			})

			t.subscribers = append(t.subscribers, s)
		case s, ok := <-t.unsubEvents:
			if !ok {
				break outer
			}

			t.subscribers = slices.DeleteFunc(t.subscribers, func(e *Subscriber[T]) bool {
				return e.id == s.id
			})

			if len(t.subscribers) == 0 {
				t.id = 0
			}
		case m, ok := <-t.recv:
			if !ok {
				break outer
			}

			t.history.Push(m)

			for _, s := range t.subscribers {
				s.send(m)
			}
		}
	}

	for _, s := range t.subscribers {
		close(s.recv)
	}
}

func (t *topic[T]) publish(msg T) error {
	err := t.ctx.Err()
	if err != nil {
		return err
	}

	t.recv <- msg
	runtime.Gosched()

	return nil
}

func (t *topic[T]) close() {
	t.once.Do(func() {
		t.cancel()
		close(t.recv)
		close(t.subEvents)
		close(t.unsubEvents)
	})
}

func (t *topic[T]) subscribe() *Subscriber[T] {
	sub := newSubscriber(t)
	t.subEvents <- sub
	return sub
}

func (t *topic[T]) unsubscribe(sub *Subscriber[T]) {
	t.unsubEvents <- sub
}
