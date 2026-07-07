//go:build !windows

package main

import "errors"

const startupShortcutName = "ClipForVRChat.lnk"

func startupShortcutStatus() StartupShortcutStatus {
	return StartupShortcutStatus{
		Supported:    false,
		ShortcutName: startupShortcutName,
		Error:        "PC起動時自動起動はWindows版でのみ設定できます",
	}
}

func setStartupShortcut(enabled bool) (StartupShortcutStatus, error) {
	status := startupShortcutStatus()
	if enabled {
		return status, errors.New(status.Error)
	}
	return status, nil
}
