package appcore

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AutoPhotoEvent struct {
	Path   string `json:"path"`
	Result Result `json:"result"`
	Error  string `json:"error"`
}

type AutoPhotoWatcher struct {
	Config   Config
	Interval time.Duration
	Handler  func(AutoPhotoEvent)
	seen     map[string]time.Time
}

func (w *AutoPhotoWatcher) Run(ctx context.Context) {
	if w.Interval <= 0 {
		w.Interval = 2 * time.Second
	}
	w.seen = scanPhotoFiles(w.Config.AutoPhoto.PhotoDirectory)
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
	current := scanPhotoFiles(w.Config.AutoPhoto.PhotoDirectory)
	for path, modTime := range current {
		if _, ok := w.seen[path]; ok {
			continue
		}
		if !fileLooksStable(path) {
			continue
		}
		w.seen[path] = modTime
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
	cfg := w.Config
	cfg.Output.UploadDiscord = true
	if strings.TrimSpace(cfg.AutoPhoto.WebhookURL) != "" {
		cfg.Discord.WebhookURL = cfg.AutoPhoto.WebhookURL
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

func scanPhotoFiles(dir string) map[string]time.Time {
	files := make(map[string]time.Time)
	if strings.TrimSpace(dir) == "" {
		return files
	}
	_ = filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !isSupportedPhotoPath(path) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		files[path] = info.ModTime()
		return nil
	})
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
