package main

import (
	"path/filepath"
	"testing"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

func TestSettingsDraftSaveLoadRemove(t *testing.T) {
	dir := t.TempDir()
	old := singleInstanceDirFunc
	singleInstanceDirFunc = func() (string, error) {
		return dir, nil
	}
	t.Cleanup(func() { singleInstanceDirFunc = old })

	cfg := appcore.DefaultConfig()
	cfg.Image.Suffix = "_draft"
	if err := saveSettingsDraft(cfg); err != nil {
		t.Fatal(err)
	}
	if _, err := filepath.Abs(filepath.Join(dir, settingsDraftFile)); err != nil {
		t.Fatal(err)
	}
	got, ok, err := loadSettingsDraft()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("draft not found")
	}
	if got.Image.Suffix != "_draft" {
		t.Fatalf("suffix = %q, want _draft", got.Image.Suffix)
	}
	if err := removeSettingsDraft(); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := loadSettingsDraft(); err != nil || ok {
		t.Fatalf("after remove ok=%t err=%v", ok, err)
	}
}

func TestSettingsDraftForConfigIgnoresUnchangedNormalizedDraft(t *testing.T) {
	dir := t.TempDir()
	old := singleInstanceDirFunc
	singleInstanceDirFunc = func() (string, error) {
		return dir, nil
	}
	t.Cleanup(func() { singleInstanceDirFunc = old })

	configPath := filepath.Join(t.TempDir(), "config.json")
	baseline := appcore.DefaultConfig()
	baseline.Normalize()
	draft := baseline
	if err := saveSettingsDraftForConfig(configPath, draft, &baseline); err != nil {
		t.Fatal(err)
	}

	_, ok, reason, diffPaths, err := loadSettingsDraftForConfig(configPath, baseline, false)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("unchanged draft was restored: diff=%v", diffPaths)
	}
	if reason == "" {
		t.Fatal("expected ignored reason")
	}
}

func TestSettingsDraftForConfigIgnoresLegacyDraftOnFirstLaunch(t *testing.T) {
	dir := t.TempDir()
	old := singleInstanceDirFunc
	singleInstanceDirFunc = func() (string, error) {
		return dir, nil
	}
	t.Cleanup(func() { singleInstanceDirFunc = old })

	baseline := appcore.DefaultConfig()
	baseline.Normalize()
	draft := baseline
	draft.AutoCapture.OSC.SendPort = 9100
	if err := saveSettingsDraft(draft); err != nil {
		t.Fatal(err)
	}

	_, ok, reason, _, err := loadSettingsDraftForConfig(filepath.Join(t.TempDir(), "config.json"), baseline, false)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("legacy first-launch draft was restored")
	}
	if reason == "" {
		t.Fatal("expected ignored reason")
	}
}

func TestConfigDiffPathsReportsOSCForwardTargets(t *testing.T) {
	before := appcore.DefaultConfig()
	before.Normalize()
	after := before
	after.AutoCapture.OSC.Forward.Targets = append(after.AutoCapture.OSC.Forward.Targets, appcore.OSCForwardTarget{Host: "127.0.0.1", Port: 9101})

	diffPaths := configDiffPaths(before, after)
	if len(diffPaths) != 1 || diffPaths[0] != "autoCapture.osc.forward.targets" {
		t.Fatalf("diffPaths = %#v, want autoCapture.osc.forward.targets", diffPaths)
	}
}
