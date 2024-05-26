package pubsub

import (
	"runtime"
	"slices"
	"sync"
)

type Topic[T any] struct {
	id int64

	recv        chan T
	subscribers []*Subscriber[T]

	subEvents   chan *Subscriber[T]
	unsubEvents chan *Subscriber[T]

	once sync.Once
}

func newTopic[T any](name string) *Topic[T] {
	t := &Topic[T]{
		recv:        make(chan T),
		subEvents:   make(chan *Subscriber[T]),
		unsubEvents: make(chan *Subscriber[T]),
	}

	go t.run()

	rw.Lock()
	topics[name] = t
	rw.Unlock()

	return t
}

func (t *Topic[T]) run() {
outer:
	for {
		select {
		case s, ok := <-t.subEvents:
			if !ok {
				break outer
			}

			t.id += 1
			s.id = t.id

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

			for _, s := range t.subscribers {
				s.send(m)
			}
		}
	}

	for _, s := range t.subscribers {
		close(s.recv)
	}
}

func (t *Topic[T]) publish(msg T) {
	t.recv <- msg
	runtime.Gosched()
}

func (t *Topic[T]) close() {
	t.once.Do(func() {
		close(t.recv)
		close(t.subEvents)
		close(t.unsubEvents)
	})
}

func (t *Topic[T]) subscribe() *Subscriber[T] {
	sub := newSubscriber(t)
	t.subEvents <- sub
	return sub
}

func (t *Topic[T]) unsubscribe(sub *Subscriber[T]) {
	t.unsubEvents <- sub
}
