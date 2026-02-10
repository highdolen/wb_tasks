package main

import (
	"fmt"
	"time"
)

func main() {
	//создаем небуферезированный канал
	ch := make(chan int)
	//создаем счетчик воркеров
	var countWorkers int
	fmt.Println("Введите количество воркеров:")
	fmt.Scan(&countWorkers)

	//создаем горутины воркеров
	for i := 0; i < countWorkers; i++ {
		go worker(ch, i)
	}

	//создаем тикер, который будет сигнализировать каждую секунду
	ticker := time.NewTicker(time.Second)
	//останавливаем тикер при завершении main
	defer ticker.Stop()

	i := 1
	//бесконечный цикл
	for range ticker.C {
		//отправка чисел в канал
		ch <- i
		//инкремент числа
		i++
	}
}

// горутина воркеров, принимаем на вход канал и id воркера
func worker(jobs <-chan int, id int) {
	//проходимся по каналу, каждый воркер принимает определенное число и обрабатывает
	for nums := range jobs {
		fmt.Printf("Worker %v - %v\n", id, nums)
	}
}
