package appcore

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
)

type Processor struct {
	Config Config
}

func (p Processor) ProcessPaths(paths []string) ([]Result, error) {
	if len(paths) == 0 {
		return p.ProcessClipboard()
	}
	if !p.Config.Output.SaveLocal && !p.Config.Output.UploadDiscord && len(paths) > 1 {
		return nil, fmt.Errorf("複数画像はクリップボードへ直接保持できません。ローカル保存またはDiscord投稿をONにするか、1枚ずつ処理してください。")
	}
	results := make([]Result, 0, len(paths))
	for _, path := range paths {
		result := p.processFile(path)
		results = append(results, result)
	}
	return results, nil
}

func (p Processor) ProcessClipboard() ([]Result, error) {
	data, err := ReadClipboardImage()
	if err != nil {
		return nil, err
	}
	img, format, err := DecodeImage(data, "clipboard.png")
	if err != nil {
		return nil, err
	}
	img = TrimClipboardArtifactBorder(img)
	img = FlattenTransparentImage(img, color.White)
	result := p.processImage(img, format, "clipboard.png", true)
	return []Result{result}, nil
}

func (p Processor) processFile(path string) Result {
	img, format, err := DecodeImageFile(path)
	if err != nil {
		return Result{SourcePath: path, Name: filepath.Base(path), Error: fmt.Sprintf("画像を読み込めませんでした: %v", err)}
	}
	return p.processImage(img, format, path, false)
}

func (p Processor) processImage(img image.Image, format string, sourcePath string, clipboardInput bool) Result {
	resized := ResizeToFit(img, p.Config.Image.MaxWidth, p.Config.Image.MaxHeight)
	encoded, err := EncodeImage(resized, format, p.Config.Image.JPEGQuality)
	result := Result{SourcePath: sourcePath, Name: filepath.Base(sourcePath)}
	if err != nil {
		result.Error = fmt.Sprintf("画像を書き出せませんでした: %v", err)
		return result
	}

	if thumb, err := ThumbnailDataURL(resized); err == nil {
		result.Thumbnail = thumb
	}

	if p.Config.Output.SaveLocal {
		outputPath, err := SaveEncodedImage(encoded, sourcePath, p.Config, clipboardInput)
		if err != nil {
			result.Error = fmt.Sprintf("縮小画像を保存できませんでした: %v", err)
			return result
		}
		result.OutputPath = outputPath
	}

	if p.Config.Output.UploadDiscord {
		name := filepath.Base(sourcePath)
		name = name[:len(name)-len(filepath.Ext(name))] + encoded.Extension
		url, err := UploadDiscord(p.Config.Discord.WebhookURL, name, encoded)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.URL = url
	} else if !p.Config.Output.SaveLocal {
		if err := WriteClipboardImage(encoded.Data); err != nil {
			result.Error = fmt.Sprintf("画像をクリップボードへコピーできませんでした: %v", err)
			return result
		}
	}

	return result
}

func CopySingleURLIfNeeded(cfg Config, results []Result) error {
	if !cfg.Output.CopySingleURLToClipboard || len(results) != 1 || results[0].URL == "" || results[0].Error != "" {
		return nil
	}
	return WriteClipboardText(results[0].URL)
}
