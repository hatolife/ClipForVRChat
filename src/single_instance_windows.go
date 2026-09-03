//go:build windows

package main

import (
	"io"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
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
