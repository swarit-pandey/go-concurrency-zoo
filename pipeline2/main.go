package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func readURLs(path string) (<-chan string, error) {
	out := make(chan string)

	// Step 1: Read file from path
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Step 2: Start a goroutine
	go func() {
		defer file.Close()

		// Step 3: Start scanning the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()
			out <- url
		}
		close(out)
	}()

	return out, nil
}

func fetchContent(urls <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		var wg sync.WaitGroup
		for url := range urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()

				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("failed to fetch url %s: %v\n", url, err)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("failed to read data from response: %v", err)
					return
				}

				out <- string(body)
			}(url)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func processContent(contents <-chan string) <-chan int {
	out := make(chan int)

	go func() {
		for content := range contents {
			wordCount := len(strings.Fields(content))
			out <- wordCount
		}
		close(out)
	}()
	return out
}

func aggregateResult(wordCount <-chan int) int {
	totalWords := 0
	for count := range wordCount {
		totalWords += count
	}
	return totalWords
}

func main() {
	const filePath = "urls.txt"

	urls, err := readURLs(filePath)
	if err != nil {
		fmt.Printf("failed to read URLs from file: %v\n", err)
		return
	}

	contents := fetchContent(urls)

	wordCount := processContent(contents)

	totalWords := aggregateResult(wordCount)

	fmt.Printf("Total number of words from all the URLS: %d\n", totalWords)
}
