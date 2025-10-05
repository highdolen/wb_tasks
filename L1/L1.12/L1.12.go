package main

import "fmt"

func doSet(str []string) []string {
	//инициализируем слайс, где будет результирующее множество
	resultSet := []string{}
	//инициализируем мапу, где слова будут храниться, как ключ
	m := make(map[string]struct{})

	//проходимся по нашим словам(слайсу строк)
	for _, v := range str {
		//добавляем в мапу слово, как ключ
		m[v] = struct{}{}
	}

	//проходимся по мапе
	for key, _ := range m {
		//добавляем ключи к нашему результирующему слайсу
		resultSet = append(resultSet, key)
	}
	return resultSet
}

func main() {
	//инициализируем слайс слов
	words := []string{"cat", "cat", "dog", "cat", "tree"}

	//печатаем множество слов
	fmt.Println(doSet(words))
}
