package main

import (
	"sync"
)

func main() {
	const n = 1000
	var wg sync.WaitGroup
	var mu sync.Mutex
	current := 1

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			for {
				mu.Lock()
				if current == num {
					current++
					mu.Unlock()
					break
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
}
