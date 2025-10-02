package main

import (
	"fmt"
	"sync"
)

func main() {
	//создание слайса чисел
	nums := []int{2, 4, 6, 8, 10}
	//инициализируем вэйтгруппу для синхронизации
	wg := &sync.WaitGroup{}
	//проходимся по нашим числам
	for _, value := range nums {
		wg.Add(1) //добавляем горутину в вэйтгруппу, чтобы не было гонки данных
		go func() {
			defer wg.Done()
			//берем квадрат от исходного числа
			newValue := value * value
			fmt.Println(newValue)
		}()
	}
	//ждем выполнение всех горутин, чтобы main-горутина не закончилась раньше, чем выполнятся остальные горутины
	wg.Wait()
}
