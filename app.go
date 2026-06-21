package main

import (
	"context"
	"os"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

type App struct {
	ctx        context.Context
	configPath string
	state      appcore.UIState
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
