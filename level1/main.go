package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

type track struct {
	cache map[any]any
	mu    sync.Mutex
}

func NewTrack() *track {
	return &track{
		cache: make(map[any]any),
	}
}

func (t *track) write(put <-chan map[any]any, wg *sync.WaitGroup) {
	defer wg.Done()

	for kvMap := range put {
		t.mu.Lock()
		for k, v := range kvMap {
			t.cache[k] = v
		}
		t.mu.Unlock()
	}
}

func (t *track) read(get <-chan any, wg *sync.WaitGroup) chan any {
	resultChan := make(chan any)

	go func() {
		defer close(resultChan)
		defer wg.Done()

		for val := range get {
			t.mu.Lock()
			resultChan <- t.cache[val]
			t.mu.Unlock()
		}
	}()

	return resultChan
}

func main() {
	bigMap := generateRandomKeyVal()
	t := NewTrack()
	write := make(chan map[any]any)
	read := make(chan any)

	var wg sync.WaitGroup
	wg.Add(2)

	go t.write(write, &wg)
	result := t.read(read, &wg)

	go func() {
		for r := range result {
			fmt.Println(r)
		}
	}()

	for key, val := range bigMap {
		write <- map[any]any{key: val}
		read <- key
	}

	close(write)
	close(read)

	wg.Wait()
	slog.Info("everything is done")
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomKeyVal() map[any]any {
	const numkeys = 100_000
	randomMap := make(map[any]any, numkeys)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numkeys; i++ {
		keyLength := rand.Intn(20) + 1
		key := make([]byte, keyLength)
		for j := range key {
			key[j] = letterBytes[rand.Intn(len(letterBytes))]
		}

		valueLength := rand.Intn(20) + 1
		value := make([]byte, valueLength)
		for j := range value {
			value[j] = letterBytes[rand.Intn(len(letterBytes))]
		}

		randomMap[string(key)] = string(value)
	}

	slog.Info("created random map")
	return randomMap
}
