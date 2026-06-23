package appcore

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	MaxAutoPhotoScanFiles      = 5000
	MaxAutoPhotoProcessPerTick = 5
)

type AutoPhotoEvent struct {
	Path   string `json:"path"`
	Result Result `json:"result"`
	Error  string `json:"error"`
}

type AutoPhotoWatcher struct {
	Config     Config
	Directory  string
	WebhookURL string
	Interval   time.Duration
	Handler    func(AutoPhotoEvent)
	Process    func(string) Result
	seen       map[string]time.Time
}

func (w *AutoPhotoWatcher) Run(ctx context.Context) {
	if w.Interval <= 0 {
		w.Interval = 2 * time.Second
	}
	w.seen = scanPhotoFiles(w.directory())
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick()
		}
	}
}

func (w *AutoPhotoWatcher) tick() {
	current := scanPhotoFiles(w.directory())
	processed := 0
	for _, path := range sortedPhotoPaths(current) {
		modTime := current[path]
		if _, ok := w.seen[path]; ok {
			continue
		}
		if processed >= MaxAutoPhotoProcessPerTick {
			break
		}
		if !fileLooksStable(path) {
			continue
		}
		w.seen[path] = modTime
		processed++
		result := w.process(path)
		event := AutoPhotoEvent{Path: path, Result: result}
		if result.Error != "" {
			event.Error = result.Error
		}
		if w.Handler != nil {
			w.Handler(event)
		}
	}
	for path := range w.seen {
		if _, ok := current[path]; !ok {
			delete(w.seen, path)
		}
	}
}

func (w *AutoPhotoWatcher) process(path string) Result {
	if w.Process != nil {
		return w.Process(path)
	}
	cfg := w.Config
	cfg.Output.UploadDiscord = true
	if strings.TrimSpace(w.webhookURL()) != "" {
		cfg.Discord.WebhookURL = w.webhookURL()
	}
	results, err := Processor{Config: cfg}.ProcessPaths([]string{path})
	if err != nil {
		return Result{SourcePath: path, Name: filepath.Base(path), Error: err.Error()}
	}
	if len(results) == 0 {
		return Result{SourcePath: path, Name: filepath.Base(path), Error: "画像を処理できませんでした。"}
	}
	return results[0]
}

func (w *AutoPhotoWatcher) directory() string {
	if strings.TrimSpace(w.Directory) != "" {
		return w.Directory
	}
	return w.Config.AutoPhoto.PhotoDirectory
}

func (w *AutoPhotoWatcher) webhookURL() string {
	if strings.TrimSpace(w.WebhookURL) != "" {
		return w.WebhookURL
	}
	return w.Config.AutoPhoto.WebhookURL
}

func scanPhotoFiles(dir string) map[string]time.Time {
	return scanPhotoFilesLimited(dir, MaxAutoPhotoScanFiles)
}

func scanPhotoFilesLimited(dir string, maxFiles int) map[string]time.Time {
	files := make(map[string]time.Time)
	if strings.TrimSpace(dir) == "" {
		return files
	}
	if maxFiles <= 0 {
		return files
	}
	paths := make([]string, 0, maxFiles)
	_ = filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !isSupportedPhotoPath(path) {
			return nil
		}
		if len(paths) >= maxFiles {
			return filepath.SkipAll
		}
		paths = append(paths, path)
		return nil
	})
	sort.Strings(paths)
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		files[path] = info.ModTime()
	}
	return files
}

func isSupportedPhotoPath(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".webp":
		return true
	default:
		return false
	}
}

func sortedPhotoPaths(files map[string]time.Time) []string {
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func fileLooksStable(path string) bool {
	first, err := os.Stat(path)
	if err != nil || first.Size() == 0 {
		return false
	}
	time.Sleep(700 * time.Millisecond)
	second, err := os.Stat(path)
	if err != nil {
		return false
	}
	return first.Size() == second.Size() && first.ModTime().Equal(second.ModTime())
}
