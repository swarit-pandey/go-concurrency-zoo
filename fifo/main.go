package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func readURLs(filePath string) (<-chan string, error) {
	out := make(chan string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()
			out <- url
		}
		close(out)
	}()
	return out, nil
}

func fetchData(data <-chan string, result chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range data {
		resp, err := http.Get(val)
		if err != nil {
			fmt.Printf("failed to get data via http: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read http content: %v", err)
			resp.Body.Close()
			continue
		}

		result <- string(body)
	}
}

// fanIn takes variadic argument for result chan and then
// spawns a new goroutine for each of the chan
func fanIn(resultChans ...<-chan string) <-chan string {
	out := make(chan string)
	var wg *sync.WaitGroup
	for _, ch := range resultChans {
		go func(c <-chan string) {
			defer wg.Done()
			for res := range c {
				out <- res
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	const numWorkers = 4
	const filePath = "urls.txt"

	// Stage 1: Read URLs from a file
	urls, err := readURLs(filePath)
	if err != nil {
		fmt.Printf("Failed to read URLs from file: %v\n", err)
		return
	}

	// Create channels for distributing URLs to workers and collecting results
	urlChan := make(chan string)
	resultChans := make([]chan string, numWorkers)
	for i := range resultChans {
		resultChans[i] = make(chan string)
	}

	// Create a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Stage 2: Start worker goroutines to fetch content (Fan-Out)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go fetchData(urlChan, resultChans[i], &wg)
	}

	// Stage 3: Distribute URLs to the workers
	go func() {
		for url := range urls {
			urlChan <- url
		}
		close(urlChan)
	}()

	// Convert resultChans to a slice of receive-only channels for fanIn
	receiveOnlyResultChans := make([]<-chan string, numWorkers)
	for i, ch := range resultChans {
		receiveOnlyResultChans[i] = ch
	}

	// Stage 4: Collect results from all workers (Fan-In)
	combinedResults := fanIn(receiveOnlyResultChans...)

	// Stage 5: Process the combined content
	go func() {
		wg.Wait()
		for _, ch := range resultChans {
			close(ch)
		}
	}()

	for result := range combinedResults {
		fmt.Println("Fetched content length:", len(result))
	}
}
