package main

import "fmt"

func reverseWord(str string) string {
	//переводим строку в руны
	strRune := []rune(str)
	//создаем слайс рун(результирующий)
	resultRune := []rune{}

	//проходимся по слову в обратном порядке
	for i := len(strRune) - 1; i >= 0; i-- {
		//добавляем в результирующий слайс руны по порядку
		resultRune = append(resultRune, strRune[i])
	}

	//переводим слайс рун в строку и возвращаем ее
	return string(resultRune)
}

func main() {
	//вводим строку
	fmt.Print("Введите строку: ")
	var word string
	fmt.Scan(&word)

	//печатаем перевернутую строку
	fmt.Println(reverseWord(word))
}
