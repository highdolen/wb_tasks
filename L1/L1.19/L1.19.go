package main

import "fmt"

func reverseWord(str string) string {
	strRune := []rune(str)
	resultRune := []rune{}

	for i := len(strRune) - 1; i >= 0; i-- {
		resultRune = append(resultRune, strRune[i])
	}

	return string(resultRune)
}

func main() {
	fmt.Print("Введите строку: ")
	var word string
	fmt.Scan(&word)

	fmt.Println(reverseWord(word))
}
