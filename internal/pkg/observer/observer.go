package observer

import (
	"context"
	"time"
)

type EventHandler[T any] func(ctx context.Context, events []T)

type Observer[T any] struct {
	queue   *queue[T]
	handler *handler[T]
}

func NewObserver[T any](fn EventHandler[T], tickDuration *time.Duration) *Observer[T] {
	queue := newQueue[T]()

	o := &Observer[T]{
		queue:   queue,
		handler: newHandler(queue, fn).withTick(*tickDuration),
	}
	return o
}

func (o *Observer[T]) Start(ctx context.Context) {
	go o.handler.listen(ctx)
}

func (o *Observer[T]) Dispatch(event T) {
	o.queue.Enqueue(event)
}

func (o *Observer[T]) Flush() {
	o.handler.flush()
}

func (o *Observer[T]) Wait(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		o.handler.flushAndWait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-done:
		return
	}
}
