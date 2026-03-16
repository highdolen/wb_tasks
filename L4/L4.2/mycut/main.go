package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// Флаги командной строки для настройки работы утилиты
var (
	fields    = flag.String("f", "", "fields")     // номера колонок, которые нужно вывести
	delimiter = flag.String("d", "", "delimiter")  // пользовательский разделитель колонок
	separated = flag.Bool("s", false, "separated") // выводить только строки с разделителем

	serverMode = flag.Bool("server", false, "run as server") // запуск программы в режиме сервера
	port       = flag.String("port", "9001", "server port")  // порт для запуска сервера

	nodes = flag.String("nodes", "", "node list") // список серверов для распределённой обработки
)

// main - точка входа в программу
func main() {
	flag.Parse()

	// если указан флаг server — запускаем сервер
	if *serverMode {
		StartServer(*port, *fields, *delimiter, *separated)
		return
	}

	// если серверы не указаны — выполняем локальную обработку
	if *nodes == "" {
		runLocal()
		return
	}

	// иначе запускаем распределённую обработку
	runDistributed()
}

// runLocal - выполняет обработку входного потока локально без использования серверов
func runLocal() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		line := scanner.Text()

		// обработка строки через функцию cut
		result := ProcessLine(line, *fields, *delimiter, *separated)

		// вывод результата
		if result != "" {
			fmt.Println(result)
		}
	}
}

// runDistributed - распределяет обработку строк между несколькими серверами
func runDistributed() {

	// преобразуем строку серверов в список
	nodeList := strings.Split(*nodes, ",")

	conns := []net.Conn{}
	writers := []*bufio.Writer{}

	// подключаемся к каждому серверу
	for _, node := range nodeList {

		conn, err := net.Dial("tcp", node)

		// если соединение не удалось — пропускаем сервер
		if err != nil {
			continue
		}

		conns = append(conns, conn)
		writers = append(writers, bufio.NewWriter(conn))

		// запускаем горутину для чтения результатов от сервера
		go readResults(conn)
	}

	// если ни к одному серверу не удалось подключиться
	if len(conns) == 0 {
		fmt.Println("No servers available")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	i := 0

	// читаем строки из входного потока
	for scanner.Scan() {

		line := scanner.Text()

		// выбираем сервер по принципу round-robin
		server := i % len(writers)

		// обрабатываем строку
		result := ProcessLine(line, *fields, *delimiter, *separated)

		// отправляем результат на сервер
		fmt.Fprintln(writers[server], result)

		i++
	}

	// очищаем буферы и отправляем данные
	for _, w := range writers {
		w.Flush()
	}
}

// readResults - читает строки, которые возвращает сервер, и выводит их в stdout
func readResults(conn net.Conn) {

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
