package main

//импортируем нужные библиотеки
import (
	"fmt"
	"math"
)

// создаем структуру Point
type Point struct {
	x float64
	y float64
}

// функция-конструктор, создающая новую точку
func NewPoint(x, y float64) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

// метод, который вычисляет расстояние между точками
func (p1 *Point) Distance(p2 Point) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	//возвращаем квадрат результата сложения
	return math.Sqrt(dx*dx + dy*dy)
}
func main() {

	//инициализируем точки с помощью функции-конструктора
	point1 := NewPoint(1, 2)
	point2 := NewPoint(4, 6)

	//вызываем метод, который вычисляет расстояние между точками и печатаем результат
	fmt.Println(point1.Distance(*point2))
}
