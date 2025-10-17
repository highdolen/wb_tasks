package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Command интерфейс для встроенных команд
type Command interface {
	Run(output io.Writer, args []string) error
}

type commandCd struct{}

func (c *commandCd) Run(output io.Writer, args []string) error {
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
	if err := os.Chdir(dir); err != nil {
		return err
	}
	wd, _ := os.Getwd()
	fmt.Fprintln(output, "cd:", wd)
	return nil
}

type commandPwd struct{}

func (c *commandPwd) Run(output io.Writer, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(output, wd)
	return nil
}

type commandEcho struct{}

func (c *commandEcho) Run(output io.Writer, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(output)
		return nil
	}
	noNewline := false
	if args[0] == "-n" {
		noNewline = true
		args = args[1:]
	}
	text := strings.Join(args, " ")
	if noNewline {
		fmt.Fprint(output, text)
	} else {
		fmt.Fprintln(output, text)
	}
	return nil
}

type commandKill struct{}

func (c *commandKill) Run(output io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kill: pid required")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid pid")
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Kill(); err != nil {
		return err
	}
	fmt.Fprintf(output, "process %d killed\n", pid)
	return nil
}

type commandPs struct{}

func (c *commandPs) Run(output io.Writer, args []string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", "tasklist")
	} else {
		cmd = exec.Command("ps", "aux")
	}
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
