package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	var countWorkers int
	fmt.Println("Введите количество воркеров:")
	fmt.Scan(&countWorkers)

	for i := 0; i < countWorkers; i++ {
		go worker(ch, i)
	}

	for i := 1; ; i++ {
		ch <- i
		time.Sleep(1 * time.Second)
	}
}

func worker(jobs <-chan int, id int) {
	for nums := range jobs {
		fmt.Printf("Worker %v - %v\n", id, nums)
	}
}
