package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	fields    = flag.String("f", "", "Номера колонок, которые нужно вывести")
	delimiter = flag.String("d", "", "Использовать другой разделитель(символ)")
	separated = flag.Bool("s", false, "Вывод строк только с разделителями")
)

func toInt(s string) int {
	n := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			n = n*10 + int(r-'0')
		}
	}
	return n
}

func selectCollumns(nums string) []int {
	numColumns := []int{}
	parts := strings.Split(nums, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) != 2 {
				continue
			}
			start := toInt(bounds[0])
			end := toInt(bounds[1])
			if start > 0 && end > 0 {
				if start > end {
					start, end = end, start
				}
				for n := start; n <= end; n++ {
					numColumns = append(numColumns, n)
				}
			}
		} else {
			numColumns = append(numColumns, toInt(part))
		}
	}

	return numColumns
}

func Collumns(lines []string) {
	for _, line := range lines {
		str := []string{}
		if *delimiter == "" {
			str = strings.Split(line, "\t")
		} else {
			str = strings.Split(line, *delimiter)
		}

		if *separated && len(str) == 1 {
			continue
		}

		newStr := []string{}
		if *fields != "" {
			numCollumns := selectCollumns(*fields)
			for _, v := range numCollumns {
				if v-1 < len(str) && v-1 >= 0 {
					newStr = append(newStr, str[v-1])
				}
			}
			fmt.Println(strings.Join(newStr, "\t"))
		}
		if *fields == "" {
			fmt.Println(strings.Join(str, "\t"))
		}
	}

}

func main() {
	flag.Parse()
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println("Вывод:")
	Collumns(lines)
}
