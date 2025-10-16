package main

import (
	"fmt"
	"sync"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	resultChannel := make(chan interface{})

	go func() {
		defer close(resultChannel)
		var wg sync.WaitGroup
		for _, ch := range channels {
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case _, ok := <-ch:
					if !ok {
						close(resultChannel)
					}
				}
			}()
		}
		wg.Wait()
	}()

	return resultChannel
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(2*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
