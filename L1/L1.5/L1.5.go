package main

import (
	"fmt"
	"time"
)

func main() {
	//создаем канал
	ch := make(chan int)

	//задаем время работы программы
	var n int
	fmt.Print("Введите время работы программы(сек):")
	fmt.Scan(&n)
	//задаем таймаут(сработает через n секунд)
	timeout := time.After(time.Duration(n) * time.Second)

	//запускаем горутину для обработки числа
	go analyseNum(ch)

	//запускаем бесконечный цикл
	for i := 0; ; i++ {
		select {
		//если сработал наш таймаут, закрываем канал и пишем что программа
		//остановлена по истечении n секунд
		case <-timeout:
			close(ch)
			fmt.Printf("Программа остановлена по истечениии %v сек", n)
			return
		//пока таймаут не сработал, отправляем числа в канал
		case ch <- i:
			time.Sleep(1 * time.Second)
		}
	}
}

// функция принимает на вход канал чтения
func analyseNum(nums <-chan int) {
	//читаем числа из канала и печатаем их
	for num := range nums {
		fmt.Printf("Число %v обработано\n", num)
	}
}
