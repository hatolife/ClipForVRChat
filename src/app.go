package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	configPath string
	state      appcore.UIState
	autoCancel context.CancelFunc
	mu         sync.Mutex
}

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	GitHub  string `json:"github"`
}

type OSSLicense struct {
	Name      string `json:"name"`
	License   string `json:"license"`
	Copyright string `json:"copyright"`
	URL       string `json:"url"`
}

func NewApp(configPath string, initial appcore.UIState) *App {
	return &App{
		configPath: configPath,
		state:      initial,
	}
}

func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ctx = ctx
	a.restartAutoPhotoWatcher(a.state.Config)
}

func (a *App) GetInitialState() appcore.UIState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.refreshHistory()
	return a.state
}

func (a *App) GetAppInfo() AppInfo {
	return AppInfo{
		Name:    "ClipForVRChat",
		Version: appVersion(),
		GitHub:  githubURL,
	}
}

func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

func (a *App) CheckForUpdate() (appcore.UpdateInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	return appcore.CheckLatestRelease(ctx, nil, appVersion(), appReleaseTime())
}

func (a *App) GetOSSLicenses() []OSSLicense {
	return []OSSLicense{
		{Name: "Wails", License: "MIT", Copyright: "Copyright (c) 2018-Present Lea Anthony", URL: "https://github.com/wailsapp/wails"},
		{Name: "Vue.js", License: "MIT", Copyright: "Copyright (c) 2018-present, Yuxi (Evan) You", URL: "https://github.com/vuejs/core"},
		{Name: "Vite", License: "MIT", Copyright: "Copyright (c) 2019-present, VoidZero Inc. and Vite contributors", URL: "https://github.com/vitejs/vite"},
		{Name: "imaging", License: "MIT", Copyright: "Copyright (c) 2012 Grigory Dryapak", URL: "https://github.com/disintegration/imaging"},
		{Name: "golang.design/x/clipboard", License: "MIT", Copyright: "Copyright (c) 2021 Changkun Ou", URL: "https://github.com/golang-design/clipboard"},
		{Name: "golang.org/x/image", License: "BSD-3-Clause", Copyright: "Copyright (c) The Go Authors", URL: "https://cs.opensource.google/go/x/image"},
	}
}

func (a *App) SelectOutputDirectory(current string) (string, error) {
	return a.selectDirectory("出力先フォルダを選択", current)
}

func (a *App) SelectPhotoDirectory(current string) (string, error) {
	return a.selectDirectory("VRChat写真フォルダを選択", current)
}

func (a *App) SelectScreenshotDirectory(current string) (string, error) {
	return a.selectDirectory("スクリーンショットフォルダを選択", current)
}

func (a *App) selectDirectory(title string, current string) (string, error) {
	dir := strings.Trim(strings.TrimSpace(current), `"`)
	if dir != "" && !filepath.IsAbs(dir) {
		exe, err := os.Executable()
		if err == nil {
			dir = filepath.Join(filepath.Dir(exe), dir)
		}
	}
	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		dir = ""
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: dir,
	})
}

func (a *App) ClearResults() appcore.UIState {
	a.mu.Lock()
	defer a.mu.Unlock()
	var ids []string
	for _, result := range a.state.Results {
		if result.HistoryID != "" {
			ids = append(ids, result.HistoryID)
		}
	}
	if len(ids) > 0 {
		if history, err := appcore.MarkHistoryCleared(appcore.HistoryPath(a.configPath), ids); err == nil {
			a.state.History = history
		}
	}
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state
}

func (a *App) GetHistory() ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.LoadHistory(appcore.HistoryPath(a.configPath))
	if err != nil {
		return nil, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) MarkHistoryCleared(ids []string) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.MarkHistoryCleared(appcore.HistoryPath(a.configPath), ids)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) SetHistoryPinned(id string, pinned bool) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.SetHistoryPinned(appcore.HistoryPath(a.configPath), id, pinned)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) DeleteDiscordHistoryEntries(ids []string) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.LoadHistory(appcore.HistoryPath(a.configPath))
	if err != nil {
		return history, err
	}
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	for i := range history {
		if !idSet[history[i].ID] || history[i].Pinned {
			continue
		}
		if err := appcore.DeleteDiscordMessage(history[i].DiscordWebhookID, history[i].DiscordToken, history[i].DiscordMessageID); err != nil {
			return history, err
		}
		history[i].DiscordDeleted = true
		history[i].DeletedAt = nowRFC3339()
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(a.configPath), history); err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) PurgeDeletedHistoryEntries() ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return a.state.History, err
	}
	history, _, err := appcore.PurgeDiscordDeletedHistory(appcore.HistoryPath(a.configPath), cfg.Output.DeleteOutputOnHistoryPurge)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) LoadConfig() (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return appcore.LoadConfig(a.configPath)
}

func (a *App) SaveConfig(cfg appcore.Config) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.saveConfigLocked(cfg)
}

func (a *App) saveConfigLocked(cfg appcore.Config) error {
	if err := appcore.SaveConfig(a.configPath, cfg); err != nil {
		return err
	}
	a.state.Config = cfg
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.restartAutoPhotoWatcher(cfg)
	return nil
}

func (a *App) CloseSettings() appcore.UIState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state
}

func (a *App) OpenSettings(path string) (appcore.UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if strings.TrimSpace(path) != "" {
		a.configPath = path
	}
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return a.state, err
	}
	a.state.Mode = appcore.ModeSettings
	a.state.Message = ""
	a.state.ConfigPath = a.configPath
	a.state.Config = cfg
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state, nil
}

func (a *App) SaveConfigAndProcess(cfg appcore.Config, paths []string) (appcore.UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.saveConfigLocked(cfg); err != nil {
		return a.state, err
	}
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
	if err != nil {
		a.state.Mode = appcore.ModeError
		a.state.Message = err.Error()
		return a.state, nil
	}
	_ = appcore.CopySingleURLIfNeeded(cfg, results)
	if history, historyErr := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results); historyErr == nil {
		a.state.History = history
	}
	a.state.Results = results
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	if hasResultErrors(results) {
		a.state.Mode = appcore.ModeError
		a.state.Message = "処理中にエラーが発生しました。内容を確認してください。"
	} else {
		a.state.Mode = appcore.ModeResults
		a.state.Message = resultMessage(cfg, results)
	}
	return a.state, nil
}

func (a *App) ProcessToState(paths []string) (appcore.UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if hasJSONPath(paths) {
		a.state.Mode = appcore.ModeError
		a.state.Message = "画像ファイルと設定ファイルが混在しています。設定編集と画像処理は別々に実行してください。"
		a.state.Results = nil
		return a.state, nil
	}
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return a.state, err
	}
	results, err := appcore.Processor{Config: cfg}.ProcessPathsWithProgress(paths, func(event appcore.ProgressEvent) {
		runtime.EventsEmit(a.ctx, "process:progress", event)
	})
	if err != nil {
		a.state.Mode = appcore.ModeError
		a.state.Message = err.Error()
		a.state.Results = nil
		return a.state, nil
	}
	_ = appcore.CopySingleURLIfNeeded(cfg, results)
	if history, historyErr := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results); historyErr == nil {
		a.state.History = history
	}
	a.state.Config = cfg
	a.state.Results = results
	if hasResultErrors(results) {
		a.state.Mode = appcore.ModeError
		a.state.Message = "処理中にエラーが発生しました。内容を確認してください。"
	} else {
		a.state.Mode = appcore.ModeResults
		a.state.Message = resultMessage(cfg, results)
	}
	return a.state, nil
}

func (a *App) Process(paths []string) ([]appcore.Result, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return nil, err
	}
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
	if err != nil {
		return nil, err
	}
	_ = appcore.CopySingleURLIfNeeded(cfg, results)
	_, _ = appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results)
	return results, nil
}

func (a *App) CopyText(text string) error {
	return appcore.WriteClipboardText(text)
}

func defaultConfigPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	return appcore.ConfigPath(exe)
}

func hasResultErrors(results []appcore.Result) bool {
	for _, result := range results {
		if result.Error != "" {
			return true
		}
	}
	return false
}

func resultMessage(cfg appcore.Config, results []appcore.Result) string {
	if len(results) == 1 && results[0].URL != "" && cfg.Output.CopySingleURLToClipboard {
		return "画像URLをクリップボードにコピーしました。"
	}
	if len(results) == 1 && results[0].URL == "" && results[0].Error == "" {
		var parts []string
		if results[0].OutputPath != "" {
			parts = append(parts, "縮小画像を保存しました")
		}
		if len(parts) == 0 {
			parts = append(parts, "縮小画像をクリップボードにコピーしました")
		}
		return strings.Join(parts, "。") + "。"
	}
	return ""
}

func hasJSONPath(paths []string) bool {
	for _, path := range paths {
		if strings.EqualFold(filepath.Ext(path), ".json") {
			return true
		}
	}
	return false
}

func (a *App) refreshHistory() {
	if history, err := appcore.LoadHistory(appcore.HistoryPath(a.configPath)); err == nil {
		a.state.History = history
	}
}

func nowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func (a *App) restartAutoPhotoWatcher(cfg appcore.Config) {
	if a.autoCancel != nil {
		a.autoCancel()
		a.autoCancel = nil
	}
	if a.ctx == nil || (!cfg.AutoPhoto.Enabled && !cfg.ScreenshotAutoPost.Enabled) {
		return
	}
	ctx, cancel := context.WithCancel(a.ctx)
	a.autoCancel = cancel
	handler := func(event appcore.AutoPhotoEvent) {
		a.mu.Lock()
		defer a.mu.Unlock()
		if event.Result.URL != "" && event.Result.Error == "" {
			results := []appcore.Result{event.Result}
			if history, err := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results); err == nil {
				a.state.History = history
				event.Result = results[0]
			}
		}
		a.state.Results = append([]appcore.Result{event.Result}, a.state.Results...)
		a.state.Mode = appcore.ModeResults
		runtime.EventsEmit(a.ctx, "auto-photo:result", event)
	}
	interval := time.Duration(cfg.AutoPhoto.ScanIntervalSeconds) * time.Second
	if cfg.AutoPhoto.Enabled {
		watcher := appcore.AutoPhotoWatcher{
			Config:     cfg,
			Directory:  cfg.AutoPhoto.PhotoDirectory,
			WebhookURL: cfg.AutoPhoto.WebhookURL,
			Interval:   interval,
			Handler:    handler,
		}
		go watcher.Run(ctx)
	}
	if cfg.ScreenshotAutoPost.Enabled {
		watcher := appcore.AutoPhotoWatcher{
			Config:    cfg,
			Directory: cfg.ScreenshotAutoPost.ScreenshotDirectory,
			Interval:  interval,
			Handler:   handler,
		}
		go watcher.Run(ctx)
	}
}
