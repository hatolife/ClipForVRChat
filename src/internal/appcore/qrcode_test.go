package appcore

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
	"testing"

	qrcode "github.com/skip2/go-qrcode"
)

func TestDetectQRCodeURLsFindsMultipleURLs(t *testing.T) {
	img := compositeQRImage(t, []string{
		"https://example.com/first",
		"https://example.com/second",
	})

	urls := DetectQRCodeURLs(img)
	if len(urls) != 2 {
		t.Fatalf("len(urls) = %d, want 2: %+v", len(urls), urls)
	}
	if urls[0] != "https://example.com/first" || urls[1] != "https://example.com/second" {
		t.Fatalf("urls = %+v", urls)
	}
}

func TestDetectQRCodeURLsSkipsNonURLs(t *testing.T) {
	img := compositeQRImage(t, []string{"not a url"})
	if urls := DetectQRCodeURLs(img); len(urls) != 0 {
		t.Fatalf("urls = %+v, want empty", urls)
	}
}

func TestQRDiscordContentFormatsMultipleURLs(t *testing.T) {
	got := qrDiscordContent([]string{"https://example.com/1", "https://example.com/2"})
	if !strings.Contains(got, "QRコードのURL:") ||
		!strings.Contains(got, "https://example.com/1") ||
		!strings.Contains(got, "https://example.com/2") {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestQRImageVariantsSkipsUpscalesWhenScaledPixelsExceedLimit(t *testing.T) {
	img := staticTestImage{rect: image.Rect(0, 0, 2, MaxImagePixels/2)}

	variants := qrImageVariants(img)
	if len(variants) != 1 {
		t.Fatalf("len(variants) = %d, want 1: %+v", len(variants), variantNames(variants))
	}
	if variants[0].name != "original" {
		t.Fatalf("variants[0].name = %q, want original", variants[0].name)
	}
}

func TestQRImageVariantsKeepsUpscalesWithinLimit(t *testing.T) {
	img := staticTestImage{rect: image.Rect(0, 0, 16, 16)}

	variants := qrImageVariants(img)
	names := variantNames(variants)
	for _, want := range []string{"upscale2", "upscale3", "upscale4"} {
		if !containsString(names, want) {
			t.Fatalf("variant names = %+v, want %q", names, want)
		}
	}
}

func compositeQRImage(t *testing.T, values []string) image.Image {
	t.Helper()
	const qrSize = 180
	const gap = 20
	width := len(values)*qrSize + (len(values)+1)*gap
	height := qrSize + 2*gap
	out := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(out, out.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	for i, value := range values {
		data, err := qrcode.Encode(value, qrcode.Medium, qrSize)
		if err != nil {
			t.Fatal(err)
		}
		qr, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		x := gap + i*(qrSize+gap)
		draw.Draw(out, image.Rect(x, gap, x+qrSize, gap+qrSize), qr, image.Point{}, draw.Src)
	}
	return out
}

func writeImagePNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
}

type staticTestImage struct {
	rect image.Rectangle
}

func (s staticTestImage) ColorModel() color.Model { return color.NRGBAModel }

func (s staticTestImage) Bounds() image.Rectangle { return s.rect }

func (s staticTestImage) At(x, y int) color.Color { return color.White }

func variantNames(variants []qrImageVariant) []string {
	names := make([]string, 0, len(variants))
	for _, variant := range variants {
		names = append(names, variant.name)
	}
	return names
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
