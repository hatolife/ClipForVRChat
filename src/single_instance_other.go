//go:build !windows

package main

import (
	"fmt"
	"io"
	"os"
)

func rejectElevatedProcess() error {
	return nil
}

func chooseExistingSingleInstanceAction(state singleInstanceState, stderr io.Writer) singleInstanceChoice {
	if stderr != nil {
		if state.ExecutablePath != "" {
			fmt.Fprintf(stderr, "ClipForVRChat はすでに起動しています: %s\n", state.ExecutablePath)
		} else {
			fmt.Fprintln(stderr, "ClipForVRChat はすでに起動しています。")
		}
	}
	return singleInstanceChoiceActivate
}

func showNativeMessage(title string, message string) {
	if title == "" {
		fmt.Fprintln(os.Stderr, message)
		return
	}
	fmt.Fprintln(os.Stderr, title+": "+message)
}
