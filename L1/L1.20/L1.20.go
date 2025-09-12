package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func reverseStr(s string) string {
	s = strings.TrimSpace(s)
	wordsSlice := strings.Split(s, " ")
	for i, j := 0, len(wordsSlice)-1; i < j; i, j = i+1, j-1 {
		wordsSlice[i], wordsSlice[j] = wordsSlice[j], wordsSlice[i]
	}
	return strings.Join(wordsSlice, " ")
}

func main() {
	fmt.Println("Введите строку")
	r := bufio.NewReader(os.Stdin)
	str, err := r.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(reverseStr(str))
}
