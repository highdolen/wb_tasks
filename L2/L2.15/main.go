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
	shell := NewShell(os.Stdout)

	// Настраиваем обработку Ctrl+C для Unix
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT) // Используем Unix сигнал
	go shell.HandleSignals(sigs)

	fmt.Println("minishell (use 'exit' or Ctrl+D to quit)")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("minishell> ")
		if !scanner.Scan() {
			// Ctrl+D или EOF - нормальный выход для Unix
			if scanner.Err() == nil {
				fmt.Println("\nexit")
			} else {
				fmt.Fprintf(os.Stderr, "read error: %v\n", scanner.Err())
			}
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := shell.ExecuteLine(line); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}
