package main

import (
	"fmt"
	"strings"
)

func Check(str string) bool {
	m := make(map[string]int)      //создаем мапу, где ключ - буква, значение - количество раз, когда эта буква встретилась в строке
	newStr := strings.ToLower(str) //все буквы переводим в нижний регистр
	for _, v := range newStr {     //итерируемся по моей строке с нижний регистром
		m[string(v)]++ //добавляем наш ключ и увеличиваем значение на одиг
	}
	for _, value := range m { //итерируемся по мапе
		if value > 1 { //если значение больше одного, то возвращаем false
			return false
		}
	}

	return true //если значения были везде 1, значит возвращаем true
}

func main() {
	var word string //создаем строку
	fmt.Print("Введите строку: ")
	fmt.Scan(&word)          //ввели строку
	fmt.Println(Check(word)) //вызываем функцию Check
}
