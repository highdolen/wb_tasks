package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := []int{2, 4, 6, 8, 10}
	wg := &sync.WaitGroup{}
	for _, value := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newValue := value * value
			fmt.Println(newValue)
		}()
	}
	wg.Wait()
}
