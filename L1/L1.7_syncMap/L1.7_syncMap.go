package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	ch := make(chan int)
	numWorkers := 3
	wg := &sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workerMap(ch, i, wg, &m)
	}

	for i := 7; i < 10; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	fmt.Println("Итоговое содержимое мапы:")
	m.Range(func(key any, value any) bool {
		fmt.Printf("Worker %v - %v\n", key, value)
		return true
	})
}

func workerMap(jobs <-chan int, id int, wg *sync.WaitGroup, m *sync.Map) {
	defer wg.Done()
	for job := range jobs {
		m.Store(id, job)
		fmt.Printf("Worker %v - %v\n", id, job)
	}
}
