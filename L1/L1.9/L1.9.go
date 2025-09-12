package main

import "fmt"

func main() {
	arr := [7]int{2, 4, 6, 7, 8, 9, 11}

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for _, v := range arr {
			ch1 <- v
		}
		close(ch1)
	}()

	go func(nums <-chan int) {
		for num := range nums {
			numDouble := num * 2
			ch2 <- numDouble
		}
		close(ch2)
	}(ch1)

	for numDouble := range ch2 {
		fmt.Println(numDouble)
	}
}
