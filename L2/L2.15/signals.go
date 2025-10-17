package main

import (
	"fmt"
	"os"
	"syscall"
)

// Обрабатывает Ctrl+C
func (s *Shell) HandleSignals(sigChan chan os.Signal) {
	for sig := range sigChan {
		if sig == os.Interrupt {
			if s.active != nil {
				fmt.Println("\ninterrupting process...")
				if err := s.active.Process.Signal(syscall.SIGINT); err != nil {
					fmt.Println("failed to send SIGINT:", err)
				}
			} else {
				fmt.Println("\n(press Ctrl+D to exit)")
			}
		}
	}
}
