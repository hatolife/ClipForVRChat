package appcore

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"strings"
)

type Processor struct {
	Config Config
}

func (p Processor) ProcessPaths(paths []string) ([]Result, error) {
	return p.ProcessPathsWithProgress(paths, nil)
}

func (p Processor) ProcessPathsWithProgress(paths []string, progress func(ProgressEvent)) ([]Result, error) {
	if len(paths) == 0 {
		if progress != nil {
			progress(ProgressEvent{Index: 0, Total: 1, Stage: "processing", Result: Result{Name: "clipboard.png", SourcePath: "clipboard.png"}})
		}
		results, err := p.ProcessClipboard()
		if progress != nil && err == nil && len(results) == 1 {
			progress(ProgressEvent{Index: 0, Total: 1, Stage: "done", Result: results[0]})
		}
		return results, err
	}
	if !p.Config.Output.SaveLocal && !p.Config.Output.UploadDiscord && len(paths) > 1 {
		return nil, fmt.Errorf("複数画像はクリップボードへ直接保持できません。ローカル保存またはDiscord投稿をONにするか、1枚ずつ処理してください。")
	}
	results := make([]Result, 0, len(paths))
	for i, path := range paths {
		if progress != nil {
			progress(ProgressEvent{Index: i, Total: len(paths), Stage: "processing", Result: Result{SourcePath: path, Name: filepath.Base(path)}})
		}
		result := p.processFile(path)
		results = append(results, result)
		if progress != nil {
			progress(ProgressEvent{Index: i, Total: len(paths), Stage: "done", Result: result})
		}
	}
	return results, nil
}

func (p Processor) ProcessClipboard() ([]Result, error) {
	data, err := ReadClipboardImage()
	if err != nil {
		return nil, err
	}
	img, format, err := DecodeImageWithLimit(data, "clipboard.png", p.Config.Image.MaxInputMB)
	if err != nil {
		return nil, err
	}
	img = TrimClipboardArtifactBorder(img)
	img = FlattenTransparentImage(img, color.White)
	result := p.processImage(img, format, "clipboard.png", true)
	return []Result{result}, nil
}

func (p Processor) processFile(path string) Result {
	img, format, err := DecodeImageFileWithLimit(path, p.Config.Image.MaxInputMB)
	if err != nil {
		return Result{SourcePath: path, Name: filepath.Base(path), Error: fmt.Sprintf("画像を読み込めませんでした: %v", err)}
	}
	return p.processImage(img, format, path, false)
}

func (p Processor) processImage(img image.Image, format string, sourcePath string, clipboardInput bool) Result {
	result := Result{SourcePath: sourcePath, Name: filepath.Base(sourcePath)}
	inputBounds := img.Bounds()
	p.logf("process start source=%q format=%s clipboard=%t input=%dx%d qr_enabled=%t", sourcePath, format, clipboardInput, inputBounds.Dx(), inputBounds.Dy(), p.Config.Output.DetectQRCodeURLs)
	if p.Config.Output.DetectQRCodeURLs {
		qrResult := DetectQRCodeURLsWithDiagnostics(img)
		result.QRURLs = qrResult.URLs
		p.logf("qr result source=%q urls=%d details=%s", sourcePath, len(qrResult.URLs), strings.Join(qrResult.Diagnostics, " | "))
	}
	resized := ResizeToFit(img, p.Config.Image.MaxWidth, p.Config.Image.MaxHeight)
	resizedBounds := resized.Bounds()
	p.logf("resize result source=%q output=%dx%d", sourcePath, resizedBounds.Dx(), resizedBounds.Dy())
	encoded, err := EncodeImage(resized, p.Config.Image.OutputFormat, p.Config.Image.JPEGQuality)
	if err != nil {
		result.Error = fmt.Sprintf("画像を書き出せませんでした: %v", err)
		p.logf("encode error source=%q error=%q", sourcePath, result.Error)
		return result
	}

	if thumb, err := ThumbnailDataURL(resized); err == nil {
		result.Thumbnail = thumb
	}

	if p.Config.Output.SaveLocal {
		outputPath, err := SaveEncodedImage(encoded, sourcePath, p.Config, clipboardInput)
		if err != nil {
			result.Error = fmt.Sprintf("縮小画像を保存できませんでした: %v", err)
			p.logf("save error source=%q error=%q", sourcePath, result.Error)
			return result
		}
		result.OutputPath = outputPath
		p.logf("save result source=%q output_path=%q bytes=%d", sourcePath, outputPath, len(encoded.Data))
	}

	if p.Config.Output.UploadDiscord {
		name := filepath.Base(sourcePath)
		name = name[:len(name)-len(filepath.Ext(name))] + encoded.Extension
		uploaded, err := UploadDiscordWithContent(p.Config.Discord.WebhookURL, name, encoded, qrDiscordContent(result.QRURLs))
		if err != nil {
			result.Error = err.Error()
			p.logf("discord error source=%q error=%q", sourcePath, result.Error)
			return result
		}
		result.URL = uploaded.URL
		result.DiscordMessageID = uploaded.MessageID
		result.DiscordWebhookID = uploaded.WebhookID
		result.DiscordToken = uploaded.Token
		p.logf("discord result source=%q url=%q message_id=%q", sourcePath, result.URL, result.DiscordMessageID)
	} else if !p.Config.Output.SaveLocal {
		if err := WriteClipboardImage(encoded.Data); err != nil {
			result.Error = fmt.Sprintf("画像をクリップボードへコピーできませんでした: %v", err)
			p.logf("clipboard write error source=%q error=%q", sourcePath, result.Error)
			return result
		}
		p.logf("clipboard write result source=%q bytes=%d", sourcePath, len(encoded.Data))
	}

	return result
}

func (p Processor) logf(format string, args ...any) {
	AppendDiagnosticLog(p.Config.DiagnosticLogPath, format, args...)
}

func CopySingleURLIfNeeded(cfg Config, results []Result) error {
	if !cfg.Output.CopySingleURLToClipboard || len(results) != 1 || results[0].URL == "" || results[0].Error != "" {
		return nil
	}
	return WriteClipboardText(results[0].URL)
}
