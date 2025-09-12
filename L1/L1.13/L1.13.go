package main

import "fmt"

func changeNums(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b
	return a, b
}

func main() {
	var a, b int
	fmt.Print("Введите число а: ")
	fmt.Scan(&a)
	fmt.Print("Введите число b: ")
	fmt.Scan(&b)
	fmt.Println(changeNums(a, b))
}
