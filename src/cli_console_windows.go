//go:build windows

package main

import (
	"io"
	"os"
	"strings"
	"syscall"
)

const attachParentProcess = ^uint32(0)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procAttachConsole   = kernel32.NewProc("AttachConsole")
	procGetConsoleCP    = kernel32.NewProc("GetConsoleCP")
	procGetConsoleMode  = kernel32.NewProc("GetConsoleMode")
	procGetStdHandle    = kernel32.NewProc("GetStdHandle")
	procSetConsoleMode  = kernel32.NewProc("SetConsoleMode")
	procSetConsoleCP    = kernel32.NewProc("SetConsoleCP")
	procSetConsoleOutCP = kernel32.NewProc("SetConsoleOutputCP")
)

func cliOutputWriters(args []string, stdout io.Writer, stderr io.Writer) (io.Writer, io.Writer, func()) {
	if !needsConsoleOutput(args) {
		return stdout, stderr, func() {}
	}
	_, _, _ = procAttachConsole.Call(uintptr(attachParentProcess))
	out, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
	if err != nil {
		return stdout, stderr, func() {}
	}
	return out, out, func() {
		_ = out.Close()
	}
}

func needsConsoleOutput(args []string) bool {
	for _, value := range args {
		switch strings.TrimSpace(value) {
		case "--version", "--help", "-h":
			return true
		}
	}
	return false
}
