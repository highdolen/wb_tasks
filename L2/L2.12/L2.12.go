package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	extraStringsAfter  = flag.Int("A", 0, "Дополнительный вывод N строк после найденной строки")
	extraStringsBefore = flag.Int("B", 0, "Дополнительный вывод N строк до найденной строки")
	extraStringsAround = flag.Int("C", 0, "Дополнительный вывод N строк до и после найденной строки")
	count              = flag.Bool("c", false, "Вывод количества строк, совпадающих с шаблоном")
	ignored            = flag.Bool("i", false, "Игнорировать регистр")
	invers             = flag.Bool("v", false, "Выводить строки, не содержащие шаблон")
	fixed              = flag.Bool("F", false, "Воспринимать шаблон, как фиксированную строку")
	num                = flag.Bool("n", false, "Выводить номер строки перед каждой найденной строкой")
)

func normalizePattern(pattern string) string {
	if *ignored {
		strings.ToLower(pattern)
	}
	return pattern
}

func matches(line string, pattern string) bool {
	if *ignored {
		strings.ToLower(line)
	}
	if *fixed {
		found := strings.Contains(pattern, line)
	} else {
		flags := ""
	}

}

func formatPrint(res []string) {
	if *count {
		fmt.Println(len(res))
		return
	}

	if *num {
		for i, v := range res {
			fmt.Println(i, ".", v)
		}
		return
	}

	fmt.Println(res)
}

func Filtred(lines []string, pattern string) {
	result := []string{}

	if *extraStringsAfter != 0 {
		for i, line := range lines {
			if matches(line, pattern) {
				for j := 0; j < *extraStringsAfter; j++ {
					result = append(result, lines[i+j])
				}
			}
		}
	}

	if *extraStringsBefore != 0 {
		for i, line := range lines {
			if matches(line, pattern) {
				for j := *extraStringsBefore; j > 0; j-- {
					result = append(result, lines[i-j])
				}
			}
		}
	}

	if *extraStringsAround != 0 {
		for i, line := range lines {
			if matches(line, pattern) {
				for j := *extraStringsAround; j > 0; j-- {
					result = append(result, lines[i-j])
				}
				for j := 1; j < *extraStringsAround; j++ {
					result = append(result, lines[i+j])
				}

			}
		}
	}

	formatPrint(result)
}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	pattern := "Error"
	pattern = normalizePattern(pattern)

	Filtred(lines, pattern)
}
