package main

import (
	"path/filepath"
	"testing"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

func TestAcquireInstanceLockPreventsSecondLock(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	first, err := acquireInstanceLock(configPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = first.Unlock()
	}()

	second, err := acquireInstanceLock(configPath)
	if err == nil {
		_ = second.Unlock()
		t.Fatal("expected second lock to fail")
	}
}

func TestAcquireInstanceLockAllowsAfterUnlock(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	first, err := acquireInstanceLock(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Unlock(); err != nil {
		t.Fatal(err)
	}

	second, err := acquireInstanceLock(configPath)
	if err != nil {
		t.Fatal(err)
	}
	_ = second.Unlock()
}

func TestShouldExitWithoutUI(t *testing.T) {
	cfg := appcore.DefaultConfig()
	cfg.Output.ShowUI = "auto"
	cfg.Output.CopySingleURLToClipboard = true
	results := []appcore.Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}
	if !shouldExitWithoutUI(cfg, results) {
		t.Fatal("auto mode with copied single URL should exit without UI")
	}

	cfg.Output.ShowUI = "always"
	if shouldExitWithoutUI(cfg, results) {
		t.Fatal("always mode should keep UI open")
	}

	cfg.Output.ShowUI = "never"
	if !shouldExitWithoutUI(cfg, results) {
		t.Fatal("never mode with successful single result should exit")
	}

	if shouldExitWithoutUI(cfg, []appcore.Result{{Error: "x"}}) {
		t.Fatal("error result should keep UI open")
	}
}

func TestHasErrors(t *testing.T) {
	if hasErrors([]appcore.Result{{}, {Error: "x"}}) != true {
		t.Fatal("expected hasErrors true")
	}
	if hasErrors([]appcore.Result{{}}) != false {
		t.Fatal("expected hasErrors false")
	}
}
