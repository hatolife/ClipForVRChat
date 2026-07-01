package appcore

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCapturedImageRejectsBlankFrames(t *testing.T) {
	tests := []struct {
		name string
		fill color.RGBA
		want string
	}{
		{name: "white", fill: color.RGBA{R: 255, G: 255, B: 255, A: 255}, want: "ほぼ白"},
		{name: "black", fill: color.RGBA{A: 255}, want: "ほぼ黒"},
		{name: "transparent", fill: color.RGBA{}, want: "ほぼ透明"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "capture.png")
			writeSolidPNG(t, path, tt.fill)
			_, err := validateCapturedImage(path)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateCapturedImageAcceptsVariedFrame(t *testing.T) {
	path := filepath.Join(t.TempDir(), "capture.png")
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 8), G: uint8(y * 8), B: 80, A: 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := validateCapturedImage(path); err != nil {
		t.Fatal(err)
	}
}

func writeSolidPNG(t *testing.T, path string, fill color.RGBA) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.SetRGBA(x, y, fill)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
