package appcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigNormalizeAppliesDefaultsAndTrimsQuotes(t *testing.T) {
	cfg := Config{
		Image: ImageConfig{
			MaxWidth:        -1,
			MaxHeight:       0,
			MaxInputMB:      0,
			Suffix:          "",
			OutputFormat:    "gif",
			JPEGQuality:     500,
			OutputDirectory: ` "./quoted-output" `,
		},
		Output: OutputConfig{ShowUI: "sometimes"},
		AutoPhoto: AutoPhotoConfig{
			PhotoDirectory: ` "C:\VRChat Photos" `,
			WebhookURL:     ` "https://discord.com/api/webhooks/1/token" `,
		},
		ScreenshotAutoPost: ScreenshotAutoPostConfig{
			ScreenshotDirectory: ` "C:\Users\test\Pictures\Screenshots" `,
			WebhookURL:          ` "https://discord.com/api/webhooks/2/screenshot" `,
		},
	}

	cfg.Normalize()

	if cfg.Image.MaxWidth != 2048 || cfg.Image.MaxHeight != 2048 {
		t.Fatalf("unexpected max size: %dx%d", cfg.Image.MaxWidth, cfg.Image.MaxHeight)
	}
	if cfg.Image.MaxInputMB != DefaultMaxImageInputMB {
		t.Fatalf("MaxInputMB = %d, want %d", cfg.Image.MaxInputMB, DefaultMaxImageInputMB)
	}
	if cfg.Image.Suffix != "_2048" {
		t.Fatalf("Suffix = %q, want _2048", cfg.Image.Suffix)
	}
	if cfg.Image.OutputFormat != "png" {
		t.Fatalf("OutputFormat = %q, want png", cfg.Image.OutputFormat)
	}
	if cfg.Image.JPEGQuality != 92 {
		t.Fatalf("JPEGQuality = %d, want 92", cfg.Image.JPEGQuality)
	}
	if cfg.Image.OutputDirectory != "./quoted-output" {
		t.Fatalf("OutputDirectory = %q", cfg.Image.OutputDirectory)
	}
	if cfg.Output.ShowUI != "auto" {
		t.Fatalf("ShowUI = %q, want auto", cfg.Output.ShowUI)
	}
	if cfg.AutoPhoto.PhotoDirectory != `C:\VRChat Photos` {
		t.Fatalf("PhotoDirectory = %q", cfg.AutoPhoto.PhotoDirectory)
	}
	if cfg.AutoPhoto.WebhookURL != "https://discord.com/api/webhooks/1/token" {
		t.Fatalf("WebhookURL = %q", cfg.AutoPhoto.WebhookURL)
	}
	if cfg.AutoPhoto.ScanIntervalSeconds != 2 {
		t.Fatalf("ScanIntervalSeconds = %d, want 2", cfg.AutoPhoto.ScanIntervalSeconds)
	}
	if cfg.AutoCapture.PlayerLocal.BasisSource != PlayerLocalBasisSourceAvatarOSC {
		t.Fatalf("PlayerLocal.BasisSource = %q, want avatar_osc", cfg.AutoCapture.PlayerLocal.BasisSource)
	}
	if cfg.AutoCapture.PlayerLocal.AvatarOSC.ParameterPrefix != "coord" {
		t.Fatalf("PlayerLocal.AvatarOSC.ParameterPrefix = %q", cfg.AutoCapture.PlayerLocal.AvatarOSC.ParameterPrefix)
	}
	if cfg.AutoCapture.PlayerLocal.AvatarOSC.PositionScale != 1000 {
		t.Fatalf("PlayerLocal.AvatarOSC.PositionScale = %v, want 1000", cfg.AutoCapture.PlayerLocal.AvatarOSC.PositionScale)
	}
	if cfg.ScreenshotAutoPost.ScreenshotDirectory != `C:\Users\test\Pictures\Screenshots` {
		t.Fatalf("ScreenshotDirectory = %q", cfg.ScreenshotAutoPost.ScreenshotDirectory)
	}
	if cfg.ScreenshotAutoPost.WebhookURL != "https://discord.com/api/webhooks/2/screenshot" {
		t.Fatalf("Screenshot WebhookURL = %q", cfg.ScreenshotAutoPost.WebhookURL)
	}
	if cfg.ScreenshotAutoPost.ScanIntervalSeconds != 2 {
		t.Fatalf("Screenshot ScanIntervalSeconds = %d, want 2", cfg.ScreenshotAutoPost.ScanIntervalSeconds)
	}
}

func TestLoadConfigReturnsDefaultWithoutCreatingFileWhenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Image.MaxWidth != 2048 {
		t.Fatalf("MaxWidth = %d, want 2048", cfg.Image.MaxWidth)
	}
	if ConfigExists(path) {
		t.Fatal("LoadConfig should not create config file")
	}
}

func TestSaveAndLoadConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := DefaultConfig()
	want.Image.MaxWidth = 1024
	want.Image.OutputFormat = "jpg"
	want.Image.MaxInputMB = 64
	want.Output.ShowUI = "always"
	want.AutoCapture.Stream.SpoutHelperPath = `C:\tools\spout-capture.exe`

	if err := SaveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Image.MaxWidth != want.Image.MaxWidth ||
		got.Image.OutputFormat != want.Image.OutputFormat ||
		got.Image.MaxInputMB != want.Image.MaxInputMB ||
		got.Output.ShowUI != want.Output.ShowUI ||
		got.Update.CheckEnabled != want.Update.CheckEnabled ||
		got.Update.NotificationEnabled != want.Update.NotificationEnabled ||
		got.AutoCapture.Stream.SpoutHelperPath != want.AutoCapture.Stream.SpoutHelperPath {
		t.Fatalf("loaded config mismatch: %+v", got)
	}
}

func TestConfigNormalizeTrimsSpoutHelperPathWithoutResettingCustomPath(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.Stream.SpoutHelperPath = ` "C:\tools\custom-spout-capture.exe" `
	cfg.Normalize()
	if cfg.AutoCapture.Stream.SpoutHelperPath != `C:\tools\custom-spout-capture.exe` {
		t.Fatalf("SpoutHelperPath = %q", cfg.AutoCapture.Stream.SpoutHelperPath)
	}
}

func TestOSCForwardConfigNormalize(t *testing.T) {
	cfg := OSCForwardConfig{
		Enabled: true,
		Mode:    "unknown",
		Targets: []OSCForwardTarget{
			{Host: ` " " `, Port: 9101},
			{Host: ` "127.0.0.2" `, Port: 9102},
		},
	}
	cfg.Normalize()
	if cfg.Mode != OSCForwardModeAll {
		t.Fatalf("Mode = %q, want %q", cfg.Mode, OSCForwardModeAll)
	}
	if cfg.Targets[0].Host != "127.0.0.1" {
		t.Fatalf("Targets[0].Host = %q, want default loopback", cfg.Targets[0].Host)
	}
	if cfg.Targets[1].Host != "127.0.0.2" {
		t.Fatalf("Targets[1].Host = %q", cfg.Targets[1].Host)
	}
}

func TestAutoCapturePlayerLocalNormalizeDefaultsAndSourceFallback(t *testing.T) {
	cfg := AutoCapturePlayerLocalConfig{
		AvatarOSC: AutoCapturePlayerLocalAvatarOSCConfig{
			ParameterPrefix:       ` "  " `,
			PositionScale:         0,
			InvertMagnitude:       true,
			PositiveFlagThreshold: 0.75,
			MaxAbsPosition:        -1,
			MaxAbsForward:         0,
			FreshnessSec:          0,
		},
	}
	cfg.Normalize()
	if cfg.BasisSource != PlayerLocalBasisSourceAvatarOSC {
		t.Fatalf("BasisSource = %q, want avatar_osc", cfg.BasisSource)
	}

	cfg.BasisSource = "unknown"
	cfg.Normalize()
	if cfg.BasisSource != PlayerLocalBasisSourceManual {
		t.Fatalf("BasisSource = %q, want manual", cfg.BasisSource)
	}
	if cfg.AvatarOSC.ParameterPrefix != "coord" {
		t.Fatalf("ParameterPrefix = %q", cfg.AvatarOSC.ParameterPrefix)
	}
	if cfg.AvatarOSC.PositionScale != 1000 {
		t.Fatalf("PositionScale = %v, want 1000", cfg.AvatarOSC.PositionScale)
	}
	if cfg.AvatarOSC.PositiveFlagThreshold != 0.75 {
		t.Fatalf("PositiveFlagThreshold = %v, want 0.75", cfg.AvatarOSC.PositiveFlagThreshold)
	}
	if cfg.AvatarOSC.MaxAbsPosition != 10000 || cfg.AvatarOSC.MaxAbsForward != 2000 {
		t.Fatalf("unexpected max abs limits: %+v", cfg.AvatarOSC)
	}
	if cfg.AvatarOSC.FreshnessSec != 3 {
		t.Fatalf("FreshnessSec = %d, want 3", cfg.AvatarOSC.FreshnessSec)
	}
}

func TestLoadConfigDefaultsUpdateSettingsWhenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"image":{"maxWidth":1024}}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Update.CheckEnabled || !got.Update.NotificationEnabled {
		t.Fatalf("update defaults = %+v, want both enabled", got.Update)
	}
}

func TestConfigPathUsesExecutableDirectory(t *testing.T) {
	got := ConfigPath(filepath.Join("C:", "tools", "ClipForVRChat.exe"))
	if filepath.Base(got) != "config.json" {
		t.Fatalf("ConfigPath base = %q, want config.json", filepath.Base(got))
	}
}

func TestDefaultVRChatPhotoDirectoryUsesUserProfile(t *testing.T) {
	t.Setenv("USERPROFILE", filepath.Join("C:", "Users", "test"))
	got := DefaultVRChatPhotoDirectory()
	want := filepath.Join("C:", "Users", "test", "Pictures", "VRChat")
	if got != want {
		t.Fatalf("DefaultVRChatPhotoDirectory = %q, want %q", got, want)
	}
}

func TestDefaultScreenshotsDirectoryUsesUserProfile(t *testing.T) {
	t.Setenv("USERPROFILE", filepath.Join("C:", "Users", "test"))
	got := DefaultScreenshotsDirectory()
	want := filepath.Join("C:", "Users", "test", "Pictures", "Screenshots")
	if got != want {
		t.Fatalf("DefaultScreenshotsDirectory = %q, want %q", got, want)
	}
}

func TestWritePrivateFileCreatesReadableFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "private.txt")
	if err := WritePrivateFile(path, []byte("secret")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "secret" {
		t.Fatalf("data = %q, want secret", string(data))
	}
}
