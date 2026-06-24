package appcore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanPhotoFilesLimitsFileCount(t *testing.T) {
	dir := t.TempDir()
	const maxFiles = 3
	for i := 0; i < maxFiles+3; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%05d.png", i))
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	files := scanPhotoFilesLimited(dir, maxFiles)
	if len(files) != maxFiles {
		t.Fatalf("len(files) = %d, want %d", len(files), maxFiles)
	}
	_, status := scanPhotoFilesLimitedWithStatus(dir, maxFiles)
	if !status.LimitReached {
		t.Fatal("expected limit reached status")
	}
}

func TestScanPhotoFilesReportsMissingDirectory(t *testing.T) {
	_, status := scanPhotoFilesLimitedWithStatus(filepath.Join(t.TempDir(), "missing"), 3)
	if status.Error == "" {
		t.Fatal("expected missing directory error")
	}
}

func TestAutoPhotoWatcherEmitsScanStatus(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	var events []AutoPhotoEvent
	watcher := AutoPhotoWatcher{
		Directory: missing,
		seen:      map[string]time.Time{},
		Handler: func(event AutoPhotoEvent) {
			events = append(events, event)
		},
	}

	watcher.tick()

	if len(events) != 1 || events[0].Error == "" || events[0].Result.Error == "" {
		t.Fatalf("events = %+v, want scan error event", events)
	}
}

func TestAutoPhotoWatcherTickLimitsProcessing(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < MaxAutoPhotoProcessPerTick+3; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%05d.png", i))
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	processed := 0
	watcher := AutoPhotoWatcher{
		Config: Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: dir}},
		seen:   map[string]time.Time{},
		Process: func(path string) Result {
			processed++
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}

	watcher.tick()

	if processed != MaxAutoPhotoProcessPerTick {
		t.Fatalf("processed = %d, want %d", processed, MaxAutoPhotoProcessPerTick)
	}
}

func TestAutoPhotoWatcherUsesExplicitDirectory(t *testing.T) {
	vrchatDir := t.TempDir()
	screenshotDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(vrchatDir, "vrchat.png"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	screenshotPath := filepath.Join(screenshotDir, "screenshot.png")
	if err := os.WriteFile(screenshotPath, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	var processed []string
	watcher := AutoPhotoWatcher{
		Config:    Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: vrchatDir}},
		Directory: screenshotDir,
		seen:      map[string]time.Time{},
		Process: func(path string) Result {
			processed = append(processed, path)
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}

	watcher.tick()

	if len(processed) != 1 || processed[0] != screenshotPath {
		t.Fatalf("processed = %v, want only %q", processed, screenshotPath)
	}
}
