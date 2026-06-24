package appcore

import (
	"context"
	"fmt"
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

type autoPhotoScanStatus struct {
	Error        string
	LimitReached bool
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
	var status autoPhotoScanStatus
	w.seen, status = scanPhotoFilesWithStatus(w.directory())
	w.emitScanStatus(status)
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
	current, status := scanPhotoFilesWithStatus(w.directory())
	w.emitScanStatus(status)
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

func (w *AutoPhotoWatcher) emitScanStatus(status autoPhotoScanStatus) {
	if w.Handler == nil {
		return
	}
	message := strings.TrimSpace(status.Error)
	if message == "" && status.LimitReached {
		message = fmt.Sprintf("自動投稿の監視対象が%d件を超えたため、一部のファイルはこのスキャンで確認されませんでした。", MaxAutoPhotoScanFiles)
	}
	if message == "" {
		return
	}
	name := "自動投稿"
	if dir := w.directory(); strings.TrimSpace(dir) != "" {
		name = filepath.Base(dir)
	}
	w.Handler(AutoPhotoEvent{
		Path: w.directory(),
		Result: Result{
			SourcePath: w.directory(),
			Name:       name,
			Error:      message,
		},
		Error: message,
	})
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
	files, _ := scanPhotoFilesLimitedWithStatus(dir, MaxAutoPhotoScanFiles)
	return files
}

func scanPhotoFilesLimited(dir string, maxFiles int) map[string]time.Time {
	files, _ := scanPhotoFilesLimitedWithStatus(dir, maxFiles)
	return files
}

func scanPhotoFilesWithStatus(dir string) (map[string]time.Time, autoPhotoScanStatus) {
	return scanPhotoFilesLimitedWithStatus(dir, MaxAutoPhotoScanFiles)
}

func scanPhotoFilesLimitedWithStatus(dir string, maxFiles int) (map[string]time.Time, autoPhotoScanStatus) {
	files := make(map[string]time.Time)
	if strings.TrimSpace(dir) == "" {
		return files, autoPhotoScanStatus{Error: "自動投稿の監視フォルダが未設定です。"}
	}
	if maxFiles <= 0 {
		return files, autoPhotoScanStatus{Error: "自動投稿のスキャン上限が0以下です。"}
	}
	info, err := os.Stat(dir)
	if err != nil {
		return files, autoPhotoScanStatus{Error: fmt.Sprintf("自動投稿の監視フォルダを確認できません: %v", err)}
	}
	if !info.IsDir() {
		return files, autoPhotoScanStatus{Error: "自動投稿の監視フォルダにファイルが指定されています。"}
	}
	status := autoPhotoScanStatus{}
	paths := make([]string, 0, maxFiles)
	_ = filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			if status.Error == "" {
				status.Error = fmt.Sprintf("自動投稿の監視フォルダを走査できません: %v", err)
			}
			return nil
		}
		if entry.IsDir() || !isSupportedPhotoPath(path) {
			return nil
		}
		if len(paths) >= maxFiles {
			status.LimitReached = true
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
	return files, status
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
