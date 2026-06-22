package appcore

import (
	"fmt"
	"image"
	"image/color"
	"net/url"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/makiuchi-d/gozxing"
	multiqrcode "github.com/makiuchi-d/gozxing/multi/qrcode"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type qrDetectionResult struct {
	URLs        []string
	Diagnostics []string
}

func DetectQRCodeURLs(img image.Image) []string {
	return DetectQRCodeURLsWithDiagnostics(img).URLs
}

func DetectQRCodeURLsWithDiagnostics(img image.Image) qrDetectionResult {
	seen := map[string]bool{}
	var urls []string
	var diagnostics []string
	for _, variant := range qrImageVariants(img) {
		texts, errText := decodeQRCodeTexts(variant.image)
		if errText != "" {
			diagnostics = append(diagnostics, fmt.Sprintf("%s: %s", variant.name, errText))
		}
		for _, text := range texts {
			raw := strings.TrimSpace(text)
			if !isHTTPURL(raw) {
				diagnostics = append(diagnostics, fmt.Sprintf("%s: ignored non-url payload %q", variant.name, truncateForUserMessage(raw, 80)))
				continue
			}
			if seen[raw] {
				continue
			}
			seen[raw] = true
			urls = append(urls, raw)
			diagnostics = append(diagnostics, fmt.Sprintf("%s: detected %q", variant.name, raw))
		}
	}
	if len(urls) == 0 {
		diagnostics = append(diagnostics, "no QR URL detected")
	}
	return qrDetectionResult{URLs: urls, Diagnostics: diagnostics}
}

type qrImageVariant struct {
	name  string
	image image.Image
}

func qrImageVariants(img image.Image) []qrImageVariant {
	variants := []qrImageVariant{{name: "original", image: img}}
	b := img.Bounds()
	width := b.Dx()
	if width <= 0 {
		return variants
	}
	if width < 1800 {
		up2 := imaging.Resize(img, width*2, 0, imaging.NearestNeighbor)
		variants = append(variants,
			qrImageVariant{name: "upscale2", image: up2},
			qrImageVariant{name: "upscale2-threshold120", image: thresholdGray(up2, 120)},
			qrImageVariant{name: "upscale2-threshold140", image: thresholdGray(up2, 140)},
		)
	}
	if width < 900 {
		up3 := imaging.Resize(img, width*3, 0, imaging.NearestNeighbor)
		variants = append(variants,
			qrImageVariant{name: "upscale3", image: up3},
			qrImageVariant{name: "upscale3-threshold120", image: thresholdGray(up3, 120)},
			qrImageVariant{name: "upscale3-threshold140", image: thresholdGray(up3, 140)},
		)
		up4 := imaging.Resize(img, width*4, 0, imaging.NearestNeighbor)
		variants = append(variants,
			qrImageVariant{name: "upscale4", image: up4},
			qrImageVariant{name: "upscale4-threshold140", image: thresholdGray(up4, 140)},
		)
	}
	return variants
}

func decodeQRCodeTexts(img image.Image) ([]string, string) {
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil, err.Error()
	}
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}
	seen := map[string]bool{}
	var texts []string
	var errors []string

	multiReader := multiqrcode.NewQRCodeMultiReader()
	results, err := multiReader.DecodeMultiple(bitmap, hints)
	if err != nil {
		errors = append(errors, err.Error())
	}
	for _, result := range results {
		text := result.GetText()
		if text == "" || seen[text] {
			continue
		}
		seen[text] = true
		texts = append(texts, text)
	}

	singleReader := qrcode.NewQRCodeReader()
	result, err := singleReader.Decode(bitmap, hints)
	if err != nil {
		errors = append(errors, err.Error())
	} else if text := result.GetText(); text != "" && !seen[text] {
		texts = append(texts, text)
	}

	return texts, strings.Join(errors, "; ")
}

func thresholdGray(img image.Image, threshold uint8) image.Image {
	b := img.Bounds()
	out := image.NewGray(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			if gray < threshold {
				gray = 0
			} else {
				gray = 255
			}
			out.SetGray(x-b.Min.X, y-b.Min.Y, color.Gray{Y: gray})
		}
	}
	return out
}

func isHTTPURL(raw string) bool {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return false
	}
	return (parsed.Scheme == "https" || parsed.Scheme == "http") && parsed.Hostname() != ""
}

func qrDiscordContent(urls []string) string {
	if len(urls) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("QRコードのURL:\n")
	for _, u := range urls {
		b.WriteString(u)
		b.WriteByte('\n')
	}
	return strings.TrimSpace(b.String())
}
