package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

const settingsDraftFile = "settings-draft.json"

type settingsDraftState struct {
	Config         appcore.Config  `json:"config"`
	ConfigPath     string          `json:"configPath,omitempty"`
	BaselineConfig *appcore.Config `json:"baselineConfig,omitempty"`
	UpdatedAt      string          `json:"updatedAt"`
}

func settingsDraftPath() (string, error) {
	dir, err := singleInstanceDirFunc()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsDraftFile), nil
}

func loadSettingsDraft() (appcore.Config, bool, error) {
	draft, ok, _, _, err := loadSettingsDraftForConfig("", appcore.Config{}, true)
	return draft, ok, err
}

func loadSettingsDraftForConfig(configPath string, baseline appcore.Config, configExists bool) (appcore.Config, bool, string, []string, error) {
	path, err := settingsDraftPath()
	if err != nil {
		return appcore.Config{}, false, "", nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return appcore.Config{}, false, "", nil, nil
		}
		return appcore.Config{}, false, "", nil, fmt.Errorf("設定一時変更を読み込めませんでした: %w", err)
	}
	var draft settingsDraftState
	if err := json.Unmarshal(data, &draft); err != nil {
		return appcore.Config{}, false, "", nil, fmt.Errorf("設定一時変更が壊れています: %w", err)
	}
	draft.Config.Normalize()
	if strings.TrimSpace(configPath) != "" {
		if strings.TrimSpace(draft.ConfigPath) == "" && !configExists {
			_ = removeSettingsDraft()
			return appcore.Config{}, false, "legacy draft without config path for first launch", nil, nil
		}
		if strings.TrimSpace(draft.ConfigPath) != "" && !sameSettingsDraftConfigPath(draft.ConfigPath, configPath) {
			_ = removeSettingsDraft()
			return appcore.Config{}, false, fmt.Sprintf("config path mismatch draft=%q current=%q", draft.ConfigPath, configPath), nil, nil
		}
		diffPaths := configDiffPaths(baseline, draft.Config)
		if len(diffPaths) == 0 {
			_ = removeSettingsDraft()
			return appcore.Config{}, false, "draft has no changes from current config", nil, nil
		}
		return draft.Config, true, "", diffPaths, nil
	}
	return draft.Config, true, "", nil, nil
}

func saveSettingsDraft(cfg appcore.Config) error {
	return saveSettingsDraftForConfig("", cfg, nil)
}

func saveSettingsDraftForConfig(configPath string, cfg appcore.Config, baseline *appcore.Config) error {
	path, err := settingsDraftPath()
	if err != nil {
		return err
	}
	cfg.Normalize()
	draft := settingsDraftState{
		Config:     cfg,
		ConfigPath: configPath,
		UpdatedAt:  time.Now().Format(time.RFC3339Nano),
	}
	if baseline != nil {
		base := *baseline
		base.Normalize()
		draft.BaselineConfig = &base
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

func sameSettingsDraftConfigPath(a string, b string) bool {
	clean := func(value string) string {
		value = strings.TrimSpace(value)
		if abs, err := filepath.Abs(value); err == nil {
			value = abs
		}
		return filepath.Clean(value)
	}
	return strings.EqualFold(clean(a), clean(b))
}

func configDiffPaths(before appcore.Config, after appcore.Config) []string {
	before.Normalize()
	after.Normalize()
	var beforeJSON any
	var afterJSON any
	beforeData, err := json.Marshal(before)
	if err != nil {
		return []string{"config"}
	}
	afterData, err := json.Marshal(after)
	if err != nil {
		return []string{"config"}
	}
	if err := json.Unmarshal(beforeData, &beforeJSON); err != nil {
		return []string{"config"}
	}
	if err := json.Unmarshal(afterData, &afterJSON); err != nil {
		return []string{"config"}
	}
	return diffJSONPaths(beforeJSON, afterJSON, "")
}

func diffJSONPaths(before any, after any, prefix string) []string {
	if reflect.DeepEqual(before, after) {
		return nil
	}
	beforeMap, beforeOK := before.(map[string]any)
	afterMap, afterOK := after.(map[string]any)
	if !beforeOK || !afterOK {
		if prefix == "" {
			return []string{"config"}
		}
		return []string{prefix}
	}
	keySet := make(map[string]bool, len(beforeMap)+len(afterMap))
	for key := range beforeMap {
		keySet[key] = true
	}
	for key := range afterMap {
		keySet[key] = true
	}
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	paths := make([]string, 0)
	for _, key := range keys {
		child := key
		if prefix != "" {
			child = prefix + "." + key
		}
		paths = append(paths, diffJSONPaths(beforeMap[key], afterMap[key], child)...)
	}
	return paths
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
