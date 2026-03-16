package main

import (
	"strings"
	"unicode"
)

// toInt конвертирует строку в число
func toInt(s string) int {
	n := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			n = n*10 + int(r-'0')
		}
	}
	return n
}

// selectCollumns разбивает строку с номерами колонок и возвращает список индексов
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

// ProcessLine обрабатывает строку line по указанным колонкам и разделителю
func ProcessLine(line string, fields string, delimiter string, separated bool) string {

	var str []string
	if delimiter == "" {
		str = strings.Split(line, "\t")
	} else {
		str = strings.Split(line, delimiter)
	}

	if separated && len(str) == 1 {
		return ""
	}

	if fields == "" {
		return strings.Join(str, "\t")
	}

	numColumns := selectCollumns(fields)
	newStr := []string{}

	for _, v := range numColumns {
		if v-1 < len(str) && v-1 >= 0 {
			newStr = append(newStr, str[v-1])
		}
	}

	return strings.Join(newStr, "\t")
}
