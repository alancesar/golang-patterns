package fetcher

import (
	"sync"
)

type (
	incoming struct {
		index   int
		payload interface{}
	}

	ProducerFn func() interface{}
	ConsumerFn func(interface{})

	Fetcher struct {
		producers []<-chan incoming
		consumers []ConsumerFn
	}
)

func New() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) With(producerFn ProducerFn, consumerFn ConsumerFn) *Fetcher {
	producer := f.prepare(producerFn)
	f.producers = append(f.producers, producer)
	f.consumers = append(f.consumers, consumerFn)
	return f
}

func (f Fetcher) Fetch(target sync.Locker) {
	consumers := f.fanIn(f.producers)

	for input := range consumers {
		f.safeMerge(target, input)
	}
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
		c <- incoming{
			index:   index,
			payload: fn(),
		}
	}()

	return c
}

func (Fetcher) fanIn(producers []<-chan incoming) <-chan incoming {
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
	}()

	return out
}
