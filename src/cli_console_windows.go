//go:build windows

package main

import (
	"io"
	"os"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const attachParentProcess = ^uint32(0)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	procAttachConsole = kernel32.NewProc("AttachConsole")
	procWriteConsoleW = kernel32.NewProc("WriteConsoleW")
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
	writer := windowsConsoleWriter{file: out}
	return writer, writer, func() {
		_ = out.Close()
	}
}

type windowsConsoleWriter struct {
	file *os.File
}

func (w windowsConsoleWriter) Write(data []byte) (int, error) {
	text := string(data)
	runes := []rune(text)
	wide := utf16.Encode(runes)
	for len(wide) > 0 {
		chunk := wide
		if len(chunk) > 32767 {
			chunk = chunk[:32767]
		}
		var written uint32
		ret, _, err := procWriteConsoleW.Call(
			w.file.Fd(),
			uintptr(unsafe.Pointer(&chunk[0])),
			uintptr(len(chunk)),
			uintptr(unsafe.Pointer(&written)),
			0,
		)
		if ret == 0 {
			return 0, err
		}
		wide = wide[len(chunk):]
	}
	return len(data), nil
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
