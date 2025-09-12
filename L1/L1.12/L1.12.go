package main

import "fmt"

func doSet(str []string) []string {
	resultSet := []string{}
	m := make(map[string]struct{})

	for _, v := range str {
		m[v] = struct{}{}
	}

	for key, _ := range m {
		resultSet = append(resultSet, key)
	}
	return resultSet
}

func main() {
	words := []string{"cat", "cat", "dog", "cat", "tree"}

	fmt.Println(doSet(words))
}
