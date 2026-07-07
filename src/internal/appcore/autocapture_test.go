package appcore

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
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
	if !cfg.AutoCapture.Schedule.CaptureOnStart {
		t.Fatal("capture on start should be enabled by default")
	}
	if cfg.AutoCapture.OSC.Host != "127.0.0.1" || cfg.AutoCapture.OSC.SendPort != 9000 || cfg.AutoCapture.OSC.ReceivePort != 9001 {
		t.Fatalf("unexpected osc defaults: %+v", cfg.AutoCapture.OSC)
	}
	if cfg.AutoCapture.PlayerLocal.BasisSource != PlayerLocalBasisSourceAvatarOSC {
		t.Fatalf("PlayerLocal.BasisSource = %q, want avatar_osc", cfg.AutoCapture.PlayerLocal.BasisSource)
	}
	if len(cfg.AutoCapture.Views) != 3 {
		t.Fatalf("default views = %d, want 3", len(cfg.AutoCapture.Views))
	}
	if cfg.AutoCapture.Capture.Mode != "stream" || cfg.AutoCapture.Stream.SpoutHelperPath == "" || !cfg.AutoCapture.Stream.SpoutAutoSelect {
		t.Fatalf("unexpected stream defaults: capture=%+v stream=%+v", cfg.AutoCapture.Capture, cfg.AutoCapture.Stream)
	}
	if !cfg.AutoCapture.Restore.Enabled || !cfg.AutoCapture.Restore.PreferSnapshot || cfg.AutoCapture.Restore.SnapshotFreshnessSec != 10 {
		t.Fatalf("unexpected restore defaults: %+v", cfg.AutoCapture.Restore)
	}
	if cfg.AutoCapture.Capture.AutoLevelRollBeforeShot == nil || !*cfg.AutoCapture.Capture.AutoLevelRollBeforeShot {
		t.Fatalf("AutoLevelRollBeforeShot should be enabled by default: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Capture.OpenCameraBeforeBatch || cfg.AutoCapture.Capture.CloseCameraAfterBatch {
		t.Fatalf("camera auto open/close should default off: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Idle.Enabled {
		t.Fatalf("idle camera should default off: %+v", cfg.AutoCapture.Idle)
	}
	if cfg.AutoCapture.Idle.View.CoordinateSpace != "player_local" || cfg.AutoCapture.Idle.View.Pose.Position.Y != -5 {
		t.Fatalf("unexpected idle camera default: %+v", cfg.AutoCapture.Idle)
	}
	if cfg.AutoCapture.Capture.PreplacedLocalAnchorEnabled() || cfg.AutoCapture.Capture.AutoEnablePreplaced || cfg.AutoCapture.Capture.AutoDisablePreplaced {
		t.Fatalf("fallback mode and auto fallback should default off: %+v", cfg.AutoCapture.Capture)
	}
	if !cfg.AutoCapture.Restore.Fallback.AutoLevelRoll {
		t.Fatalf("restore fallback AutoLevelRoll should be enabled by default: %+v", cfg.AutoCapture.Restore.Fallback)
	}
	if cfg.AutoCapture.Views[0].ID != "front" || cfg.AutoCapture.Views[0].Calibrated || cfg.AutoCapture.Views[0].Zoom == nil {
		t.Fatalf("unexpected first view: %+v", cfg.AutoCapture.Views[0])
	}
	if *cfg.AutoCapture.Views[0].Zoom != 45 {
		t.Fatalf("front view zoom = %v, want 45", *cfg.AutoCapture.Views[0].Zoom)
	}
	if cfg.AutoCapture.Views[0].Exposure == nil || *cfg.AutoCapture.Views[0].Exposure != 0 {
		t.Fatalf("front view exposure = %v, want 0", cfg.AutoCapture.Views[0].Exposure)
	}
	if cfg.AutoCapture.Views[0].RemotePlayer == nil || !*cfg.AutoCapture.Views[0].RemotePlayer {
		t.Fatalf("front view RemotePlayer should default on: %+v", cfg.AutoCapture.Views[0])
	}
	if cfg.AutoCapture.Views[0].CoordinateSpace != "player_local" || cfg.AutoCapture.Views[0].Pose.Position.Z != 1.0 {
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
	if cfg.AutoCapture.Capture.OpenCameraBeforeBatch || cfg.AutoCapture.Capture.CloseCameraAfterBatch {
		t.Fatalf("camera auto open/close should default off: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Capture.AutoLevelRollBeforeShot == nil || !*cfg.AutoCapture.Capture.AutoLevelRollBeforeShot {
		t.Fatalf("AutoLevelRollBeforeShot should default on: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Capture.PreplacedLocalAnchorEnabled() {
		t.Fatalf("PreplacedLocalAnchor should default off: %+v", cfg.AutoCapture.Capture)
	}
	if cfg.AutoCapture.Idle.Enabled || cfg.AutoCapture.Idle.View.CoordinateSpace != "player_local" || cfg.AutoCapture.Idle.View.Pose.Position.Y != -5 {
		t.Fatalf("idle camera normalize failed: %+v", cfg.AutoCapture.Idle)
	}
	if cfg.AutoCapture.Capture.AutoEnablePreplaced || cfg.AutoCapture.Capture.AutoDisablePreplaced {
		t.Fatalf("auto fallback controls should default off: %+v", cfg.AutoCapture.Capture)
	}
	if len(cfg.AutoCapture.Views) != 3 {
		t.Fatalf("default views = %d, want 3", len(cfg.AutoCapture.Views))
	}
	if cfg.AutoCapture.Stream.SpoutHelperPath != "spout-capture.exe" || !cfg.AutoCapture.Stream.SpoutAutoSelect || cfg.AutoCapture.Stream.CaptureTimeoutMS != 10000 {
		t.Fatalf("stream normalize failed: %+v", cfg.AutoCapture.Stream)
	}
	if !cfg.AutoCapture.Restore.Enabled || !cfg.AutoCapture.Restore.PreferSnapshot || cfg.AutoCapture.Restore.Fallback.Zoom != 45 || cfg.AutoCapture.Restore.Fallback.Exposure != 0 {
		t.Fatalf("restore normalize failed: %+v", cfg.AutoCapture.Restore)
	}
	if !cfg.AutoCapture.Restore.Fallback.AutoLevelRoll {
		t.Fatalf("restore fallback AutoLevelRoll should default on: %+v", cfg.AutoCapture.Restore.Fallback)
	}
}

func TestAutoCaptureConfigNormalizeMigratesOldCameraViewZoom(t *testing.T) {
	oldZoom := 1.0
	validZoom := 72.5
	cfg := DefaultConfig()
	cfg.AutoCapture.Views = []CameraViewConfig{
		{ID: "old", Name: "旧Zoom", Zoom: &oldZoom},
		{ID: "valid", Name: "有効Zoom", Zoom: &validZoom},
	}
	cfg.Normalize()
	if cfg.AutoCapture.Views[0].Zoom == nil || *cfg.AutoCapture.Views[0].Zoom != 45 {
		t.Fatalf("old zoom = %v, want 45", cfg.AutoCapture.Views[0].Zoom)
	}
	if cfg.AutoCapture.Views[1].Zoom == nil || *cfg.AutoCapture.Views[1].Zoom != 72.5 {
		t.Fatalf("valid zoom = %v, want 72.5", cfg.AutoCapture.Views[1].Zoom)
	}
}

func TestPreplacedLocalAnchorSkipsPoseResolve(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(true)
	view := DefaultCameraViews()[0]
	view.CoordinateSpace = "player_local"
	cfg.AutoCapture.PlayerLocal.Calibrated = false
	runner := AutoCaptureRunner{Config: cfg}
	if err := runner.applyCameraViewAndOptions(oscClient{}, view); err != nil {
		t.Fatalf("applyCameraViewAndOptions with preplaced local anchor returned error: %v", err)
	}
}

func TestApplyIdleCameraStateSendsPlayerLocalDefaultPose(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Idle.Enabled = true
	cfg.AutoCapture.PlayerLocal.BasisSource = PlayerLocalBasisSourceManual
	cfg.AutoCapture.PlayerLocal.Calibrated = true
	cfg.AutoCapture.PlayerLocal.BasisPose = CameraPoseConfig{}
	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := (AutoCaptureRunner{Config: cfg}).applyIdleCameraState(client); err != nil {
		t.Fatal(err)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	if len(samples) < 3 {
		t.Fatalf("packet count = %d, want at least 3: %+v", len(samples), samples)
	}
	if samples[0].Address != "/usercamera/Mode" || !samples[0].HasInt || samples[0].Int != 2 {
		t.Fatalf("mode sample = %+v, want stream mode", samples[0])
	}
	var poseSample *oscPacketSample
	for i := range samples {
		if samples[i].Address == "/usercamera/Pose" {
			poseSample = &samples[i]
			break
		}
	}
	if poseSample == nil || !poseSample.HasPose {
		t.Fatalf("pose sample missing: %+v", samples)
	}
	if poseSample.Pose.Position.X != 0 || poseSample.Pose.Position.Y != -5 || poseSample.Pose.Position.Z != 0 {
		t.Fatalf("idle pose = %+v, want player-local default resolved to Y=-5", poseSample.Pose)
	}
}

func TestApplyCameraViewTemporarilyEnablesFlyingAroundPose(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	view := cfg.AutoCapture.Views[0]
	view.CoordinateSpace = "world"
	view.Calibrated = true

	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := (AutoCaptureRunner{Config: cfg}).applyCameraView(client, view); err != nil {
		t.Fatal(err)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	if len(samples) != 3 {
		t.Fatalf("packet count = %d, want 3: %+v", len(samples), samples)
	}
	if samples[0].Address != "/usercamera/Flying" || !samples[0].HasBool || !samples[0].Bool {
		t.Fatalf("packet[0] = %+v, want Flying=true", samples[0])
	}
	if samples[1].Address != "/usercamera/Pose" || !samples[1].HasPose {
		t.Fatalf("packet[1] = %+v, want Pose", samples[1])
	}
	if samples[2].Address != "/usercamera/Flying" || !samples[2].HasBool || samples[2].Bool {
		t.Fatalf("packet[2] = %+v, want Flying=false", samples[2])
	}
}

func TestApplyCameraViewDisablesFlyingWhenPoseSendFails(t *testing.T) {
	cfg := DefaultConfig()
	view := cfg.AutoCapture.Views[0]
	view.CoordinateSpace = "world"
	view.Calibrated = true
	conn := &addressFailConn{failAddress: "/usercamera/Pose"}
	client := oscClient{host: "127.0.0.1", port: 9000, conn: conn}

	err := (AutoCaptureRunner{Config: cfg}).applyCameraView(client, view)
	if err == nil {
		t.Fatal("applyCameraView returned nil, want pose send error")
	}

	got := make([]string, 0, len(conn.addresses))
	for _, address := range conn.addresses {
		if address == avatarBeaconVersionOSCAddress {
			continue
		}
		got = append(got, address)
	}
	want := []string{"/usercamera/Flying", "/usercamera/Pose", "/usercamera/Flying"}
	if len(got) != len(want) {
		t.Fatalf("addresses = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("addresses = %v, want %v", got, want)
		}
	}
}

func TestMoveUserCameraToViewRejectsPreplacedLocalAnchorBeforeOSCOpen(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(true)
	cfg.AutoCapture.OSC.Host = "bad host name"
	err := MoveUserCameraToView(context.Background(), cfg, cfg.AutoCapture.Views[0].ID)
	if err == nil || !strings.Contains(err.Error(), "フォールバックモード") {
		t.Fatalf("MoveUserCameraToView error = %v, want fallback mode rejection", err)
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

func TestResolveSpoutHelperPathUsesExplicitExistingPath(t *testing.T) {
	dir := t.TempDir()
	helper := filepath.Join(dir, "custom-spout-capture.exe")
	if err := os.WriteFile(helper, []byte("helper"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveSpoutHelperPath(helper)
	if err != nil {
		t.Fatal(err)
	}
	if got != helper {
		t.Fatalf("path = %q, want %q", got, helper)
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

func TestResetUserCameraOSCUsesStreamingCompatAndKeepsOtherSettingsUntouched(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port

	if err := ResetUserCameraOSC(nil, cfg); err != nil {
		t.Fatal(err)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	if len(samples) != 4 {
		t.Fatalf("packet count = %d, want 4: %+v", len(samples), samples)
	}

	want := []struct {
		address string
		boolVal *bool
		intVal  *int
	}{
		{address: "/usercamera/Capture", boolVal: boolPtr(false)},
		{address: "/usercamera/Streaming", boolVal: boolPtr(false)},
		{address: "/usercamera/Streaming", intVal: intPtr(0)},
		{address: "/usercamera/Mode", intVal: intPtr(0)},
	}
	for i, sample := range samples {
		if sample.Address != want[i].address {
			t.Fatalf("packet[%d].address = %q, want %q", i, sample.Address, want[i].address)
		}
		if want[i].boolVal != nil {
			if !sample.HasBool || sample.Bool != *want[i].boolVal {
				t.Fatalf("packet[%d] bool = %+v, want %t", i, sample, *want[i].boolVal)
			}
			continue
		}
		if !sample.HasInt || sample.Int != *want[i].intVal {
			t.Fatalf("packet[%d] int = %+v, want %d", i, sample, *want[i].intVal)
		}
	}
}

func TestSendCameraButtonReleasesOnCancellation(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := sendCameraButton(ctx, client, "/usercamera/Capture", 25, "", "cancel_test"); err == nil {
		t.Fatal("expected cancellation error")
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	if len(samples) != 2 {
		t.Fatalf("packet count = %d, want 2: %+v", len(samples), samples)
	}
	if samples[0].Address != "/usercamera/Capture" || !samples[0].HasBool || !samples[0].Bool {
		t.Fatalf("first packet = %+v, want press true", samples[0])
	}
	if samples[1].Address != "/usercamera/Capture" || !samples[1].HasBool || samples[1].Bool {
		t.Fatalf("second packet = %+v, want release false", samples[1])
	}
}

func TestOSCTraceHandlerCapturesSendPackets(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	events := make(chan OSCTraceEvent, 4)
	cleanup := SetOSCTraceHandler(func(event OSCTraceEvent) {
		events <- event
	})
	defer cleanup()

	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := client.sendBool("/usercamera/Streaming", true); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-events:
		if event.Direction != "send" || event.Address != "/usercamera/Streaming" || event.TypeTags != ",T" || event.Status != "ok" {
			t.Fatalf("unexpected trace event: %+v", event)
		}
		if event.Target != fmt.Sprintf("127.0.0.1:%d", port) {
			t.Fatalf("target = %q, want port %d", event.Target, port)
		}
	case <-time.After(time.Second):
		t.Fatal("trace event was not captured")
	}
}

func TestParseDebugOSCLine(t *testing.T) {
	address, typeTags, appendArgs, err := parseDebugOSCLine("/avatar/parameters/debug true i:3 f:1.5 s:note 7 2.25 label")
	if err != nil {
		t.Fatal(err)
	}
	if address != "/avatar/parameters/debug" {
		t.Fatalf("address = %q", address)
	}
	if typeTags != ",Tifsifs" {
		t.Fatalf("typeTags = %q", typeTags)
	}
	packet := buildOSCPacket(address, typeTags, appendArgs)
	gotAddress, gotTypes, payload, ok := ParseOSCPacket(packet)
	if !ok {
		t.Fatal("ParseOSCPacket failed")
	}
	if gotAddress != address || gotTypes != typeTags {
		t.Fatalf("parsed = %q %q, want %q %q", gotAddress, gotTypes, address, typeTags)
	}
	if len(payload) == 0 {
		t.Fatal("payload is empty")
	}
}

func TestAutoCaptureRunnerRunOnceReleasesStreamingOnCancellation(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Restore.Enabled = false
	cfg.AutoCapture.Capture.Mode = "stream"
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(false)
	cfg.AutoCapture.Capture.OpenCameraBeforeBatch = true
	cfg.AutoCapture.Capture.CloseCameraAfterBatch = false
	cfg.AutoCapture.Capture.SettleDelayMS = 1500
	cfg.AutoCapture.Stream.StartDelayMS = 0
	cfg.AutoCapture.Stream.SpoutHelperPath = filepath.Join(t.TempDir(), "missing-spout-capture.exe")
	cfg.AutoCapture.Output.Directory = t.TempDir()
	cfg.AutoCapture.Presence.WatchOutputLog = false
	cfg.AutoCapture.Views = []CameraViewConfig{{
		ID:              "front",
		Name:            "front",
		Enabled:         true,
		CoordinateSpace: "world",
		Calibrated:      true,
		SettleDelayMS:   2000,
	}, {
		ID:              "side",
		Name:            "side",
		Enabled:         true,
		CoordinateSpace: "world",
		Calibrated:      true,
	}}

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(2700*time.Millisecond, cancel)

	results, err := (AutoCaptureRunner{Config: cfg}).RunOnce(ctx)
	if err == nil {
		t.Fatalf("RunOnce error = nil, results=%+v", results)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one shot result before cancellation")
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	hasStreamStart := false
	hasStreamStop := false
	hasStreamStartInt := false
	hasStreamStopInt := false
	for _, sample := range samples {
		if sample.Address != "/usercamera/Streaming" {
			continue
		}
		if sample.HasBool {
			if sample.Bool {
				hasStreamStart = true
			} else {
				hasStreamStop = true
			}
		}
		if sample.HasInt {
			switch sample.Int {
			case 1:
				hasStreamStartInt = true
			case 0:
				hasStreamStopInt = true
			}
		}
	}
	if !hasStreamStart {
		t.Fatalf("stream start packet not found: %+v", samples)
	}
	if !hasStreamStop {
		t.Fatalf("stream stop packet not found: %+v", samples)
	}
	if !hasStreamStartInt {
		t.Fatalf("stream start compat int packet not found: %+v", samples)
	}
	if !hasStreamStopInt {
		t.Fatalf("stream stop compat int packet not found: %+v", samples)
	}
}

func TestAutoCaptureRunnerRunOnceSkipsCameraAutoOpenWhenDisabled(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Restore.Enabled = false
	cfg.AutoCapture.Capture.Mode = "stream"
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(false)
	cfg.AutoCapture.Capture.OpenCameraBeforeBatch = false
	cfg.AutoCapture.Capture.CloseCameraAfterBatch = true
	cfg.AutoCapture.Stream.SpoutHelperPath = filepath.Join(t.TempDir(), "missing-spout-capture.exe")
	cfg.AutoCapture.Output.Directory = t.TempDir()
	cfg.AutoCapture.Presence.WatchOutputLog = false
	cfg.AutoCapture.Views = []CameraViewConfig{{
		ID:              "front",
		Name:            "front",
		Enabled:         true,
		CoordinateSpace: "world",
		Calibrated:      true,
		SettleDelayMS:   1,
	}}

	results, err := (AutoCaptureRunner{Config: cfg}).RunOnce(context.Background())
	if err != nil {
		t.Fatalf("RunOnce error = %v", err)
	}
	if len(results) != 1 || results[0].Error == "" {
		t.Fatalf("results = %+v, want one failed stream shot", results)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	hasAutoLevelRoll := false
	for _, sample := range samples {
		switch sample.Address {
		case "/usercamera/Mode", "/usercamera/SmoothMovement", "/usercamera/Streaming":
			t.Fatalf("camera auto-open packet should not be sent when disabled: %+v all=%+v", sample, samples)
		case "/usercamera/Close":
			t.Fatalf("camera close packet should not be sent even when legacy closeCameraAfterBatch is true: %+v all=%+v", sample, samples)
		}
		if sample.Address == "/usercamera/AutoLevelRoll" && sample.HasBool && sample.Bool {
			hasAutoLevelRoll = true
		}
	}
	if !hasAutoLevelRoll {
		t.Fatalf("AutoLevelRoll=true packet not found: %+v", samples)
	}
}

func TestAutoCaptureRunnerRunOnceCanDisableAutoLevelRollBeforeShot(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Restore.Enabled = false
	cfg.AutoCapture.Capture.Mode = "stream"
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(false)
	cfg.AutoCapture.Capture.OpenCameraBeforeBatch = false
	cfg.AutoCapture.Capture.CloseCameraAfterBatch = false
	cfg.AutoCapture.Capture.AutoLevelRollBeforeShot = boolPtr(false)
	cfg.AutoCapture.Stream.SpoutHelperPath = filepath.Join(t.TempDir(), "missing-spout-capture.exe")
	cfg.AutoCapture.Output.Directory = t.TempDir()
	cfg.AutoCapture.Presence.WatchOutputLog = false
	cfg.AutoCapture.Views = []CameraViewConfig{{
		ID:              "front",
		Name:            "front",
		Enabled:         true,
		CoordinateSpace: "world",
		Calibrated:      true,
		SettleDelayMS:   1,
	}}

	results, err := (AutoCaptureRunner{Config: cfg}).RunOnce(context.Background())
	if err != nil {
		t.Fatalf("RunOnce error = %v", err)
	}
	if len(results) != 1 || results[0].Error == "" {
		t.Fatalf("results = %+v, want one failed stream shot", results)
	}

	samples := readOSCPacketSamples(t, conn)
	for _, sample := range samples {
		if sample.Address == "/usercamera/AutoLevelRoll" {
			t.Fatalf("AutoLevelRoll packet should not be sent when disabled: %+v all=%+v", sample, samples)
		}
	}
}

func TestRecoverEmptySpoutSenderListTogglesStreaming(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	originalList := autoCaptureListSpoutSenders
	defer func() { autoCaptureListSpoutSenders = originalList }()
	calls := 0
	autoCaptureListSpoutSenders = func(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) (SpoutListResult, error) {
		calls++
		return SpoutListResult{OK: true, Senders: nil}, nil
	}

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(false)
	cfg.AutoCapture.Capture.OpenCameraBeforeBatch = true
	runner := AutoCaptureRunner{Config: cfg}
	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := runner.recoverEmptySpoutSenderList(context.Background(), client, "front"); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("sender list calls = %d, want 2", calls)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	want := []struct {
		boolVal *bool
		intVal  *int
	}{
		{boolVal: boolPtr(false)},
		{intVal: intPtr(0)},
		{boolVal: boolPtr(true)},
		{intVal: intPtr(1)},
	}
	if len(samples) != len(want) {
		t.Fatalf("packet count = %d, want %d: %+v", len(samples), len(want), samples)
	}
	for i, sample := range samples {
		if sample.Address != "/usercamera/Streaming" {
			t.Fatalf("packet[%d].address = %q, want /usercamera/Streaming", i, sample.Address)
		}
		if want[i].boolVal != nil {
			if !sample.HasBool || sample.Bool != *want[i].boolVal {
				t.Fatalf("packet[%d] bool = %+v, want %t", i, sample, *want[i].boolVal)
			}
			continue
		}
		if !sample.HasInt || sample.Int != *want[i].intVal {
			t.Fatalf("packet[%d] int = %+v, want %d", i, sample, *want[i].intVal)
		}
	}
}

func TestRecoverEmptySpoutSenderListUsesOnOnlyWhenAutoOpenDisabledAndSenderReturns(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	originalList := autoCaptureListSpoutSenders
	defer func() { autoCaptureListSpoutSenders = originalList }()
	calls := 0
	autoCaptureListSpoutSenders = func(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) (SpoutListResult, error) {
		calls++
		if calls == 2 {
			return SpoutListResult{OK: true, Senders: []SpoutSenderInfo{{Name: "VRCSender1", Width: 1920, Height: 1080}}}, nil
		}
		return SpoutListResult{OK: true, Senders: nil}, nil
	}

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(false)
	cfg.AutoCapture.Capture.OpenCameraBeforeBatch = false
	runner := AutoCaptureRunner{Config: cfg}
	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := runner.recoverEmptySpoutSenderList(context.Background(), client, "front"); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("sender list calls = %d, want 2", calls)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	want := []struct {
		boolVal *bool
		intVal  *int
	}{
		{boolVal: boolPtr(true)},
		{intVal: intPtr(1)},
	}
	if len(samples) != len(want) {
		t.Fatalf("packet count = %d, want %d: %+v", len(samples), len(want), samples)
	}
	for i, sample := range samples {
		if sample.Address != "/usercamera/Streaming" {
			t.Fatalf("packet[%d].address = %q, want /usercamera/Streaming", i, sample.Address)
		}
		if sample.HasBool && !sample.Bool {
			t.Fatalf("packet[%d] stopped streaming: %+v all=%+v", i, sample, samples)
		}
		if sample.HasInt && sample.Int == 0 {
			t.Fatalf("packet[%d] stopped streaming by int compat: %+v all=%+v", i, sample, samples)
		}
		if want[i].boolVal != nil {
			if !sample.HasBool || sample.Bool != *want[i].boolVal {
				t.Fatalf("packet[%d] bool = %+v, want %t", i, sample, *want[i].boolVal)
			}
			continue
		}
		if !sample.HasInt || sample.Int != *want[i].intVal {
			t.Fatalf("packet[%d] int = %+v, want %d", i, sample, *want[i].intVal)
		}
	}
}

func TestRecoverEmptySpoutSenderListEscalatesToToggleWhenAutoOpenDisabled(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	originalList := autoCaptureListSpoutSenders
	defer func() { autoCaptureListSpoutSenders = originalList }()
	calls := 0
	autoCaptureListSpoutSenders = func(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) (SpoutListResult, error) {
		calls++
		return SpoutListResult{OK: true, Senders: nil}, nil
	}

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	cfg.AutoCapture.Capture.PreplacedLocalAnchor = boolPtr(false)
	cfg.AutoCapture.Capture.OpenCameraBeforeBatch = false
	runner := AutoCaptureRunner{Config: cfg}
	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := runner.recoverEmptySpoutSenderList(context.Background(), client, "front"); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("sender list calls = %d, want 3", calls)
	}

	samples := withoutVersionNoticePackets(readOSCPacketSamples(t, conn))
	want := []struct {
		boolVal *bool
		intVal  *int
	}{
		{boolVal: boolPtr(true)},
		{intVal: intPtr(1)},
		{boolVal: boolPtr(false)},
		{intVal: intPtr(0)},
		{boolVal: boolPtr(true)},
		{intVal: intPtr(1)},
	}
	if len(samples) != len(want) {
		t.Fatalf("packet count = %d, want %d: %+v", len(samples), len(want), samples)
	}
	for i, sample := range samples {
		if sample.Address != "/usercamera/Streaming" {
			t.Fatalf("packet[%d].address = %q, want /usercamera/Streaming", i, sample.Address)
		}
		if want[i].boolVal != nil {
			if !sample.HasBool || sample.Bool != *want[i].boolVal {
				t.Fatalf("packet[%d] bool = %+v, want %t", i, sample, *want[i].boolVal)
			}
			continue
		}
		if !sample.HasInt || sample.Int != *want[i].intVal {
			t.Fatalf("packet[%d] int = %+v, want %d", i, sample, *want[i].intVal)
		}
	}
}

func TestRecoverEmptySpoutSenderListSkipsWhenSenderExists(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	originalList := autoCaptureListSpoutSenders
	defer func() { autoCaptureListSpoutSenders = originalList }()
	autoCaptureListSpoutSenders = func(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) (SpoutListResult, error) {
		return SpoutListResult{OK: true, Senders: []SpoutSenderInfo{{Name: "VRCSender1", Width: 1920, Height: 1080}}}, nil
	}

	cfg := DefaultConfig()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.SendPort = port
	runner := AutoCaptureRunner{Config: cfg}
	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()

	if err := runner.recoverEmptySpoutSenderList(context.Background(), client, "front"); err != nil {
		t.Fatal(err)
	}
	if samples := readOSCPacketSamples(t, conn); len(samples) != 0 {
		t.Fatalf("unexpected recovery packets when sender exists: %+v", samples)
	}
}

func listenOSCUserCameraPackets(t *testing.T) (net.PacketConn, int) {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	return conn, addr.Port
}

type oscPacketSample struct {
	Address  string
	TypeTags string
	Bool     bool
	HasBool  bool
	Int      int
	HasInt   bool
	Float    float64
	HasFloat bool
	Pose     CameraPoseConfig
	HasPose  bool
	String   string
	HasStr   bool
}

func readOSCPacketSamples(t *testing.T, conn net.PacketConn) []oscPacketSample {
	t.Helper()
	samples := make([]oscPacketSample, 0, 16)
	buf := make([]byte, 2048)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(150 * time.Millisecond)); err != nil {
			t.Fatal(err)
		}
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			t.Fatal(err)
		}
		address, typeTags, payload, ok := ParseOSCPacket(append([]byte(nil), buf[:n]...))
		if !ok {
			t.Fatalf("failed to parse OSC packet: %v", buf[:n])
		}
		sample := oscPacketSample{Address: address, TypeTags: typeTags}
		typeTag, _ := firstOSCTag(typeTags)
		if typeTag != 0 {
			switch typeTag {
			case 'T':
				sample.HasBool = true
				sample.Bool = true
			case 'F':
				sample.HasBool = true
				sample.Bool = false
			case 'i':
				sample.HasInt = true
				if len(payload) < 4 {
					t.Fatalf("OSC int packet too short: %v", buf[:n])
				}
				sample.Int = int(int32(binary.BigEndian.Uint32(payload[:4])))
			case 'f':
				sample.HasFloat = true
				if len(payload) < 4 {
					t.Fatalf("OSC float packet too short: %v", buf[:n])
				}
				sample.Float = float64(math.Float32frombits(binary.BigEndian.Uint32(payload[:4])))
			case 's':
				sample.HasStr = true
				str, _, ok := readOSCString(payload, 0)
				if !ok {
					t.Fatalf("OSC string packet could not be decoded: %v", buf[:n])
				}
				sample.String = str
			}
		}
		if pose, ok := ParseOSCPose(append([]byte(nil), buf[:n]...)); ok {
			sample.Pose = pose
			sample.HasPose = true
		}
		samples = append(samples, sample)
	}
	return samples
}

func withoutVersionNoticePackets(samples []oscPacketSample) []oscPacketSample {
	filtered := make([]oscPacketSample, 0, len(samples))
	for _, sample := range samples {
		if sample.Address == avatarBeaconVersionOSCAddress {
			continue
		}
		filtered = append(filtered, sample)
	}
	return filtered
}

type addressFailConn struct {
	failAddress string
	addresses   []string
}

func (c *addressFailConn) Read(_ []byte) (int, error) {
	return 0, errors.New("read is not supported")
}

func (c *addressFailConn) Write(packet []byte) (int, error) {
	address, _, _, ok := ParseOSCPacket(packet)
	if ok {
		c.addresses = append(c.addresses, address)
	}
	if address == c.failAddress {
		return 0, errors.New("forced OSC write failure")
	}
	return len(packet), nil
}

func (c *addressFailConn) Close() error {
	return nil
}

func (c *addressFailConn) LocalAddr() net.Addr {
	return dummyAddr("local")
}

func (c *addressFailConn) RemoteAddr() net.Addr {
	return dummyAddr("remote")
}

func (c *addressFailConn) SetDeadline(_ time.Time) error {
	return nil
}

func (c *addressFailConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *addressFailConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

type dummyAddr string

func (a dummyAddr) Network() string {
	return string(a)
}

func (a dummyAddr) String() string {
	return string(a)
}

func boolPtr(value bool) *bool {
	v := value
	return &v
}

func intPtr(value int) *int {
	v := value
	return &v
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
	if view.CoordinateSpace != "player_local" || view.Pose.Position.Z != 1.0 || view.Zoom == nil {
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

func TestOSCClientSendsVersionNoticeOnce(t *testing.T) {
	conn, port := listenOSCUserCameraPackets(t)
	defer conn.Close()

	SetOSCVersionNotice("v9.8.7-b6")
	client := oscClient{host: "127.0.0.1", port: port}
	if err := client.open(); err != nil {
		t.Fatal(err)
	}
	defer client.close()
	if err := client.sendInt("/usercamera/Mode", 1); err != nil {
		t.Fatal(err)
	}
	if err := client.sendBool("/usercamera/Streaming", true); err != nil {
		t.Fatal(err)
	}

	samples := readOSCPacketSamples(t, conn)
	versionPackets := 0
	modePackets := 0
	streamingPackets := 0
	for _, sample := range samples {
		switch sample.Address {
		case avatarBeaconVersionOSCAddress:
			versionPackets++
			if !sample.HasStr || sample.String != "v9.8.7-b6" {
				t.Fatalf("version sample = %+v, want string v9.8.7-b6", sample)
			}
		case "/usercamera/Mode":
			modePackets++
		case "/usercamera/Streaming":
			streamingPackets++
		}
	}
	if versionPackets != 1 || modePackets != 1 || streamingPackets != 1 {
		t.Fatalf("packets version=%d mode=%d streaming=%d samples=%+v", versionPackets, modePackets, streamingPackets, samples)
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

func TestParseOSCUserCameraSample(t *testing.T) {
	modePacket := buildOSCPacket("/usercamera/Mode", ",i", func(buf []byte) []byte {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], uint32(2))
		return append(buf, raw[:]...)
	})
	mode, ok := ParseOSCUserCameraSample(modePacket)
	if !ok || !mode.HasInt || mode.Int != 2 {
		t.Fatalf("mode sample = %+v ok=%t", mode, ok)
	}

	streamingPacket := buildOSCPacket("/usercamera/Streaming", ",T", func(buf []byte) []byte { return buf })
	streaming, ok := ParseOSCUserCameraSample(streamingPacket)
	if !ok || !streaming.HasBool || !streaming.Bool {
		t.Fatalf("streaming sample = %+v ok=%t", streaming, ok)
	}

	zoomPacket := buildOSCPacket("/usercamera/Zoom", ",f", func(buf []byte) []byte {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], math.Float32bits(72.5))
		return append(buf, raw[:]...)
	})
	zoom, ok := ParseOSCUserCameraSample(zoomPacket)
	if !ok || !zoom.HasFloat || math.Abs(zoom.Float-72.5) > 0.001 {
		t.Fatalf("zoom sample = %+v ok=%t", zoom, ok)
	}
}

func TestMergeUserCameraRestoreStatePrefersSnapshot(t *testing.T) {
	restore := defaultAutoCaptureRestoreConfig()
	restore.Snapshot = AutoCaptureUserCameraState{
		Mode:      intStatePtr(2),
		Streaming: boolStatePtr(true),
		Zoom:      floatStatePtr(80),
	}
	target := mergeUserCameraRestoreState(restore)
	if target.Mode == nil || *target.Mode != 2 {
		t.Fatalf("mode = %v, want snapshot 2", target.Mode)
	}
	if target.Streaming == nil || !*target.Streaming {
		t.Fatalf("streaming = %v, want snapshot true", target.Streaming)
	}
	if target.Zoom == nil || *target.Zoom != 80 {
		t.Fatalf("zoom = %v, want snapshot 80", target.Zoom)
	}
	if target.Exposure == nil || *target.Exposure != 0 {
		t.Fatalf("exposure = %v, want fallback 0", target.Exposure)
	}
}

func TestMergeUserCameraRestoreStateCanIgnoreSnapshot(t *testing.T) {
	restore := defaultAutoCaptureRestoreConfig()
	restore.PreferSnapshot = false
	restore.Snapshot = AutoCaptureUserCameraState{Mode: intStatePtr(2)}
	target := mergeUserCameraRestoreState(restore)
	if target.Mode == nil || *target.Mode != 0 {
		t.Fatalf("mode = %v, want fallback 0", target.Mode)
	}
}

func TestSuppressFallbackCameraActivationStateKeepsOnlySnapshotValues(t *testing.T) {
	restore := defaultAutoCaptureRestoreConfig()
	target := mergeUserCameraRestoreState(restore)
	suppressFallbackCameraActivationState(&target, restore.Snapshot)
	if target.Mode != nil {
		t.Fatalf("mode = %v, want nil when snapshot is missing", target.Mode)
	}
	if target.Streaming != nil {
		t.Fatalf("streaming = %v, want nil when snapshot is missing", target.Streaming)
	}

	streaming := true
	mode := 2
	restore.Snapshot = AutoCaptureUserCameraState{
		Mode:      &mode,
		Streaming: &streaming,
	}
	target = mergeUserCameraRestoreState(restore)
	suppressFallbackCameraActivationState(&target, restore.Snapshot)
	if target.Mode == nil || *target.Mode != 2 {
		t.Fatalf("mode = %v, want snapshot 2", target.Mode)
	}
	if target.Streaming == nil || !*target.Streaming {
		t.Fatalf("streaming = %v, want snapshot true", target.Streaming)
	}
}

func TestFinalizeAutoCaptureImageCreatesThumbnail(t *testing.T) {
	path := filepath.Join(t.TempDir(), "capture.png")
	writeTestPNG(t, path, 20, 10)
	cfg := DefaultConfig()
	cfg.AutoCapture.Output.WriteEXIF = false
	cfg.AutoCapture.Output.WriteSidecarJSON = false
	cfg.AutoCapture.Discord.Enabled = false
	cfg.Output.UploadDiscord = false

	result := (AutoCaptureRunner{Config: cfg}).finalizeAutoCaptureImage(path, "batch-test", "shot-test", cfg.AutoCapture.Views[0], nil, nil, "unknown", AutoCaptureVRChatMetadata{}, SpoutCaptureResult{})
	if result.Error != "" {
		t.Fatalf("result error = %q", result.Error)
	}
	if result.Thumbnail == "" || !strings.HasPrefix(result.Thumbnail, "data:image/png;base64,") {
		t.Fatalf("thumbnail = %q, want data URL", result.Thumbnail)
	}
}

func TestAutoCaptureDiscordUploadEnabledFollowsPrimaryUploadSetting(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoCapture.Discord.Enabled = false
	cfg.Output.UploadDiscord = true
	if !autoCaptureDiscordUploadEnabled(cfg) {
		t.Fatal("auto-capture Discord upload should be enabled by primary upload setting")
	}
	cfg.Output.UploadDiscord = false
	cfg.AutoCapture.Discord.Enabled = true
	if !autoCaptureDiscordUploadEnabled(cfg) {
		t.Fatal("auto-capture Discord upload should be enabled by auto-capture setting")
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
2026.06.30 20:00:10 Log - Loading avatar avtr_11111111-2222-3333-4444-555555555555.
2026.06.30 21:00:00 Log - Joining wrld_ffffffff-bbbb-cccc-dddd-eeeeeeeeeeee:67890~private(usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa)~nonce(123456)~region(us).
`
	meta := parseVRChatWorldMetadata(logText)
	if meta.WorldID != "wrld_ffffffff-bbbb-cccc-dddd-eeeeeeeeeeee" || meta.InstanceID != "67890~private(usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa)~nonce(123456)~region(us)" {
		t.Fatalf("meta = %+v", meta)
	}
	if meta.AvatarID != "avtr_11111111-2222-3333-4444-555555555555" {
		t.Fatalf("AvatarID = %q", meta.AvatarID)
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
