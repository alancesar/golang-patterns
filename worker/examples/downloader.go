package main

import (
	"context"
	"fmt"
	"golang-patterns/worker"
	"log"
	"math/rand"
	"time"
)

type (
	Download struct {
		url string
	}
)

func randomSleep() {
	ms := rand.Intn(50) * 100
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func (d Download) Download(_ context.Context) {
	log.Printf("starting download from %s\n", d.url)
	randomSleep()
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
