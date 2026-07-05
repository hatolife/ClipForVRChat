//go:build windows

package appcore

import (
	"os/exec"
	"syscall"
)

func hideSpoutHelperWindow(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
