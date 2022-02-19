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

func NewDownload(url string) Download {
	return Download{
		url: url,
	}
}

func main() {
	n := 20
	downloads := make([]Download, n)

	for i := 0; i < n; i++ {
		downloads[i] = NewDownload(fmt.Sprintf("https://some-item.com?id=%d", i+1))
	}

	downloader := func(ctx context.Context, download Download) {
		download.Download(ctx)
	}

	w := worker.New(downloader, 5)
	w.Work(context.Background(), downloads)
}
