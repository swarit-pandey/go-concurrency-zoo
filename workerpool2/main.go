package main

import (
	"fmt"
	"sync"
)

type Result struct {
	Number  int
	IsPrime bool
}

func checkPrime(n int) bool {
	if n <= 1 {
		return false
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func worker(id int, numbers <-chan int, result chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for nums := range numbers {
		prime := checkPrime(nums)
		res := Result{
			Number:  nums,
			IsPrime: prime,
		}
		fmt.Printf("Goroutine working, id: %v, result, %v\n", id, res)
		result <- res
	}
}

func main() {
	const numWorkers = 4
	numbers := []int{2, 3, 4, 5, 10, 13, 17, 19, 23, 24, 29, 31, 37, 41, 43, 47, 53, 59}

	numberChan := make(chan int, len(numbers))
	resultChan := make(chan Result, len(numbers))

	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, numberChan, resultChan, &wg)
	}

	for _, number := range numbers {
		numberChan <- number
	}
	close(numberChan)

	wg.Wait()
	close(resultChan)

	fmt.Println("Got numbers")
	for res := range resultChan {
		fmt.Printf("Got %v\n", res)
	}
}
