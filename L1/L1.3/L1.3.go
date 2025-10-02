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

	//бесконечный цикл чисел
	for i := 1; ; i++ {
		//отправка числа в канал
		ch <- i
		//горутина засыпает на секунду, чтобы не было слишком частой отправки в канал
		time.Sleep(1 * time.Second)
	}
}

// горутина воркеров, принимаем на вход канал и id воркера
func worker(jobs <-chan int, id int) {
	//проходимся по каналу, каждый воркер принимает определенное число и обрабатывает
	for nums := range jobs {
		fmt.Printf("Worker %v - %v\n", id, nums)
	}
}
