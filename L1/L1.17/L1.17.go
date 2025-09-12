package main

import "fmt"

func binarySearch(n int, slice []int, left int, right int) int {
	if left > right {
		return -1
	}
	middle := (left + right) / 2

	if slice[middle] == n {
		return middle
	} else if slice[middle] < n {
		return binarySearch(n, slice, middle+1, right)
	} else {
		return binarySearch(n, slice, left, middle-1)
	}
}

func main() {
	sortedSlice := []int{3, 5, 7, 12, 14, 22, 32, 56, 78, 99, 115, 121}
	var num int
	fmt.Print("Введите искомый элемент: ")
	fmt.Scan(&num)
	fmt.Println(binarySearch(num, sortedSlice, 0, len(sortedSlice)-1))
}
