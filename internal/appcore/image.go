package appcore

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"golang.org/x/image/webp"
)

type EncodedImage struct {
	Image     image.Image
	Format    string
	Extension string
	Mime      string
	Data      []byte
}

func DecodeImage(data []byte, sourceName string) (image.Image, string, error) {
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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return DecodeImage(data, path)
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

func EncodeImage(img image.Image, sourceFormat string, quality int) (EncodedImage, error) {
	var buf bytes.Buffer
	out := EncodedImage{Image: img}
	if sourceFormat == "jpeg" || sourceFormat == "jpg" {
		out.Format = "jpeg"
		out.Extension = ".jpg"
		out.Mime = "image/jpeg"
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
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
	outputDir := cfg.Image.OutputDirectory
	if outputDir == "" {
		if clipboardInput {
			exe, err := os.Executable()
			if err != nil {
				return "", err
			}
			outputDir = filepath.Join(filepath.Dir(exe), "output")
		} else {
			outputDir = filepath.Dir(sourcePath)
		}
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	base := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	if base == "" || base == "." {
		base = "clipboard_" + time.Now().Format("20060102_150405")
	}
	target := filepath.Join(outputDir, base+cfg.Image.Suffix+encoded.Extension)
	if !cfg.Image.Overwrite {
		target = nextAvailablePath(target)
	}
	return target, os.WriteFile(target, encoded.Data, 0644)
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
