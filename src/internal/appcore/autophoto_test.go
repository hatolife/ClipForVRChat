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

func TestScanPhotoFilesAllWithExcludesStatusIgnoresScanCap(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 6; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%05d.png", i))
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	limited, limitedStatus := scanPhotoFilesLimitedWithExcludesStatus(dir, 3, nil)
	if len(limited) != 3 {
		t.Fatalf("len(limited) = %d, want 3", len(limited))
	}
	if !limitedStatus.LimitReached {
		t.Fatal("expected limited scan to report limit reached")
	}
	all, status := scanPhotoFilesAllWithExcludesStatus(dir, nil)
	if status.Error != "" {
		t.Fatalf("unexpected scan status error: %s", status.Error)
	}
	if status.LimitReached {
		t.Fatal("expected full scan to avoid limit reached status")
	}
	if len(all) != 6 {
		t.Fatalf("len(all) = %d, want 6", len(all))
	}
}

func TestAutoPhotoWatcherFullScanBaselineDoesNotResurfaceDeletedOldFiles(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 8; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%05d.png", i))
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	watcher := AutoPhotoWatcher{
		Config: Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: dir}},
		Process: func(path string) Result {
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}
	seen, status := scanPhotoFilesAllWithExcludesStatus(dir, nil)
	if status.Error != "" {
		t.Fatalf("unexpected scan status error: %s", status.Error)
	}
	if len(seen) != 8 {
		t.Fatalf("len(seen) = %d, want 8", len(seen))
	}
	watcher.seen = seen

	if err := os.Remove(filepath.Join(dir, "00000.png")); err != nil {
		t.Fatal(err)
	}

	newPath := filepath.Join(dir, "99999.png")
	if err := os.WriteFile(newPath, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	var processed []string
	watcher.Process = func(path string) Result {
		processed = append(processed, path)
		return Result{Name: filepath.Base(path), SourcePath: path}
	}
	watcher.tick()

	if len(processed) != 1 || processed[0] != newPath {
		t.Fatalf("processed = %v, want only %q", processed, newPath)
	}
}

func TestAutoPhotoWatcherFullScanBaselineProcessesLateArrivingFile(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 8; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%05d.png", i))
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	watcher := AutoPhotoWatcher{
		Config: Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: dir}},
		Process: func(path string) Result {
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}
	seen, status := scanPhotoFilesAllWithExcludesStatus(dir, nil)
	if status.Error != "" {
		t.Fatalf("unexpected scan status error: %s", status.Error)
	}
	watcher.seen = seen

	newPath := filepath.Join(dir, "99999.png")
	if err := os.WriteFile(newPath, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	var processed []string
	watcher.Process = func(path string) Result {
		processed = append(processed, path)
		return Result{Name: filepath.Base(path), SourcePath: path}
	}
	watcher.tick()

	if len(processed) != 1 || processed[0] != newPath {
		t.Fatalf("processed = %v, want only %q", processed, newPath)
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

func TestAutoPhotoWatcherSuppressesRepeatedScanStatus(t *testing.T) {
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
	watcher.tick()

	if len(events) != 1 {
		t.Fatalf("events = %+v, want repeated scan status suppressed", events)
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

func TestAutoPhotoWatcherCountsSkippedFilesTowardPerTickLimit(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < MaxAutoPhotoProcessPerTick+1; i++ {
		path := filepath.Join(dir, fmt.Sprintf("%05d.png", i))
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	processed := 0
	skipped := 0
	watcher := AutoPhotoWatcher{
		Config: Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: dir}},
		seen:   map[string]time.Time{},
		ShouldSkip: func(path string) bool {
			skipped++
			return true
		},
		Process: func(path string) Result {
			processed++
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}

	watcher.tick()

	if skipped != MaxAutoPhotoProcessPerTick {
		t.Fatalf("skipped = %d, want %d", skipped, MaxAutoPhotoProcessPerTick)
	}
	if processed != 0 {
		t.Fatalf("processed = %d, want 0", processed)
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

func TestAutoPhotoWatcherExcludesAutoCaptureOutputDirectory(t *testing.T) {
	vrchatDir := t.TempDir()
	autoCaptureDir := filepath.Join(vrchatDir, "VRC-AutoCapture")
	if err := os.MkdirAll(autoCaptureDir, 0700); err != nil {
		t.Fatal(err)
	}
	vrchatPath := filepath.Join(vrchatDir, "vrchat.png")
	autoCapturePath := filepath.Join(autoCaptureDir, "auto-capture.png")
	for _, path := range []string{vrchatPath, autoCapturePath} {
		if err := os.WriteFile(path, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	var processed []string
	watcher := AutoPhotoWatcher{
		Config:             Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: vrchatDir}},
		ExcludeDirectories: []string{autoCaptureDir},
		seen:               map[string]time.Time{},
		Process: func(path string) Result {
			processed = append(processed, path)
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}

	watcher.tick()

	if len(processed) != 1 || processed[0] != vrchatPath {
		t.Fatalf("processed = %v, want only %q", processed, vrchatPath)
	}
}

func TestAutoPhotoWatcherShouldSkipMarksPathSeen(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "auto-capture-photo.png")
	if err := os.WriteFile(source, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	processed := 0
	skipped := 0
	watcher := AutoPhotoWatcher{
		Config: Config{AutoPhoto: AutoPhotoConfig{PhotoDirectory: dir}},
		seen:   map[string]time.Time{},
		ShouldSkip: func(path string) bool {
			if path == source {
				skipped++
				return true
			}
			return false
		},
		Process: func(path string) Result {
			processed++
			return Result{Name: filepath.Base(path), SourcePath: path}
		},
	}

	watcher.tick()
	watcher.ShouldSkip = nil
	watcher.tick()

	if skipped != 1 {
		t.Fatalf("skipped = %d, want 1", skipped)
	}
	if processed != 0 {
		t.Fatalf("processed = %d, want 0", processed)
	}
}

func TestAutoPhotoWatcherProcessDoesNotForceDiscordUpload(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.png")
	writeTestPNG(t, source, 2, 2)

	cfg := DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")
	cfg.Output.SaveLocal = true
	cfg.Output.UploadDiscord = false
	cfg.Discord.WebhookURL = ""

	watcher := AutoPhotoWatcher{Config: cfg}
	result := watcher.process(source)
	if result.Error != "" {
		t.Fatalf("unexpected result error: %s", result.Error)
	}
	if result.URL != "" {
		t.Fatalf("URL = %q, want empty when Discord upload is disabled", result.URL)
	}
	if result.OutputPath == "" {
		t.Fatal("expected local output path")
	}
	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Fatal(err)
	}
}

func TestAutoPhotoWatcherWebhookURLUsesExplicitOverrideOnly(t *testing.T) {
	watcher := AutoPhotoWatcher{
		WebhookURL: "",
		Config: Config{
			AutoPhoto: AutoPhotoConfig{WebhookURL: "https://discord.com/api/webhooks/auto/token"},
		},
	}
	if got := watcher.webhookURL(); got != "" {
		t.Fatalf("webhookURL() = %q, want empty without explicit override", got)
	}
	watcher.WebhookURL = " https://discord.com/api/webhooks/explicit/token "
	if got := watcher.webhookURL(); got != "https://discord.com/api/webhooks/explicit/token" {
		t.Fatalf("webhookURL() = %q, want explicit override", got)
	}
}
