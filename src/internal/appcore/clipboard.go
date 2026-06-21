package appcore

import (
	"errors"

	"golang.design/x/clipboard"
)

var clipboardReady bool

func InitClipboard() error {
	if clipboardReady {
		return nil
	}
	if err := clipboard.Init(); err != nil {
		return err
	}
	clipboardReady = true
	return nil
}

func ReadClipboardImage() ([]byte, error) {
	if err := InitClipboard(); err != nil {
		return nil, err
	}
	if data, err := readNativeClipboardPNG(); err == nil && len(data) > 0 {
		return data, nil
	}
	data := clipboard.Read(clipboard.FmtImage)
	if len(data) == 0 {
		return nil, errors.New("クリップボードに画像がありません。画像をコピーしてから再実行するか、画像ファイルを exe にドラッグしてください。")
	}
	return data, nil
}

func WriteClipboardText(text string) error {
	if err := InitClipboard(); err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}

func WriteClipboardImage(pngData []byte) error {
	if err := InitClipboard(); err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtImage, pngData)
	return nil
}
