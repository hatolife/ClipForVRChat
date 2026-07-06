package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

func TestAppGetInitialStateRefreshesHistory(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	historyPath := appcore.HistoryPath(configPath)
	if err := appcore.SaveHistory(historyPath, []appcore.HistoryEntry{{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}); err != nil {
		t.Fatal(err)
	}

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	state := app.GetInitialState()
	if len(state.History) != 1 || state.History[0].ID != "1" {
		t.Fatalf("history = %+v", state.History)
	}
}

func appBoolPtr(value bool) *bool {
	return &value
}

func TestAppClearResultsMarksHistoryCleared(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{{ID: "h1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{
		Mode:    appcore.ModeResults,
		Results: []appcore.Result{{HistoryID: "h1", URL: history[0].URL}},
	})

	state := app.ClearResults()
	if len(state.Results) != 0 {
		t.Fatalf("Results = %+v, want empty", state.Results)
	}
	loaded, err := appcore.LoadHistory(appcore.HistoryPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || !loaded[0].Cleared {
		t.Fatalf("history = %+v, want cleared", loaded)
	}
}

func TestAppSaveConfigAndOpenSettings(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	cfg := appcore.DefaultConfig()
	cfg.Image.MaxWidth = 777

	if err := app.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	state, err := app.OpenSettings("")
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != appcore.ModeSettings || state.Config.Image.MaxWidth != 777 {
		t.Fatalf("state = %+v", state)
	}
	closed := app.CloseSettings()
	if closed.Mode != appcore.ModeResults {
		t.Fatalf("closed mode = %s, want results", closed.Mode)
	}
}

func TestAppSaveConfigWritesDiagnosticConfigSummary(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeSettings})
	cfg := appcore.DefaultConfig()
	cfg.Discord.WebhookURL = "https://discord.com/api/webhooks/primary-secret"
	cfg.AutoPhoto.WebhookURL = ""
	cfg.ScreenshotAutoPost.WebhookURL = ""
	cfg.AutoCapture.Discord.WebhookURL = ""

	if err := app.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range []string{
		"settings saved config=",
		`"discord":{"webhookConfigured":true}`,
		`"fallbackToPrimaryWebhook":true`,
		`"effectiveWebhookConfigured":true`,
		`"effectiveWebhookSource":"primary"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("diagnostic log = %q, want %q", text, want)
		}
	}
	if strings.Contains(text, "primary-secret") || strings.Contains(text, "discord.com/api/webhooks") {
		t.Fatalf("diagnostic log leaked webhook URL: %q", text)
	}
}

func TestAppSaveConfigLogsFailure(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	if err := os.Mkdir(configPath, 0700); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeSettings})

	if err := app.SaveConfig(appcore.DefaultConfig()); err == nil {
		t.Fatal("expected save error")
	}

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if text := string(data); !strings.Contains(text, "settings save error: err=") {
		t.Fatalf("diagnostic log = %q, want settings save error", text)
	}
}

func TestAppOpenSettingsKeepsConfigPathOnLoadError(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	badConfigPath := filepath.Join(dir, "bad.json")
	writeTextFile(t, badConfigPath, "{")

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults, ConfigPath: configPath})
	if _, err := app.OpenSettings(badConfigPath); err == nil {
		t.Fatal("expected invalid config error")
	}
	if app.configPath != configPath {
		t.Fatalf("configPath = %q, want %q", app.configPath, configPath)
	}
}

func TestAppSaveConfigAndProcessHandlesDecodeError(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	source := filepath.Join(dir, "bad.png")
	writeTextFile(t, source, "not image")

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeSettings, ProcessOnSave: true})
	cfg := appcore.DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")
	cfg.Output.UploadDiscord = false

	state, err := app.SaveConfigAndProcess(cfg, []string{source})
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != appcore.ModeError || !strings.Contains(state.Message, "処理中にエラー") {
		t.Fatalf("state = %+v", state)
	}
}

func TestAppProcessToStateHidesNoWorkResults(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	source := filepath.Join(dir, "image.png")
	writeTestPNG(t, source, 2, 2)
	cfg := appcore.DefaultConfig()
	cfg.Output.UploadDiscord = false
	cfg.Output.SaveLocal = false
	cfg.Output.DetectQRCodeURLs = false
	if err := appcore.SaveConfig(configPath, cfg); err != nil {
		t.Fatal(err)
	}

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	state, err := app.ProcessToState([]string{source})
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Results) != 0 {
		t.Fatalf("results = %+v, want empty", state.Results)
	}
	if !strings.Contains(state.Message, "実行された処理はありません") {
		t.Fatalf("message = %q", state.Message)
	}
	history, err := appcore.LoadHistory(appcore.HistoryPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 0 {
		t.Fatalf("history = %+v, want empty", history)
	}
}

func TestAppStartupStartsAutoPhotoWatcherWithoutDiscordUpload(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Output.UploadDiscord = false
	cfg.AutoPhoto.Enabled = true
	cfg.AutoPhoto.PhotoDirectory = t.TempDir()
	cfg.AutoPhoto.ScanIntervalSeconds = 3600

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults, Config: cfg})
	ctx, cancel := context.WithCancel(context.Background())
	app.startup(ctx)
	defer func() {
		cancel()
		app.shutdown(context.Background())
		time.Sleep(50 * time.Millisecond)
	}()

	if app.autoCancel == nil {
		t.Fatal("auto photo watcher was not started")
	}
}

func TestAppStartupDoesNotStartAutoPhotoWatcherWhileReviewingSettings(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "imported.json")
	cfg := appcore.DefaultConfig()
	cfg.Output.UploadDiscord = true
	cfg.AutoPhoto.Enabled = true
	cfg.AutoPhoto.PhotoDirectory = t.TempDir()
	cfg.AutoPhoto.WebhookURL = "https://example.com/api/webhooks/attacker/token"
	cfg.AutoPhoto.ScanIntervalSeconds = 3600

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeSettings, Config: cfg, ConfigPath: configPath})
	ctx, cancel := context.WithCancel(context.Background())
	app.startup(ctx)
	defer func() {
		cancel()
		app.shutdown(context.Background())
		time.Sleep(50 * time.Millisecond)
	}()

	if app.autoCancel != nil {
		t.Fatal("auto photo watcher started while imported settings were only open for review")
	}
}

func TestExplorerSelectArgsRejectsMissingPath(t *testing.T) {
	if _, err := explorerRevealPath("", ""); err == nil {
		t.Fatal("expected empty path error")
	}
	if _, err := explorerRevealPath(filepath.Join(t.TempDir(), "missing.png"), ""); err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestExplorerSelectArgsRejectsDirectory(t *testing.T) {
	if _, err := explorerRevealPath(t.TempDir(), ""); err == nil {
		t.Fatal("expected directory error")
	}
}

func TestExplorerSelectArgsResolvesRelativePathFromBaseDir(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "out.png"), []byte("image"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := explorerRevealPath("output/out.png", dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(outputDir, "out.png") {
		t.Fatalf("path = %q, want resolved relative output path", got)
	}
}

func TestAppLogUserActionWritesDiagnosticLog(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	app.LogUserAction("button_click", "history_delete_entries selected=1")

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, `ui action="button_click"`) || !strings.Contains(text, `history_delete_entries`) {
		t.Fatalf("diagnostic log = %q", text)
	}
	if filepath.Dir(appcore.DiagnosticLogPath(configPath)) != filepath.Join(dir, "logs") {
		t.Fatalf("diagnostic log path = %q, want dated log under log dir", appcore.DiagnosticLogPath(configPath))
	}
}

func TestAppLifecycleLogsCloseAndShutdown(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	app.domReady(context.Background())
	if prevent := app.beforeClose(context.Background()); prevent {
		t.Fatal("beforeClose should allow the window to close")
	}
	app.shutdown(context.Background())

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range []string{
		"ui lifecycle dom_ready",
		"ui lifecycle before_close: allowing close",
		"ui lifecycle shutdown begin",
		"ui lifecycle shutdown complete",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("diagnostic log = %q, want %q", text, want)
		}
	}
}

func TestTrustedExternalURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		ok   bool
	}{
		{name: "github", raw: "https://github.com/hatolife/ClipForVRChat/releases", ok: true},
		{name: "booth", raw: "https://hatolife.booth.pm/items/8531663", ok: true},
		{name: "discord help", raw: "https://support.discord.com/hc/ja/articles/228383668", ok: true},
		{name: "twitter", raw: "https://x.com/hato_poppo_life", ok: true},
		{name: "http", raw: "http://github.com/hatolife/ClipForVRChat", ok: false},
		{name: "file", raw: "file:///C:/Windows/win.ini", ok: false},
		{name: "unknown host", raw: "https://example.com", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := trustedExternalURL(tt.raw)
			if tt.ok && err != nil {
				t.Fatalf("trustedExternalURL() error = %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("trustedExternalURL() should reject URL")
			}
		})
	}
}

func TestIsGoTestBinary(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "/tmp/app.test", want: true},
		{path: `C:\Users\runner\AppData\Local\Temp\app.test.exe`, want: true},
		{path: `C:\app\ClipForVRChat.exe`, want: false},
		{path: "/tmp/not-test", want: false},
	}
	for _, tt := range tests {
		if got := isGoTestBinary(tt.path); got != tt.want {
			t.Fatalf("isGoTestBinary(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestAppStartupWritesVersionHashAndRedactedConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Discord.WebhookURL = "https://discord.com/api/webhooks/secret"
	cfg.AutoPhoto.WebhookURL = "https://discord.com/api/webhooks/auto-secret"
	app := NewApp(configPath, appcore.UIState{Config: cfg})

	ctx, cancel := context.WithCancel(context.Background())
	app.startup(ctx)
	defer func() {
		cancel()
		app.shutdown(context.Background())
		time.Sleep(50 * time.Millisecond)
	}()

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "startup app_version=") || !strings.Contains(text, "exe_sha256=") {
		t.Fatalf("diagnostic log = %q, want startup version and hash", text)
	}
	if !strings.Contains(text, `"webhookConfigured":true`) {
		t.Fatalf("diagnostic log = %q, want redacted webhook configured flag", text)
	}
	if strings.Contains(text, "secret") {
		t.Fatalf("diagnostic log leaked webhook URL: %q", text)
	}
}

func TestAppLatestPlayerLocalBasisManual(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceManual
	app.state.Config.AutoCapture.PlayerLocal.Calibrated = true
	app.state.Config.AutoCapture.PlayerLocal.BasisPose = appcore.CameraPoseConfig{
		Position: appcore.CameraVector3Config{X: 1, Y: 2, Z: 3},
		Rotation: appcore.CameraVector3Config{Y: 45},
	}

	got := app.GetLatestPlayerLocalBasis()
	if got.Source != "manual" || got.Status != "ready" || !got.Fresh {
		t.Fatalf("snapshot = %+v", got)
	}
	if got.Pose.Position.X != 1 || got.Pose.Rotation.Y != 45 {
		t.Fatalf("pose = %+v", got.Pose)
	}
}

func TestAppLatestPlayerLocalBasisAvatarOSC(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.PositionScale = 1000
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.InvertMagnitude = true
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.PositiveFlagThreshold = 0
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.MaxAbsPosition = 2000
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.MaxAbsForward = 2000
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	app.avatarOSCBasisSamples = map[string]avatarOSCBasisSample{
		"coord/x":       {Float: 0.877, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/xSign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/y":       {Float: 0.544, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/ySign":   {Float: 0, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/z":       {Float: 0.211, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/zSign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/x":     {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/xSign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/y":     {Float: 0.5, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/ySign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/z":     {Float: 0, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/zSign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
	}

	got := app.latestPlayerLocalBasisLocked(app.state.Config, now)
	if got.Source != "avatar_osc" || got.Status != "ready" || !got.Fresh {
		t.Fatalf("snapshot = %+v", got)
	}
	if got.ParameterPrefix != "coord" || got.RawSampleCount != 12 || got.LastReceivedAddress == "" {
		t.Fatalf("unexpected raw OSC status: %+v", got)
	}
	if math.Abs(got.Pose.Position.X-123) > 0.000001 || math.Abs(got.Pose.Position.Y+456) > 0.000001 || math.Abs(got.Pose.Position.Z-789) > 0.000001 {
		t.Fatalf("pose = %+v", got.Pose)
	}
	if math.Abs(got.Pose.Rotation.Y) > 0.000001 {
		t.Fatalf("rotation = %+v", got.Pose.Rotation)
	}
}

func TestAppAvatarOSCBasisStatusIgnoresManualBasisMissing(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceManual
	app.state.Config.AutoCapture.PlayerLocal.Calibrated = false

	got := app.GetAvatarOSCBasisStatus()
	if got.Source != "avatar_osc" || got.Status != "missing" {
		t.Fatalf("snapshot = %+v, want missing avatar_osc", got)
	}
	if strings.Contains(got.Error, "プレイヤー基準Pose") {
		t.Fatalf("error = %q, want avatar OSC diagnostic", got.Error)
	}
	if !strings.Contains(got.Error, "avatar OSC parameter") || !strings.Contains(got.Error, "coord/* / forward/*") {
		t.Fatalf("error = %q, want AvatarBeacon basis guidance", got.Error)
	}
}

func TestAppAvatarOSCBasisStatusResolvesEvenWhenManualBasisMissing(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceManual
	app.state.Config.AutoCapture.PlayerLocal.Calibrated = false
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	app.avatarOSCBasisSamples = map[string]avatarOSCBasisSample{
		"coord/x":       {Float: 0.8, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/xSign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/y":       {Float: 0.5, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/ySign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/z":       {Float: 0.2, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"coord/zSign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/x":     {Float: 0, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/xSign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/y":     {Float: 0, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/ySign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/z":     {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"forward/zSign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
	}

	got := app.latestAvatarOSCBasisSnapshotLocked(app.state.Config, now)
	if got.Source != "avatar_osc" || got.Status != "ready" || !got.Fresh {
		t.Fatalf("snapshot = %+v, want ready avatar_osc", got)
	}
	if strings.Contains(got.Error, "プレイヤー基準Pose") {
		t.Fatalf("error = %q, want no manual basis error", got.Error)
	}
}

func TestAppAvatarOSCBasisStatusExplainsStaleWithNonBasisLastParameter(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	now := time.Date(2026, 7, 4, 4, 35, 50, 0, time.UTC)
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	staleAt := now.Add(-31 * time.Second)
	app.avatarOSCBasisSamples = map[string]avatarOSCBasisSample{
		"coord/x":       {Float: 0.8, HasFloat: true, ReceivedAt: staleAt},
		"coord/xSign":   {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"coord/y":       {Float: 0.5, HasFloat: true, ReceivedAt: staleAt},
		"coord/ySign":   {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"coord/z":       {Float: 0.2, HasFloat: true, ReceivedAt: staleAt},
		"coord/zSign":   {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"forward/x":     {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"forward/xSign": {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"forward/y":     {Float: 0, HasFloat: true, ReceivedAt: staleAt},
		"forward/ySign": {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"forward/z":     {Float: 0, HasFloat: true, ReceivedAt: staleAt},
		"forward/zSign": {Float: 1, HasFloat: true, ReceivedAt: staleAt},
		"Ahoge_Angle":   {Float: 0.4, HasFloat: true, ReceivedAt: now},
	}

	got := app.latestAvatarOSCBasisSnapshotLocked(app.state.Config, now)
	if got.Status != "stale" || got.Fresh {
		t.Fatalf("snapshot = %+v, want stale", got)
	}
	if !got.AutoFallback {
		t.Fatalf("autoFallback = false, want true for stale AvatarBeacon basis")
	}
	if got.LastReceivedAddress != "Ahoge_Angle" {
		t.Fatalf("last address = %q, want Ahoge_Angle", got.LastReceivedAddress)
	}
	for _, want := range []string{"coord/* / forward/*", "AvatarBeaconのbasis parameterではありません", "Ahoge_Angle"} {
		if !strings.Contains(got.Error, want) {
			t.Fatalf("error = %q, want %q", got.Error, want)
		}
	}
}

func TestAppRebuildAvatarOSCBasisDoesNotRepeatPartialLog(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	logPath := appcore.DiagnosticLogPath(configPath)
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	app.avatarOSCBasisSamples = map[string]avatarOSCBasisSample{
		"coord/x": {Float: 0.8, HasFloat: true, ReceivedAt: now},
	}

	app.rebuildAvatarOSCBasisLocked(now, logPath)
	app.rebuildAvatarOSCBasisLocked(now.Add(10*time.Millisecond), logPath)

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), "avatar osc basis resolve status"); got != 1 {
		t.Fatalf("partial log count = %d, want 1; log=%s", got, string(data))
	}
}

func TestOSCForwarderForwardsPacket(t *testing.T) {
	receiver := listenUDPForTest(t)
	cfg := appcore.DefaultConfig().AutoCapture.OSC
	cfg.Forward.Enabled = true
	cfg.Forward.Mode = appcore.OSCForwardModeAll
	cfg.Forward.Targets = []appcore.OSCForwardTarget{{Host: "127.0.0.1", Port: receiver.LocalAddr().(*net.UDPAddr).Port}}

	forwarder := newOSCForwarder(cfg, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9001}, appcore.DiagnosticLogPath(filepath.Join(t.TempDir(), "config.json")), nil)
	defer forwarder.Close()

	packet := []byte("/avatar/parameters/test\x00\x00\x00\x00,T\x00\x00")
	forwarder.Forward(packet, "/avatar/parameters/test", true)

	got := readUDPForTest(t, receiver)
	if string(got) != string(packet) {
		t.Fatalf("forwarded packet = %q, want %q", got, packet)
	}
}

func TestOSCForwarderUnhandledOnlySkipsHandledPacket(t *testing.T) {
	receiver := listenUDPForTest(t)
	cfg := appcore.DefaultConfig().AutoCapture.OSC
	cfg.Forward.Enabled = true
	cfg.Forward.Mode = appcore.OSCForwardModeUnhandledOnly
	cfg.Forward.Targets = []appcore.OSCForwardTarget{{Host: "127.0.0.1", Port: receiver.LocalAddr().(*net.UDPAddr).Port}}

	forwarder := newOSCForwarder(cfg, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9001}, appcore.DiagnosticLogPath(filepath.Join(t.TempDir(), "config.json")), nil)
	defer forwarder.Close()

	forwarder.Forward([]byte("/avatar/parameters/test\x00\x00\x00\x00,T\x00\x00"), "/avatar/parameters/test", true)
	_ = receiver.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	buf := make([]byte, 256)
	if n, _, err := receiver.ReadFromUDP(buf); err == nil {
		t.Fatalf("unexpected forwarded packet bytes=%d", n)
	}
}

func TestOSCForwarderSkipsSelfTarget(t *testing.T) {
	cfg := appcore.DefaultConfig().AutoCapture.OSC
	cfg.Forward.Enabled = true
	cfg.Forward.Targets = []appcore.OSCForwardTarget{{Host: "127.0.0.1", Port: 9001}}

	forwarder := newOSCForwarder(cfg, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9001}, appcore.DiagnosticLogPath(filepath.Join(t.TempDir(), "config.json")), nil)
	defer forwarder.Close()

	if len(forwarder.targets) != 0 {
		t.Fatalf("targets = %d, want 0", len(forwarder.targets))
	}
}

func TestAppOSCLogEntriesAreBounded(t *testing.T) {
	app := NewApp(filepath.Join(t.TempDir(), "config.json"), appcore.UIState{Mode: appcore.ModeResults})
	for i := 0; i < 550; i++ {
		app.recordOSCLogEvent(OSCLogEntry{Direction: "receive", Address: fmt.Sprintf("/test/%d", i)})
	}
	got := app.GetOSCLogEntries()
	if len(got) != 500 {
		t.Fatalf("entries = %d, want 500", len(got))
	}
	if got[0].Address != "/test/50" || got[len(got)-1].Address != "/test/549" {
		t.Fatalf("unexpected ring contents: first=%q last=%q", got[0].Address, got[len(got)-1].Address)
	}
}

func listenUDPForTest(t *testing.T) *net.UDPConn {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})
	return conn
}

func freeUDPPortForTest(t *testing.T) int {
	t.Helper()
	conn := listenUDPForTest(t)
	port := conn.LocalAddr().(*net.UDPAddr).Port
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	return port
}

func readUDPForTest(t *testing.T, conn *net.UDPConn) []byte {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	buf := make([]byte, 256)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatal(err)
	}
	return append([]byte(nil), buf[:n]...)
}

func TestAppRestartCameraPoseReceiverKeepsSameEndpoint(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	logPath := appcore.DiagnosticLogPath(configPath)
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.ctx = context.Background()
	cfg := appcore.DefaultConfig()
	cfg.AutoCapture.OSC.ReceivePort = freeUDPPortForTest(t)
	cancelCalled := false
	app.oscCancel = func() {
		cancelCalled = true
	}
	app.oscReceiverHost = "127.0.0.1"
	app.oscReceiverPort = cfg.AutoCapture.OSC.ReceivePort
	app.oscReceiverLogPath = logPath
	app.oscReceiverForwardKey = oscForwardConfigKey(cfg.AutoCapture.OSC.Forward)
	app.oscReceiverSeq = 7

	app.restartCameraPoseReceiverLocked(cfg)

	if cancelCalled {
		t.Fatal("receiver was canceled for the same endpoint")
	}
	if app.oscCancel == nil || app.oscReceiverHost != "127.0.0.1" || app.oscReceiverPort != cfg.AutoCapture.OSC.ReceivePort || app.oscReceiverSeq != 7 {
		t.Fatalf("receiver state changed: cancel=%v host=%q port=%d seq=%d", app.oscCancel != nil, app.oscReceiverHost, app.oscReceiverPort, app.oscReceiverSeq)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "auto-capture osc receiver keep") {
		t.Fatalf("keep log missing: %s", string(data))
	}
}

func TestAppSaveCurrentCameraPoseToViewCapturesFreshZoom(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Normalize()
	cfg.AutoCapture.Views[0].CoordinateSpace = "world"
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults, Config: cfg})
	now := time.Now()
	app.latestPose = appcore.CameraPoseConfig{
		Position: appcore.CameraVector3Config{X: 1, Y: 2, Z: 3},
		Rotation: appcore.CameraVector3Config{X: 4, Y: 5, Z: 6},
	}
	app.poseAt = now
	app.userCameraSamples = map[string]userCameraOSCSample{
		"/usercamera/Zoom": {
			UserCameraOSCSample: appcore.UserCameraOSCSample{
				Address:  "/usercamera/Zoom",
				Float:    72.5,
				HasFloat: true,
			},
			ReceivedAt: now,
		},
	}

	got, err := app.SaveCurrentCameraPoseToView("front")
	if err != nil {
		t.Fatal(err)
	}
	view, ok := findAutoCaptureViewByID(got.AutoCapture.Views, "front")
	if !ok {
		t.Fatal("front view missing")
	}
	if view.Zoom == nil || *view.Zoom != 72.5 {
		t.Fatalf("zoom = %v, want 72.5", view.Zoom)
	}
	if view.Pose.Position.X != 1 || view.Pose.Rotation.Z != 6 || !view.Calibrated {
		t.Fatalf("view pose/calibrated = %+v", view)
	}
}

func TestAppSaveCurrentCameraPoseToViewKeepsZoomWithoutFreshSample(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Normalize()
	cfg.AutoCapture.Views[0].CoordinateSpace = "world"
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults, Config: cfg})
	now := time.Now()
	app.latestPose = appcore.CameraPoseConfig{
		Position: appcore.CameraVector3Config{X: 1},
		Rotation: appcore.CameraVector3Config{Z: 2},
	}
	app.poseAt = now
	app.userCameraSamples = map[string]userCameraOSCSample{
		"/usercamera/Zoom": {
			UserCameraOSCSample: appcore.UserCameraOSCSample{
				Address:  "/usercamera/Zoom",
				Float:    72.5,
				HasFloat: true,
			},
			ReceivedAt: now.Add(-10 * time.Minute),
		},
	}

	got, err := app.SaveCurrentCameraPoseToView("front")
	if err != nil {
		t.Fatal(err)
	}
	view, ok := findAutoCaptureViewByID(got.AutoCapture.Views, "front")
	if !ok {
		t.Fatal("front view missing")
	}
	if view.Zoom == nil || *view.Zoom != 45 {
		t.Fatalf("zoom = %v, want existing 45", view.Zoom)
	}
}

func TestAppRestartCameraPoseReceiverRestartsWhenForwardConfigChanges(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	logPath := appcore.DiagnosticLogPath(configPath)
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.ctx = context.Background()
	cfg := appcore.DefaultConfig()
	cfg.AutoCapture.OSC.ReceivePort = freeUDPPortForTest(t)
	cancelCalled := false
	app.oscCancel = func() {
		cancelCalled = true
	}
	app.oscReceiverHost = "127.0.0.1"
	app.oscReceiverPort = cfg.AutoCapture.OSC.ReceivePort
	app.oscReceiverLogPath = logPath
	app.oscReceiverForwardKey = oscForwardConfigKey(appcore.OSCForwardConfig{Enabled: false, Mode: appcore.OSCForwardModeAll})
	app.oscReceiverSeq = 7

	cfg.AutoCapture.OSC.Forward.Enabled = true
	cfg.AutoCapture.OSC.Forward.Targets = []appcore.OSCForwardTarget{{Host: "127.0.0.1", Port: 9101}}
	app.restartCameraPoseReceiverLocked(cfg)
	defer stopCameraPoseReceiverForTest(app)

	if !cancelCalled {
		t.Fatal("receiver was not canceled after forward config changed")
	}
	if app.oscReceiverForwardKey != oscForwardConfigKey(cfg.AutoCapture.OSC.Forward) {
		t.Fatal("forward config key was not updated")
	}
}

func stopCameraPoseReceiverForTest(app *App) {
	app.mu.Lock()
	done := app.oscReceiverDone
	app.stopCameraPoseReceiverLocked()
	app.mu.Unlock()
	if done == nil {
		return
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		panic("camera pose receiver did not stop")
	}
}

func TestAppAvatarOSCBasisAddressAffectsBasisOnlyForConfiguredAxes(t *testing.T) {
	app := NewApp(filepath.Join(t.TempDir(), "config.json"), appcore.UIState{Mode: appcore.ModeResults})
	cfg := appcore.DefaultConfig()

	for _, address := range []string{
		"/avatar/parameters/coord/x",
		"/avatar/parameters/coord/xSign",
		"/avatar/parameters/forward/z",
		"/avatar/parameters/forward/zSign",
	} {
		if !app.avatarOSCBasisAddressAffectsBasisLocked(address, cfg) {
			t.Fatalf("%s should affect AvatarBeacon basis", address)
		}
	}
	for _, address := range []string{
		"/avatar/parameters/GestureLeft",
		"/avatar/parameters/Viseme",
		"/avatar/parameters/coord/debug",
	} {
		if app.avatarOSCBasisAddressAffectsBasisLocked(address, cfg) {
			t.Fatalf("%s should not affect AvatarBeacon basis", address)
		}
	}
}

func TestAppFormatAvatarOSCBasisValuesOmitsDebugPing(t *testing.T) {
	app := NewApp(filepath.Join(t.TempDir(), "config.json"), appcore.UIState{Mode: appcore.ModeResults})
	cfg := appcore.DefaultConfig()
	got := app.formatAvatarOSCBasisValuesLocked(cfg)

	if strings.Contains(got, "avatar_beacon/debug/ping") {
		t.Fatalf("summary values include removed debug ping: %q", got)
	}
	for _, want := range []string{"coord/x=<missing>", "forward/zSign=<missing>"} {
		if !strings.Contains(got, want) {
			t.Fatalf("summary values = %q, want %q", got, want)
		}
	}
}

func TestLatestAvatarOSCSampleTimeUsesLatestBasisParameter(t *testing.T) {
	now := time.Date(2026, 7, 4, 6, 0, 0, 0, time.UTC)
	got := latestAvatarOSCSampleTime(map[string]avatarOSCBasisSample{
		"coord/x":     {Float: 1, HasFloat: true, ReceivedAt: now.Add(-30 * time.Second)},
		"coord/y":     {Float: 1, HasFloat: true, ReceivedAt: now.Add(-20 * time.Second)},
		"coord/z":     {Float: 1, HasFloat: true, ReceivedAt: now.Add(-10 * time.Second)},
		"forward/y":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-time.Second)},
		"Ahoge_Angle": {Float: 1, HasFloat: true, ReceivedAt: now},
	}, "coord", "forward", "Sign")
	if !got.Equal(now.Add(-time.Second)) {
		t.Fatalf("latestAvatarOSCSampleTime = %s, want %s", got, now.Add(-time.Second))
	}
}

func TestFormatOSCLogValuesIncludesUnsupportedAndRemainingPayload(t *testing.T) {
	payload := make([]byte, 0, 24)
	payload = binary.BigEndian.AppendUint32(payload, math.Float32bits(0.4))
	negativeTwo := int32(-2)
	payload = binary.BigEndian.AppendUint32(payload, uint32(negativeTwo))
	payload = append(payload, []byte("hello\x00\x00\x00")...)
	payload = append(payload, 0xaa, 0xbb)

	got := formatOSCLogValues(",fisTz", payload)
	for _, want := range []string{"0.4", "-2", `"hello"`, "true", "<unsupported:z>", "remaining_hex=aabb"} {
		if !strings.Contains(got, want) {
			t.Fatalf("formatOSCLogValues = %q, want %q", got, want)
		}
	}
}

func TestHexPreviewTruncatesLongPayload(t *testing.T) {
	got := hexPreview([]byte{0x00, 0x11, 0x22, 0x33}, 2)
	if got != "0011..." {
		t.Fatalf("hexPreview = %q", got)
	}
}

func TestAppPrepareAutoCaptureConfigKeepsFallbackOffAfterAvatarBeaconTimeout(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	now := time.Now()
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	app.state.Config.AutoCapture.Capture.PreplacedLocalAnchor = appBoolPtr(false)
	app.avatarOSCBasisSamples = map[string]avatarOSCBasisSample{
		"coord/x":       {Float: 0.8, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"coord/xSign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"coord/y":       {Float: 0.5, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"coord/ySign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"coord/z":       {Float: 0.2, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"coord/zSign":   {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"forward/x":     {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"forward/xSign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"forward/y":     {Float: 0, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"forward/ySign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"forward/z":     {Float: 0, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
		"forward/zSign": {Float: 1, HasFloat: true, ReceivedAt: now.Add(-31 * time.Second)},
	}

	app.mu.Lock()
	cfg, err := app.prepareAutoCaptureConfigForRunLocked(app.state.Config)
	app.mu.Unlock()
	if err == nil {
		t.Fatal("expected stale avatar basis error when fallback mode is off")
	}
	if cfg.AutoCapture.Capture.PreplacedLocalAnchorEnabled() {
		t.Fatal("fallback mode should remain off after avatar beacon timeout")
	}
}

func TestAppPrepareAutoCaptureConfigKeepsFallbackOnWhenAvatarBeaconReceived(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	app.state.Config.AutoCapture.Capture.PreplacedLocalAnchor = appBoolPtr(true)
	app.avatarOSCBasisSamples = readyAvatarOSCBasisSamples(time.Now().Add(-5 * time.Second))

	app.mu.Lock()
	cfg, err := app.prepareAutoCaptureConfigForRunLocked(app.state.Config)
	app.mu.Unlock()
	if err != nil {
		t.Fatalf("prepareAutoCaptureConfigForRunLocked error = %v", err)
	}
	if !cfg.AutoCapture.Capture.PreplacedLocalAnchorEnabled() {
		t.Fatal("fallback mode should remain on when AvatarBeacon basis was received")
	}
}

func TestAppMoveCameraToViewUsesAvatarOSCBasis(t *testing.T) {
	conn, port := listenAppOSCPackets(t)
	defer conn.Close()

	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	cfg := app.state.Config
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Capture.Mode = "photo"
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = appBoolPtr(false)
	cfg.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	cfg.AutoCapture.PlayerLocal.Calibrated = false
	cfg.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	view := appcore.CameraViewConfig{
		ID:              "local",
		Name:            "ローカル",
		CoordinateSpace: "player_local",
		Pose: appcore.CameraPoseConfig{
			Position: appcore.CameraVector3Config{X: 1, Y: 2, Z: 3},
			Rotation: appcore.CameraVector3Config{X: 4, Y: 5, Z: 6},
		},
	}
	cfg.AutoCapture.Views = []appcore.CameraViewConfig{view}
	app.state.Config = cfg
	app.avatarOSCBasisSamples = readyAvatarOSCBasisSamples(time.Now())

	snapshot := app.latestPlayerLocalBasisLocked(app.state.Config, time.Now())
	if snapshot.Status != "ready" || !snapshot.Fresh {
		t.Fatalf("snapshot = %+v, want ready", snapshot)
	}
	wantPose := appcore.TransformPlayerLocalPose(snapshot.Pose, view.Pose)

	if err := app.MoveCameraToView("local"); err != nil {
		t.Fatal(err)
	}

	packets := readAppOSCPackets(t, conn)
	var posePacket appOSCPacket
	for _, packet := range packets {
		if packet.Address == "/usercamera/Pose" {
			posePacket = packet
			break
		}
	}
	if posePacket.Address == "" {
		t.Fatalf("/usercamera/Pose was not sent: %+v", packets)
	}
	if len(posePacket.Floats) != 6 {
		t.Fatalf("pose floats = %+v, want 6 values", posePacket.Floats)
	}
	gotPose := appcore.CameraPoseConfig{
		Position: appcore.CameraVector3Config{X: float64(posePacket.Floats[0]), Y: float64(posePacket.Floats[1]), Z: float64(posePacket.Floats[2])},
		Rotation: appcore.CameraVector3Config{X: float64(posePacket.Floats[3]), Y: float64(posePacket.Floats[4]), Z: float64(posePacket.Floats[5])},
	}
	assertAppPoseClose(t, gotPose, wantPose)
}

func readyAvatarOSCBasisSamples(now time.Time) map[string]avatarOSCBasisSample {
	return map[string]avatarOSCBasisSample{
		"coord/x":       {Float: 0.8, HasFloat: true, ReceivedAt: now},
		"coord/xSign":   {Float: 1, HasFloat: true, ReceivedAt: now},
		"coord/y":       {Float: 0.5, HasFloat: true, ReceivedAt: now},
		"coord/ySign":   {Float: 1, HasFloat: true, ReceivedAt: now},
		"coord/z":       {Float: 0.2, HasFloat: true, ReceivedAt: now},
		"coord/zSign":   {Float: 1, HasFloat: true, ReceivedAt: now},
		"forward/x":     {Float: 1, HasFloat: true, ReceivedAt: now},
		"forward/xSign": {Float: 1, HasFloat: true, ReceivedAt: now},
		"forward/y":     {Float: 0, HasFloat: true, ReceivedAt: now},
		"forward/ySign": {Float: 1, HasFloat: true, ReceivedAt: now},
		"forward/z":     {Float: 0, HasFloat: true, ReceivedAt: now},
		"forward/zSign": {Float: 1, HasFloat: true, ReceivedAt: now},
	}
}

func listenAppOSCPackets(t *testing.T) (net.PacketConn, int) {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	return conn, addr.Port
}

type appOSCPacket struct {
	Address  string
	TypeTags string
	Ints     []int32
	Floats   []float32
}

func readAppOSCPackets(t *testing.T, conn net.PacketConn) []appOSCPacket {
	t.Helper()
	packets := make([]appOSCPacket, 0, 8)
	buf := make([]byte, 2048)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond)); err != nil {
			t.Fatal(err)
		}
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			t.Fatal(err)
		}
		address, typeTags, payload, ok := appcore.ParseOSCPacket(append([]byte(nil), buf[:n]...))
		if !ok {
			t.Fatalf("failed to parse OSC packet: %v", buf[:n])
		}
		packet := appOSCPacket{Address: address, TypeTags: typeTags}
		payloadOffset := 0
		for _, tag := range strings.TrimPrefix(typeTags, ",") {
			switch tag {
			case 'i':
				if payloadOffset+4 > len(payload) {
					t.Fatalf("OSC int payload too short: %v", payload)
				}
				packet.Ints = append(packet.Ints, int32(binary.BigEndian.Uint32(payload[payloadOffset:payloadOffset+4])))
				payloadOffset += 4
			case 'f':
				if payloadOffset+4 > len(payload) {
					t.Fatalf("OSC float payload too short: %v", payload)
				}
				packet.Floats = append(packet.Floats, math.Float32frombits(binary.BigEndian.Uint32(payload[payloadOffset:payloadOffset+4])))
				payloadOffset += 4
			case 'T', 'F':
			default:
				t.Fatalf("unsupported OSC type tag %q in %q", tag, typeTags)
			}
		}
		packets = append(packets, packet)
	}
	return packets
}

func assertAppPoseClose(t *testing.T, got appcore.CameraPoseConfig, want appcore.CameraPoseConfig) {
	t.Helper()
	const epsilon = 1e-4
	values := []struct {
		name string
		got  float64
		want float64
	}{
		{"pos.x", got.Position.X, want.Position.X},
		{"pos.y", got.Position.Y, want.Position.Y},
		{"pos.z", got.Position.Z, want.Position.Z},
		{"rot.x", got.Rotation.X, want.Rotation.X},
		{"rot.y", got.Rotation.Y, want.Rotation.Y},
		{"rot.z", got.Rotation.Z, want.Rotation.Z},
	}
	for _, value := range values {
		if math.Abs(value.got-value.want) > epsilon {
			t.Fatalf("%s = %v, want %v; got pose %+v want %+v", value.name, value.got, value.want, got, want)
		}
	}
}

func TestAppWaitForAutoCaptureStartReadinessWaitsForAvatarOSCBasis(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	app.state.Config.AutoCapture.Capture.PreplacedLocalAnchor = appBoolPtr(false)
	app.state.Config.AutoCapture.Presence.WatchOutputLog = false
	app.avatarOSCBasisSamples = map[string]avatarOSCBasisSample{
		"coord/x": {Float: 0.8, HasFloat: true, ReceivedAt: time.Now()},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		time.Sleep(30 * time.Millisecond)
		app.mu.Lock()
		app.avatarOSCBasisSamples = readyAvatarOSCBasisSamples(time.Now())
		app.mu.Unlock()
	}()
	start := time.Now()
	if err := app.waitForAutoCaptureStartReadiness(ctx, app.state.Config, 500*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed < 30*time.Millisecond {
		t.Fatalf("wait returned too early: %s", elapsed)
	}
}

func TestAppWaitForAutoCaptureStartReadinessFallsBackWithoutAvatarOSCBasis(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	ctx := context.Background()
	start := time.Now()
	err := app.waitForAutoCaptureStartReadiness(ctx, app.state.Config, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("readiness wait error = %v", err)
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("fallback readiness took too long: %s", elapsed)
	}
}

func TestAppWaitForAutoCaptureStartReadinessTimesOutWithoutWorldMetadata(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	app.state.Config.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC
	app.state.Config.AutoCapture.PlayerLocal.AvatarOSC.FreshnessSec = 3
	app.state.Config.AutoCapture.Presence.WatchOutputLog = true
	app.state.Config.AutoCapture.Presence.OutputLogDirectory = t.TempDir()
	app.avatarOSCBasisSamples = readyAvatarOSCBasisSamples(time.Now())

	err := app.waitForAutoCaptureStartReadiness(context.Background(), app.state.Config, 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected readiness wait to time out without world metadata")
	}
	if !strings.Contains(err.Error(), "world情報") {
		t.Fatalf("err = %v, want world metadata timeout", err)
	}
}

func TestCreateEncryptedDiagnosticPackageEncryptsZip(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	userProfile := filepath.Join(dir, "Users", "alice")
	t.Setenv("USERPROFILE", userProfile)
	cfg := appcore.DefaultConfig()
	cfg.Discord.WebhookURL = "https://discord.com/api/webhooks/secret"
	cfg.AutoPhoto.PhotoDirectory = filepath.Join(userProfile, "Pictures", "VRChat")
	cfg.Image.OutputDirectory = filepath.Join(userProfile, "Pictures", "ClipForVRChatOutput")
	if err := appcore.SaveConfig(configPath, cfg); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(cfg.Image.OutputDirectory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.Image.OutputDirectory, "result.png"), []byte("real output image"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), []appcore.HistoryEntry{{ID: "1", URL: "https://example.com", SourcePath: filepath.Join(userProfile, "Pictures", "VRChat", "photo.png")}}); err != nil {
		t.Fatal(err)
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "test log path=%q", filepath.Join(userProfile, "Pictures", "VRChat", "photo.png"))

	app := NewApp(configPath, appcore.UIState{Config: cfg})
	path, err := app.CreateEncryptedDiagnosticPackage()
	if err != nil {
		t.Fatal(err)
	}
	workDir := filepath.Dir(path)
	if filepath.Base(filepath.Dir(workDir)) != "diagnostics" || filepath.Ext(path) != ".gpg" {
		t.Fatalf("path = %q, want encrypted package under diagnostics timestamp dir", path)
	}
	zipPath := strings.TrimSuffix(path, ".gpg")
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("plain zip should remain for user review: %v", err)
	}
	dataDir := filepath.Join(workDir, "data")
	if _, err := os.Stat(dataDir); !os.IsNotExist(err) {
		t.Fatalf("data dir should be removed after zip creation, err=%v", err)
	}
	plainZip, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer plainZip.Close()
	logData, err := readZipFile(t, plainZip.File, "logs/"+filepath.Base(appcore.DiagnosticLogPath(configPath)))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(logData), userProfile) || !strings.Contains(string(logData), `%USERPROFILE%`) {
		t.Fatalf("redacted log = %q, want USERPROFILE placeholder and no raw profile path", string(logData))
	}
	if !strings.Contains(string(logData), "diagnostic output directory=") || !strings.Contains(string(logData), "result.png") {
		t.Fatalf("redacted log = %q, want output directory listing", string(logData))
	}
	configData, err := readZipFile(t, plainZip.File, "config.json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(configData), userProfile) || !strings.Contains(string(configData), `%USERPROFILE%`) {
		t.Fatalf("redacted config = %q, want USERPROFILE placeholder and no raw profile path", string(configData))
	}
	historyData, err := readZipFile(t, plainZip.File, "history.json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(historyData), userProfile) || !strings.Contains(string(historyData), `%USERPROFILE%`) {
		t.Fatalf("redacted history = %q, want USERPROFILE placeholder and no raw profile path", string(historyData))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, leaked := range []string{"config.json", "history.json", "manifest.json", "secret"} {
		if strings.Contains(string(data), leaked) {
			t.Fatalf("encrypted package leaked %q", leaked)
		}
	}
}

func TestEncryptZipFileWithPublicKey(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "sample.zip")
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	writer, err := zipWriter.Create("inside.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("secret text")); err != nil {
		t.Fatal(err)
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zipPath, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}

	outputPath, err := encryptZipFileWithPublicKey(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	if outputPath != zipPath+".gpg" {
		t.Fatalf("outputPath = %q, want %q", outputPath, zipPath+".gpg")
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, leaked := range []string{"inside.txt", "secret text"} {
		if strings.Contains(string(data), leaked) {
			t.Fatalf("encrypted zip leaked %q", leaked)
		}
	}
}

func TestDiagnosticZipDoesNotIncludeOutputImages(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Image.OutputDirectory = "./output"
	if err := appcore.SaveConfig(configPath, cfg); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "real.png"), []byte("not the real image in diagnostics"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "real.jpg"), []byte("not the real jpeg in diagnostics"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), []appcore.HistoryEntry{{ID: "1", OutputPath: "output/history.png"}}); err != nil {
		t.Fatal(err)
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "test log")
	zipData, _, err := buildDiagnosticZip(configPath, cfg)
	if err != nil {
		t.Fatal(err)
	}
	plainZip, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		t.Fatalf("plain diagnostic zip is invalid: %v", err)
	}
	entries := zipEntryNames(plainZip.File)
	for _, want := range []string{"logs/" + filepath.Base(appcore.DiagnosticLogPath(configPath))} {
		if !entries[want] {
			t.Fatalf("zip entries = %+v, missing %s", entries, want)
		}
	}
	if entries["log/"+filepath.Base(appcore.DiagnosticLogPath(configPath))] {
		t.Fatalf("zip entries = %+v, should use logs/ not log/", entries)
	}
	for _, unwanted := range []string{"output/history.png", "output/real.jpg", "output/real.png"} {
		if entries[unwanted] {
			t.Fatalf("zip entries = %+v, should not include output images", entries)
		}
	}

	entity, err := openpgp.NewEntity("Diagnostic Test", "", "diagnostic@example.test", &packet.Config{
		Time: func() time.Time { return time.Now().Add(-time.Hour) },
	})
	if err != nil {
		t.Fatal(err)
	}
	encrypted, err := encryptDiagnosticZip(zipData, openpgp.EntityList{entity})
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(encrypted), "config.json") {
		t.Fatal("encrypted package leaked zip entry name")
	}
	message, err := openpgp.ReadMessage(bytes.NewReader(encrypted), openpgp.EntityList{entity}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if message.LiteralData == nil || !message.LiteralData.IsBinary {
		t.Fatalf("literal data = %+v, want binary", message.LiteralData)
	}
	decrypted, err := io.ReadAll(message.UnverifiedBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := zip.NewReader(bytes.NewReader(decrypted), int64(len(decrypted))); err != nil {
		t.Fatalf("decrypted diagnostic zip is invalid: %v", err)
	}
}

func zipEntryNames(files []*zip.File) map[string]bool {
	names := make(map[string]bool, len(files))
	for _, file := range files {
		names[file.Name] = true
	}
	return names
}

func readZipFile(t *testing.T, files []*zip.File, name string) ([]byte, error) {
	t.Helper()
	for _, file := range files {
		if file.Name != name {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	return nil, fmt.Errorf("zip entry %s not found", name)
}

func TestAppProcessToStateRejectsMixedJSON(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	state, err := app.ProcessToState([]string{"config.json", "image.png"})
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != appcore.ModeError || !strings.Contains(state.Message, "混在") {
		t.Fatalf("state = %+v", state)
	}
}

func TestAppHistoryOperations(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{
		{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png", DiscordDeleted: true},
		{ID: "2", URL: ""},
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	got, err := app.GetHistory()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len(history) = %d, want 2", len(got))
	}

	got, err = app.MarkHistoryCleared([]string{"2"})
	if err != nil {
		t.Fatal(err)
	}
	if !got[1].Cleared {
		t.Fatalf("entry 2 should be cleared: %+v", got)
	}

	got, err = app.PurgeDeletedHistoryEntries()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "2" {
		t.Fatalf("purged history = %+v", got)
	}
}

func TestAppDeleteDiscordHistoryEntriesSkipsEntriesWithoutStoredWebhookData(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	got, err := app.DeleteDiscordHistoryEntries([]string{"1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].DiscordDeleted {
		t.Fatalf("history = %+v, want unchanged", got)
	}
}

func TestAppDeleteDiscordHistoryEntriesPersistsPartialSuccess(t *testing.T) {
	oldDelete := deleteDiscordMessage
	t.Cleanup(func() {
		deleteDiscordMessage = oldDelete
	})

	calls := 0
	deleteDiscordMessage = func(webhookID, token, messageID string) error {
		calls++
		if messageID == "m2" {
			return errors.New("delete failed")
		}
		return nil
	}

	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{
		{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png", DiscordWebhookID: "w", DiscordToken: "t", DiscordMessageID: "m1"},
		{ID: "2", URL: "https://cdn.discordapp.com/attachments/1/2/b.png", DiscordWebhookID: "w", DiscordToken: "t", DiscordMessageID: "m2"},
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	if _, err := app.DeleteDiscordHistoryEntries([]string{"1", "2"}); err == nil {
		t.Fatal("expected partial delete error")
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	loaded, err := appcore.LoadHistory(appcore.HistoryPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if !loaded[0].DiscordDeleted || loaded[0].DeletedAt == "" {
		t.Fatalf("first entry should be persisted as deleted: %+v", loaded[0])
	}
	if loaded[1].DiscordDeleted {
		t.Fatalf("second entry should not be deleted: %+v", loaded[1])
	}
}

func TestResultMessage(t *testing.T) {
	cfg := appcore.DefaultConfig()
	cfg.Output.CopySingleURLToClipboard = true
	msg := resultMessage(cfg, []appcore.Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}, nil)
	if !strings.Contains(msg, "コピー") {
		t.Fatalf("message = %q", msg)
	}

	msg = resultMessage(cfg, []appcore.Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}, errors.New("clipboard busy"))
	if !strings.Contains(msg, "コピーできません") {
		t.Fatalf("copy error message = %q", msg)
	}

	msg = resultMessage(cfg, []appcore.Result{{OutputPath: "out.png"}}, nil)
	if !strings.Contains(msg, "保存") {
		t.Fatalf("message = %q", msg)
	}

	msg = resultMessage(cfg, []appcore.Result{{QRURLs: []string{"https://example.com/qr"}}}, nil)
	if !strings.Contains(msg, "QRコードURL") {
		t.Fatalf("message = %q", msg)
	}

	msg = resultMessage(cfg, nil, nil)
	if !strings.Contains(msg, "表示できる処理結果はありません") {
		t.Fatalf("message = %q", msg)
	}

	cfg.Output.UploadDiscord = false
	cfg.Output.SaveLocal = false
	cfg.Output.DetectQRCodeURLs = false
	msg = resultMessage(cfg, nil, nil)
	if !strings.Contains(msg, "すべてOFF") || !strings.Contains(msg, "設定を確認") {
		t.Fatalf("all disabled message = %q", msg)
	}

	cfg.Output.DetectQRCodeURLs = true
	msg = resultMessage(cfg, nil, nil)
	if !strings.Contains(msg, "URLを取得できません") || !strings.Contains(msg, "Discord投稿またはローカル保存") {
		t.Fatalf("qr-only no result message = %q", msg)
	}
}

func TestCheckFFmpegReportsMissingPath(t *testing.T) {
	app := NewApp(filepath.Join(t.TempDir(), "config.json"), appcore.UIState{})
	status := app.CheckFFmpeg(filepath.Join(t.TempDir(), "missing-ffmpeg.exe"))
	if status.Available {
		t.Fatalf("status = %+v, want unavailable", status)
	}
	if !strings.Contains(status.Message, "ffmpegがインストールされていないかPATHにありません") {
		t.Fatalf("message = %q", status.Message)
	}
}

func writeTextFile(t *testing.T, path string, text string) {
	t.Helper()
	if err := appcore.WritePrivateFile(path, []byte(text)); err != nil {
		t.Fatal(err)
	}
}

func writeTestPNG(t *testing.T, path string, width int, height int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 80, G: 120, B: 180, A: 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
}
