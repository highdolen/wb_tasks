package main

import (
	"fmt"
	"sync"
)

func main() {
	//создаем синкМапу
	var m sync.Map
	//создаем канал
	ch := make(chan int)
	//число воркеров
	numWorkers := 3
	//создаем вэйтгруппу
	wg := &sync.WaitGroup{}
	//запускаем воркеры
	for i := 0; i < numWorkers; i++ {
		//добавляем задачу
		//увеличиваем счетчик вэйтгрупп
		wg.Add(1)
		go workerMap(ch, i, wg, &m)
	}

	//запускаем цикл для передачи чисел в канал
	for i := 7; i < 10; i++ {
		//передаем числа в канал
		ch <- i
	}

	//закрываем канал(сообщаем о том, что задач больше не будет)
	close(ch)

	//ждем, пока все горутины выполнят свою задачу
	wg.Wait()
	fmt.Println("Итоговое содержимое мапы:")
	//обходим синкМапу и выводим все пары ключ-значение
	m.Range(func(key any, value any) bool {
		fmt.Printf("Worker %v - %v\n", key, value)
		return true
	})
}

// в горутину передаем канал для чтения, айди воркера, вэйтгруппу и синкМапу
func workerMap(jobs <-chan int, id int, wg *sync.WaitGroup, m *sync.Map) {
	//сообщаем о вэйтгруппе о том, что задача выполнена(откладываем о том, что задача выполнена на конец(с помощью defer))
	//уменьшаем счетчик вэйтгрупп
	defer wg.Done()
	//читаем данные из канала
	for job := range jobs {
		//сохраняем данные в синкМапу(Store безопасен для конкурентного доступа)
		m.Store(id, job)
		fmt.Printf("Worker %v - %v\n", id, job)
	}
}
