package main

import "fmt"

func main() {
	slice := []int{1, 2, 3, 4, 5}
	var index int
	fmt.Print("Напиши индекс для удаления")
	fmt.Scan(&index)
	copy(slice[index:], slice[index+1:])
	slice = slice[:len(slice)-1]
	fmt.Println(slice)
}
