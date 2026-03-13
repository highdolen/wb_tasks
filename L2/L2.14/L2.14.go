package main

import (
	"fmt"
	"sync"
	"time"
)

// or принимает на вход каналы и возвращает результирующий канал(через него подается сигнал о заверщении работы)
func or(channels ...<-chan interface{}) <-chan interface{} {
	//Если в функцию or не поступает сигналов, тогда возвращаем nil
	if len(channels) == 0 {
		return nil
	}
	//создаем результирующий канал
	resultChannel := make(chan interface{})
	//создаем once для того, чтобы закрыть канал только один раз
	var once sync.Once
	//читаем данные с каналов, пока они не закрыты
	for _, ch := range channels {
		//запускаем горутину
		go func() {
			select {
			//читаем значение из канала
			case <-ch:
				//закрываем канал один раз
				once.Do(func() {
					close(resultChannel)
				})
			}
		}()
	}

	return resultChannel
}

func main() {
	//инициализируем сигнал(который будет сигнализировать о том, что прошло время и нужно завершать работу) через функцию
	sig := func(after time.Duration) <-chan interface{} {
		//создаем канал для сигнала о завершении работы
		c := make(chan interface{})
		//горутина, которая закрывает канал по истечении времени
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		//возвращаем канал
		return c
	}

	//инициализируем переменную, которая хранит время в данный момент
	start := time.Now()
	//передаем если поступил сигнал из любого из каналов о том, что нужно завершить программу, тогда завершаем ее
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	//фиксируем время работы программы с ее начала до завершения через канал
	fmt.Printf("done after %v", time.Since(start))
}
