package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Выполняет пайплайн команд (упрощенная версия)
func (s *Shell) ExecutePipeline(cmdLines []string) error {
	n := len(cmdLines)
	if n == 0 {
		return fmt.Errorf("empty pipeline")
	}

	// Запускаем команды последовательно, соединяя пайпами
	var prevOutput io.Reader = os.Stdin

	for i, line := range cmdLines {
		args := strings.Fields(line)
		if len(args) == 0 {
			return fmt.Errorf("empty command in pipeline")
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = prevOutput
		cmd.Stderr = os.Stderr

		// Если это последняя команда - выводим в stdout shell
		if i == n-1 {
			cmd.Stdout = s.output

			// Запускаем и ждем завершения
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			// Для промежуточных команд создаем пайп
			reader, writer := io.Pipe()
			cmd.Stdout = writer

			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start %s: %v", args[0], err)
			}

			// Закрываем writer после завершения команды
			go func() {
				cmd.Wait()
				writer.Close()
			}()

			prevOutput = reader
		}
	}

	return nil
}
