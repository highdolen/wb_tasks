package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Выход из горутины 1 по условию")
	//запускаем горутину
	go func() {
		//запускаем бесконечный цикл чисел
		for i := 0; ; i++ {
			//условие, при котором завершаем работу нашей горутины
			if i > 5 {
				fmt.Println("Завершение горутины 1 по условию")
				return
			}
			//пишем что наше число обработано
			fmt.Printf("Число %v обработано\n", i)
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(7 * time.Second)

	fmt.Println("Выход из горутины 2 через канал уведомлений")
	//создаем канал через который будет приходить уведомление о завершении работы
	done := make(chan struct{})
	//запускаем горутину
	go func() {
		//запускаем бесконечный цикл
		for {
			select {
			//если пришло уведомление и закрытии канала, то завершаем работу горутины
			case <-done:
				fmt.Println("Завершение горутины 2 через канал уведомлений")
				return
			//пока уведомление не пришло, горутина в работе
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
	//создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//запускаем горутину
	go func() {
		//запускаем бесконечный цикл
		for {
			select {
			//если пришло уведомление о завершении работы через контекст, завершаем работу горутины
			case <-ctx.Done():
				fmt.Println("Завершение горутины 3 через контекст")
				//завершаем работу
				cancel()
				return
			//пока уведомление не пришло, горутина работает
			default:
				fmt.Println("Горутина 3 в процессе")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	time.Sleep(5 * time.Second)

	fmt.Println("Выход из горутины 4 через runtime.Goexit()")
	//запускаем горутину
	go func() {
		fmt.Println("Горутина 4 в процессе")
		defer fmt.Println("Завершение горутины 4 через runtime.Goexit()")
		//вызываем Goexit, работа горутины сразу же завершается
		runtime.Goexit()
		fmt.Println("Эта строчка не выполнится")
	}()
	time.Sleep(5 * time.Second)

	fmt.Println("Выход из горутины 5 через timeout")

	//создаем таймаут(через 5 сек)
	timeout := time.After(5 * time.Second)
	//запускаем горутину
	go func() {
		//создаем бесконечный цикл
		for {
			select {
			//если пришло уведомление через timeout, выходим из горутины
			case <-timeout:
				fmt.Println("Горутина 5 прекращает работу через timeout")
				return
			//пока уведомление не пришло, горутина в работе
			default:
				fmt.Println("Горутина 5 в процессе")
				time.Sleep(time.Second)
			}
		}
	}()
	time.Sleep(7 * time.Second)

}
