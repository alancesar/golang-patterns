package worker

import (
	"context"
	"sync"
)

type (
	Fn[T any] func(context.Context, T)

	Worker[T any] struct {
		fn     Fn[T]
		buffer int
	}
)

func New[T any](workerFn Fn[T], bufferSize int) *Worker[T] {
	return &Worker[T]{
		fn:     workerFn,
		buffer: bufferSize,
	}
}

func (w Worker[T]) Work(ctx context.Context, items []T) {
	var wg sync.WaitGroup
	wg.Add(len(items))

	jobs := make(chan T)

	for i := 1; i <= w.buffer; i++ {
		go func() {
			for j := range jobs {
				w.fn(ctx, j)
				wg.Done()
			}
		}()
	}

	for _, item := range items {
		jobs <- item
	}

	close(jobs)
	wg.Wait()
}
