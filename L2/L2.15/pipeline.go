package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Выполняет пайплайн команд
func (s *Shell) ExecutePipeline(cmdLines []string) error {
	n := len(cmdLines)
	cmds := make([]*exec.Cmd, n)

	for i, line := range cmdLines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			return fmt.Errorf("empty command in pipeline")
		}
		cmds[i] = exec.Command(parts[0], parts[1:]...)
	}

	// Соединяем пайпы
	for i := 0; i < n-1; i++ {
		r, w := io.Pipe()
		cmds[i].Stdout = w
		cmds[i+1].Stdin = r
		cmds[i].Stderr = os.Stderr
	}

	// Последняя команда выводит в stdout
	cmds[n-1].Stdout = s.output
	cmds[n-1].Stderr = os.Stderr

	// Запускаем
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	// Ждём
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}
