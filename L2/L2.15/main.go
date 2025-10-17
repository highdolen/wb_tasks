package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func main() {
	shell := NewShell(os.Stdout)

	// Настраиваем обработку Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go shell.HandleSignals(sigs)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("minishell> ")
		if !scanner.Scan() {
			fmt.Println("\nexit")
			break // Ctrl+D
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
