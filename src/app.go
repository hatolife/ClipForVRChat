package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
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

func (a *App) LoadConfig() (appcore.Config, error) {
	return appcore.LoadConfig(a.configPath)
}

func (a *App) SaveConfig(cfg appcore.Config) error {
	if err := appcore.SaveConfig(a.configPath, cfg); err != nil {
		return err
	}
	a.state.Config = cfg
	return nil
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
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
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
