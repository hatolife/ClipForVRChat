package appcore

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultAutoCaptureConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.AutoCapture.Schedule.Enabled {
		t.Fatal("auto capture should be disabled by default")
	}
	if cfg.AutoCapture.OSC.Host != "127.0.0.1" || cfg.AutoCapture.OSC.SendPort != 9000 || cfg.AutoCapture.OSC.ReceivePort != 9001 {
		t.Fatalf("unexpected osc defaults: %+v", cfg.AutoCapture.OSC)
	}
	if len(cfg.AutoCapture.Views) != 3 {
		t.Fatalf("default views = %d, want 3", len(cfg.AutoCapture.Views))
	}
	if cfg.AutoCapture.Capture.Mode != "stream" || cfg.AutoCapture.Stream.SpoutHelperPath == "" || !cfg.AutoCapture.Stream.SpoutAutoSelect {
		t.Fatalf("unexpected stream defaults: capture=%+v stream=%+v", cfg.AutoCapture.Capture, cfg.AutoCapture.Stream)
	}
	if cfg.AutoCapture.Views[0].ID != "front" || cfg.AutoCapture.Views[0].Calibrated || cfg.AutoCapture.Views[0].Zoom == nil {
		t.Fatalf("unexpected first view: %+v", cfg.AutoCapture.Views[0])
	}
	if cfg.AutoCapture.Views[0].CoordinateSpace != "world" || cfg.AutoCapture.Views[0].Pose.Position.Z == 0 {
		t.Fatalf("default front view pose was not initialized: %+v", cfg.AutoCapture.Views[0])
	}
}

func TestAutoCaptureConfigNormalize(t *testing.T) {
	cfg := Config{AutoCapture: AutoCaptureConfig{
		OSC:      AutoCaptureOSCConfig{SendPort: -1, ReceivePort: 70000},
		Schedule: AutoCaptureScheduleConfig{CaptureIntervalSec: 1, InitialDelaySec: -1, MaxBatches: -1},
		Capture:  AutoCaptureCaptureConfig{Mode: "bad", ConcurrentMode: "bad", RequestedCameraCount: 10},
		Output:   AutoCaptureOutputConfig{ImageFormat: "gif"},
		Discord:  AutoCaptureDiscordConfig{PostMode: "bad"},
	}}
	cfg.Normalize()
	if cfg.AutoCapture.Schedule.CaptureIntervalSec != 10 {
		t.Fatalf("CaptureIntervalSec = %d, want 10", cfg.AutoCapture.Schedule.CaptureIntervalSec)
	}
	if cfg.AutoCapture.Capture.Mode != "stream" || cfg.AutoCapture.Capture.ConcurrentMode != "sequential" {
		t.Fatalf("capture normalize failed: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Capture.RequestedCameraCount != 4 {
		t.Fatalf("RequestedCameraCount = %d, want 4", cfg.AutoCapture.Capture.RequestedCameraCount)
	}
	if len(cfg.AutoCapture.Views) != 3 {
		t.Fatalf("default views = %d, want 3", len(cfg.AutoCapture.Views))
	}
	if cfg.AutoCapture.Stream.SpoutHelperPath != "spout-capture.exe" || !cfg.AutoCapture.Stream.SpoutAutoSelect || cfg.AutoCapture.Stream.CaptureTimeoutMS != 10000 {
		t.Fatalf("stream normalize failed: %+v", cfg.AutoCapture.Stream)
	}
}

func TestAutoCaptureConfigNormalizeMigratesOldDesktopFFmpegInputToLegacy(t *testing.T) {
	for _, oldArgs := range []string{oldDesktopFFmpegInputArgs, oldTitleFFmpegInputArgs} {
		cfg := DefaultConfig()
		cfg.AutoCapture.Stream.LegacyInputArgs = oldArgs
		cfg.Normalize()
		if cfg.AutoCapture.Stream.LegacyInputArgs != DefaultAutoCaptureFFmpegInputArgs() {
			t.Fatalf("legacy stream input args = %q, want %q", cfg.AutoCapture.Stream.LegacyInputArgs, DefaultAutoCaptureFFmpegInputArgs())
		}
	}
}

func TestSplitCommandLine(t *testing.T) {
	got, err := splitCommandLine(`-f gdigrab -i "title=VRChat Window" -frames:v 1`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-f", "gdigrab", "-i", "title=VRChat Window", "-frames:v", "1"}
	if len(got) != len(want) {
		t.Fatalf("args = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("args = %#v, want %#v", got, want)
		}
	}
}

func TestResolveFFmpegPathRejectsMissingPath(t *testing.T) {
	_, err := ResolveFFmpegPath(filepath.Join(t.TempDir(), "missing-ffmpeg.exe"))
	if err == nil || !strings.Contains(err.Error(), "ffmpegがインストールされていないかPATHにありません") {
		t.Fatalf("err = %v, want missing ffmpeg message", err)
	}
}

func TestResolveSpoutHelperPathRejectsMissingPath(t *testing.T) {
	_, err := ResolveSpoutHelperPath(filepath.Join(t.TempDir(), "missing-spout-capture.exe"))
	if err == nil || !strings.Contains(err.Error(), "Spout helperが見つかりません") {
		t.Fatalf("err = %v, want missing spout helper message", err)
	}
}

func TestExpandFFmpegInputPlaceholdersNoWindow(t *testing.T) {
	args := []string{"-f", "gdigrab", "-i", "desktop"}
	got, err := expandFFmpegInputPlaceholders(nil, args, "")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, " ") != strings.Join(args, " ") {
		t.Fatalf("args = %#v, want %#v", got, args)
	}
}

func TestMoveUserCameraToViewMissingView(t *testing.T) {
	cfg := DefaultConfig()
	err := MoveUserCameraToView(nil, cfg, "missing")
	if err == nil || !strings.Contains(err.Error(), "構図が見つかりません") {
		t.Fatalf("err = %v, want missing view error", err)
	}
}

func TestResetUserCameraOSCRejectsBadPort(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "bad host name"
	err := ResetUserCameraOSC(nil, cfg)
	if err == nil {
		t.Fatal("expected OSC open error")
	}
}

func TestAutoCaptureOutputPath(t *testing.T) {
	cfg := DefaultConfig().AutoCapture
	cfg.Output.Directory = t.TempDir()
	path, err := autoCaptureOutputPath(cfg, "batch-test", "shot-test", 2, CameraViewConfig{ID: "front", Name: "正面"})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(path) != cfg.Output.Directory || filepath.Ext(path) != ".png" || !strings.Contains(filepath.Base(path), "02") {
		t.Fatalf("output path = %q", path)
	}
}

func TestAutoCaptureConfigNormalizeMigratesDefaultTemplateViews(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.Views = []CameraViewConfig{
		{ID: "front", Name: "正面", Enabled: true, CoordinateSpace: "template_relative"},
	}
	cfg.Normalize()
	view := cfg.AutoCapture.Views[0]
	if view.CoordinateSpace != "world" || view.Pose.Position.Z == 0 || view.Zoom == nil {
		t.Fatalf("default template view was not migrated: %+v", view)
	}
}

func TestAutoCapturePhotoDirectoryUsesAutoPhotoSetting(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoPhoto.PhotoDirectory = filepath.Join("C:", "VRChat", "Photos")
	if got := autoCapturePhotoDirectory(cfg); got != cfg.AutoPhoto.PhotoDirectory {
		t.Fatalf("photo dir = %q, want %q", got, cfg.AutoPhoto.PhotoDirectory)
	}
}

func TestEnabledCameraViewsUsesEnabledToggleOnly(t *testing.T) {
	views := []CameraViewConfig{
		{ID: "template", Enabled: true, CoordinateSpace: "template_relative", Calibrated: false, SortOrder: 1},
		{ID: "disabled", Enabled: false, CoordinateSpace: "world", Calibrated: true, SortOrder: 2},
		{ID: "world", Enabled: true, CoordinateSpace: "world", Calibrated: true, SortOrder: 3},
	}
	got := enabledCameraViews(views)
	if len(got) != 2 || got[0].ID != "template" || got[1].ID != "world" {
		t.Fatalf("enabled views = %+v, want enabled views regardless of calibration", got)
	}
}

func TestAppendOSCStringPadsToFourBytes(t *testing.T) {
	got := appendOSCString(nil, "/x")
	if len(got)%4 != 0 {
		t.Fatalf("OSC string length = %d, want multiple of 4", len(got))
	}
	if string(got[:2]) != "/x" || got[2] != 0 {
		t.Fatalf("unexpected OSC string bytes: %v", got)
	}
}

func TestBuildOSCButtonPacketUsesBoolTypeTag(t *testing.T) {
	got := buildOSCPacket("/usercamera/Capture", ",T", func(buf []byte) []byte { return buf })
	want := appendOSCString(nil, "/usercamera/Capture")
	want = appendOSCString(want, ",T")
	if string(got) != string(want) {
		t.Fatalf("button packet = %v, want %v", got, want)
	}
}

func TestParseOSCPose(t *testing.T) {
	packet := buildOSCPacket("/usercamera/Pose", ",ffffff", func(buf []byte) []byte {
		for _, value := range []float32{1.25, 2.5, -3.75, 10, 20, 30} {
			var raw [4]byte
			binary.BigEndian.PutUint32(raw[:], math.Float32bits(value))
			buf = append(buf, raw[:]...)
		}
		return buf
	})
	pose, ok := ParseOSCPose(packet)
	if !ok {
		t.Fatal("ParseOSCPose failed")
	}
	if pose.Position.X != 1.25 || pose.Position.Y != 2.5 || pose.Position.Z != -3.75 || pose.Rotation.X != 10 || pose.Rotation.Y != 20 || pose.Rotation.Z != 30 {
		t.Fatalf("pose = %+v", pose)
	}
}

func TestNewPhotoCandidatesSortsByModTimeAndFiltersOldFiles(t *testing.T) {
	base := time.Date(2026, 6, 30, 4, 28, 0, 0, time.Local)
	files := map[string]time.Time{
		"old-before-map.png": base.Add(5 * time.Second),
		"old-time.png":       base.Add(-1 * time.Second),
		"newer.png":          base.Add(20 * time.Second),
		"new.png":            base.Add(10 * time.Second),
	}
	before := map[string]time.Time{
		"old-before-map.png": base.Add(5 * time.Second),
	}
	got := newPhotoCandidates(files, before, base)
	want := []string{"newer.png", "new.png"}
	if len(got) != len(want) {
		t.Fatalf("candidates = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("candidates = %v, want %v", got, want)
		}
	}
}

func TestParseVRChatPresenceLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output_log_2026-06-29_12-00-00.txt")
	log := "" +
		"2026.06.29 12:00:00 Log - OnPlayerJoined displayName: Alice usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa\n" +
		"2026.06.29 12:01:00 Log - OnPlayerJoined displayName: Bob usr_bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb\n" +
		"2026.06.29 12:02:00 Log - OnPlayerLeft displayName: Bob usr_bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb\n" +
		"2026.06.29 12:03:00 Debug      -  [Behaviour] OnPlayerJoined はとぽ_ (usr_dc4f8eca-e074-443a-b271-21ef533c9c3e)\n"
	if err := os.WriteFile(path, []byte(log), 0600); err != nil {
		t.Fatal(err)
	}
	users, ok := parseVRChatPresenceLog(path)
	if !ok {
		t.Fatal("parse failed")
	}
	if len(users) != 2 {
		t.Fatalf("users = %d, want 2: %+v", len(users), users)
	}
	user := users["usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"]
	if user.DisplayName != "Alice" || user.Confidence != "confirmed" {
		t.Fatalf("unexpected user: %+v", user)
	}
	vrcUser := users["usr_dc4f8eca-e074-443a-b271-21ef533c9c3e"]
	if vrcUser.DisplayName != "はとぽ_" || vrcUser.Confidence != "confirmed" {
		t.Fatalf("unexpected VRChat user: %+v", vrcUser)
	}
}

func TestWriteAutoCaptureSidecar(t *testing.T) {
	imagePath := filepath.Join(t.TempDir(), "photo.png")
	if err := os.WriteFile(imagePath, []byte("image"), 0600); err != nil {
		t.Fatal(err)
	}
	sidecar := AutoCaptureSidecar{
		SchemaVersion:   1,
		BatchID:         "batch-test",
		ShotID:          "shot-test",
		CapturedAtLocal: time.Now().Format(time.RFC3339),
		CapturedAtUTC:   time.Now().UTC().Format(time.RFC3339),
		CaptureMode:     "photo",
		View:            DefaultCameraViews()[0],
		VRChat:          AutoCaptureVRChatMetadata{UsersSource: "output_log", UsersConfidence: "partial"},
		Users:           []PresenceUser{{DisplayName: "Alice", UserID: "usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Status: "present"}},
	}
	if err := WriteAutoCaptureSidecar(imagePath, sidecar); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(imagePath + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var got AutoCaptureSidecar
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Files.ImagePath != imagePath || got.Files.SHA256 == "" || len(got.Users) != 1 {
		t.Fatalf("unexpected sidecar: %+v", got)
	}
}

func TestAutoCaptureUserIDOutputsAreIndependent(t *testing.T) {
	users := []PresenceUser{{
		DisplayName: "Alice",
		UserID:      "usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Status:      "present",
		Confidence:  "confirmed",
	}}
	cfg := DefaultAutoCaptureConfig()
	cfg.Presence.IncludeUserIDsInSidecar = false
	cfg.Presence.IncludeUserIDsInDiscord = true
	cfg.Presence.IncludeDisplayNamesInDiscord = true
	sidecarUsers := autoCaptureSidecarUsers(cfg, users)
	if len(sidecarUsers) != 1 || sidecarUsers[0].UserID != "" {
		t.Fatalf("sidecar users = %+v, want user ID removed", sidecarUsers)
	}
	content := autoCaptureDiscordContent(cfg, cfg.Views[0], users)
	if !strings.Contains(content, "usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa") {
		t.Fatalf("discord content = %q, want user ID", content)
	}

	cfg.Presence.IncludeUserIDsInSidecar = true
	cfg.Presence.IncludeUserIDsInDiscord = false
	sidecarUsers = autoCaptureSidecarUsers(cfg, users)
	if sidecarUsers[0].UserID == "" {
		t.Fatalf("sidecar users = %+v, want user ID preserved", sidecarUsers)
	}
	content = autoCaptureDiscordContent(cfg, cfg.Views[0], users)
	if strings.Contains(content, "usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa") {
		t.Fatalf("discord content = %q, want user ID omitted", content)
	}
}

func TestParseVRChatWorldMetadata(t *testing.T) {
	logText := `
2026.06.30 20:00:00 Log - Joining wrld_aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee:12345~region(jp)
2026.06.30 21:00:00 Log - Joining wrld_ffffffff-bbbb-cccc-dddd-eeeeeeeeeeee:67890~region(us)
`
	meta := parseVRChatWorldMetadata(logText)
	if meta.WorldID != "wrld_ffffffff-bbbb-cccc-dddd-eeeeeeeeeeee" || meta.InstanceID != "67890~region" {
		t.Fatalf("meta = %+v", meta)
	}
}

func TestWaitCaptureDelayCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	runner := AutoCaptureRunner{Config: Config{DiagnosticLogPath: ""}}
	view := DefaultCameraViews()[0]
	view.CaptureDelayMS = 1000
	if runner.waitCaptureDelay(ctx, view, view.Name) {
		t.Fatal("waitCaptureDelay should report cancellation")
	}
}

func TestWaitCaptureDelayPositiveWaits(t *testing.T) {
	runner := AutoCaptureRunner{Config: Config{DiagnosticLogPath: ""}}
	view := DefaultCameraViews()[0]
	view.CaptureDelayMS = 5
	start := time.Now()
	if !runner.waitCaptureDelay(context.Background(), view, view.Name) {
		t.Fatal("waitCaptureDelay should succeed")
	}
	if elapsed := time.Since(start); elapsed < 5*time.Millisecond {
		t.Fatalf("waitCaptureDelay elapsed = %s, want at least 5ms", elapsed)
	}
}

func TestWaitCaptureDelayZero(t *testing.T) {
	runner := AutoCaptureRunner{Config: Config{DiagnosticLogPath: ""}}
	view := DefaultCameraViews()[0]
	view.CaptureDelayMS = 0
	if !runner.waitCaptureDelay(context.Background(), view, view.Name) {
		t.Fatal("waitCaptureDelay should succeed without waiting when delay is zero")
	}
}
