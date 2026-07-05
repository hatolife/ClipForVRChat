//go:build !windows

package appcore

import "os/exec"

func hideSpoutHelperWindow(cmd *exec.Cmd) {
}
