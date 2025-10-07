package main

import "fmt"

// NewCup — это новый интерфейс, который ожидает клиент.
// Он определяет метод SystemPoints, работающий с новой системой начисления очков.
type NewCup interface {
	SystemPoints(msg string)
}

// OldCup — это старая структура, которая использует старую систему начисления очков.
// Она несовместима с новым интерфейсом, т.к. у нее другой метод — OldSystemPoints.
type OldCup struct {
	points int
}

// OldSystemPoints — метод старой системы подсчета очков.
// Он возвращает количество очков в зависимости от результата игры.
func (o *OldCup) OldSystemPoints(msg string) int {
	if msg == "draw" {
		o.points++ // за ничью 1 очко
	} else if msg == "win" {
		o.points += 3 // за победу 3 очка
	} else if msg == "lose" {
		o.points += 0 // за поражение 0 очков
	} else {
		fmt.Println("Unknown result") // результат неизвестен
		return 0
	}
	return o.points
}

// AdapterCup — это адаптер, который "переводит" вызовы нового интерфейса (NewCup)
// в формат, понятный старой реализации (OldCup).
type AdapterCup struct {
	Old *OldCup // ссылка на экземпляр старой структуры
}

// SystemPoints — реализация метода нового интерфейса.
// Здесь адаптер вызывает соответствующий метод старого объекта.
func (a *AdapterCup) SystemPoints(msg string) {
	fmt.Println("Adapter in work") // сообщение о работе адаптера
	fmt.Printf("Points %v", a.Old.OldSystemPoints(msg))
}

func main() {
	// Создаем объект типа NewCup, но фактически он — адаптер над старой системой.
	var cup NewCup
	cup = &AdapterCup{
		Old: &OldCup{
			points: 5, // начальное количество очков
		},
	}

	// Пользователь вводит результат матча.
	fmt.Println("Введите результат (win/draw/lose):")
	var result string
	fmt.Scan(&result)

	// Клиентский код работает с новым интерфейсом,
	// не зная, что внутри используется старая система через адаптер.
	cup.SystemPoints(result)
}
