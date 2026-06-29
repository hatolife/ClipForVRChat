package appcore

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	if cfg.AutoCapture.Views[0].ID != "front" || cfg.AutoCapture.Views[0].Calibrated {
		t.Fatalf("unexpected first view: %+v", cfg.AutoCapture.Views[0])
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
	if cfg.AutoCapture.Capture.Mode != "photo" || cfg.AutoCapture.Capture.ConcurrentMode != "sequential" {
		t.Fatalf("capture normalize failed: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Capture.RequestedCameraCount != 4 {
		t.Fatalf("RequestedCameraCount = %d, want 4", cfg.AutoCapture.Capture.RequestedCameraCount)
	}
	if len(cfg.AutoCapture.Views) != 3 {
		t.Fatalf("default views = %d, want 3", len(cfg.AutoCapture.Views))
	}
}

func TestAutoCapturePhotoDirectoryUsesAutoPhotoSetting(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoPhoto.PhotoDirectory = filepath.Join("C:", "VRChat", "Photos")
	if got := autoCapturePhotoDirectory(cfg); got != cfg.AutoPhoto.PhotoDirectory {
		t.Fatalf("photo dir = %q, want %q", got, cfg.AutoPhoto.PhotoDirectory)
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

func TestBuildOSCActionPacketUsesEmptyTypeTag(t *testing.T) {
	got := buildOSCPacket("/usercamera/Capture", ",", func(buf []byte) []byte { return buf })
	want := appendOSCString(nil, "/usercamera/Capture")
	want = appendOSCString(want, ",")
	if string(got) != string(want) {
		t.Fatalf("action packet = %v, want %v", got, want)
	}
}

func TestParseVRChatPresenceLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output_log_2026-06-29_12-00-00.txt")
	log := "" +
		"2026.06.29 12:00:00 Log - OnPlayerJoined displayName: Alice usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa\n" +
		"2026.06.29 12:01:00 Log - OnPlayerJoined displayName: Bob usr_bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb\n" +
		"2026.06.29 12:02:00 Log - OnPlayerLeft displayName: Bob usr_bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb\n"
	if err := os.WriteFile(path, []byte(log), 0600); err != nil {
		t.Fatal(err)
	}
	users, ok := parseVRChatPresenceLog(path)
	if !ok {
		t.Fatal("parse failed")
	}
	if len(users) != 1 {
		t.Fatalf("users = %d, want 1: %+v", len(users), users)
	}
	user := users["usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"]
	if user.DisplayName != "Alice" || user.Confidence != "confirmed" {
		t.Fatalf("unexpected user: %+v", user)
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
		View:            defaultCameraView("front", "正面", 0),
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
