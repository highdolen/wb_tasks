package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Определение глобальных флагов
var (
	keyColumn   = flag.Int("k", 0, "Сортировать по столбцу N (1-based)")
	numeric     = flag.Bool("n", false, "Сортировать по числовому значению")
	reverse     = flag.Bool("r", false, "Сортировать в обратном порядке")
	unique      = flag.Bool("u", false, "Выводить только уникальные строки")
	monthSort   = flag.Bool("M", false, "Сортировать по названию месяца (Jan, Feb...)")
	ignoreSpace = flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	checkSorted = flag.Bool("c", false, "Проверить, отсортированы ли данные")
	humanSort   = flag.Bool("h", false, "Сортировать с учётом суффиксов (K, M, G)")
)

// Мапа для месяцев
var months = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

// humanReadableRegexp для парсинга размеров
var humanReadableRegexp = regexp.MustCompile(`(?i)^([0-9.]+)([KMGTP]?)$`)

// parseHumanValue преобразует строку с суффиксом K/M/G в число
func parseHumanValue(s string) float64 {
	matches := humanReadableRegexp.FindStringSubmatch(s)
	if matches == nil {
		return 0
	}
	value, _ := strconv.ParseFloat(matches[1], 64)
	switch strings.ToUpper(matches[2]) {
	case "K":
		value *= 1024
	case "M":
		value *= 1024 * 1024
	case "G":
		value *= 1024 * 1024 * 1024
	case "T":
		value *= 1024 * 1024 * 1024 * 1024
	case "P":
		value *= 1024 * 1024 * 1024 * 1024 * 1024
	}
	return value
}

// getColumn возвращает ключ для сортировки по колонке
func getColumn(line string) string {
	cols := strings.Split(line, "\t")
	if *keyColumn <= 0 || *keyColumn > len(cols) {
		return line
	}
	return cols[*keyColumn-1]
}

// compare строки с учётом флагов
func compare(a, b string) bool {
	aKey := getColumn(a)
	bKey := getColumn(b)

	if *ignoreSpace {
		aKey = strings.TrimRight(aKey, " ")
		bKey = strings.TrimRight(bKey, " ")
	}

	if *numeric {
		af, _ := strconv.ParseFloat(aKey, 64)
		bf, _ := strconv.ParseFloat(bKey, 64)
		if *reverse {
			return af > bf
		}
		return af < bf
	}

	if *humanSort {
		af := parseHumanValue(aKey)
		bf := parseHumanValue(bKey)
		if *reverse {
			return af > bf
		}
		return af < bf
	}

	if *monthSort {
		af := months[aKey]
		bf := months[bKey]
		if *reverse {
			return af > bf
		}
		return af < bf
	}

	if *reverse {
		return aKey > bKey
	}
	return aKey < bKey
}

// sortLines сортирует срез строк
func sortLines(lines []string) {
	sort.Slice(lines, func(i, j int) bool {
		return compare(lines[i], lines[j])
	})
	if *unique {
		uniq := lines[:0]
		for i, line := range lines {
			if i == 0 || line != lines[i-1] {
				uniq = append(uniq, line)
			}
		}
		lines = uniq
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}

// checkSortedFunc проверяет отсортированность
func checkSortedFunc(lines []string) {
	for i := 1; i < len(lines); i++ {
		if !compare(lines[i-1], lines[i]) {
			fmt.Println("Файл не отсортирован")
			os.Exit(1)
		}
	}
	fmt.Println("Файл отсортирован")
}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}

	if *checkSorted {
		checkSortedFunc(lines)
		return
	}

	sortLines(lines)
}
