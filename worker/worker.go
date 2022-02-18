package worker

import (
	"context"
	"sync"
)

type (
	Fn func(context.Context, interface{})

	Worker struct {
		fn     Fn
		buffer int
	}
)

func New(workerFn Fn, bufferSize int) *Worker {
	return &Worker{
		fn:     workerFn,
		buffer: bufferSize,
	}
}

func (w Worker) Work(ctx context.Context, items []interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(items))

	jobs := make(chan interface{})

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
