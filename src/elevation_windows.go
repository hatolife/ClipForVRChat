//go:build windows && !ciguismoke

package main

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	errElevatedProcess      = errors.New("管理者権限では起動できません。通常のユーザー権限で起動してください。")
	errElevationCheckFailed = errors.New("管理者権限の確認に失敗しました")
)

func rejectElevatedProcess() error {
	elevated, err := isProcessElevated()
	if err != nil {
		showNativeMessage("ClipForVRChat", errElevationCheckFailed.Error()+": "+err.Error())
		return fmt.Errorf("%w: %v", errElevationCheckFailed, err)
	}
	if !elevated {
		return nil
	}
	showNativeMessage("ClipForVRChat", errElevatedProcess.Error())
	return errElevatedProcess
}

func isProcessElevated() (bool, error) {
	var token windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token); err != nil {
		return false, err
	}
	defer token.Close()

	var elevated uint32
	var outLen uint32
	err := windows.GetTokenInformation(
		token,
		windows.TokenElevation,
		(*byte)(unsafe.Pointer(&elevated)),
		uint32(unsafe.Sizeof(elevated)),
		&outLen,
	)
	if err != nil {
		return false, err
	}
	return elevated != 0, nil
}
