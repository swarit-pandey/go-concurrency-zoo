package main

import (
	"fmt"
	"sync"
)

func worker(id int, numbers <-chan int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range numbers {
		res := num * num
		fmt.Printf("Goroutine %d working on number %v\n", id, num)
		result <- res
	}
}

func main() {
	const workers = 3
	nums := []int{2, 4, 6, 8, 10}

	numChan := make(chan int, len(nums))
	resChan := make(chan int, len(nums))
	var wg sync.WaitGroup

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(i, numChan, resChan, &wg)
	}

    // Push to the numChan which holds the numbers
    // that we need to process
	for _, number := range nums {
		numChan <- number
	}
	// Since the producer has finished
	// sending data we can close the numchan
	close(numChan)

	wg.Wait()
	// after all goroutines have finished working
	// close the result chan
	close(resChan)

    // Read from resChan which holds the results 
	fmt.Println("Results: ")
	for res := range resChan {
		fmt.Println(res)
	}
}
