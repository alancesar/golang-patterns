package fetcher

import (
	"sync"
)

type (
	Event struct {
		Name string
		Data interface{}
	}

	ReducerFn func(locker sync.Locker, event Event)

	Fetcher struct {
		reducerFn   ReducerFn
		dispatchers []<-chan Event
	}
)

func New(reducer ReducerFn) *Fetcher {
	return &Fetcher{
		reducerFn: reducer,
	}
}

func (f *Fetcher) AddProducer(name string, fn func() interface{}) *Fetcher {
	producer := f.createProducer(name, fn)
	f.dispatchers = append(f.dispatchers, producer)
	return f
}

func (f Fetcher) Fetch(locker sync.Locker) {
	incoming := f.merge(f.dispatchers)

	reduce := func(event Event) {
		defer locker.Unlock()
		locker.Lock()
		f.reducerFn(locker, event)
	}

	for input := range incoming {
		reduce(input)
	}
}

func (Fetcher) createProducer(name string, fn func() interface{}) <-chan Event {
	c := make(chan Event)

	go func() {
		defer close(c)
		c <- Event{
			Name: name,
			Data: fn(),
		}
	}()

	return c
}

func (Fetcher) merge(dispatchers []<-chan Event) <-chan Event {
	var wg sync.WaitGroup
	wg.Add(len(dispatchers))
	dispatcher := make(chan Event)

	fn := func(incoming <-chan Event) {
		for in := range incoming {
			dispatcher <- in
		}
		wg.Done()
	}

	for _, c := range dispatchers {
		go fn(c)
	}

	go func() {
		wg.Wait()
		close(dispatcher)
	}()

	return dispatcher
}
