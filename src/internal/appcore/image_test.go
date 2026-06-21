package appcore

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecodeImageAcceptsSmallPNG(t *testing.T) {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.White)
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}

	decoded, format, err := DecodeImage(buf.Bytes(), "small.png")
	if err != nil {
		t.Fatal(err)
	}
	if format != "png" {
		t.Fatalf("format = %q, want png", format)
	}
	if decoded.Bounds().Dx() != 2 || decoded.Bounds().Dy() != 2 {
		t.Fatalf("bounds = %v, want 2x2", decoded.Bounds())
	}
}

func TestDecodeImageFileRejectsLargeFileBeforeRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(maxImageInputBytes(1) + 1); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if _, _, err := DecodeImageFileWithLimit(path, 1); err == nil {
		t.Fatal("expected large file to be rejected")
	}
}

func TestDecodeImageRejectsLargeDimensions(t *testing.T) {
	data := pngHeaderOnly(100_000, 100_000)
	if _, _, err := DecodeImage(data, "huge.png"); err == nil {
		t.Fatal("expected large dimensions to be rejected")
	}
}

func TestResizeToFitKeepsSmallImageAndResizesLargeImage(t *testing.T) {
	small := solidImage(4, 3)
	if got := ResizeToFit(small, 10, 10); got != small {
		t.Fatal("small image should be returned unchanged")
	}

	large := solidImage(20, 10)
	got := ResizeToFit(large, 5, 5)
	if got.Bounds().Dx() != 5 || got.Bounds().Dy() != 2 {
		t.Fatalf("resized bounds = %v, want 5x2", got.Bounds())
	}
}

func TestEncodeImageFormats(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 255, A: 128})

	pngOut, err := EncodeImage(img, "png", 92)
	if err != nil {
		t.Fatal(err)
	}
	if pngOut.Extension != ".png" || pngOut.Mime != "image/png" || len(pngOut.Data) == 0 {
		t.Fatalf("unexpected png output: %+v", pngOut)
	}

	jpgOut, err := EncodeImage(img, "jpg", 80)
	if err != nil {
		t.Fatal(err)
	}
	if jpgOut.Extension != ".jpg" || jpgOut.Mime != "image/jpeg" || len(jpgOut.Data) == 0 {
		t.Fatalf("unexpected jpg output: %+v", jpgOut)
	}
}

func TestSaveEncodedImageUsesSuffixAndAvoidsOverwrite(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "avatar.png")
	writeTestPNG(t, source, 2, 2)
	encoded, err := EncodeImage(solidImage(2, 2), "png", 92)
	if err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")
	cfg.Image.Suffix = "_vrchat"
	cfg.Image.Overwrite = false

	first, err := SaveEncodedImage(encoded, source, cfg, false)
	if err != nil {
		t.Fatal(err)
	}
	second, err := SaveEncodedImage(encoded, source, cfg, false)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(first) != "avatar_vrchat.png" {
		t.Fatalf("first = %q", first)
	}
	if filepath.Base(second) != "avatar_vrchat_2.png" {
		t.Fatalf("second = %q", second)
	}
}

func TestSaveEncodedImageUsesClipboardName(t *testing.T) {
	dir := t.TempDir()
	encoded, err := EncodeImage(solidImage(2, 2), "png", 92)
	if err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")

	path, err := SaveEncodedImage(encoded, "clipboard.png", cfg, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(filepath.Base(path), "clipboard_") {
		t.Fatalf("clipboard output base = %q", filepath.Base(path))
	}
}

func TestThumbnailDataURL(t *testing.T) {
	dataURL, err := ThumbnailDataURL(solidImage(10, 10))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		t.Fatalf("unexpected data URL prefix: %q", dataURL)
	}
}

func pngHeaderOnly(width, height uint32) []byte {
	var out bytes.Buffer
	out.Write([]byte{137, 80, 78, 71, 13, 10, 26, 10})

	var ihdr bytes.Buffer
	_ = binary.Write(&ihdr, binary.BigEndian, width)
	_ = binary.Write(&ihdr, binary.BigEndian, height)
	ihdr.Write([]byte{8, 2, 0, 0, 0})
	writePNGChunk(&out, "IHDR", ihdr.Bytes())
	writePNGChunk(&out, "IEND", nil)
	return out.Bytes()
}

func writePNGChunk(out *bytes.Buffer, kind string, data []byte) {
	_ = binary.Write(out, binary.BigEndian, uint32(len(data)))
	out.WriteString(kind)
	out.Write(data)
	crc := crc32.NewIEEE()
	_, _ = crc.Write([]byte(kind))
	_, _ = crc.Write(data)
	_ = binary.Write(out, binary.BigEndian, crc.Sum32())
}

func writeTestPNG(t *testing.T, path string, width, height int) {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, solidImage(width, height)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
}

func solidImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	return img
}
