package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	count int
}

func main() {
	wg := &sync.WaitGroup{}
	workerPool := 3

	nums := &Counter{}

	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go worker(nums, wg)
	}
	wg.Wait()
	fmt.Println(nums.count)
}

func worker(c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}
