package main

import "fmt"

func changeNums(a, b int) (int, int) {
	a = a + b   // 13 + 14 = 27
	b = a - b   //27-13 = 14
	a = a - b   //27 - 14 = 13
	return a, b //возвращаем числа, которые мы поменяли местами(точнее разные стали значения переменных)
}

func main() {
	//инициализируем 2 переменные
	var a, b int
	fmt.Print("Введите число а: ")
	fmt.Scan(&a)
	fmt.Print("Введите число b: ")
	fmt.Scan(&b)
	//печатаем числа, которые поменяны местами
	fmt.Println(changeNums(a, b))
}
