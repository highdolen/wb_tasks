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
	if n == 0 {
		return fmt.Errorf("empty pipeline")
	}

	cmds := make([]*exec.Cmd, n)

	// Создаем команды
	for i, line := range cmdLines {
		args := strings.Fields(line)
		if len(args) == 0 {
			return fmt.Errorf("empty command in pipeline")
		}

		// Для Unix все команды запускаем напрямую
		cmds[i] = exec.Command(args[0], args[1:]...)
	}

	// Соединяем команды пайпами
	for i := 0; i < n-1; i++ {
		r, w := io.Pipe()
		cmds[i].Stdout = w
		cmds[i+1].Stdin = r
		cmds[i].Stderr = os.Stderr
	}

	// Настраиваем первую и последнюю команду
	cmds[0].Stdin = os.Stdin
	cmds[n-1].Stdout = s.output
	cmds[n-1].Stderr = os.Stderr

	// Запускаем все команды
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %v", err)
		}
	}

	// Закрываем stdout предыдущих команд после запуска следующих
	for i := 0; i < n-1; i++ {
		if writer, ok := cmds[i].Stdout.(*io.PipeWriter); ok {
			go func(idx int, w *io.PipeWriter) {
				cmds[idx].Wait()
				w.Close()
			}(i, writer)
		}
	}

	// Ждем завершения всех команд
	var lastError error
	for i, cmd := range cmds {
		err := cmd.Wait()
		if i == n-1 {
			// Сохраняем ошибку только от последней команды
			lastError = err
		}
	}

	return lastError
}

// Упрощенная версия для встроенных команд в пайплайнах
func (s *Shell) ExecutePipelineSimple(cmdLines []string) error {
	n := len(cmdLines)
	if n == 0 {
		return fmt.Errorf("empty pipeline")
	}

	cmds := make([]*exec.Cmd, n)

	// Создаем команды
	for i, line := range cmdLines {
		args := strings.Fields(line)
		if len(args) == 0 {
			return fmt.Errorf("empty command in pipeline")
		}

		cmdName := args[0]

		// Для встроенных команд используем системные аналоги
		switch cmdName {
		case "echo":
			cmds[i] = exec.Command("echo", args[1:]...)
		case "pwd":
			cmds[i] = exec.Command("pwd")
		case "ps":
			cmds[i] = exec.Command("ps", "aux")
		default:
			cmds[i] = exec.Command(args[0], args[1:]...)
		}
	}

	// Соединяем команды пайпами
	for i := 0; i < n-1; i++ {
		r, w := io.Pipe()
		cmds[i].Stdout = w
		cmds[i+1].Stdin = r
		cmds[i].Stderr = os.Stderr
	}

	// Настраиваем первую и последнюю команду
	cmds[0].Stdin = os.Stdin
	cmds[n-1].Stdout = s.output
	cmds[n-1].Stderr = os.Stderr

	// Запускаем все команды
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %v", err)
		}
	}

	// Ждем завершения всех команд
	for i, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			if i < n-1 {
				// Для промежуточных команд ошибки могут быть из-за закрытия пайпов
				continue
			}
			return err
		}
	}

	return nil
}
