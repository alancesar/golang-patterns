package main

import (
	"context"
	"fmt"
	"golang-patterns/internal/sleep"
	"golang-patterns/worker"
	"log"
)

type (
	Download struct {
		url string
	}
)

func (d Download) Download(_ context.Context) {
	log.Printf("starting download from %s\n", d.url)
	sleep.Random()
	log.Printf("dowloaded from %s successfuly\n", d.url)
}

func NewDownloader(url string) Download {
	return Download{
		url: url,
	}
}

func main() {
	n := 20
	items := make([]interface{}, n)

	for i := 0; i < n; i++ {
		items[i] = NewDownloader(fmt.Sprintf("https://some-item.com?id=%d", i+1))
	}

	workerFn := func(ctx context.Context, item interface{}) {
		item.(Download).Download(ctx)
	}

	w := worker.New(workerFn, 5)
	w.Work(context.Background(), items)
}
