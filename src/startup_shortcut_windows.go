//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const startupShortcutName = "ClipForVRChat.lnk"

func startupShortcutPath() (string, error) {
	appData := strings.TrimSpace(os.Getenv("APPDATA"))
	if appData == "" {
		return "", errors.New("APPDATA が見つかりません")
	}
	return filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup", startupShortcutName), nil
}

func currentStartupShortcutExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(exe)
}

func setStartupShortcut(enabled bool) (StartupShortcutStatus, error) {
	shortcutPath, err := startupShortcutPath()
	if err != nil {
		return StartupShortcutStatus{Supported: true, ShortcutPath: shortcutPath, Error: err.Error()}, err
	}
	if err := os.Remove(shortcutPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		status := startupShortcutStatus()
		status.Error = err.Error()
		return status, err
	}
	if !enabled {
		return startupShortcutStatus(), nil
	}
	exePath, err := currentStartupShortcutExecutablePath()
	if err != nil {
		status := startupShortcutStatus()
		status.Error = err.Error()
		return status, err
	}
	if err := os.MkdirAll(filepath.Dir(shortcutPath), 0700); err != nil {
		status := startupShortcutStatus()
		status.Error = err.Error()
		return status, err
	}
	workingDir := filepath.Dir(exePath)
	script := strings.Join([]string{
		"$shortcutPath = $args[0]",
		"$targetPath = $args[1]",
		"$workingDirectory = $args[2]",
		"$shell = New-Object -ComObject WScript.Shell",
		"$shortcut = $shell.CreateShortcut($shortcutPath)",
		"$shortcut.TargetPath = $targetPath",
		"$shortcut.WorkingDirectory = $workingDirectory",
		"$shortcut.Save()",
	}, "; ")
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script, shortcutPath, exePath, workingDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		status := startupShortcutStatus()
		status.Error = strings.TrimSpace(string(output))
		if status.Error == "" {
			status.Error = err.Error()
		}
		return status, fmt.Errorf("startup shortcut create failed: %w", err)
	}
	return startupShortcutStatus(), nil
}

func startupShortcutStatus() StartupShortcutStatus {
	shortcutPath, err := startupShortcutPath()
	status := StartupShortcutStatus{
		Supported:     true,
		ShortcutPath:  shortcutPath,
		ShortcutName:  startupShortcutName,
		CurrentTarget: "",
	}
	if err != nil {
		status.Error = err.Error()
		return status
	}
	if _, statErr := os.Stat(shortcutPath); statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			status.Error = statErr.Error()
		}
		return status
	}
	status.Enabled = true
	status.CurrentTarget = readStartupShortcutTarget(shortcutPath)
	if exePath, exeErr := currentStartupShortcutExecutablePath(); exeErr == nil {
		status.CurrentExe = exePath
		status.TargetMatchesCurrentExe = sameWindowsPath(status.CurrentTarget, exePath)
	}
	return status
}

func readStartupShortcutTarget(shortcutPath string) string {
	script := "$shell = New-Object -ComObject WScript.Shell; $shortcut = $shell.CreateShortcut($args[0]); [Console]::Write($shortcut.TargetPath)"
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script, shortcutPath)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func sameWindowsPath(a string, b string) bool {
	if a == "" || b == "" {
		return false
	}
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA == nil {
		a = absA
	}
	if errB == nil {
		b = absB
	}
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
