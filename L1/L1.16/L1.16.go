package main

import "fmt"

func quicksort(nums []int) []int {
	if len(nums) < 2 {
		return nums
	}

	pivot := nums[len(nums)/2]
	left := []int{}
	right := []int{}
	sortedSlice := []int{}
	for i, v := range nums {
		if i == len(nums)/2 {
			continue
		} else if v < pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	sortedSlice = append(quicksort(left), pivot)
	sortedSlice = append(sortedSlice, quicksort(right)...)
	return sortedSlice
}

func main() {
	slice := []int{2, 4, 1, 13, 44, 22, 45, 111, 10, 21, 30}
	fmt.Println(quicksort(slice))
}
