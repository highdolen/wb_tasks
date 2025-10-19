package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	//создаем shell с выводом в Stdout
	shell := NewShell(os.Stdout)

	//настраиваем обработку Ctrl+C для Unix
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT) // Используем Unix сигнал
	//запускаем горутину, которая обрабатывает сигнал
	go shell.HandleSignals(sigs)

	//печетаем инфо для пользователя
	fmt.Println("minishell (use 'exit' or Ctrl+D to quit)")

	//создаем сканер для чтения из Stdin
	scanner := bufio.NewScanner(os.Stdin)
	//основной цикл работы shell
	for {
		//выводим приглашение к вводу
		fmt.Print("minishell> ")
		//если не смогли прочитать строчку из Stdin
		if !scanner.Scan() {
			//проверяем причину прекращения чтения
			if scanner.Err() == nil {
				//EOF (Ctrl+D) - нормальный выход
				fmt.Println("\nexit")
			} else {
				//произошла ошибка чтения
				fmt.Fprintf(os.Stderr, "read error: %v\n", scanner.Err())
			}
			break
		}
		//удаляем лишние пробелы из строки
		line := strings.TrimSpace(scanner.Text())
		//если строчка пустая, shell продолжает работу
		if line == "" {
			continue
		}

		//выполняем команду и проверяем ее на ошибку
		if err := shell.ExecuteLine(line); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}
