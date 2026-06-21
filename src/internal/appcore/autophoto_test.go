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
