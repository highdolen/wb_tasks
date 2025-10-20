package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// создаем структуру Shell
// 1 поле - выход
// 2 поле - команда, которая активна(выполняется)
type Shell struct {
	output io.Writer
	active *exec.Cmd
}

// NewShell - конструктор, который заполняет поле output
func NewShell(output io.Writer) *Shell {
	return &Shell{output: output}
}

// ExecuteLine - метод, который выполняет команду
func (s *Shell) ExecuteLine(line string) error {
	line = strings.TrimSpace(line)

	//добавляем поддержку && и ||
	tokens := splitByLogicalOperators(line)
	if len(tokens) > 1 {
		return s.executeLogical(tokens)
	}

	//логика пайплайнов
	commands := strings.Split(line, "|")
	for i := range commands {
		commands[i] = strings.TrimSpace(commands[i])
	}

	if len(commands) == 1 {
		return s.ExecuteSingleCommand(commands[0])
	}
	return s.ExecutePipeline(commands)
}

// Разбивает строку на команды с учетом && и ||
func splitByLogicalOperators(line string) []string {
	// делаем простейший парсинг
	line = strings.ReplaceAll(line, "&&", " && ")
	line = strings.ReplaceAll(line, "||", " || ")
	return strings.Fields(line)
}

// выполняет команды по && и ||
func (s *Shell) executeLogical(tokens []string) error {
	var lastErr error       // последняя возвращённая ошибка (результат последней выполненной команды)
	expectCmd := true       // флаг: ожидаем ли сейчас команду (true значит: перед последним оператором)
	var op string           // последний увиденный логический оператор ("&&" или "||"), применяется к следующей команде
	var cmdBuilder []string // накопитель токенов для текущей команды

	runCmd := func() error {
		if len(cmdBuilder) == 0 {
			return nil
		}
		cmd := strings.Join(cmdBuilder, " ")
		cmdBuilder = nil
		return s.ExecuteLine(cmd) // рекурсивно выполняем
	}

	for _, tok := range tokens {
		switch tok {
		case "&&", "||":
			if expectCmd {
				return fmt.Errorf("syntax error near '%s'", tok)
			}
			// Выполнить накопленную команду
			err := runCmd()
			lastErr = err
			expectCmd = true
			op = tok
		default:
			cmdBuilder = append(cmdBuilder, tok)
			expectCmd = false
		}
	}

	// выполняем последнюю команду
	if len(cmdBuilder) > 0 {
		if op == "&&" && lastErr != nil {
			return nil // предыдущая упала → пропускаем
		}
		if op == "||" && lastErr == nil {
			return nil // предыдущая прошла → пропускаем
		}
		return runCmd()
	}
	return lastErr
}

// Выполнение одиночной команды (встроенной или внешней)
func (s *Shell) ExecuteSingleCommand(cmdLine string) error {
	//делим строку на слова(пробелы используются как раздилители)
	args := strings.Fields(cmdLine)
	//если строка пустая, тогда возвращаем nil
	if len(args) == 0 {
		return nil
	}
	//именем функции(командой) будет являться 1ый элемент строки
	name := args[0]
	//аргументами будет являться то, что идет дальше
	args = args[1:]

	//смотрим на то, какая команда нам пришла и выполняем ее
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
		//закрываем shell
		os.Exit(0)
	//если name ни одна из этих команд
	//тогда запускается внешняя команда
	default:
		// Запускаем внешнюю команду
		cmd := exec.Command(name, args...)
		cmd.Stdout = s.output
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		//активная команда будет являться той, которая выполняется в данный момент
		s.active = cmd
		//выполняем команду
		err := cmd.Run()
		//когда выполнили, поле снова становится пустым
		s.active = nil
		//возвращаем ошибку(либо nil)
		return err
	}
	return nil
}

// Встроенные команды как методы Shell для упрощения
func (s *Shell) commandCd(args []string) error {
	var dir string
	//если после cd ничего не ввели, тогда переходим в папку HOME
	if len(args) == 0 {
		dir = os.Getenv("HOME")
	} else {
		//в случае, если после cd указан путь, тогда приравниваем его к переменной dir
		dir = args[0]

		if dir == "~" {
			dir = os.Getenv("HOME")
		}
	}
	//если вдруг невалидный запрос, тогда возвращаем ошибку
	if dir == "" {
		return fmt.Errorf("cd: path required")
	}
	//возвращаем изменененную директорию
	return os.Chdir(dir)
}

func (s *Shell) commandPwd(args []string) error {
	//получаем текущую директорию
	wd, err := os.Getwd()
	//если произошла ошибка при получении директории, то возвращаем ее
	if err != nil {
		return err
	}
	//печатаем в Stdout текущую директорию
	fmt.Fprintln(s.output, wd)
	return nil
}

func (s *Shell) commandEcho(args []string) error {
	//если не пришло никаких аргументов, значит ничего не печатаем
	if len(args) == 0 {
		fmt.Fprintln(s.output)
		return nil
	}
	//объединяем текст по пробелу
	text := strings.Join(args, " ")
	//печатаем текст в Stdout
	fmt.Fprintln(s.output, text)
	return nil
}

func (s *Shell) commandKill(args []string) error {
	//если аргументов не пришло, тогда пишем, что pid не указан
	if len(args) == 0 {
		return fmt.Errorf("kill: pid required")
	}
	//парсим pid
	pid, err := parsePid(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid pid")
	}

	//находим процесс по pid
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	//используем SIGTERM вместо Kill для более чистого завершения
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
	for sig := range sigs {
		if sig != syscall.SIGINT {
			continue
		}

		if s.active != nil && s.active.Process != nil {
			// Прерываем только активный процесс
			_ = s.active.Process.Signal(syscall.SIGINT)
		} else {
			// Shell сам не завершается — просто новая строка
			fmt.Fprintln(s.output)
		}
	}
}

// Вспомогательная функция для парсинга PID
func parsePid(pidStr string) (int, error) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid pid: %v", err)
	}
	return pid, nil
}
