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

// пределение глобальных флагов
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

// мапа для месяцев
var months = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

// humanReadableRegexp для парсинга размеров
var humanReadableRegexp = regexp.MustCompile(`(?i)^([0-9.]+)([KMGT]?)$`)

// parseHumanValue преобразует строку с суффиксом K/M/G в число
func parseHumanValue(s string) float64 {
	matches := humanReadableRegexp.FindStringSubmatch(s)
	//если не нашли значений для парсинга, тогда возвращаем 0
	if matches == nil {
		return 0
	}
	//парсим значение(это число, потому что относится к первой группе)
	value, _ := strconv.ParseFloat(matches[1], 64)
	//далее переводим все буквы(второй группы) в верхний регистр
	switch strings.ToUpper(matches[2]) {
	//если K(килобайты), то умножаем число на 1024
	case "K":
		value *= 1024
	//если М(мегабайты), то умножаем число на 1024^2
	case "M":
		value *= 1024 * 1024
	//если G(гигабайты), то умножаем число на 1024^3
	case "G":
		value *= 1024 * 1024 * 1024
	//если T(терабайты), то умножаем число на 1024^4
	case "T":
		value *= 1024 * 1024 * 1024 * 1024
	}
	//возвращаем нужное значение
	return value
}

// getColumn возвращает ключ для сортировки по колонке, функци получает строку
func getColumn(line string) string {
	//инициализируем переменную, которая делит строку по знаку табуляции(по пробелу)
	cols := strings.Split(line, "\t")
	//если у нас один столбец или введено значение больше чем столбцов есть на самом деле, тогда просто возвращаем строку
	if *keyColumn <= 0 || *keyColumn > len(cols) {
		return line
	}
	//возвращаем столбец по указанному номеру
	return cols[*keyColumn-1]
}

// ompare строки с учётом флагов(проверяем 2 рядом стоящие строки)
func compare(a, b string) bool {
	//получаем столбец для проверки первой строки
	aKey := getColumn(a)
	//получаем столбец для проверки второй строки
	bKey := getColumn(b)

	//если поступил флаг по игнорированию хвостовых пробелов(-b), тогда игнорируем их
	if *ignoreSpace {
		//удаляем лишние пробелы справа
		aKey = strings.TrimRight(aKey, " ")
		bKey = strings.TrimRight(bKey, " ")
	}
	//если поступил флаг по сортировке по числовому значению(-n), значит сортируем по числовому значению
	if *numeric {
		//парсим string в float64
		af, _ := strconv.ParseFloat(aKey, 64)
		bf, _ := strconv.ParseFloat(bKey, 64)
		//если поступил флаг по сортировке в обратном направлении(-r), значит сортируем числа в обратном направлении
		if *reverse {
			//если первый элемент больше второго, тогда выведем true(т.к. реверс)
			return af > bf
		}
		//если второй элемент больше первого, значит все отсортированно верно
		return af < bf
	}

	//если поступил флаг сортировки с учетом суффиксов
	if *humanSort {
		//парсим переменные в соотвествии с условием
		af := parseHumanValue(aKey)
		bf := parseHumanValue(bKey)
		//если стоит флаг реверса, то реверсим результат
		if *reverse {
			return af > bf
		}
		return af < bf
	}

	//если поступил флаг сортировки по месяцам
	if *monthSort {
		//присваивам переменной значение из мапы в соотвествии с месяцем
		af := months[aKey]
		bf := months[bKey]
		//если стоит флаг реверса, то переворачиваем результат
		if *reverse {
			return af > bf
		}
		return af < bf
	}

	//если стоит флаг реверса, то переворачиваем результат
	if *reverse {
		return aKey > bKey
	}
	return aKey < bKey
}

// сортирует срез строк, на вхож получает данный срез
func sortLines(lines []string) {
	//вызываем функцию sort.Slice, которая принимает срез строк и анонимную функцию
	//анонимная функция возвращает bool(0 или 1). в зависимости от того, отсортированны ли два рядом стоящих элемента
	sort.Slice(lines, func(i, j int) bool {
		return compare(lines[i], lines[j])
	})
	//если был флаг, который говорит о том, что нужно выводить только уникальные строки, тогда входим в условие
	if *unique {
		//уникальный срез(ссылается на тот же базовый массив)
		uniq := lines[:0]
		//проходимся по lines
		for i, line := range lines {
			//если у нас первая строка или стоящий между собой элементы не равны, тогда добавляем в слайс уникальных чисел
			if i == 0 || line != lines[i-1] {
				//добавляем в слайс уникальных строк
				uniq = append(uniq, line)
			}
		}
		//приравниваем изначальный слайс к массиву уникальных значений
		lines = uniq
	}
	//проходимся по слайсу и печатаем значения
	for _, line := range lines {
		fmt.Println(line)
	}
}

// checkSortedFunc проверяет отсортированность
func checkSortedFunc(lines []string) {
	//проходимся по срезу строк
	for i := 1; i < len(lines); i++ {
		//если элементы стоят в неправильном порядке, значит слайс не отсортирован
		if !compare(lines[i-1], lines[i]) {
			//пишем что слайс не отсортирован
			fmt.Println("Входные данные не отсортированы")
			//завершаем работу программы
			os.Exit(1)
		}
	}
	//если все в порядке пишем, что входные данные отсортированны
	fmt.Println("Входные данные отсортированы")
}

func main() {
	//парсим флаги
	flag.Parse()
	//создаем буфер чтения из Stdin
	scanner := bufio.NewScanner(os.Stdin)
	//создаем срез строк
	lines := []string{}
	//сканируем входные строку и добавляем их к срезу lines
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	//если произошла ошибка при чтении из буфера, тогда завершаем программу и пишем, что произошла ошибка при чтении
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}

	//проверяем отсортированны ли наши слова(если получили флаг -c)
	if *checkSorted {
		//передаем наш срез в функцию проверки на сортировку
		checkSortedFunc(lines)
		//завершаем main-горутину
		return
	}

	fmt.Println("Отсортированные значения")
	//сортируем входные данные
	sortLines(lines)
}
