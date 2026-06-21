//go:build windows

package appcore

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	winUser32                     = syscall.NewLazyDLL("user32.dll")
	winOpenClipboard              = winUser32.NewProc("OpenClipboard")
	winCloseClipboard             = winUser32.NewProc("CloseClipboard")
	winGetClipboardData           = winUser32.NewProc("GetClipboardData")
	winIsClipboardFormatAvailable = winUser32.NewProc("IsClipboardFormatAvailable")
	winRegisterClipboardFormatA   = winUser32.NewProc("RegisterClipboardFormatA")

	winKernel32     = syscall.NewLazyDLL("kernel32.dll")
	winGlobalLock   = winKernel32.NewProc("GlobalLock")
	winGlobalUnlock = winKernel32.NewProc("GlobalUnlock")
	winGlobalSize   = winKernel32.NewProc("GlobalSize")
)

func readNativeClipboardPNG() ([]byte, error) {
	for _, formatName := range []string{"PNG", "image/png"} {
		data, err := readRegisteredClipboardFormat(formatName)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}
	return nil, nil
}

func readRegisteredClipboardFormat(formatName string) ([]byte, error) {
	format, err := registerClipboardFormat(formatName)
	if err != nil {
		return nil, err
	}
	available, _, _ := winIsClipboardFormatAvailable.Call(format)
	if available == 0 {
		return nil, nil
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	opened, _, openErr := winOpenClipboard.Call(0)
	if opened == 0 {
		return nil, openErr
	}
	defer winCloseClipboard.Call()

	handle, _, dataErr := winGetClipboardData.Call(format)
	if handle == 0 {
		return nil, dataErr
	}
	ptr, _, lockErr := winGlobalLock.Call(handle)
	if ptr == 0 {
		return nil, lockErr
	}
	defer winGlobalUnlock.Call(handle)

	size, _, _ := winGlobalSize.Call(handle)
	if size == 0 {
		return nil, errors.New("clipboard PNG data is empty")
	}

	data := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(size))
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func registerClipboardFormat(name string) (uintptr, error) {
	raw, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, err
	}
	format, _, callErr := winRegisterClipboardFormatA.Call(uintptr(unsafe.Pointer(raw)))
	if format == 0 {
		return 0, callErr
	}
	return format, nil
}
