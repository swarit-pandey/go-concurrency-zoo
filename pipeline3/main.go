package main

import (
	"fmt"
	"sync"
)

func populateMap(m *sync.Map, keys []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, key := range keys {
		m.Store(key, len(keys))
	}
}

func processData(m *sync.Map, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	m.Range(func(key, value any) bool {
		strKey := key.(string)
		intValue := value.(int)
		result := intValue * 2
		fmt.Printf("processing key: %s, value %d, result: %d\n", strKey, intValue, result)
		resultChan <- result
		return true
	})
}

func furtherProcessing(m *sync.Map, processedChan <-chan int, finalChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	total := 0

	for result := range processedChan {
		total += result
	}

	m.Range(func(key, _ any) bool {
		strKey := key.(string)
		finalChan <- fmt.Sprintf("Key: %s, total processed result: %d", strKey, total)
		return true
	})
}

func main() {
	var sharedMap sync.Map

	keys := []string{"apple", "banana", "cherry", "date", "elderberry"}

	resultChan := make(chan int)
	finalChan := make(chan string)

	var wg sync.WaitGroup

	wg.Add(1)
	go populateMap(&sharedMap, keys, &wg)

	wg.Add(1)
	go processData(&sharedMap, resultChan, &wg)

	// close the result chan after processing is done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	wg.Add(1)
	go furtherProcessing(&sharedMap, resultChan, finalChan, &wg)

	go func() {
		wg.Wait()
		close(finalChan)
	}()

	for result := range finalChan {
		fmt.Println(result)
	}
}
