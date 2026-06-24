package main

import (
	"bytes"
	"path/filepath"
	"strings"
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
	if !shouldExitWithoutUI(cfg, results, nil) {
		t.Fatal("auto mode with copied single URL should exit without UI")
	}

	cfg.Output.ShowUI = "always"
	if shouldExitWithoutUI(cfg, results, nil) {
		t.Fatal("always mode should keep UI open")
	}

	cfg.Output.ShowUI = "never"
	if !shouldExitWithoutUI(cfg, results, nil) {
		t.Fatal("never mode with successful single result should exit")
	}

	if shouldExitWithoutUI(cfg, []appcore.Result{{Error: "x"}}, nil) {
		t.Fatal("error result should keep UI open")
	}

	cfg.Output.ShowUI = "auto"
	if shouldExitWithoutUI(cfg, results, errTestCopyFailed{}) {
		t.Fatal("copy failure should keep UI open")
	}
}

type errTestCopyFailed struct{}

func (errTestCopyFailed) Error() string {
	return "copy failed"
}

func TestHasErrors(t *testing.T) {
	if hasErrors([]appcore.Result{{}, {Error: "x"}}) != true {
		t.Fatal("expected hasErrors true")
	}
	if hasErrors([]appcore.Result{{}}) != false {
		t.Fatal("expected hasErrors false")
	}
}

func TestHandleCLIArgsVersion(t *testing.T) {
	oldVersion := version
	oldRevision := revision
	t.Cleanup(func() {
		version = oldVersion
		revision = oldRevision
	})
	version = "v1.2.3"
	revision = "abcdef0"

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	handled, code := handleCLIArgs([]string{"--version"}, &stdout, &stderr)

	if !handled || code != 0 {
		t.Fatalf("handled=%t code=%d", handled, code)
	}
	if got := stdout.String(); got != "ClipForVRChat v1.2.3.abcdef0\n" {
		t.Fatalf("stdout = %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestHandleCLIArgsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	handled, code := handleCLIArgs([]string{"--help"}, &stdout, &stderr)

	if !handled || code != 0 {
		t.Fatalf("handled=%t code=%d", handled, code)
	}
	if got := stdout.String(); !strings.Contains(got, "--version") || !strings.Contains(got, "--help") {
		t.Fatalf("help output = %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestHandleCLIArgsLeavesExistingPositionalArgsAlone(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	handled, code := handleCLIArgs([]string{"image.png"}, &stdout, &stderr)

	if handled || code != 0 {
		t.Fatalf("handled=%t code=%d", handled, code)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}
