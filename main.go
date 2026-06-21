package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	configPath := defaultConfigPath()
	args := os.Args[1:]

	state := appcore.UIState{
		Mode:       appcore.ModeResults,
		ConfigPath: configPath,
	}

	cfg, err := appcore.LoadConfig(configPath)
	if err != nil {
		state.Mode = appcore.ModeError
		state.Message = fmt.Sprintf("設定を読み込めませんでした: %v", err)
		state.Config = appcore.DefaultConfig()
		runUI(configPath, state)
		return
	}
	state.Config = cfg

	if len(args) == 1 && strings.EqualFold(filepath.Ext(args[0]), ".json") {
		configPath = args[0]
		cfg, err := appcore.LoadConfig(configPath)
		if err != nil {
			state.Mode = appcore.ModeError
			state.Message = fmt.Sprintf("設定を読み込めませんでした: %v", err)
		} else {
			state.Mode = appcore.ModeSettings
			state.Config = cfg
			state.ConfigPath = configPath
		}
		runUI(configPath, state)
		return
	}

	for _, arg := range args {
		if strings.EqualFold(filepath.Ext(arg), ".json") {
			state.Mode = appcore.ModeError
			state.Message = "画像ファイルと設定ファイルが混在しています。設定編集と画像処理は別々に起動してください。"
			runUI(configPath, state)
			return
		}
	}

	results, err := appcore.Processor{Config: cfg}.ProcessPaths(args)
	if err != nil {
		state.Mode = appcore.ModeError
		state.Message = err.Error()
		runUI(configPath, state)
		return
	}
	_ = appcore.CopySingleURLIfNeeded(cfg, results)

	state.Results = results
	if shouldExitWithoutUI(cfg, results) {
		return
	}
	if hasErrors(results) {
		state.Mode = appcore.ModeError
		state.Message = "処理中にエラーが発生しました。内容を確認してください。"
	} else {
		state.Mode = appcore.ModeResults
	}
	runUI(configPath, state)
}

func shouldExitWithoutUI(cfg appcore.Config, results []appcore.Result) bool {
	if cfg.Output.ShowUI == "always" {
		return false
	}
	if len(results) != 1 || hasErrors(results) {
		return false
	}
	if cfg.Output.ShowUI == "never" {
		return true
	}
	return results[0].URL != "" && cfg.Output.CopySingleURLToClipboard
}

func hasErrors(results []appcore.Result) bool {
	for _, result := range results {
		if result.Error != "" {
			return true
		}
	}
	return false
}

func runUI(configPath string, state appcore.UIState) {
	app := NewApp(configPath, state)
	err := wails.Run(&options.App{
		Title:  "ClipForVRChat",
		Width:  900,
		Height: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
