package appcore

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessorProcessPathsSavesLocalImage(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.png")
	writeTestPNG(t, source, 4, 4)

	cfg := DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")
	cfg.Image.Suffix = "_small"
	cfg.Image.OutputFormat = "jpg"
	cfg.Output.UploadDiscord = false
	cfg.Output.SaveLocal = true

	results, err := Processor{Config: cfg}.ProcessPaths([]string{source})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Error != "" {
		t.Fatalf("unexpected result error: %s", results[0].Error)
	}
	if !strings.HasSuffix(results[0].OutputPath, "_small.jpg") {
		t.Fatalf("OutputPath = %q", results[0].OutputPath)
	}
	if _, err := os.Stat(results[0].OutputPath); err != nil {
		t.Fatal(err)
	}
	if results[0].Thumbnail == "" {
		t.Fatal("expected thumbnail data URL")
	}
}

func TestProcessorRejectsMultipleClipboardOnlyOutputs(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Output.SaveLocal = false
	cfg.Output.UploadDiscord = false

	_, err := Processor{Config: cfg}.ProcessPaths([]string{"a.png", "b.png"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "複数画像") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessorReportsDecodeErrorPerFile(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "not-image.png")
	if err := os.WriteFile(source, []byte("not an image"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	cfg.Output.UploadDiscord = false

	results, err := Processor{Config: cfg}.ProcessPaths([]string{source})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Error == "" {
		t.Fatalf("expected per-file decode error, got %+v", results)
	}
}

func TestProcessorEmitsProgressEvents(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.png")
	second := filepath.Join(dir, "second.png")
	writeTestPNG(t, first, 2, 2)
	writeTestPNG(t, second, 2, 2)

	cfg := DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")
	cfg.Output.UploadDiscord = false
	var events []ProgressEvent
	_, err := Processor{Config: cfg}.ProcessPathsWithProgress([]string{first, second}, func(event ProgressEvent) {
		events = append(events, event)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 4 {
		t.Fatalf("len(events) = %d, want 4", len(events))
	}
	if events[0].Stage != "processing" || events[1].Stage != "done" {
		t.Fatalf("unexpected first events: %+v", events[:2])
	}
}

func TestProcessorUploadDiscordRequiresWebhook(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.png")
	writeTestPNG(t, source, 2, 2)
	cfg := DefaultConfig()
	cfg.Output.SaveLocal = false
	cfg.Output.UploadDiscord = true
	cfg.Discord.WebhookURL = ""

	results, err := Processor{Config: cfg}.ProcessPaths([]string{source})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !strings.Contains(results[0].Error, "Webhook URLが未設定") {
		t.Fatalf("expected webhook error, got %+v", results)
	}
}

func TestCopySingleURLIfNeededSkipsNonSingleOrError(t *testing.T) {
	cfg := DefaultConfig()
	if err := CopySingleURLIfNeeded(cfg, []Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}, {URL: "https://cdn.discordapp.com/attachments/1/2/b.png"}}); err != nil {
		t.Fatal(err)
	}
	if err := CopySingleURLIfNeeded(cfg, []Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png", Error: "x"}}); err != nil {
		t.Fatal(err)
	}
}

func TestProcessImageResizesToFit(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Output.SaveLocal = false
	cfg.Output.UploadDiscord = true
	cfg.Discord.WebhookURL = ""
	img := image.NewRGBA(image.Rect(0, 0, 20, 10))
	cfg.Image.MaxWidth = 5
	cfg.Image.MaxHeight = 5

	result := Processor{Config: cfg}.processImage(img, "png", "source.png", false)
	if !strings.Contains(result.Error, "Webhook URLが未設定") {
		t.Fatalf("expected webhook error after encode path, got %+v", result)
	}
}
