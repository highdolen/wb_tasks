package main

import (
	"bufio"
	"fmt"
	"net"
)

// StartServer запускает TCP сервер на указанном порту
func StartServer(port string, fields string, delimiter string, separated bool) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server started on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, fields, delimiter, separated)
	}
}

// handleConnection обрабатывает каждое соединение с клиентом
func handleConnection(conn net.Conn, fields string, delimiter string, separated bool) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("connection close error:", err)
		}
	}()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		result := ProcessLine(line, fields, delimiter, separated)
		if result != "" {
			if _, err := fmt.Fprintln(conn, result); err != nil {
				fmt.Println("send error:", err)
				return
			}
		}
	}
}
