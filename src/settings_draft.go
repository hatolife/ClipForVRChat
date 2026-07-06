package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

const settingsDraftFile = "settings-draft.json"

type settingsDraftState struct {
	Config    appcore.Config `json:"config"`
	UpdatedAt string         `json:"updatedAt"`
}

func settingsDraftPath() (string, error) {
	dir, err := singleInstanceDirFunc()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsDraftFile), nil
}

func loadSettingsDraft() (appcore.Config, bool, error) {
	path, err := settingsDraftPath()
	if err != nil {
		return appcore.Config{}, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return appcore.Config{}, false, nil
		}
		return appcore.Config{}, false, fmt.Errorf("設定一時変更を読み込めませんでした: %w", err)
	}
	var draft settingsDraftState
	if err := json.Unmarshal(data, &draft); err != nil {
		return appcore.Config{}, false, fmt.Errorf("設定一時変更が壊れています: %w", err)
	}
	draft.Config.Normalize()
	return draft.Config, true, nil
}

func saveSettingsDraft(cfg appcore.Config) error {
	path, err := settingsDraftPath()
	if err != nil {
		return err
	}
	cfg.Normalize()
	draft := settingsDraftState{
		Config:    cfg,
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
	}
	data, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("設定一時変更フォルダを作成できませんでした: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("設定一時変更を書き込めませんでした: %w", err)
	}
	return nil
}

func removeSettingsDraft() error {
	path, err := settingsDraftPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("設定一時変更を削除できませんでした: %w", err)
	}
	return nil
}
