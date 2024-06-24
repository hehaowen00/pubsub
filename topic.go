package pubsub

import (
	"context"
	"runtime"
	"slices"
	"sync"
)

type topic[T any] struct {
	id      int64
	history ICache[T]

	recv        chan T
	subscribers []*Subscriber[T]

	ctx    context.Context
	cancel context.CancelFunc

	subEvents   chan *Subscriber[T]
	unsubEvents chan *Subscriber[T]

	once sync.Once
}

func newTopic[T any](name string, historyStrategy ...ICache[T]) error {
	rw.Lock()

	history := NewEmptyCache[T]()

	if len(historyStrategy) == 1 {
		history = historyStrategy[0]
	}

	_, ok := topics[name]
	if !ok {
		ctx, cancel := context.WithCancel(context.Background())
		t := &topic[T]{
			history: history,
			recv:    make(chan T),

			ctx:    ctx,
			cancel: cancel,

			subEvents:   make(chan *Subscriber[T]),
			unsubEvents: make(chan *Subscriber[T]),
		}

		go t.run()
		topics[name] = t

		rw.Unlock()

		return nil
	}

	rw.Unlock()

	return ErrTopicExists
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

			s.mu.Lock()
			t.history.Iter(func(i T) {
				s.out = append(s.out, i)
			})
			s.mu.Unlock()

			t.subscribers = append(t.subscribers, s)
		case s, ok := <-t.unsubEvents:
			if !ok {
				break outer
			}

			index := slices.IndexFunc(t.subscribers, func(sub *Subscriber[T]) bool {
				return sub.id == s.id
			})

			if index == -1 {
				continue
			}

			t.subscribers = slices.Delete(t.subscribers, index, 1)

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
