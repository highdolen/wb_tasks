package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	//создаем переменные, где будет храниться адрес флагов
	extraStringsAfter  = flag.Int("A", 0, "Дополнительный вывод N строк после найденной строки")
	extraStringsBefore = flag.Int("B", 0, "Дополнительный вывод N строк до найденной строки")
	extraStringsAround = flag.Int("C", 0, "Дополнительный вывод N строк до и после найденной строки")
	count              = flag.Bool("c", false, "Вывод количества строк, совпадающих с шаблоном")
	ignored            = flag.Bool("i", false, "Игнорировать регистр")
	invers             = flag.Bool("v", false, "Выводить строки, не содержащие шаблон")
	fixed              = flag.Bool("F", false, "Воспринимать шаблон, как фиксированную строку")
	num                = flag.Bool("n", false, "Выводить номер строки перед каждой найденной строкой")
)

// normalizePattern возвращает значение в соотвествии с указанными флагами
func normalizePattern(pattern string) string {
	if *ignored {
		pattern = strings.ToLower(pattern)
	}
	return pattern
}

func matches(line string, pattern string) bool {
	// если нужно игнорировать регистр приводим строку к нижнему
	if *ignored {
		line = strings.ToLower(line)
	}
	//изначально ставим флаг false(нашли нужную строку или нет)
	found := false

	//если пришел флаг -F, то просто проверяем содержится ли наш шаблон в строке(точное совпадение)
	if *fixed {
		// точное совпадение подстроки
		found = strings.Contains(line, pattern) //found = true
		//если данный флаг не приходил, тогда ищем строки по стандартному алгоритму
	} else {
		// используем регулярное выражение
		flags := ""
		//если пришел флаг для игнорирования регистра, тогда добавляем в регулярное выражение флаги для игнора регистра
		if *ignored {
			flags = "(?i)" // флаг для игнорирования регистра
		}
		//вот эту строчку честно говоря не особо понимаю, ну я так понимаю, что мы можем работать с каким то протатипом, который нам поступил + флаги
		re, err := regexp.Compile(flags + pattern)
		//если произошла ошибка, то печатаем ее
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка в регулярном выражении:", err)
			//завершаем программу с кодом "1"
			os.Exit(1)
		}
		//говорим что нам строка удовлетворяет
		found = re.MatchString(line)
	}

	// инверсия результата, если пришел флаг -v
	if *invers {
		return !found
	}
	return found
}

// formatPrint выводит результат
func formatPrint(res []string) {
	//если пришел флаг -c, тогда только выводим количество выведенных строк
	if *count {
		fmt.Println(len(res))
		return
	}
	//если пришел флаг -n, то выводим номер строки перед выводом
	if *num {
		for i, v := range res {
			fmt.Printf("%v: %s\n", i+1, v)
		}
		return
	}

	//если не приходило флагов по форматированию вывода, тогда просто выводим нужные строки
	for _, v := range res {
		fmt.Println(v)
	}
}

// Filtred выбирает нужные строки
func Filtred(lines []string, pattern string) {
	//создаем результирующий слайс
	result := []string{}
	//создаем слайс с индексами нужных строк
	matchIndexes := []int{}

	// сначала находим индексы совпадений
	for i, line := range lines {
		//если нам подходит строка, тогда мы добавляем ее индекс в matchIndexes
		if matches(line, pattern) {
			matchIndexes = append(matchIndexes, i)
		}
	}

	//если длина слайса 0, значит не нашли подходящих строк
	if len(matchIndexes) == 0 {
		formatPrint([]string{})
		return
	}

	// если есть контекст (-A, -B, -C)
	after := *extraStringsAfter
	before := *extraStringsBefore
	//если получили контекст -С, значит значения before и after будут равны количеству, которое пришло нам от флага
	if *extraStringsAround > 0 {
		after = *extraStringsAround
		before = *extraStringsAround
	}

	seen := make(map[int]bool) // чтобы не дублировать строки

	//проходимся по индексам
	for _, idx := range matchIndexes {
		//определяем с какого индекса начинаем(если поступил флаг before, значит начинаем с (i - before))
		start := idx - before
		//если начальный индекс получился меньше 0, значит просто зануляем индекс и идем с начала
		if start < 0 {
			start = 0
		}
		//конечным индексом будет индекс + указанное количество строк после выбранной строки
		end := idx + after
		// если последний индекс больше или равен длине слайсу строк, то для того, чтобы не выйти за пределы цикла, последняя строка будет последним элементом
		if end >= len(lines) {
			end = len(lines) - 1
		}

		//проходимся по слайсу строк
		for i := start; i <= end; i++ {
			//если уже встречали этот индекс, тогда не выводим его
			if !seen[i] {
				//добавляем строку в результирующий слайс
				result = append(result, lines[i])
				//через мапу сообщаем, что элемент уже был добавлен в результирующий слайс
				seen[i] = true
			}
		}
	}

	formatPrint(result)
}

func main() {
	//парсим флаги
	flag.Parse()

	//инициализируем буфер для чтения из stdin
	scanner := bufio.NewScanner(os.Stdin)
	//создаем слайс строк
	lines := []string{}
	for scanner.Scan() {
		//добавляем слайс строк из stdin
		lines = append(lines, scanner.Text())
	}

	// шаблон передаётся последним аргументом
	args := flag.Args()
	//если длина переданного шаблона равна 0, значит шаблон не был указан
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: не указан шаблон для поиска")
		//выводим ошибку
		os.Exit(1)
	}
	//инициализируем шаблон и приводим его к нужному виду
	pattern := normalizePattern(args[0])

	//запускаем функцию, которая будет фильтровать строки в соотвествии с шаблоном
	fmt.Println("Вывод:")
	Filtred(lines, pattern)
}
