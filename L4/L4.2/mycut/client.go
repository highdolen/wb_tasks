package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// SendToNode отправляет строки на сервер node, обрабатывает их через ProcessLine
func SendToNode(node string, lines []string, fields string, delimiter string, separated bool, ch chan string, chResults chan string) {
	conn, err := net.Dial("tcp", node)
	if err != nil {
		ch <- "" // если соединение не удалось, отправляем пустой ответ
		return
	}

	scanner := bufio.NewScanner(conn)

	// отправляем строки серверу
	for _, line := range lines {
		result := ProcessLine(line, fields, delimiter, separated)
		fmt.Fprintln(conn, result)
	}

	// закрываем запись, чтобы сервер понял, что данные закончились
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	// читаем ответы сервера и отправляем их в chResults
	for scanner.Scan() {
		chResults <- scanner.Text()
	}

	conn.Close()
	ch <- "ok" // сигнализируем, что задача выполнена
}

// SplitLines делит строки lines на примерно равные куски по числу серверов nodes
func SplitLines(lines []string, nodes []string) [][]string {
	chunkSize := (len(lines) + len(nodes) - 1) / len(nodes)
	chunks := [][]string{}

	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, lines[i:end])
	}

	return chunks
}

// ParseNodes преобразует строку с адресами серверов через запятую в срез адресов
func ParseNodes(nodes string) []string {
	return strings.Split(nodes, ",")
}
