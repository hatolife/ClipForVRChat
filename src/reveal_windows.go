//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	coinitApartmentThreaded = 0x2
	rpcEChangedMode         = 0x80010106
	sOK                     = 0x0
	sFalse                  = 0x1
)

var (
	shell32                        = syscall.NewLazyDLL("shell32.dll")
	ole32                          = syscall.NewLazyDLL("ole32.dll")
	procILCreateFromPathW          = shell32.NewProc("ILCreateFromPathW")
	procILFree                     = shell32.NewProc("ILFree")
	procSHOpenFolderAndSelectItems = shell32.NewProc("SHOpenFolderAndSelectItems")
	procCoInitializeEx             = ole32.NewProc("CoInitializeEx")
	procCoUninitialize             = ole32.NewProc("CoUninitialize")
)

func revealFileInExplorer(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	hr, _, _ := procCoInitializeEx.Call(0, coinitApartmentThreaded)
	if hr == sOK || hr == sFalse {
		defer procCoUninitialize.Call()
	} else if uint32(hr) != rpcEChangedMode {
		return fmt.Errorf("Windows Shellを初期化できません: HRESULT 0x%08x", uint32(hr))
	}

	pidl, _, _ := procILCreateFromPathW.Call(uintptr(unsafe.Pointer(pathPtr)))
	if pidl == 0 {
		return fmt.Errorf("Explorerで表示するファイルを解決できません: %s", path)
	}
	defer procILFree.Call(pidl)

	hr, _, _ = procSHOpenFolderAndSelectItems.Call(pidl, 0, 0, 0)
	if hr != 0 {
		return fmt.Errorf("Explorerでファイルを選択表示できません: HRESULT 0x%08x", uint32(hr))
	}
	return nil
}
