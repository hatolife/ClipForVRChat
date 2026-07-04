package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

func TestAcquireInstanceLockPreventsSecondLock(t *testing.T) {
	withSingleInstanceTestDir(t)
	first, err := acquireInstanceLock()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = first.Unlock()
	}()

	second, err := acquireInstanceLock()
	if err == nil {
		_ = second.Unlock()
		t.Fatal("expected second lock to fail")
	}
}

func TestAcquireInstanceLockAllowsAfterUnlock(t *testing.T) {
	withSingleInstanceTestDir(t)
	first, err := acquireInstanceLock()
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Unlock(); err != nil {
		t.Fatal(err)
	}

	second, err := acquireInstanceLock()
	if err != nil {
		t.Fatal(err)
	}
	_ = second.Unlock()
}

func TestSingleInstanceStateRoundTrip(t *testing.T) {
	dir := withSingleInstanceTestDir(t)
	path := filepath.Join(dir, singleInstanceStateFile)
	want := singleInstanceState{
		PID:            123,
		ExecutablePath: "/tmp/ClipForVRChat.exe",
		Version:        "v1.2.3",
		Revision:       "abcdef0",
		Endpoint:       "127.0.0.1:12345",
		Token:          strings.Repeat("a", 64),
		StartedAt:      "2026-07-04T00:00:00Z",
	}
	if err := writeSingleInstanceState(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := readSingleInstanceState(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("state = %+v, want %+v", got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("state file mode = %v, want regular file", info.Mode())
	}
}

func TestSingleInstanceServerCommands(t *testing.T) {
	server, err := startSingleInstanceServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	var activated bool
	var shutdown bool
	server.SetHandlers(func() error {
		activated = true
		return nil
	}, func() error {
		shutdown = true
		return nil
	})
	state := singleInstanceState{Endpoint: server.Endpoint(), Token: server.Token()}
	if err := sendSingleInstanceCommand(state, "ping", time.Second); err != nil {
		t.Fatal(err)
	}
	if err := sendSingleInstanceCommand(state, "activate", time.Second); err != nil {
		t.Fatal(err)
	}
	if err := sendSingleInstanceCommand(state, "shutdown", time.Second); err != nil {
		t.Fatal(err)
	}
	if !activated || !shutdown {
		t.Fatalf("activated=%t shutdown=%t", activated, shutdown)
	}
}

func withSingleInstanceTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	old := singleInstanceDirFunc
	singleInstanceDirFunc = func() (string, error) {
		return dir, nil
	}
	t.Cleanup(func() {
		singleInstanceDirFunc = old
	})
	return dir
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
	oldBuildChannel := buildChannel
	t.Cleanup(func() {
		version = oldVersion
		revision = oldRevision
		buildChannel = oldBuildChannel
	})
	version = "v1.2.3"
	revision = "abcdef0"
	buildChannel = "release"

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

func TestFrontendAssetSummaryReportsEmbeddedIndex(t *testing.T) {
	got := frontendAssetSummary()
	if !strings.Contains(got, "index_html=true") {
		t.Fatalf("frontendAssetSummary() = %q, want embedded index.html", got)
	}
}
