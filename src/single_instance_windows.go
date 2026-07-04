//go:build windows

package main

import (
	"errors"
	"fmt"
	"io"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW         = user32.NewProc("MessageBoxW")
	errElevatedProcess      = errors.New("管理者権限では起動できません。通常のユーザー権限で起動してください。")
	errElevationCheckFailed = errors.New("管理者権限の確認に失敗しました")
)

const (
	mbOK             = 0x00000000
	mbYesNoCancel    = 0x00000003
	mbIconWarning    = 0x00000030
	mbIconError      = 0x00000010
	mbDefaultButton2 = 0x00000100
	idCancel         = 2
	idYes            = 6
	idNo             = 7
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

func chooseExistingSingleInstanceAction(state singleInstanceState, _ io.Writer) singleInstanceChoice {
	message := "ClipForVRChat はすでに起動しています。\n\n" +
		"はい: 既存のClipForVRChatを閉じて、このClipForVRChatを起動します。\n" +
		"いいえ: 既存のClipForVRChatを表示して、この起動を終了します。\n" +
		"キャンセル: 何もせず終了します。"
	if state.ExecutablePath != "" {
		message += "\n\n既存プロセス:\n" + state.ExecutablePath
	}
	ret := messageBox("ClipForVRChat", message, mbYesNoCancel|mbIconWarning|mbDefaultButton2)
	switch ret {
	case idYes:
		return singleInstanceChoiceReplace
	case idNo:
		return singleInstanceChoiceActivate
	default:
		return singleInstanceChoiceCancel
	}
}

func showNativeMessage(title string, message string) {
	_ = messageBox(title, message, mbOK|mbIconError)
}

func messageBox(title string, message string, flags uintptr) uintptr {
	titlePtr, _ := windows.UTF16PtrFromString(title)
	messagePtr, _ := windows.UTF16PtrFromString(message)
	ret, _, _ := procMessageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		flags,
	)
	return ret
}
