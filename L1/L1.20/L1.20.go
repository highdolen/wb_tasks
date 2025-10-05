package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func reverseStr(s string) string {
	//удаляем лишние пробелы в начале и конце строки
	s = strings.TrimSpace(s)
	//разделяем строку по пробелам на отдельные слова
	wordsSlice := strings.Split(s, " ")
	//создаме цикл, где меняем местами слова
	for i, j := 0, len(wordsSlice)-1; i < j; i, j = i+1, j-1 {
		wordsSlice[i], wordsSlice[j] = wordsSlice[j], wordsSlice[i]
	}
	//возвращаем объединенную строку(с помощью join, через пробел)
	return strings.Join(wordsSlice, " ")
}

func main() {
	//вводим строку
	fmt.Println("Введите строку")
	//создаем "ридер" для чтения строки из стандартного ввода
	r := bufio.NewReader(os.Stdin)
	//читаем строку и проверяем на ошибку
	str, err := r.ReadString('\n')
	//если ошибка, то завершаем программу с ее кодом (1)
	if err != nil {
		os.Exit(1)
	}
	//печатаем перевернутую строку(обратный порядок слов)
	fmt.Println(reverseStr(str))
}
