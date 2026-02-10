package main

import (
	"fmt"
	"sync"
)

type SafeMap struct {
	mu sync.Mutex
	m  map[int]int
}

// NewSafeMap - конструктор мапы с мьютексом
func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[int]int),
	}
}

// Set - добавление в мапу ключа и значения
func (s *SafeMap) Set(key int, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//добавление ключа и значения
	s.m[key] = value

}

// Print - печать мапы
func (s *SafeMap) Print() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println(s.m)
}

func main() {
	safeMap := NewSafeMap()
	//создаем канал
	ch := make(chan int)
	//создаем переменную, отвечающую за количество работающих воркеров
	numWorkers := 3
	//создаем вэйтгруппу, для конкурентной работы
	wg := &sync.WaitGroup{}
	//запускаем воркеров
	for i := 0; i < numWorkers; i++ {
		//увеличиваем счетчик для каждой задачи
		wg.Add(1)
		go worker(ch, i, safeMap, wg)
	}
	//запускаем цикл для записи чисел в мапу
	for i := 7; i < 10; i++ {
		//передаем число в канал
		ch <- i
	}
	close(ch)

	//ждем пока работа воркеров завершится
	wg.Wait()
	//печатаем мапу
	safeMap.Print()
}

// worker - принимает канал чтения, свой айди, мапу(с мьютексом) и вейтгруппу
func worker(jobs <-chan int, id int, s *SafeMap, wg *sync.WaitGroup) {
	defer wg.Done()
	//читаем данные из канала
	for job := range jobs {
		s.Set(id, job)
		fmt.Printf("Worker %v - %v\n", id, job)
	}
}
