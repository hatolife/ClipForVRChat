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
