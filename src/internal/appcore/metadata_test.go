package appcore

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteAutoCaptureEmbeddedMetadataPNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "capture.png")
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	metadata := AutoCaptureEmbeddedMetadata{SchemaVersion: 1, App: "ClipForVRChat", BatchID: "batch", ShotID: "shot"}
	if err := WriteAutoCaptureEmbeddedMetadata(path, metadata); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(embeddedMetadataKeyword)) || !bytes.Contains(data, []byte(`"batch_id":"batch"`)) {
		t.Fatalf("PNG metadata payload was not embedded")
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("embedded PNG should remain decodable: %v", err)
	}
}

func TestWriteAutoCaptureEmbeddedMetadataJPEG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "capture.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{G: 255, A: 255})
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	metadata := AutoCaptureEmbeddedMetadata{SchemaVersion: 1, App: "ClipForVRChat", BatchID: "batch", ShotID: "shot"}
	if err := WriteAutoCaptureEmbeddedMetadata(path, metadata); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("Exif\x00\x00")) || !bytes.Contains(data, []byte(`"shot_id":"shot"`)) {
		t.Fatalf("JPEG EXIF payload was not embedded")
	}
	if _, err := jpeg.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("embedded JPEG should remain decodable: %v", err)
	}
}

func TestBuildAutoCaptureEmbeddedMetadataMasksUserIDs(t *testing.T) {
	cfg := DefaultAutoCaptureConfig()
	cfg.Output.WriteEXIF = true
	cfg.Output.WriteUserListToEXIF = true
	cfg.Output.WriteUserIDsToEXIF = false
	users := []PresenceUser{{DisplayName: "Alice", UserID: "usr_secret"}}
	metadata := BuildAutoCaptureEmbeddedMetadata(cfg, "batch", "shot", cfg.Views[0], users, "confirmed", SpoutCaptureResult{})
	if len(metadata.Users) != 1 || metadata.Users[0].UserID != "" {
		t.Fatalf("metadata users = %+v, want user ID masked", metadata.Users)
	}
	cfg.Output.WriteUserIDsToEXIF = true
	metadata = BuildAutoCaptureEmbeddedMetadata(cfg, "batch", "shot", cfg.Views[0], users, "confirmed", SpoutCaptureResult{})
	if !strings.Contains(metadata.Users[0].UserID, "usr_secret") {
		t.Fatalf("metadata users = %+v, want user ID preserved", metadata.Users)
	}
}
