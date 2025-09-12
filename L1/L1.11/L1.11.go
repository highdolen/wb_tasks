package main

import "fmt"

func check(a, b []int) []int {
	m := make(map[int]int)
	resultSet := []int{}
	for _, v := range a {
		if m[v] > 1 {
			continue
		}
		m[v]++
	}
	for _, v := range b {
		if m[v] == 1 {
			m[v]++
		}
	}

	for key, _ := range m {
		if m[key] >= 2 {
			resultSet = append(resultSet, key)
		}
	}

	return resultSet
}

func main() {
	aSet := []int{1, 2, 3}
	bSet := []int{2, 3, 4}

	fmt.Println(check(aSet, bSet))
}
