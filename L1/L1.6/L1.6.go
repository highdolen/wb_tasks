package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Выход из горутины 1 по условию")
	go func() {
		for i := 0; ; i++ {
			if i > 5 {
				fmt.Println("Завершение горутины 1 по условию")
				return
			}
			fmt.Printf("Число %v обработано\n", i)
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(7 * time.Second)

	fmt.Println("Выход из горутины 2 через канал уведомлений")
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Завершение горутины 2 через канал уведомлений")
				return
			default:
				fmt.Println("Горутина 2 в работе")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	time.Sleep(3 * time.Second)
	close(done)
	time.Sleep(5 * time.Second)

	fmt.Println("Выход из горутины 3 через контекст")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Завершение горутины 3 через контекст")
				cancel()
				return
			default:
				fmt.Println("Горутина 3 в процессе")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	time.Sleep(5 * time.Second)

	fmt.Println("Выход из горутины 4 через runtime.Goexit()")
	go func() {
		fmt.Println("Горутина 4 в процессе")
		defer fmt.Println("Завершение горутины 4 через runtime.Goexit()")
		runtime.Goexit()
		fmt.Println("Эта строчка не выполнится")
	}()
	time.Sleep(5 * time.Second)

	fmt.Println("Выход из горутины 5 через timeout")

	timeout := time.After(5 * time.Second)
	go func() {
		for {
			select {
			case <-timeout:
				fmt.Println("Горутина 5 прекращает работу через timeout")
				return
			default:
				fmt.Println("Горутина 5 в процессе")
				time.Sleep(time.Second)
			}
		}
	}()
	time.Sleep(7 * time.Second)

}
