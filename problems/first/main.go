package main

import (
	"log/slog"
	"sync"
)

// Create a simple program that launches multiple goroutines to calculate the square of numbers and send the results through a channel.
func square(wg *sync.WaitGroup, number int, squareChan chan<- int) {
	defer wg.Done()

	res := number * number

	squareChan <- res
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	squareChan := make(chan int, len(numbers))
	// This is vim mode
	var wg sync.WaitGroup
	for _, number := range numbers {
		wg.Add(1)
		go square(&wg, number, squareChan)
	}

	go func() {
		wg.Wait()
		close(squareChan)
	}()

	for val := range squareChan {
		slog.Info("squared", slog.Int("val", val))
	}
}
