package main

import "fmt"

func getNumbers(numbers []int) <-chan int {
	out := make(chan int)

	go func() {
		for _, nums := range numbers {
			out <- nums
		}
		close(out)
	}()

	return out
}

func squareNumbers(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for nums := range in {
			out <- nums * nums
		}
		close(out)
	}()

	return out
}

func sumNumbers(in <-chan int) int {
	sum := 0

	for nums := range in {
		sum += nums
	}

	return sum
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}

	// Stage 1: Generate numbers
	numChan := getNumbers(numbers)

	// Stage 2: Square numbers
	squareChan := squareNumbers(numChan)

	// Stage 3: Sum numbers
	final := sumNumbers(squareChan)

	fmt.Printf("Got final sum: %v", final)
}
