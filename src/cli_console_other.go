//go:build !windows

package main

import "io"

func cliOutputWriters(_ []string, stdout io.Writer, stderr io.Writer) (io.Writer, io.Writer, func()) {
	return stdout, stderr, func() {}
}
