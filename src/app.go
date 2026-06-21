package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	configPath string
	state      appcore.UIState
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
	a.ctx = ctx
}

func (a *App) GetInitialState() appcore.UIState {
	return a.state
}

func (a *App) GetAppInfo() AppInfo {
	return AppInfo{
		Name:    "ClipForVRChat",
		Version: version,
		GitHub:  githubURL,
	}
}

func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

func (a *App) GetOSSLicenses() []OSSLicense {
	return []OSSLicense{
		{Name: "ClipForVRChat", License: "MIT", Copyright: "Copyright (c) 2026 ClipForVRChat contributors", URL: githubURL},
		{Name: "Wails", License: "MIT", Copyright: "Copyright (c) 2018-Present Lea Anthony", URL: "https://github.com/wailsapp/wails"},
		{Name: "Vue.js", License: "MIT", Copyright: "Copyright (c) 2018-present, Yuxi (Evan) You", URL: "https://github.com/vuejs/core"},
		{Name: "Vite", License: "MIT", Copyright: "Copyright (c) 2019-present, VoidZero Inc. and Vite contributors", URL: "https://github.com/vitejs/vite"},
		{Name: "imaging", License: "MIT", Copyright: "Copyright (c) 2012 Grigory Dryapak", URL: "https://github.com/disintegration/imaging"},
		{Name: "golang.design/x/clipboard", License: "MIT", Copyright: "Copyright (c) 2021 Changkun Ou", URL: "https://github.com/golang-design/clipboard"},
		{Name: "golang.org/x/image", License: "BSD-3-Clause", Copyright: "Copyright (c) The Go Authors", URL: "https://cs.opensource.google/go/x/image"},
	}
}

func (a *App) ClearResults() appcore.UIState {
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state
}

func (a *App) LoadConfig() (appcore.Config, error) {
	return appcore.LoadConfig(a.configPath)
}

func (a *App) SaveConfig(cfg appcore.Config) error {
	if err := appcore.SaveConfig(a.configPath, cfg); err != nil {
		return err
	}
	a.state.Config = cfg
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	return nil
}

func (a *App) CloseSettings() appcore.UIState {
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state
}

func (a *App) OpenSettings(path string) (appcore.UIState, error) {
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
	if err := a.SaveConfig(cfg); err != nil {
		return a.state, err
	}
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
	if err != nil {
		a.state.Mode = appcore.ModeError
		a.state.Message = err.Error()
		return a.state, nil
	}
	_ = appcore.CopySingleURLIfNeeded(cfg, results)
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
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return nil, err
	}
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
	if err != nil {
		return nil, err
	}
	_ = appcore.CopySingleURLIfNeeded(cfg, results)
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
