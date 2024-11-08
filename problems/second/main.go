package main

import (
	"bufio"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

// Implement a worker pool that processes a large number of tasks concurrently using
// a fixed number of worker goroutines.

type result struct {
	url  string
	size int
	err  error
}

type worker struct {
	filename string
	size     int
	tasks    chan string
	results  chan result
	wg       sync.WaitGroup
}

func newPool(filename string, size int) *worker {
	return &worker{
		filename: filename,
		size:     size,
		results:  make(chan result),
		tasks:    make(chan string),
	}
}

func (w *worker) worker() {
	defer w.wg.Done()

	for url := range w.tasks {
		size, err := fetchURL(url)
		w.results <- result{url: url, size: size, err: err}
	}
}

func (w *worker) run() {
	// Add workers and start fetching from URL
	for i := 0; i < w.size; i++ {
		w.wg.Add(1)
		go w.worker()
	}

	// Start filling the tasks chan async
	go func() {
		file, err := os.Open(w.filename)
		if err != nil {
			slog.Error(err.Error())
			close(w.tasks)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			w.tasks <- scanner.Text()
		}
		close(w.tasks)
	}()

	go func() {
		w.wg.Wait()
		close(w.results)
	}()

	for result := range w.results {
		slog.Info("", slog.Any("result", result))
	}
}

func fetchURL(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return len(body), nil
}

func main() {
	const filename = "urls.txt"
	const numWorkers = 100 

	pool := newPool(filename, numWorkers)
	pool.run()
}
