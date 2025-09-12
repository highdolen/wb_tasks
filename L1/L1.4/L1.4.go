package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		fmt.Println("\nОстанавливаем программу.")
		cancel()
	}()

	ch := make(chan int)
	var countWorkers int
	fmt.Println("Введите количество воркеров:")
	fmt.Scan(&countWorkers)

	for i := 0; i < countWorkers; i++ {
		go worker(ctx, ch, i)
	}

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			close(ch)
			return
		case ch <- i:
			time.Sleep(1 * time.Second)
		}
	}
}
func worker(ctx context.Context, jobs <-chan int, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %v завершает работу\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("Worker %v - %v\n", id, job)
		}
	}
}
