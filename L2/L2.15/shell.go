package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Shell struct {
	output io.Writer
	active *exec.Cmd
}

func NewShell(output io.Writer) *Shell {
	return &Shell{output: output}
}

func (s *Shell) ExecuteLine(line string) error {
	// Разделяем на пайплайн
	commands := strings.Split(line, "|")
	for i := range commands {
		commands[i] = strings.TrimSpace(commands[i])
	}

	// Если это одиночная команда — выполнить напрямую
	if len(commands) == 1 {
		return s.ExecuteSingleCommand(commands[0])
	}
	// Иначе — создать пайплайн
	return s.ExecutePipelineSimple(commands)
}

// Выполнение одиночной команды (встроенной или внешней)
func (s *Shell) ExecuteSingleCommand(cmdLine string) error {
	args := strings.Fields(cmdLine)
	if len(args) == 0 {
		return nil
	}
	name := args[0]
	args = args[1:]

	switch name {
	case "cd":
		return s.commandCd(args)
	case "pwd":
		return s.commandPwd(args)
	case "echo":
		return s.commandEcho(args)
	case "kill":
		return s.commandKill(args)
	case "ps":
		return s.commandPs(args)
	case "exit", "quit":
		fmt.Fprintln(s.output, "exit")
		os.Exit(0)
	default:
		// Запускаем внешнюю команду
		cmd := exec.Command(name, args...)
		cmd.Stdout = s.output
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		s.active = cmd
		err := cmd.Run()
		s.active = nil
		return err
	}
	return nil
}

// Встроенные команды как методы Shell для упрощения
func (s *Shell) commandCd(args []string) error {
	var dir string
	if len(args) == 0 {
		dir = os.Getenv("HOME")
	} else {
		dir = args[0]
		if dir == "~" {
			dir = os.Getenv("HOME")
		}
	}
	if dir == "" {
		return fmt.Errorf("cd: path required")
	}
	return os.Chdir(dir)
}

func (s *Shell) commandPwd(args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(s.output, wd)
	return nil
}

func (s *Shell) commandEcho(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(s.output)
		return nil
	}
	text := strings.Join(args, " ")
	fmt.Fprintln(s.output, text)
	return nil
}

func (s *Shell) commandKill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kill: pid required")
	}
	pid, err := parsePid(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid pid")
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// Используем SIGTERM вместо Kill для более чистого завершения
	return proc.Signal(syscall.SIGTERM)
}

func (s *Shell) commandPs(args []string) error {
	// Всегда используем ps aux для Unix
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = s.output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Shell) HandleSignals(sigs <-chan os.Signal) {
	go func() {
		for sig := range sigs {
			if sig == syscall.SIGINT {
				if s.active != nil {
					if process := s.active.Process; process != nil {
						// Используем SIGINT для Unix
						process.Signal(syscall.SIGINT)
					}
				} else {
					fmt.Fprintln(s.output, "\nType 'exit' to quit")
				}
			}
		}
	}()
}

// Вспомогательная функция для парсинга PID
func parsePid(pidStr string) (int, error) {
	var pid int
	_, err := fmt.Sscanf(pidStr, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}
