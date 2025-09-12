package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	var n int
	fmt.Print("Введите время работы программы(сек):")
	fmt.Scan(&n)
	timeout := time.After(time.Duration(n) * time.Second)

	go analyseNum(ch)

	for i := 0; ; i++ {
		select {
		case <-timeout:
			close(ch)
			fmt.Printf("Программа остановлена по истечениии %v сек", n)
			return
		case ch <- i:
			time.Sleep(1 * time.Second)
		}
	}
}

func analyseNum(nums <-chan int) {
	for num := range nums {
		fmt.Printf("Число %v обработано\n", num)
	}
}
