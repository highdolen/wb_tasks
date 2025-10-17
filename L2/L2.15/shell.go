package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	return s.ExecutePipeline(commands)
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
		return (&commandCd{}).Run(s.output, args)
	case "pwd":
		return (&commandPwd{}).Run(s.output, args)
	case "echo":
		return (&commandEcho{}).Run(s.output, args)
	case "kill":
		return (&commandKill{}).Run(s.output, args)
	case "ps":
		return (&commandPs{}).Run(s.output, args)
	case "exit", "quit":
		fmt.Fprintln(s.output, "exit")
		os.Exit(0)
	default:
		// Попробуем запустить внешнюю команду
		cmd := exec.Command(name, args...)
		cmd.Stdout = s.output
		cmd.Stderr = os.Stderr
		s.active = cmd
		err := cmd.Run()
		s.active = nil
		return err
	}
	return nil
}
