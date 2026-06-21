package appcore

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"golang.org/x/image/webp"
)

const (
	DefaultMaxImageInputMB = 32
	MaxImagePixels         = 64_000_000
)

type EncodedImage struct {
	Image     image.Image
	Format    string
	Extension string
	Mime      string
	Data      []byte
}

func DecodeImage(data []byte, sourceName string) (image.Image, string, error) {
	return DecodeImageWithLimit(data, sourceName, DefaultMaxImageInputMB)
}

func DecodeImageWithLimit(data []byte, sourceName string, maxInputMB int) (image.Image, string, error) {
	if err := ValidateImageBytes(data, maxInputMB); err != nil {
		return nil, "", err
	}
	format, err := ValidateImageDimensions(data, sourceName)
	if err != nil {
		return nil, "", err
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err == nil {
		return img, strings.ToLower(format), nil
	}
	if strings.EqualFold(filepath.Ext(sourceName), ".webp") {
		webpImg, webpErr := webp.Decode(bytes.NewReader(data))
		if webpErr == nil {
			return webpImg, "webp", nil
		}
	}
	return nil, "", err
}

func DecodeImageFile(path string) (image.Image, string, error) {
	return DecodeImageFileWithLimit(path, DefaultMaxImageInputMB)
}

func DecodeImageFileWithLimit(path string, maxInputMB int) (image.Image, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}
	maxBytes := maxImageInputBytes(maxInputMB)
	if info.Size() > maxBytes {
		return nil, "", fmt.Errorf("画像ファイルが大きすぎます。%dMB以下の画像を指定してください。", maxBytes/1024/1024)
	}
	data, err := os.ReadFile(path) // #nosec G304 -- path comes from explicit user file drop or configured photo folder.
	if err != nil {
		return nil, "", err
	}
	return DecodeImageWithLimit(data, path, maxInputMB)
}

func ValidateImageBytes(data []byte, maxInputMB int) error {
	if len(data) == 0 {
		return errors.New("画像データが空です。")
	}
	maxBytes := maxImageInputBytes(maxInputMB)
	if int64(len(data)) > maxBytes {
		return fmt.Errorf("画像データが大きすぎます。%dMB以下の画像を指定してください。", maxBytes/1024/1024)
	}
	return nil
}

func maxImageInputBytes(maxInputMB int) int64 {
	if maxInputMB <= 0 {
		maxInputMB = DefaultMaxImageInputMB
	}
	return int64(maxInputMB) * 1024 * 1024
}

func ValidateImageDimensions(data []byte, sourceName string) (string, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil && strings.EqualFold(filepath.Ext(sourceName), ".webp") {
		var webpErr error
		cfg, webpErr = webp.DecodeConfig(bytes.NewReader(data))
		if webpErr == nil {
			format = "webp"
			err = nil
		}
	}
	if err != nil {
		return "", err
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return "", errors.New("画像サイズを確認できませんでした。")
	}
	if cfg.Height > MaxImagePixels/cfg.Width {
		return "", fmt.Errorf("画像のピクセル数が大きすぎます。%dメガピクセル以下の画像を指定してください。", MaxImagePixels/1_000_000)
	}
	return strings.ToLower(format), nil
}

func ResizeToFit(img image.Image, maxWidth, maxHeight int) image.Image {
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()
	if width <= maxWidth && height <= maxHeight {
		return img
	}
	widthRatio := float64(maxWidth) / float64(width)
	heightRatio := float64(maxHeight) / float64(height)
	ratio := widthRatio
	if heightRatio < ratio {
		ratio = heightRatio
	}
	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}
	return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
}

func FlattenTransparentImage(img image.Image, bg color.Color) image.Image {
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Over)
	return out
}

func TrimClipboardArtifactBorder(img image.Image) image.Image {
	b := img.Bounds()
	if b.Dx() <= 4 || b.Dy() <= 4 {
		return img
	}
	return cropImage(img, image.Rect(b.Min.X+1, b.Min.Y+1, b.Max.X-1, b.Max.Y-1))
}

func cropImage(img image.Image, rect image.Rectangle) image.Image {
	out := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(out, out.Bounds(), img, rect.Min, draw.Src)
	return out
}

func EncodeImage(img image.Image, outputFormat string, quality int) (EncodedImage, error) {
	var buf bytes.Buffer
	out := EncodedImage{Image: img}
	if outputFormat == "jpg" || outputFormat == "jpeg" {
		out.Format = "jpeg"
		out.Extension = ".jpg"
		out.Mime = "image/jpeg"
		if err := jpeg.Encode(&buf, FlattenTransparentImage(img, color.White), &jpeg.Options{Quality: quality}); err != nil {
			return out, err
		}
	} else {
		out.Format = "png"
		out.Extension = ".png"
		out.Mime = "image/png"
		if err := png.Encode(&buf, img); err != nil {
			return out, err
		}
	}
	out.Data = buf.Bytes()
	return out, nil
}

func SaveEncodedImage(encoded EncodedImage, sourcePath string, cfg Config, clipboardInput bool) (string, error) {
	outputDir := strings.Trim(strings.TrimSpace(cfg.Image.OutputDirectory), `"`)
	if outputDir == "" {
		outputDir = "./output"
	}
	if !filepath.IsAbs(outputDir) {
		exe, err := os.Executable()
		if err != nil {
			return "", err
		}
		outputDir = filepath.Join(filepath.Dir(exe), outputDir)
	}
	if err := os.MkdirAll(outputDir, privateDirMode); err != nil {
		return "", err
	}

	base := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	if clipboardInput {
		base = "clipboard_" + time.Now().Format("20060102_150405")
	} else if base == "" || base == "." {
		base = "clipboard_" + time.Now().Format("20060102_150405")
	}
	target := filepath.Join(outputDir, base+cfg.Image.Suffix+encoded.Extension)
	if !cfg.Image.Overwrite {
		target = nextAvailablePath(target)
	}
	return target, WritePrivateFile(target, encoded.Data)
}

func nextAvailablePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

func ThumbnailDataURL(img image.Image) (string, error) {
	thumb := imaging.Fit(img, 160, 160, imaging.Lanczos)
	var buf bytes.Buffer
	if err := png.Encode(&buf, thumb); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
