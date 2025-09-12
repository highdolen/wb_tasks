package main

import (
	"fmt"
	"sync"
)

func main() {
	m := make(map[int]int)
	ch := make(chan int)
	numWorkers := 3
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ch, i, &mu, wg, m)
	}

	for i := 7; i < 10; i++ {
		ch <- i
	}
	wg.Wait()
	fmt.Println(m)
}

func worker(jobs <-chan int, id int, mu *sync.Mutex, wg *sync.WaitGroup, m map[int]int) {
	for job := range jobs {
		mu.Lock()
		fmt.Printf("Worker %v - %v\n", id, job)
		m[id] = job
		mu.Unlock()
		wg.Done()
	}
}
