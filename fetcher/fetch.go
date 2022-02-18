package fetcher

import (
	"sync"
)

type (
	incoming struct {
		index   int
		payload interface{}
	}

	ProducerFn func() (interface{}, error)
	ConsumerFn func(interface{})

	Fetcher struct {
		producers []<-chan incoming
		consumers []ConsumerFn
		errChan   chan error
	}
)

func New() *Fetcher {
	return &Fetcher{
		errChan: make(chan error, 1),
	}
}

func (f *Fetcher) With(producerFn ProducerFn, consumerFn ConsumerFn) *Fetcher {
	producer := f.prepare(producerFn)
	f.producers = append(f.producers, producer)
	f.consumers = append(f.consumers, consumerFn)
	return f
}

func (f Fetcher) Fetch(target sync.Locker) error {
	var (
		wg  sync.WaitGroup
		err error
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = <-f.errChan
	}()

	go func() {
		defer wg.Done()
		consumers := f.fanIn(f.producers)
		for input := range consumers {
			f.safeMerge(target, input)
		}
	}()

	wg.Wait()
	return err
}

func (f Fetcher) safeMerge(target sync.Locker, incoming incoming) {
	defer target.Unlock()
	target.Lock()
	f.consumers[incoming.index](incoming.payload)
}

func (f Fetcher) prepare(fn ProducerFn) <-chan incoming {
	c := make(chan incoming)
	index := len(f.producers)

	go func() {
		defer close(c)
		payload, err := fn()
		if err != nil {
			f.errChan <- err
			return
		}

		c <- incoming{
			index:   index,
			payload: payload,
		}
	}()

	return c
}

func (f Fetcher) fanIn(producers []<-chan incoming) <-chan incoming {
	var wg sync.WaitGroup
	wg.Add(len(producers))
	out := make(chan incoming)

	fn := func(incoming <-chan incoming) {
		for in := range incoming {
			out <- in
		}
		wg.Done()
	}

	for _, c := range producers {
		go fn(c)
	}

	go func() {
		wg.Wait()
		close(out)
		close(f.errChan)
	}()

	return out
}
