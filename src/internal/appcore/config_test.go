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
	if cfg.AutoCapture.PlayerLocal.AvatarOSC.ParameterPrefix != "avatar_beacon" {
		t.Fatalf("PlayerLocal.AvatarOSC.ParameterPrefix = %q", cfg.AutoCapture.PlayerLocal.AvatarOSC.ParameterPrefix)
	}
	if cfg.AutoCapture.PlayerLocal.AvatarOSC.PositionScale != 1000 {
		t.Fatalf("PlayerLocal.AvatarOSC.PositionScale = %v, want 1000", cfg.AutoCapture.PlayerLocal.AvatarOSC.PositionScale)
	}
	if cfg.AutoCapture.Capture.AutoEnablePreplacedAfterMinutes != 5 {
		t.Fatalf("AutoEnablePreplacedAfterMinutes = %d, want 5", cfg.AutoCapture.Capture.AutoEnablePreplacedAfterMinutes)
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

func TestCameraViewNormalizeRoundsPoseAndZoom(t *testing.T) {
	zoom := 44.44449
	exposure := -1.23449
	view := CameraViewConfig{
		ID:              "view",
		Name:            "View",
		CoordinateSpace: "player_local",
		Pose: CameraPoseConfig{
			Position: CameraVector3Config{X: 1.23449, Y: -2.3455, Z: 3.0004},
			Rotation: CameraVector3Config{X: 10.12349, Y: 20.5555, Z: -30.0004},
		},
		Zoom:     &zoom,
		Exposure: &exposure,
	}
	view.Normalize(0)
	if view.Pose.Position.X != 1.234 || view.Pose.Position.Y != -2.346 || view.Pose.Position.Z != 3 {
		t.Fatalf("position = %+v", view.Pose.Position)
	}
	if view.Pose.Rotation.X != 10.123 || view.Pose.Rotation.Y != 20.556 || view.Pose.Rotation.Z != -30 {
		t.Fatalf("rotation = %+v", view.Pose.Rotation)
	}
	if view.Zoom == nil || *view.Zoom != 44.444 {
		t.Fatalf("zoom = %v, want 44.444", view.Zoom)
	}
	if view.Exposure == nil || *view.Exposure != -1.234 {
		t.Fatalf("exposure = %v, want -1.234", view.Exposure)
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

func TestAutoCaptureConfigNormalizePreservesExplicitRemotePlayerFalse(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.Views = []CameraViewConfig{
		{
			ID:              "front",
			Enabled:         true,
			CoordinateSpace: "player_local",
			RemotePlayer:    boolConfigPtr(false),
		},
		{
			ID:              "back",
			Enabled:         true,
			CoordinateSpace: "player_local",
		},
	}

	cfg.Normalize()

	if cfg.AutoCapture.Views[0].RemotePlayer == nil || *cfg.AutoCapture.Views[0].RemotePlayer {
		t.Fatalf("front RemotePlayer = %v, want explicit false preserved", cfg.AutoCapture.Views[0].RemotePlayer)
	}
	if cfg.AutoCapture.Views[1].RemotePlayer == nil || !*cfg.AutoCapture.Views[1].RemotePlayer {
		t.Fatalf("back RemotePlayer = %v, want default true", cfg.AutoCapture.Views[1].RemotePlayer)
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
	if cfg.AvatarOSC.ParameterPrefix != "avatar_beacon" {
		t.Fatalf("ParameterPrefix = %q", cfg.AvatarOSC.ParameterPrefix)
	}
	for _, legacyPrefix := range []string{"coord", "/coord/", "forward", "avatar_beacon/coord", "avatar_beacon/forward"} {
		cfg.AvatarOSC.ParameterPrefix = legacyPrefix
		cfg.Normalize()
		if cfg.AvatarOSC.ParameterPrefix != "avatar_beacon" {
			t.Fatalf("legacy ParameterPrefix %q normalized to %q, want avatar_beacon", legacyPrefix, cfg.AvatarOSC.ParameterPrefix)
		}
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

func TestLoadConfigDefaultsCaptureOnStartWhenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"autoCapture":{"schedule":{"enabled":true}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.AutoCapture.Schedule.CaptureOnStart {
		t.Fatalf("CaptureOnStart = %t, want true when field is missing", got.AutoCapture.Schedule.CaptureOnStart)
	}
}

func TestLoadConfigPreservesExplicitCaptureOnStartFalse(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"autoCapture":{"schedule":{"enabled":true,"captureOnStart":false}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.AutoCapture.Schedule.CaptureOnStart {
		t.Fatalf("CaptureOnStart = %t, want false when field is explicit", got.AutoCapture.Schedule.CaptureOnStart)
	}
}

func TestLoadConfigDefaultsAutoLevelRollBeforeShotWhenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"autoCapture":{"capture":{"mode":"stream","closeCameraAfterBatch":true}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.AutoCapture.Capture.AutoLevelRollBeforeShot == nil || !*got.AutoCapture.Capture.AutoLevelRollBeforeShot {
		t.Fatalf("AutoLevelRollBeforeShot = %v, want true when field is missing", got.AutoCapture.Capture.AutoLevelRollBeforeShot)
	}
	if got.AutoCapture.Capture.CloseCameraAfterBatch {
		t.Fatalf("CloseCameraAfterBatch = %t, want false after normalize", got.AutoCapture.Capture.CloseCameraAfterBatch)
	}
	if !got.AutoCapture.Restore.Fallback.AutoLevelRoll {
		t.Fatalf("Restore fallback AutoLevelRoll = %t, want true when field is missing", got.AutoCapture.Restore.Fallback.AutoLevelRoll)
	}
}

func TestLoadConfigPreservesExplicitAutoLevelRollBeforeShotFalse(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"autoCapture":{"capture":{"autoLevelRollBeforeShot":false},"restore":{"fallback":{"autoLevelRoll":false}}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.AutoCapture.Capture.AutoLevelRollBeforeShot == nil || *got.AutoCapture.Capture.AutoLevelRollBeforeShot {
		t.Fatalf("AutoLevelRollBeforeShot = %v, want false when field is explicit", got.AutoCapture.Capture.AutoLevelRollBeforeShot)
	}
	if got.AutoCapture.Restore.Fallback.AutoLevelRoll {
		t.Fatalf("Restore fallback AutoLevelRoll = %t, want false when field is explicit", got.AutoCapture.Restore.Fallback.AutoLevelRoll)
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
