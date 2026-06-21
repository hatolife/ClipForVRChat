//go:build !windows

package appcore

func readNativeClipboardPNG() ([]byte, error) {
	return nil, nil
}
