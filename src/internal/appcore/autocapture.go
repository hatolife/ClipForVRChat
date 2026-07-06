package appcore

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type OSCTraceEvent struct {
	Direction string
	Address   string
	TypeTags  string
	Payload   []byte
	Target    string
	Status    string
	Error     string
}

var oscTraceMu sync.RWMutex
var oscTraceHandler func(OSCTraceEvent)
var oscTraceHandlerID uint64
var autoCaptureListSpoutSenders = ListSpoutSenders

func SetOSCTraceHandler(handler func(OSCTraceEvent)) func() {
	oscTraceMu.Lock()
	previous := oscTraceHandler
	previousID := oscTraceHandlerID
	oscTraceHandlerID++
	id := oscTraceHandlerID
	oscTraceHandler = handler
	oscTraceMu.Unlock()
	return func() {
		oscTraceMu.Lock()
		if oscTraceHandlerID == id {
			oscTraceHandler = previous
			oscTraceHandlerID = previousID
		}
		oscTraceMu.Unlock()
	}
}

func emitOSCTrace(event OSCTraceEvent) {
	oscTraceMu.RLock()
	handler := oscTraceHandler
	oscTraceMu.RUnlock()
	if handler != nil {
		handler(event)
	}
}

type AutoCaptureEvent struct {
	BatchID string `json:"batchId"`
	ShotID  string `json:"shotId"`
	Path    string `json:"path"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type DebugOSCSendResult struct {
	OK       bool   `json:"ok"`
	Message  string `json:"message"`
	Target   string `json:"target"`
	Address  string `json:"address"`
	TypeTags string `json:"typeTags"`
}

type debugOSCArg struct {
	tag   byte
	intV  int32
	float float32
	str   string
}

type PresenceUser struct {
	DisplayName string `json:"displayName"`
	UserID      string `json:"userId,omitempty"`
	Status      string `json:"status"`
	Source      string `json:"source"`
	Confidence  string `json:"confidence"`
	JoinedAt    string `json:"joinedAt,omitempty"`
	LeftAt      string `json:"leftAt,omitempty"`
}

type AutoCaptureRunner struct {
	Config  Config
	Handler func(AutoCaptureEvent)
}

type CameraPoseSnapshot struct {
	Pose       CameraPoseConfig `json:"pose"`
	UpdatedAt  string           `json:"updatedAt"`
	AgeMS      int64            `json:"ageMs"`
	Fresh      bool             `json:"fresh"`
	Configured bool             `json:"configured"`
}

type UserCameraOSCSample struct {
	Address  string
	Bool     bool
	HasBool  bool
	Int      int
	HasInt   bool
	Float    float64
	HasFloat bool
	Pose     CameraPoseConfig
	HasPose  bool
}

func MoveUserCameraToView(ctx context.Context, cfg Config, viewID string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg.Normalize()
	ac := cfg.AutoCapture
	viewID = strings.TrimSpace(viewID)
	var view CameraViewConfig
	found := false
	for _, candidate := range ac.Views {
		if candidate.ID == viewID {
			view = candidate
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("構図が見つかりません: %s", viewID)
	}
	if ac.Capture.PreplacedLocalAnchorEnabled() {
		return fmt.Errorf("ローカルアンカー配置済みカメラを使うフォールバックモードでは、ClipForVRChatからカメラ移動は送信しません。VRChat内でUser Cameraをローカルアンカーにして配置してください")
	}
	logPath := cfg.DiagnosticLogPath
	client := oscClient{host: ac.OSC.Host, port: ac.OSC.SendPort}
	diagAutoCapture(logPath, "move camera open begin: target=%s:%d view_id=%q", ac.OSC.Host, ac.OSC.SendPort, viewID)
	if err := client.open(); err != nil {
		diagAutoCapture(logPath, "move camera open error: target=%s:%d view_id=%q err=%v", ac.OSC.Host, ac.OSC.SendPort, viewID, err)
		return err
	}
	defer client.close()
	cameraMode := int32(1)
	if ac.Capture.Mode == "stream" {
		cameraMode = 2
	}
	if err := client.sendInt("/usercamera/Mode", cameraMode); err != nil {
		diagAutoCapture(logPath, "move camera mode error: view_id=%q mode=%d err=%v", viewID, cameraMode, err)
		return err
	}
	if !sleepContext(ctx, 500*time.Millisecond) {
		return ctx.Err()
	}
	runner := AutoCaptureRunner{Config: cfg}
	if err := runner.applyCameraView(client, view); err != nil {
		return err
	}
	sentOptions := sendOptionalFloat(client, "/usercamera/Zoom", view.Zoom) +
		sendOptionalFloat(client, "/usercamera/Exposure", view.Exposure) +
		sendOptionalFloat(client, "/usercamera/FocalDistance", view.FocalDistance) +
		sendOptionalFloat(client, "/usercamera/Aperture", view.Aperture) +
		sendOptionalBool(client, "/usercamera/LookAtMe", view.LookAtMe) +
		sendOptionalBool(client, "/usercamera/ShowUIInCamera", view.ShowUIInCamera) +
		sendOptionalBool(client, "/usercamera/LocalPlayer", view.LocalPlayer) +
		sendOptionalBool(client, "/usercamera/RemotePlayer", view.RemotePlayer) +
		sendOptionalBool(client, "/usercamera/Environment", view.Environment)
	diagAutoCapture(logPath, "move camera complete: view_id=%q mode=%d optional_params=%d", viewID, cameraMode, sentOptions)
	return nil
}

func ResetUserCameraOSC(ctx context.Context, cfg Config) error {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg.Normalize()
	ac := cfg.AutoCapture
	logPath := cfg.DiagnosticLogPath
	client := oscClient{host: ac.OSC.Host, port: ac.OSC.SendPort}
	diagAutoCapture(logPath, "osc recovery open begin: target=%s:%d", ac.OSC.Host, ac.OSC.SendPort)
	if err := client.open(); err != nil {
		diagAutoCapture(logPath, "osc recovery open error: target=%s:%d err=%v", ac.OSC.Host, ac.OSC.SendPort, err)
		return err
	}
	defer client.close()
	for _, address := range []string{"/usercamera/Capture", "/usercamera/Close"} {
		if err := client.sendBool(address, false); err != nil {
			diagAutoCapture(logPath, "osc recovery bool error: address=%q err=%v", address, err)
			return err
		}
		diagAutoCapture(logPath, "osc recovery bool success: address=%q value=false", address)
		if !sleepContext(ctx, 100*time.Millisecond) {
			return ctx.Err()
		}
	}
	if err := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", false, "osc_recovery"); err != nil {
		return err
	}
	if err := client.sendInt("/usercamera/Mode", 0); err != nil {
		diagAutoCapture(logPath, "osc recovery mode error: err=%v", err)
		return err
	}
	diagAutoCapture(logPath, "osc recovery mode success: value=0")
	return nil
}

func (r AutoCaptureRunner) RunOnce(ctx context.Context) (results []Result, err error) {
	cfg := r.Config
	cfg.Normalize()
	ac := cfg.AutoCapture
	logPath := cfg.DiagnosticLogPath
	diagAutoCapture(logPath, "run_once begin: mode=%q schedule_enabled=%t capture_on_start=%t interval_sec=%d preplaced_local_anchor=%t open_before_batch=%t close_after_batch=%t button_release_ms=%d settle_ms=%d",
		ac.Capture.Mode,
		ac.Schedule.Enabled,
		ac.Schedule.CaptureOnStart,
		ac.Schedule.CaptureIntervalSec,
		ac.Capture.PreplacedLocalAnchorEnabled(),
		ac.Capture.OpenCameraBeforeBatch,
		ac.Capture.CloseCameraAfterBatch,
		ac.Capture.ButtonReleaseDelayMS,
		ac.Capture.SettleDelayMS,
	)
	if ac.Capture.Mode != "photo" && ac.Capture.Mode != "stream" {
		diagAutoCapture(logPath, "run_once reject: unsupported_mode=%q", ac.Capture.Mode)
		return nil, fmt.Errorf("対応していない自動撮影方式です: %s", ac.Capture.Mode)
	}
	views := enabledCameraViews(ac.Views)
	if len(views) == 0 {
		diagAutoCapture(logPath, "run_once reject: enabled_views=0 total_views=%d", len(ac.Views))
		return nil, fmt.Errorf("撮影ONの構図がありません。自動撮影タブで撮影する構図を1つ以上ONにしてください。")
	}
	batchID := newBatchID(time.Now())
	users, confidence, presenceLogPath := SnapshotVRChatPresenceWithSource(ac.Presence.OutputLogDirectory)
	world := SnapshotVRChatWorld(ac.Presence.OutputLogDirectory)
	if !ac.Presence.WatchOutputLog {
		users = nil
		confidence = "unknown"
		world = AutoCaptureVRChatMetadata{}
	}
	sidecarUsers := autoCaptureSidecarUsers(ac, users)
	photoDir := autoCapturePhotoDirectory(cfg)
	before := scanAutoCapturePhotoFiles(photoDir, ac.Output.Directory)
	diagAutoCapture(logPath, "run_once prepared: batch_id=%q views=%d total_views=%d users=%d sidecar_users=%d users_confidence=%q watch_output_log=%t output_log_dir=%q output_log_path=%q photo_dir=%q output_dir=%q before_files=%d before_latest=%s",
		batchID,
		len(views),
		len(ac.Views),
		len(users),
		len(sidecarUsers),
		confidence,
		ac.Presence.WatchOutputLog,
		ac.Presence.OutputLogDirectory,
		presenceLogPath,
		photoDir,
		ac.Output.Directory,
		len(before),
		photoFileSummary(before),
	)
	client := oscClient{host: ac.OSC.Host, port: ac.OSC.SendPort}
	diagAutoCapture(logPath, "osc open begin: target=%s:%d", ac.OSC.Host, ac.OSC.SendPort)
	if err := client.open(); err != nil {
		diagAutoCapture(logPath, "osc open error: target=%s:%d err=%v", ac.OSC.Host, ac.OSC.SendPort, err)
		return nil, err
	}
	defer client.close()
	streamStarted := false
	if ac.Restore.Enabled && !ac.Capture.PreplacedLocalAnchorEnabled() {
		defer r.restoreUserCameraState(client)
	} else if ac.Capture.PreplacedLocalAnchorEnabled() {
		diagAutoCapture(logPath, "camera restore skipped: preplaced_local_anchor=true")
	}
	defer func() {
		if err == nil || !streamStarted {
			return
		}
		diagAutoCapture(logPath, "stream cleanup begin: reason=%q", "run_once_error")
		if cleanupErr := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", false, "run_once_cleanup"); cleanupErr != nil {
			diagAutoCapture(logPath, "stream cleanup failed: err=%v", cleanupErr)
			return
		}
		diagAutoCapture(logPath, "stream cleanup success: reason=%q", "run_once_error")
	}()
	diagAutoCapture(logPath, "osc open success: target=%s:%d", ac.OSC.Host, ac.OSC.SendPort)
	if ac.Capture.OpenCameraBeforeBatch && !ac.Capture.PreplacedLocalAnchorEnabled() {
		cameraMode := 1
		if ac.Capture.Mode == "stream" {
			cameraMode = 2
		}
		diagAutoCapture(logPath, "osc send begin: address=%q value=%d capture_mode=%q", "/usercamera/Mode", cameraMode, ac.Capture.Mode)
		if err := client.sendInt("/usercamera/Mode", int32(cameraMode)); err != nil {
			diagAutoCapture(logPath, "osc send error: address=%q err=%v", "/usercamera/Mode", err)
			return nil, err
		}
		diagAutoCapture(logPath, "osc send success: address=%q value=%d capture_mode=%q", "/usercamera/Mode", cameraMode, ac.Capture.Mode)
		modeWait := 2500 * time.Millisecond
		diagAutoCapture(logPath, "camera mode wait begin: duration_ms=%d", modeWait.Milliseconds())
		if !sleepContext(ctx, modeWait) {
			diagAutoCapture(logPath, "camera mode wait cancelled: err=%v", ctx.Err())
			return nil, ctx.Err()
		}
		diagAutoCapture(logPath, "camera mode wait complete")
		if ac.Capture.Mode == "stream" {
			diagAutoCapture(logPath, "osc send begin: address=%q value=%t detail=%q", "/usercamera/SmoothMovement", true, "stream_prepare")
			if err := client.sendBool("/usercamera/SmoothMovement", true); err != nil {
				diagAutoCapture(logPath, "osc send error: address=%q detail=%q err=%v", "/usercamera/SmoothMovement", "stream_prepare", err)
				return nil, err
			}
			diagAutoCapture(logPath, "osc send success: address=%q value=%t detail=%q", "/usercamera/SmoothMovement", true, "stream_prepare")
			diagAutoCapture(logPath, "osc button press begin: address=%q detail=%q", "/usercamera/Streaming", "stream_start")
			if err := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", true, "stream_start"); err != nil {
				return nil, err
			}
			streamStarted = true
			diagAutoCapture(logPath, "osc button press success: address=%q detail=%q", "/usercamera/Streaming", "stream_start")
			startDelay := time.Duration(ac.Stream.StartDelayMS) * time.Millisecond
			if startDelay > 0 {
				diagAutoCapture(logPath, "stream start wait begin: duration_ms=%d", startDelay.Milliseconds())
				if !sleepContext(ctx, startDelay) {
					diagAutoCapture(logPath, "stream start wait cancelled: err=%v", ctx.Err())
					return nil, ctx.Err()
				}
				diagAutoCapture(logPath, "stream start wait complete")
			}
		}
	} else if ac.Capture.PreplacedLocalAnchorEnabled() {
		diagAutoCapture(logPath, "camera auto open skipped: preplaced_local_anchor=true capture_mode=%q", ac.Capture.Mode)
	} else {
		diagAutoCapture(logPath, "camera auto open skipped: open_before_batch=false capture_mode=%q", ac.Capture.Mode)
	}
	results = make([]Result, 0, len(views))
	for i, view := range views {
		if err := ctx.Err(); err != nil {
			diagAutoCapture(logPath, "run_once cancelled before shot: index=%d err=%v", i+1, err)
			return results, err
		}
		shotID := fmt.Sprintf("%s-%02d", batchID, i+1)
		var result Result
		if ac.Capture.Mode == "stream" {
			result = r.captureStreamShot(ctx, client, batchID, shotID, i+1, view, sidecarUsers, users, confidence, world)
		} else {
			result = r.capturePhotoShot(ctx, client, batchID, shotID, i+1, view, photoDir, before, sidecarUsers, users, confidence, world)
		}
		results = append(results, result)
		diagAutoCapture(logPath, "shot result: batch_id=%q shot_id=%q source_path=%q error=%q", batchID, shotID, result.SourcePath, result.Error)
		r.emit(AutoCaptureEvent{BatchID: batchID, ShotID: shotID, Path: result.SourcePath, Error: result.Error, Message: result.Name})
		if result.SourcePath != "" {
			before[result.SourcePath] = time.Now()
		}
	}
	successCount := 0
	for _, result := range results {
		if result.SourcePath != "" {
			successCount++
		}
	}
	if ac.Capture.CloseCameraAfterBatch && ac.Capture.PreplacedLocalAnchorEnabled() {
		diagAutoCapture(logPath, "camera close skipped: preplaced_local_anchor=true")
	} else if ac.Capture.CloseCameraAfterBatch && (successCount > 0 || ac.Capture.Mode == "stream") {
		if ac.Capture.Mode == "stream" {
			diagAutoCapture(logPath, "osc button release begin: address=%q detail=%q", "/usercamera/Streaming", "stream_stop")
			if err := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", false, "stream_stop"); err != nil {
				diagAutoCapture(logPath, "stream stop failed: err=%v", err)
			} else {
				diagAutoCapture(logPath, "osc button release success: address=%q detail=%q", "/usercamera/Streaming", "stream_stop")
			}
		}
		if err := sendCameraButton(ctx, client, "/usercamera/Close", ac.Capture.ButtonReleaseDelayMS, logPath, "batch_close"); err != nil {
			diagAutoCapture(logPath, "camera close failed: err=%v", err)
		}
	} else if ac.Capture.CloseCameraAfterBatch {
		diagAutoCapture(logPath, "camera close skipped: reason=%q successful_shots=%d", "no_successful_shots", successCount)
	} else {
		diagAutoCapture(logPath, "camera close skipped: close_after_batch=false")
	}
	diagAutoCapture(logPath, "run_once complete: batch_id=%q results=%d", batchID, len(results))
	return results, nil
}

func (r AutoCaptureRunner) restoreUserCameraState(client oscClient) {
	cfg := r.Config.AutoCapture
	restore := cfg.Restore
	logPath := r.Config.DiagnosticLogPath
	if !restore.Enabled {
		diagAutoCapture(logPath, "camera restore skipped: enabled=false")
		return
	}
	target := mergeUserCameraRestoreState(restore)
	modeSummary := "<missing>"
	if target.Mode != nil {
		modeSummary = fmt.Sprintf("%d", *target.Mode)
	}
	streamingSummary := "<missing>"
	if target.Streaming != nil {
		streamingSummary = fmt.Sprintf("%t", *target.Streaming)
	}
	diagAutoCapture(
		logPath,
		"camera restore begin: prefer_snapshot=%t snapshot_values=%d fallback_values=%d target_values=%d mode=%s streaming=%s pose=%t snapshot=%q fallback=%q target=%q",
		restore.PreferSnapshot,
		countUserCameraStateValues(restore.Snapshot),
		countUserCameraStateValues(userCameraFallbackState(restore.Fallback)),
		countUserCameraStateValues(target),
		modeSummary,
		streamingSummary,
		target.Pose != nil,
		formatUserCameraStateValues(restore.Snapshot),
		formatUserCameraStateValues(userCameraFallbackState(restore.Fallback)),
		formatUserCameraStateValues(target),
	)
	if target.Streaming != nil && !*target.Streaming {
		sendRestoreBoolCompat(client, logPath, "/usercamera/Streaming", *target.Streaming)
	}
	if target.Mode != nil {
		sendRestoreInt(client, logPath, "/usercamera/Mode", *target.Mode)
		time.Sleep(150 * time.Millisecond)
	}
	if target.Pose != nil && shouldRestoreCameraParameters(target.Mode) {
		pose := *target.Pose
		diagAutoCapture(logPath, "camera restore send begin: address=%q", "/usercamera/Pose")
		if err := client.sendFloats("/usercamera/Pose", []float32{
			float32(pose.Position.X), float32(pose.Position.Y), float32(pose.Position.Z),
			float32(pose.Rotation.X), float32(pose.Rotation.Y), float32(pose.Rotation.Z),
		}); err != nil {
			diagAutoCapture(logPath, "camera restore send error: address=%q err=%v", "/usercamera/Pose", err)
		} else {
			diagAutoCapture(logPath, "camera restore send success: address=%q", "/usercamera/Pose")
		}
	}
	if shouldRestoreCameraParameters(target.Mode) {
		for _, item := range []struct {
			address string
			value   *float64
		}{
			{"/usercamera/Zoom", target.Zoom},
			{"/usercamera/Exposure", target.Exposure},
			{"/usercamera/FocalDistance", target.FocalDistance},
			{"/usercamera/Aperture", target.Aperture},
			{"/usercamera/Hue", target.Hue},
			{"/usercamera/Saturation", target.Saturation},
			{"/usercamera/Lightness", target.Lightness},
			{"/usercamera/LookAtMeXOffset", target.LookAtMeXOffset},
			{"/usercamera/LookAtMeYOffset", target.LookAtMeYOffset},
			{"/usercamera/FlySpeed", target.FlySpeed},
			{"/usercamera/TurnSpeed", target.TurnSpeed},
			{"/usercamera/SmoothingStrength", target.SmoothingStrength},
			{"/usercamera/PhotoRate", target.PhotoRate},
			{"/usercamera/Duration", target.Duration},
		} {
			if item.value != nil {
				sendRestoreFloat(client, logPath, item.address, *item.value)
			}
		}
		for _, item := range []struct {
			address string
			value   *bool
		}{
			{"/usercamera/SmoothMovement", target.SmoothMovement},
			{"/usercamera/ShowUIInCamera", target.ShowUIInCamera},
			{"/usercamera/Lock", target.Lock},
			{"/usercamera/LocalPlayer", target.LocalPlayer},
			{"/usercamera/RemotePlayer", target.RemotePlayer},
			{"/usercamera/Environment", target.Environment},
			{"/usercamera/GreenScreen", target.GreenScreen},
			{"/usercamera/LookAtMe", target.LookAtMe},
			{"/usercamera/AutoLevelRoll", target.AutoLevelRoll},
			{"/usercamera/AutoLevelPitch", target.AutoLevelPitch},
			{"/usercamera/Flying", target.Flying},
			{"/usercamera/TriggerTakesPhotos", target.TriggerTakesPhotos},
			{"/usercamera/DollyPathsStayVisible", target.DollyPathsStayVisible},
			{"/usercamera/CameraEars", target.CameraEars},
			{"/usercamera/ShowFocus", target.ShowFocus},
			{"/usercamera/RollWhileFlying", target.RollWhileFlying},
			{"/usercamera/OrientationIsLandscape", target.OrientationIsLandscape},
		} {
			if item.value != nil {
				sendRestoreBool(client, logPath, item.address, *item.value)
			}
		}
	}
	if target.Streaming != nil && *target.Streaming && target.Mode != nil && *target.Mode == 2 {
		sendRestoreBoolCompat(client, logPath, "/usercamera/Streaming", *target.Streaming)
	}
	diagAutoCapture(logPath, "camera restore complete")
}

func shouldRestoreCameraParameters(mode *int) bool {
	return mode == nil || *mode != 0
}

func sendRestoreInt(client oscClient, logPath string, address string, value int) {
	diagAutoCapture(logPath, "camera restore send begin: address=%q value=%d", address, value)
	if err := client.sendInt(address, int32(value)); err != nil {
		diagAutoCapture(logPath, "camera restore send error: address=%q value=%d err=%v", address, value, err)
		return
	}
	diagAutoCapture(logPath, "camera restore send success: address=%q value=%d", address, value)
}

func sendRestoreFloat(client oscClient, logPath string, address string, value float64) {
	diagAutoCapture(logPath, "camera restore send begin: address=%q value=%.4f", address, value)
	if err := client.sendFloat(address, float32(value)); err != nil {
		diagAutoCapture(logPath, "camera restore send error: address=%q value=%.4f err=%v", address, value, err)
		return
	}
	diagAutoCapture(logPath, "camera restore send success: address=%q value=%.4f", address, value)
}

func sendRestoreBool(client oscClient, logPath string, address string, value bool) {
	diagAutoCapture(logPath, "camera restore send begin: address=%q value=%t", address, value)
	if err := client.sendBool(address, value); err != nil {
		diagAutoCapture(logPath, "camera restore send error: address=%q value=%t err=%v", address, value, err)
		return
	}
	diagAutoCapture(logPath, "camera restore send success: address=%q value=%t", address, value)
}

func sendRestoreBoolCompat(client oscClient, logPath string, address string, value bool) {
	sendRestoreBool(client, logPath, address, value)
	intValue := boolOSCInt(value)
	diagAutoCapture(logPath, "camera restore compat int send begin: address=%q value=%d", address, intValue)
	if err := client.sendInt(address, intValue); err != nil {
		diagAutoCapture(logPath, "camera restore compat int send error: address=%q value=%d err=%v", address, intValue, err)
		return
	}
	diagAutoCapture(logPath, "camera restore compat int send success: address=%q value=%d", address, intValue)
}

func mergeUserCameraRestoreState(restore AutoCaptureRestoreConfig) AutoCaptureUserCameraState {
	target := userCameraFallbackState(restore.Fallback)
	if restore.PreferSnapshot {
		overlayUserCameraState(&target, restore.Snapshot)
	}
	return target
}

func userCameraFallbackState(fallback AutoCaptureUserCameraFallbackConfig) AutoCaptureUserCameraState {
	state := AutoCaptureUserCameraState{
		Mode:                   intStatePtr(fallback.Mode),
		Streaming:              boolStatePtr(fallback.Streaming),
		SmoothMovement:         boolStatePtr(fallback.SmoothMovement),
		Zoom:                   floatStatePtr(fallback.Zoom),
		Exposure:               floatStatePtr(fallback.Exposure),
		FocalDistance:          floatStatePtr(fallback.FocalDistance),
		Aperture:               floatStatePtr(fallback.Aperture),
		Hue:                    floatStatePtr(fallback.Hue),
		Saturation:             floatStatePtr(fallback.Saturation),
		Lightness:              floatStatePtr(fallback.Lightness),
		LookAtMeXOffset:        floatStatePtr(fallback.LookAtMeXOffset),
		LookAtMeYOffset:        floatStatePtr(fallback.LookAtMeYOffset),
		FlySpeed:               floatStatePtr(fallback.FlySpeed),
		TurnSpeed:              floatStatePtr(fallback.TurnSpeed),
		SmoothingStrength:      floatStatePtr(fallback.SmoothingStrength),
		PhotoRate:              floatStatePtr(fallback.PhotoRate),
		Duration:               floatStatePtr(fallback.Duration),
		ShowUIInCamera:         boolStatePtr(fallback.ShowUIInCamera),
		Lock:                   boolStatePtr(fallback.Lock),
		LocalPlayer:            boolStatePtr(fallback.LocalPlayer),
		RemotePlayer:           boolStatePtr(fallback.RemotePlayer),
		Environment:            boolStatePtr(fallback.Environment),
		GreenScreen:            boolStatePtr(fallback.GreenScreen),
		LookAtMe:               boolStatePtr(fallback.LookAtMe),
		AutoLevelRoll:          boolStatePtr(fallback.AutoLevelRoll),
		AutoLevelPitch:         boolStatePtr(fallback.AutoLevelPitch),
		Flying:                 boolStatePtr(fallback.Flying),
		TriggerTakesPhotos:     boolStatePtr(fallback.TriggerTakesPhotos),
		DollyPathsStayVisible:  boolStatePtr(fallback.DollyPathsStayVisible),
		CameraEars:             boolStatePtr(fallback.CameraEars),
		ShowFocus:              boolStatePtr(fallback.ShowFocus),
		RollWhileFlying:        boolStatePtr(fallback.RollWhileFlying),
		OrientationIsLandscape: boolStatePtr(fallback.OrientationIsLandscape),
	}
	if fallback.RestorePose {
		state.Pose = poseStatePtr(fallback.Pose)
	}
	return state
}

func overlayUserCameraState(target *AutoCaptureUserCameraState, snapshot AutoCaptureUserCameraState) {
	overlayIntState(&target.Mode, snapshot.Mode)
	overlayPoseState(&target.Pose, snapshot.Pose)
	overlayBoolState(&target.Streaming, snapshot.Streaming)
	overlayBoolState(&target.SmoothMovement, snapshot.SmoothMovement)
	overlayFloatState(&target.Zoom, snapshot.Zoom)
	overlayFloatState(&target.Exposure, snapshot.Exposure)
	overlayFloatState(&target.FocalDistance, snapshot.FocalDistance)
	overlayFloatState(&target.Aperture, snapshot.Aperture)
	overlayFloatState(&target.Hue, snapshot.Hue)
	overlayFloatState(&target.Saturation, snapshot.Saturation)
	overlayFloatState(&target.Lightness, snapshot.Lightness)
	overlayFloatState(&target.LookAtMeXOffset, snapshot.LookAtMeXOffset)
	overlayFloatState(&target.LookAtMeYOffset, snapshot.LookAtMeYOffset)
	overlayFloatState(&target.FlySpeed, snapshot.FlySpeed)
	overlayFloatState(&target.TurnSpeed, snapshot.TurnSpeed)
	overlayFloatState(&target.SmoothingStrength, snapshot.SmoothingStrength)
	overlayFloatState(&target.PhotoRate, snapshot.PhotoRate)
	overlayFloatState(&target.Duration, snapshot.Duration)
	overlayBoolState(&target.ShowUIInCamera, snapshot.ShowUIInCamera)
	overlayBoolState(&target.Lock, snapshot.Lock)
	overlayBoolState(&target.LocalPlayer, snapshot.LocalPlayer)
	overlayBoolState(&target.RemotePlayer, snapshot.RemotePlayer)
	overlayBoolState(&target.Environment, snapshot.Environment)
	overlayBoolState(&target.GreenScreen, snapshot.GreenScreen)
	overlayBoolState(&target.LookAtMe, snapshot.LookAtMe)
	overlayBoolState(&target.AutoLevelRoll, snapshot.AutoLevelRoll)
	overlayBoolState(&target.AutoLevelPitch, snapshot.AutoLevelPitch)
	overlayBoolState(&target.Flying, snapshot.Flying)
	overlayBoolState(&target.TriggerTakesPhotos, snapshot.TriggerTakesPhotos)
	overlayBoolState(&target.DollyPathsStayVisible, snapshot.DollyPathsStayVisible)
	overlayBoolState(&target.CameraEars, snapshot.CameraEars)
	overlayBoolState(&target.ShowFocus, snapshot.ShowFocus)
	overlayBoolState(&target.RollWhileFlying, snapshot.RollWhileFlying)
	overlayBoolState(&target.OrientationIsLandscape, snapshot.OrientationIsLandscape)
}

func countUserCameraStateValues(state AutoCaptureUserCameraState) int {
	count := 0
	if state.Mode != nil {
		count++
	}
	if state.Pose != nil {
		count++
	}
	for _, value := range []*bool{
		state.Streaming,
		state.SmoothMovement,
		state.ShowUIInCamera,
		state.Lock,
		state.LocalPlayer,
		state.RemotePlayer,
		state.Environment,
		state.GreenScreen,
		state.LookAtMe,
		state.AutoLevelRoll,
		state.AutoLevelPitch,
		state.Flying,
		state.TriggerTakesPhotos,
		state.DollyPathsStayVisible,
		state.CameraEars,
		state.ShowFocus,
		state.RollWhileFlying,
		state.OrientationIsLandscape,
	} {
		if value != nil {
			count++
		}
	}
	for _, value := range []*float64{
		state.Zoom,
		state.Exposure,
		state.FocalDistance,
		state.Aperture,
		state.Hue,
		state.Saturation,
		state.Lightness,
		state.LookAtMeXOffset,
		state.LookAtMeYOffset,
		state.FlySpeed,
		state.TurnSpeed,
		state.SmoothingStrength,
		state.PhotoRate,
		state.Duration,
	} {
		if value != nil {
			count++
		}
	}
	return count
}

func formatUserCameraStateValues(state AutoCaptureUserCameraState) string {
	values := make([]string, 0, countUserCameraStateValues(state))
	if state.Mode != nil {
		values = append(values, fmt.Sprintf("Mode=%d", *state.Mode))
	}
	if state.Pose != nil {
		values = append(values, "Pose=<set>")
	}
	appendBool := func(name string, value *bool) {
		if value != nil {
			values = append(values, fmt.Sprintf("%s=%t", name, *value))
		}
	}
	appendFloat := func(name string, value *float64) {
		if value != nil {
			values = append(values, fmt.Sprintf("%s=%.4f", name, *value))
		}
	}
	appendBool("Streaming", state.Streaming)
	appendBool("SmoothMovement", state.SmoothMovement)
	appendBool("ShowUIInCamera", state.ShowUIInCamera)
	appendBool("Lock", state.Lock)
	appendBool("LocalPlayer", state.LocalPlayer)
	appendBool("RemotePlayer", state.RemotePlayer)
	appendBool("Environment", state.Environment)
	appendBool("GreenScreen", state.GreenScreen)
	appendBool("LookAtMe", state.LookAtMe)
	appendBool("AutoLevelRoll", state.AutoLevelRoll)
	appendBool("AutoLevelPitch", state.AutoLevelPitch)
	appendBool("Flying", state.Flying)
	appendBool("TriggerTakesPhotos", state.TriggerTakesPhotos)
	appendBool("DollyPathsStayVisible", state.DollyPathsStayVisible)
	appendBool("CameraEars", state.CameraEars)
	appendBool("ShowFocus", state.ShowFocus)
	appendBool("RollWhileFlying", state.RollWhileFlying)
	appendBool("OrientationIsLandscape", state.OrientationIsLandscape)
	appendFloat("Zoom", state.Zoom)
	appendFloat("Exposure", state.Exposure)
	appendFloat("FocalDistance", state.FocalDistance)
	appendFloat("Aperture", state.Aperture)
	appendFloat("Hue", state.Hue)
	appendFloat("Saturation", state.Saturation)
	appendFloat("Lightness", state.Lightness)
	appendFloat("LookAtMeXOffset", state.LookAtMeXOffset)
	appendFloat("LookAtMeYOffset", state.LookAtMeYOffset)
	appendFloat("FlySpeed", state.FlySpeed)
	appendFloat("TurnSpeed", state.TurnSpeed)
	appendFloat("SmoothingStrength", state.SmoothingStrength)
	appendFloat("PhotoRate", state.PhotoRate)
	appendFloat("Duration", state.Duration)
	if len(values) == 0 {
		return "<none>"
	}
	return strings.Join(values, ",")
}

func intStatePtr(value int) *int {
	v := value
	return &v
}

func boolStatePtr(value bool) *bool {
	v := value
	return &v
}

func floatStatePtr(value float64) *float64 {
	v := value
	return &v
}

func poseStatePtr(value CameraPoseConfig) *CameraPoseConfig {
	v := value
	return &v
}

func overlayIntState(target **int, value *int) {
	if value != nil {
		*target = intStatePtr(*value)
	}
}

func overlayBoolState(target **bool, value *bool) {
	if value != nil {
		*target = boolStatePtr(*value)
	}
}

func overlayFloatState(target **float64, value *float64) {
	if value != nil {
		*target = floatStatePtr(*value)
	}
}

func overlayPoseState(target **CameraPoseConfig, value *CameraPoseConfig) {
	if value != nil {
		*target = poseStatePtr(*value)
	}
}

func (r AutoCaptureRunner) capturePhotoShot(ctx context.Context, client oscClient, batchID string, shotID string, index int, view CameraViewConfig, photoDir string, before map[string]time.Time, sidecarUsers []PresenceUser, discordUsers []PresenceUser, confidence string, world AutoCaptureVRChatMetadata) Result {
	cfg := r.Config.AutoCapture
	logPath := r.Config.DiagnosticLogPath
	name := view.Name
	if name == "" {
		name = view.ID
	}
	diagAutoCapture(logPath, "shot begin: batch_id=%q shot_id=%q index=%d view_id=%q view_name=%q coordinate_space=%q calibrated=%t settle_ms=%d capture_delay_ms=%d",
		batchID,
		shotID,
		index,
		view.ID,
		name,
		view.CoordinateSpace,
		view.Calibrated,
		view.SettleDelayMS,
		view.CaptureDelayMS,
	)
	if err := r.applyCameraViewAndOptions(client, view); err != nil {
		return Result{Name: name, Error: err.Error()}
	}
	settle := time.Duration(cfg.Capture.SettleDelayMS) * time.Millisecond
	if view.SettleDelayMS > 0 {
		settle = time.Duration(view.SettleDelayMS) * time.Millisecond
	}
	diagAutoCapture(logPath, "shot settle begin: view_id=%q duration_ms=%d", view.ID, settle.Milliseconds())
	if !sleepContext(ctx, settle) {
		diagAutoCapture(logPath, "shot settle cancelled: view_id=%q err=%v", view.ID, ctx.Err())
		return Result{Name: name, Error: "自動撮影が中断されました。"}
	}
	diagAutoCapture(logPath, "shot settle complete: view_id=%q", view.ID)
	if !r.waitCaptureDelay(ctx, view, name) {
		return Result{Name: name, Error: "自動撮影が中断されました。"}
	}
	captureNotBefore := time.Now().Add(-1 * time.Second)
	if err := sendCameraButton(ctx, client, "/usercamera/Capture", cfg.Capture.ButtonReleaseDelayMS, logPath, view.ID); err != nil {
		return Result{Name: name, Error: err.Error()}
	}
	photoPath := waitForNewPhoto(ctx, photoDir, before, 30*time.Second, captureNotBefore, logPath)
	if photoPath == "" {
		photoPath = waitForNewPhoto(ctx, cfg.Output.Directory, before, 3*time.Second, captureNotBefore, logPath)
	}
	if photoPath == "" {
		diagAutoCapture(logPath, "shot photo detection failed: view_id=%q photo_dir=%q output_dir=%q before_files=%d before_latest=%s", view.ID, photoDir, cfg.Output.Directory, len(before), photoFileSummary(before))
		return Result{Name: name, Error: "撮影後のVRChat写真ファイルを検出できませんでした。Photo方式ではUser Cameraが表示され、VRChatの写真保存先が正しい必要があります。Stream方式を使う場合はffmpeg入力設定を確認してください。"}
	}
	return r.finalizeAutoCaptureImage(photoPath, batchID, shotID, view, sidecarUsers, discordUsers, confidence, world, SpoutCaptureResult{})
}

func (r AutoCaptureRunner) applyCameraView(client oscClient, view CameraViewConfig) error {
	logPath := r.Config.DiagnosticLogPath
	pose, err := ResolveCameraViewPose(r.Config.AutoCapture, view)
	if err != nil {
		diagAutoCapture(logPath, "camera pose resolve error: view_id=%q view_name=%q coordinate_space=%q basis_source=%q manual_calibrated=%t err=%v", view.ID, view.Name, view.CoordinateSpace, r.Config.AutoCapture.PlayerLocal.BasisSource, r.Config.AutoCapture.PlayerLocal.Calibrated, err)
		return err
	}
	if view.CoordinateSpace == "player_local" {
		diagAutoCapture(logPath, "camera pose resolved: view_id=%q view_name=%q coordinate_space=%q basis_source=%q basis_pose=%+v local_pose=%+v resolved_pose=%+v", view.ID, view.Name, view.CoordinateSpace, r.Config.AutoCapture.PlayerLocal.BasisSource, r.Config.AutoCapture.PlayerLocal.BasisPose, view.Pose, pose)
	} else {
		diagAutoCapture(logPath, "camera pose resolved: view_id=%q view_name=%q coordinate_space=%q world_pose=%+v", view.ID, view.Name, view.CoordinateSpace, pose)
	}
	diagAutoCapture(logPath, "osc send begin: address=%q view_id=%q", "/usercamera/Pose", view.ID)
	if err := client.sendFloats("/usercamera/Pose", []float32{
		float32(pose.Position.X), float32(pose.Position.Y), float32(pose.Position.Z),
		float32(pose.Rotation.X), float32(pose.Rotation.Y), float32(pose.Rotation.Z),
	}); err != nil {
		diagAutoCapture(logPath, "osc send error: address=%q view_id=%q err=%v", "/usercamera/Pose", view.ID, err)
		return err
	} else {
		diagAutoCapture(logPath, "osc send success: address=%q view_id=%q coordinate_space=%q resolved_pose=%+v", "/usercamera/Pose", view.ID, view.CoordinateSpace, pose)
	}
	return nil
}

func (r AutoCaptureRunner) applyCameraViewAndOptions(client oscClient, view CameraViewConfig) error {
	logPath := r.Config.DiagnosticLogPath
	if r.Config.AutoCapture.Capture.PreplacedLocalAnchorEnabled() {
		diagAutoCapture(logPath, "camera pose/options skipped: view_id=%q preplaced_local_anchor=true", view.ID)
		return r.applyAutoLevelRollBeforeShot(client, view.ID)
	}
	if err := r.applyCameraView(client, view); err != nil {
		return err
	}
	sentOptions := sendOptionalFloat(client, "/usercamera/Zoom", view.Zoom) +
		sendOptionalFloat(client, "/usercamera/Exposure", view.Exposure) +
		sendOptionalFloat(client, "/usercamera/FocalDistance", view.FocalDistance) +
		sendOptionalFloat(client, "/usercamera/Aperture", view.Aperture) +
		sendOptionalBool(client, "/usercamera/LookAtMe", view.LookAtMe) +
		sendOptionalBool(client, "/usercamera/ShowUIInCamera", view.ShowUIInCamera) +
		sendOptionalBool(client, "/usercamera/LocalPlayer", view.LocalPlayer) +
		sendOptionalBool(client, "/usercamera/RemotePlayer", view.RemotePlayer) +
		sendOptionalBool(client, "/usercamera/Environment", view.Environment)
	if err := r.applyAutoLevelRollBeforeShot(client, view.ID); err != nil {
		return err
	}
	diagAutoCapture(logPath, "shot optional_params sent: view_id=%q count=%d", view.ID, sentOptions)
	return nil
}

func (r AutoCaptureRunner) applyAutoLevelRollBeforeShot(client oscClient, viewID string) error {
	logPath := r.Config.DiagnosticLogPath
	enabled := true
	if r.Config.AutoCapture.Capture.AutoLevelRollBeforeShot != nil {
		enabled = *r.Config.AutoCapture.Capture.AutoLevelRollBeforeShot
	}
	if !enabled {
		diagAutoCapture(logPath, "auto level roll skipped: view_id=%q auto_level_roll_before_shot=false", viewID)
		return nil
	}
	if client.conn == nil {
		diagAutoCapture(logPath, "auto level roll skipped: view_id=%q reason=%q", viewID, "osc_not_open")
		return nil
	}
	diagAutoCapture(logPath, "osc send begin: address=%q value=%t view_id=%q detail=%q", "/usercamera/AutoLevelRoll", true, viewID, "before_shot")
	if err := client.sendBool("/usercamera/AutoLevelRoll", true); err != nil {
		diagAutoCapture(logPath, "osc send error: address=%q value=%t view_id=%q detail=%q err=%v", "/usercamera/AutoLevelRoll", true, viewID, "before_shot", err)
		return err
	}
	diagAutoCapture(logPath, "osc send success: address=%q value=%t view_id=%q detail=%q", "/usercamera/AutoLevelRoll", true, viewID, "before_shot")
	return nil
}

func (r AutoCaptureRunner) captureStreamShot(ctx context.Context, client oscClient, batchID string, shotID string, index int, view CameraViewConfig, sidecarUsers []PresenceUser, discordUsers []PresenceUser, confidence string, world AutoCaptureVRChatMetadata) Result {
	cfg := r.Config.AutoCapture
	logPath := r.Config.DiagnosticLogPath
	name := view.Name
	if name == "" {
		name = view.ID
	}
	diagAutoCapture(logPath, "stream shot begin: batch_id=%q shot_id=%q index=%d view_id=%q view_name=%q", batchID, shotID, index, view.ID, name)
	if err := r.applyCameraViewAndOptions(client, view); err != nil {
		return Result{Name: name, Error: err.Error()}
	}
	settle := time.Duration(cfg.Capture.SettleDelayMS) * time.Millisecond
	if view.SettleDelayMS > 0 {
		settle = time.Duration(view.SettleDelayMS) * time.Millisecond
	}
	diagAutoCapture(logPath, "shot settle begin: view_id=%q duration_ms=%d", view.ID, settle.Milliseconds())
	if !sleepContext(ctx, settle) {
		diagAutoCapture(logPath, "shot settle cancelled: view_id=%q err=%v", view.ID, ctx.Err())
		return Result{Name: name, Error: "自動撮影が中断されました。"}
	}
	diagAutoCapture(logPath, "shot settle complete: view_id=%q", view.ID)
	if !r.waitCaptureDelay(ctx, view, name) {
		return Result{Name: name, Error: "自動撮影が中断されました。"}
	}
	if err := r.ensureStreamCameraForSpoutCapture(ctx, client, view.ID); err != nil {
		return Result{Name: name, Error: err.Error()}
	}
	if err := r.recoverEmptySpoutSenderList(ctx, client, view.ID); err != nil {
		return Result{Name: name, Error: err.Error()}
	}
	outputPath, err := autoCaptureOutputPath(cfg, batchID, shotID, index, view)
	if err != nil {
		diagAutoCapture(logPath, "stream output path error: view_id=%q err=%v", view.ID, err)
		return Result{Name: name, Error: err.Error()}
	}
	capturePath := outputPath
	if strings.EqualFold(filepath.Ext(outputPath), ".jpg") || strings.EqualFold(filepath.Ext(outputPath), ".jpeg") {
		capturePath = outputPath + ".spout.png"
	}
	streamInfo, err := captureStreamFrameWithSpout(ctx, cfg.Stream, capturePath, logPath)
	if err != nil {
		diagAutoCapture(logPath, "stream capture failed: view_id=%q output=%q capture_path=%q err=%v", view.ID, outputPath, capturePath, err)
		return Result{Name: name, Error: err.Error()}
	}
	if capturePath != outputPath {
		if err := convertAutoCaptureImage(capturePath, outputPath, cfg.Output.ImageFormat); err != nil {
			diagAutoCapture(logPath, "stream convert failed: view_id=%q capture_path=%q output=%q err=%v", view.ID, capturePath, outputPath, err)
			return Result{Name: name, Error: err.Error()}
		}
		_ = os.Remove(capturePath)
		diagAutoCapture(logPath, "stream convert success: view_id=%q capture_path=%q output=%q format=%q", view.ID, capturePath, outputPath, cfg.Output.ImageFormat)
	}
	return r.finalizeAutoCaptureImage(outputPath, batchID, shotID, view, sidecarUsers, discordUsers, confidence, world, streamInfo)
}

func (r AutoCaptureRunner) ensureStreamCameraForSpoutCapture(ctx context.Context, client oscClient, viewID string) error {
	logPath := r.Config.DiagnosticLogPath
	if r.Config.AutoCapture.Capture.PreplacedLocalAnchorEnabled() {
		diagAutoCapture(logPath, "stream camera refresh skipped: preplaced_local_anchor=true view_id=%q", viewID)
		return nil
	}
	if !r.Config.AutoCapture.Capture.OpenCameraBeforeBatch {
		diagAutoCapture(logPath, "stream camera refresh skipped: open_before_batch=false view_id=%q", viewID)
		return nil
	}
	diagAutoCapture(logPath, "stream camera refresh begin: view_id=%q", viewID)
	if err := client.sendInt("/usercamera/Mode", 2); err != nil {
		diagAutoCapture(logPath, "stream camera refresh error: address=%q view_id=%q err=%v", "/usercamera/Mode", viewID, err)
		return err
	}
	if err := client.sendBool("/usercamera/SmoothMovement", true); err != nil {
		diagAutoCapture(logPath, "stream camera refresh error: address=%q view_id=%q err=%v", "/usercamera/SmoothMovement", viewID, err)
		return err
	}
	if err := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", true, "stream_refresh:"+viewID); err != nil {
		diagAutoCapture(logPath, "stream camera refresh error: address=%q view_id=%q err=%v", "/usercamera/Streaming", viewID, err)
		return err
	}
	if !sleepContext(ctx, 500*time.Millisecond) {
		diagAutoCapture(logPath, "stream camera refresh cancelled: view_id=%q err=%v", viewID, ctx.Err())
		return ctx.Err()
	}
	diagAutoCapture(logPath, "stream camera refresh complete: view_id=%q wait_ms=%d", viewID, 500)
	return nil
}

func (r AutoCaptureRunner) recoverEmptySpoutSenderList(ctx context.Context, client oscClient, viewID string) error {
	logPath := r.Config.DiagnosticLogPath
	if r.Config.AutoCapture.Capture.PreplacedLocalAnchorEnabled() {
		diagAutoCapture(logPath, "spout sender recovery skipped: preplaced_local_anchor=true view_id=%q", viewID)
		return nil
	}
	list, err := autoCaptureListSpoutSenders(ctx, r.Config.AutoCapture.Stream, logPath)
	if err != nil {
		diagAutoCapture(logPath, "spout sender recovery skipped: list_error=%q view_id=%q", err.Error(), viewID)
		return nil
	}
	if len(list.Senders) > 0 {
		diagAutoCapture(logPath, "spout sender recovery skipped: senders=%d view_id=%q", len(list.Senders), viewID)
		return nil
	}
	diagAutoCapture(logPath, "spout sender recovery begin: reason=%q view_id=%q", "empty_sender_list", viewID)
	if err := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", false, "spout_sender_recovery_off:"+viewID); err != nil {
		diagAutoCapture(logPath, "spout sender recovery error: phase=%q view_id=%q err=%v", "off", viewID, err)
		return err
	}
	if !sleepContext(ctx, 300*time.Millisecond) {
		diagAutoCapture(logPath, "spout sender recovery cancelled: phase=%q view_id=%q err=%v", "off_wait", viewID, ctx.Err())
		return ctx.Err()
	}
	if err := sendCameraBoolCompat(client, logPath, "/usercamera/Streaming", true, "spout_sender_recovery_on:"+viewID); err != nil {
		diagAutoCapture(logPath, "spout sender recovery error: phase=%q view_id=%q err=%v", "on", viewID, err)
		return err
	}
	if !sleepContext(ctx, 800*time.Millisecond) {
		diagAutoCapture(logPath, "spout sender recovery cancelled: phase=%q view_id=%q err=%v", "on_wait", viewID, ctx.Err())
		return ctx.Err()
	}
	after, err := autoCaptureListSpoutSenders(ctx, r.Config.AutoCapture.Stream, logPath)
	if err != nil {
		diagAutoCapture(logPath, "spout sender recovery recheck error: view_id=%q err=%v", viewID, err)
		return nil
	}
	diagAutoCapture(logPath, "spout sender recovery complete: view_id=%q before_senders=%d after_senders=%d", viewID, len(list.Senders), len(after.Senders))
	return nil
}

func (r AutoCaptureRunner) waitCaptureDelay(ctx context.Context, view CameraViewConfig, name string) bool {
	delay := time.Duration(view.CaptureDelayMS) * time.Millisecond
	if delay <= 0 {
		return true
	}
	logPath := r.Config.DiagnosticLogPath
	diagAutoCapture(logPath, "shot capture_delay begin: view_id=%q duration_ms=%d", view.ID, delay.Milliseconds())
	if !sleepContext(ctx, delay) {
		diagAutoCapture(logPath, "shot capture_delay cancelled: view_id=%q view_name=%q err=%v", view.ID, name, ctx.Err())
		return false
	}
	diagAutoCapture(logPath, "shot capture_delay complete: view_id=%q", view.ID)
	return true
}

func (r AutoCaptureRunner) finalizeAutoCaptureImage(photoPath string, batchID string, shotID string, view CameraViewConfig, sidecarUsers []PresenceUser, discordUsers []PresenceUser, confidence string, world AutoCaptureVRChatMetadata, streamInfo SpoutCaptureResult) Result {
	cfg := r.Config.AutoCapture
	logPath := r.Config.DiagnosticLogPath
	resolvedPose, err := ResolveCameraViewPose(cfg, view)
	metadataWarnings := []string{}
	if cfg.Output.WriteEXIF {
		metadata := BuildAutoCaptureEmbeddedMetadata(cfg, batchID, shotID, view, discordUsers, confidence, streamInfo)
		if err == nil {
			metadata.ResolvedPose = &resolvedPose
		} else {
			diagAutoCapture(logPath, "embedded metadata resolved pose skipped: image=%q view_id=%q coordinate_space=%q err=%v", photoPath, view.ID, view.CoordinateSpace, err)
		}
		warnings, writeErr := WriteAutoCaptureEmbeddedMetadataWithWarnings(photoPath, metadata)
		metadataWarnings = append(metadataWarnings, warnings...)
		if writeErr != nil {
			warning := "埋め込みメタデータを書き込めませんでした: " + writeErr.Error()
			metadataWarnings = append(metadataWarnings, warning)
			diagAutoCapture(logPath, "embedded metadata write warning: image=%q err=%v", photoPath, writeErr)
		} else {
			diagAutoCapture(logPath, "embedded metadata write success: image=%q users=%d include_ids=%t warnings=%d", photoPath, len(metadata.Users), cfg.Output.WriteUserIDsToEXIF, len(metadataWarnings))
		}
	}
	if cfg.Output.WriteSidecarJSON {
		basisSource, basisPose, basisUpdatedAt := autoCapturePlayerLocalBasisMetadata(cfg)
		sidecar := AutoCaptureSidecar{
			SchemaVersion:             1,
			BatchID:                   batchID,
			ShotID:                    shotID,
			CapturedAtLocal:           time.Now().Format(time.RFC3339),
			CapturedAtUTC:             time.Now().UTC().Format(time.RFC3339),
			CaptureMode:               cfg.Capture.Mode,
			View:                      view,
			PlayerLocalBasisSource:    basisSource,
			PlayerLocalBasisPose:      basisPose,
			PlayerLocalBasisUpdatedAt: basisUpdatedAt,
			Stream:                    autoCaptureStreamMetadata(streamInfo),
			VRChat:                    autoCaptureVRChatMetadata(world, confidence),
			Users:                     sidecarUsers,
			MetadataWarnings:          metadataWarnings,
		}
		if err == nil {
			sidecar.ResolvedPose = &resolvedPose
		} else {
			diagAutoCapture(logPath, "sidecar resolved pose skipped: image=%q view_id=%q coordinate_space=%q err=%v", photoPath, view.ID, view.CoordinateSpace, err)
		}
		if err := WriteAutoCaptureSidecar(photoPath, sidecar); err != nil {
			diagAutoCapture(logPath, "sidecar write error: image=%q err=%v", photoPath, err)
		} else {
			diagAutoCapture(logPath, "sidecar write success: image=%q users=%d", photoPath, len(sidecarUsers))
		}
	} else {
		diagAutoCapture(logPath, "sidecar write skipped: disabled image=%q", photoPath)
	}
	result := Result{SourcePath: photoPath, OutputPath: photoPath, Name: filepath.Base(photoPath), Warnings: metadataWarnings}
	if cfg.Discord.Enabled {
		webhook := cfg.Discord.WebhookURL
		if strings.TrimSpace(webhook) == "" {
			webhook = r.Config.Discord.WebhookURL
		}
		diagAutoCapture(logPath, "discord upload begin: image=%q webhook_configured=%t users=%d", photoPath, strings.TrimSpace(webhook) != "", len(discordUsers))
		content := autoCaptureDiscordContent(cfg, view, discordUsers)
		var uploaded DiscordUpload
		var err error
		if cfg.Discord.IncludeImages {
			uploaded, err = uploadAutoCaptureDiscord(webhook, photoPath, content)
		} else {
			uploaded, err = PostDiscordContent(webhook, content)
		}
		if err != nil {
			diagAutoCapture(logPath, "discord upload error: image=%q err=%v", photoPath, err)
			result.Error = err.Error()
			return result
		}
		diagAutoCapture(logPath, "discord upload success: image=%q message_id=%q", photoPath, uploaded.MessageID)
		result.URL = uploaded.URL
		result.DiscordMessageID = uploaded.MessageID
		result.DiscordWebhookID = uploaded.WebhookID
		result.DiscordToken = uploaded.Token
	}
	return result
}

func (r AutoCaptureRunner) emit(event AutoCaptureEvent) {
	if r.Handler != nil {
		r.Handler(event)
	}
}

type AutoCaptureSidecar struct {
	SchemaVersion             int                        `json:"schema_version"`
	BatchID                   string                     `json:"batch_id"`
	ShotID                    string                     `json:"shot_id"`
	CapturedAtLocal           string                     `json:"captured_at_local"`
	CapturedAtUTC             string                     `json:"captured_at_utc"`
	CaptureMode               string                     `json:"capture_mode"`
	View                      CameraViewConfig           `json:"view"`
	ResolvedPose              *CameraPoseConfig          `json:"resolved_pose,omitempty"`
	PlayerLocalBasisSource    string                     `json:"player_local_basis_source,omitempty"`
	PlayerLocalBasisPose      *CameraPoseConfig          `json:"player_local_basis_pose,omitempty"`
	PlayerLocalBasisUpdatedAt string                     `json:"player_local_basis_updated_at,omitempty"`
	Stream                    *AutoCaptureStreamMetadata `json:"stream,omitempty"`
	VRChat                    AutoCaptureVRChatMetadata  `json:"vrchat"`
	Users                     []PresenceUser             `json:"users"`
	Files                     AutoCaptureFileMetadata    `json:"files"`
	MetadataWarnings          []string                   `json:"metadata_warnings,omitempty"`
}

type AutoCaptureStreamMetadata struct {
	Backend    string `json:"backend"`
	SenderName string `json:"sender_name,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Frame      int64  `json:"frame,omitempty"`
	CapturedAt string `json:"captured_at,omitempty"`
}

type AutoCaptureVRChatMetadata struct {
	WorldID         string `json:"world_id,omitempty"`
	InstanceID      string `json:"instance_id,omitempty"`
	UsersSource     string `json:"users_source,omitempty"`
	UsersConfidence string `json:"users_confidence,omitempty"`
}

type AutoCaptureFileMetadata struct {
	ImagePath string `json:"image_path"`
	SHA256    string `json:"sha256"`
}

func WriteAutoCaptureSidecar(imagePath string, sidecar AutoCaptureSidecar) error {
	sum, err := fileSHA256(imagePath)
	if err == nil {
		sidecar.Files.SHA256 = sum
	}
	sidecar.Files.ImagePath = imagePath
	data, err := json.MarshalIndent(sidecar, "", "  ")
	if err != nil {
		return err
	}
	return WritePrivateFile(imagePath+".json", append(data, '\n'))
}

func autoCaptureOutputPath(cfg AutoCaptureConfig, batchID string, shotID string, index int, view CameraViewConfig) (string, error) {
	dir := strings.TrimSpace(cfg.Output.Directory)
	if dir == "" {
		dir = DefaultAutoCaptureDirectory()
	}
	if err := os.MkdirAll(dir, privateDirMode); err != nil {
		return "", err
	}
	ext := strings.ToLower(strings.TrimPrefix(cfg.Output.ImageFormat, "."))
	if ext == "jpeg" {
		ext = "jpg"
	}
	if ext == "" {
		ext = "png"
	}
	name := cfg.Output.FilenameTemplate
	if strings.TrimSpace(name) == "" {
		name = "{timestamp_local}_{batch_id}_{shot_index}_{view_name}_{mode}.{ext}"
	}
	replacements := map[string]string{
		"{timestamp_local}": time.Now().Format("20060102_150405"),
		"{batch_id}":        batchID,
		"{shot_id}":         shotID,
		"{shot_index}":      fmt.Sprintf("%02d", index),
		"{view_id}":         safeFilenamePart(view.ID),
		"{view_name}":       safeFilenamePart(view.Name),
		"{mode}":            safeFilenamePart(cfg.Capture.Mode),
		"{ext}":             ext,
	}
	for old, value := range replacements {
		name = strings.ReplaceAll(name, old, value)
	}
	name = safeFilenamePart(name)
	if filepath.Ext(name) == "" {
		name += "." + ext
	}
	return filepath.Join(dir, name), nil
}

func captureStreamFrameWithFFmpeg(ctx context.Context, cfg AutoCaptureStreamConfig, outputPath string, logPath string) error {
	timeout := time.Duration(cfg.CaptureTimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if _, err := ResolveFFmpegPath(cfg.LegacyFFmpegPath); err != nil {
		diagAutoCapture(logPath, "stream ffmpeg missing: path=%q err=%v", cfg.LegacyFFmpegPath, err)
		return err
	}
	commandCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	args, err := splitCommandLine(cfg.LegacyInputArgs)
	if err != nil {
		return fmt.Errorf("ffmpeg入力引数を解釈できません: %w", err)
	}
	args, err = expandFFmpegInputPlaceholders(ctx, args, logPath)
	if err != nil {
		return err
	}
	outputInserted := false
	for i := range args {
		if strings.Contains(args[i], "{output}") {
			args[i] = strings.ReplaceAll(args[i], "{output}", outputPath)
			outputInserted = true
		}
	}
	if !outputInserted {
		args = append([]string{"-y"}, args...)
		args = append(args, "-frames:v", "1", outputPath)
	}
	diagAutoCapture(logPath, "stream ffmpeg begin: path=%q args=%q output=%q timeout_ms=%d", cfg.LegacyFFmpegPath, strings.Join(args, " "), outputPath, timeout.Milliseconds())
	cmd := exec.CommandContext(commandCtx, cfg.LegacyFFmpegPath, args...) // #nosec G204 -- user-configured local ffmpeg command for capture source.
	output, err := cmd.CombinedOutput()
	if commandCtx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("ffmpegによるStream切り出しがタイムアウトしました。ffmpegパス、入力引数、Stream Cameraの表示状態を確認してください。")
	}
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if len(trimmed) > 600 {
			trimmed = trimmed[len(trimmed)-600:]
		}
		return fmt.Errorf("ffmpegによるStream切り出しに失敗しました: %v %s", err, trimmed)
	}
	if !fileHasContent(outputPath) {
		return fmt.Errorf("ffmpegは終了しましたが画像が作成されませんでした。入力引数と出力先を確認してください。")
	}
	diagAutoCapture(logPath, "stream ffmpeg success: output=%q", outputPath)
	return nil
}

type windowRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func expandFFmpegInputPlaceholders(ctx context.Context, args []string, logPath string) ([]string, error) {
	needsWindow := false
	for _, arg := range args {
		if strings.Contains(arg, "{window_") {
			needsWindow = true
			break
		}
	}
	if !needsWindow {
		return args, nil
	}
	rect, err := detectVRChatWindowRect(ctx, logPath)
	if err != nil {
		return nil, err
	}
	replacements := map[string]string{
		"{window_x}":      fmt.Sprintf("%d", rect.X),
		"{window_y}":      fmt.Sprintf("%d", rect.Y),
		"{window_width}":  fmt.Sprintf("%d", rect.Width),
		"{window_height}": fmt.Sprintf("%d", rect.Height),
	}
	out := make([]string, len(args))
	for i, arg := range args {
		for old, value := range replacements {
			arg = strings.ReplaceAll(arg, old, value)
		}
		out[i] = arg
	}
	diagAutoCapture(logPath, "stream window rect resolved: x=%d y=%d width=%d height=%d", rect.X, rect.Y, rect.Width, rect.Height)
	return out, nil
}

func detectVRChatWindowRect(ctx context.Context, logPath string) (windowRect, error) {
	var rect windowRect
	if strings.TrimSpace(os.Getenv("OS")) != "Windows_NT" && filepath.Separator != '\\' {
		return rect, fmt.Errorf("VRChatウィンドウ範囲の自動取得はWindowsでのみ利用できます。ffmpeg入力引数を手動指定してください。")
	}
	commandCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	script := `$ErrorActionPreference='Stop'; Add-Type @"
using System;
using System.Runtime.InteropServices;
public static class Win32Rect {
  [StructLayout(LayoutKind.Sequential)] public struct RECT { public int Left; public int Top; public int Right; public int Bottom; }
  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr hWnd, out RECT rect);
}
"@; $p = Get-Process VRChat -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1; if ($null -eq $p) { throw 'VRChatのウィンドウが見つかりません。VRChatを起動し、最小化を解除してください。' }; $r = New-Object Win32Rect+RECT; if (-not [Win32Rect]::GetWindowRect($p.MainWindowHandle, [ref]$r)) { throw 'VRChatウィンドウの位置を取得できません。' }; [pscustomobject]@{x=$r.Left;y=$r.Top;width=($r.Right-$r.Left);height=($r.Bottom-$r.Top)} | ConvertTo-Json -Compress`
	cmd := exec.CommandContext(commandCtx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script) // #nosec G204 -- fixed Windows window-rect query.
	output, err := cmd.CombinedOutput()
	if commandCtx.Err() == context.DeadlineExceeded {
		return rect, fmt.Errorf("VRChatウィンドウ位置の取得がタイムアウトしました。VRChatが起動しているか確認してください。")
	}
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		diagAutoCapture(logPath, "stream window rect error: err=%v output=%q", err, trimmed)
		if trimmed != "" {
			return rect, fmt.Errorf("VRChatウィンドウ位置を取得できません: %s", trimmed)
		}
		return rect, fmt.Errorf("VRChatウィンドウ位置を取得できません: %v", err)
	}
	if err := json.Unmarshal(output, &rect); err != nil {
		diagAutoCapture(logPath, "stream window rect parse error: err=%v output=%q", err, strings.TrimSpace(string(output)))
		return rect, fmt.Errorf("VRChatウィンドウ位置の取得結果を解釈できません: %w", err)
	}
	if rect.Width <= 0 || rect.Height <= 0 {
		return rect, fmt.Errorf("VRChatウィンドウサイズが不正です: %dx%d", rect.Width, rect.Height)
	}
	return rect, nil
}

func ResolveFFmpegPath(ffmpegPath string) (string, error) {
	ffmpegPath = strings.Trim(strings.TrimSpace(ffmpegPath), `"`)
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	if strings.ContainsAny(ffmpegPath, `\/`) || filepath.IsAbs(ffmpegPath) {
		info, err := os.Stat(ffmpegPath)
		if err != nil {
			return "", fmt.Errorf("ffmpegがインストールされていないかPATHにありません。ffmpegパスを確認するか、設定画面からffmpegをインストールしてください。")
		}
		if info.IsDir() {
			return "", fmt.Errorf("ffmpegパスがフォルダを指しています。ffmpeg.exeのパスを指定してください。")
		}
		return ffmpegPath, nil
	}
	resolved, err := exec.LookPath(ffmpegPath)
	if err != nil {
		return "", fmt.Errorf("ffmpegがインストールされていないかPATHにありません。ffmpegパスを確認するか、設定画面からffmpegをインストールしてください。")
	}
	return resolved, nil
}

func splitCommandLine(input string) ([]string, error) {
	args := make([]string, 0)
	var current strings.Builder
	var quote rune
	escaped := false
	for _, r := range input {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if escaped {
		current.WriteRune('\\')
	}
	if quote != 0 {
		return nil, fmt.Errorf("クォートが閉じられていません")
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args, nil
}

func safeFilenamePart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		value = "capture"
	}
	replacer := strings.NewReplacer("\\", "_", "/", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_")
	value = replacer.Replace(value)
	value = strings.Trim(value, ". ")
	if value == "" {
		return "capture"
	}
	return value
}

func SnapshotVRChatPresence(logDir string) ([]PresenceUser, string) {
	users, confidence, _ := SnapshotVRChatPresenceWithSource(logDir)
	return users, confidence
}

func SnapshotVRChatPresenceWithSource(logDir string) ([]PresenceUser, string, string) {
	path := latestVRChatOutputLog(logDir)
	if path == "" {
		return nil, "unknown", ""
	}
	users, ok := parseVRChatPresenceLog(path)
	if !ok {
		return nil, "unknown", path
	}
	out := make([]PresenceUser, 0, len(users))
	for _, user := range users {
		out = append(out, user)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].DisplayName) < strings.ToLower(out[j].DisplayName)
	})
	return out, "partial", path
}

func SnapshotVRChatWorld(logDir string) AutoCaptureVRChatMetadata {
	path := latestVRChatOutputLog(logDir)
	if path == "" {
		return AutoCaptureVRChatMetadata{}
	}
	data, err := os.ReadFile(path) // #nosec G304 -- VRChat output log path is user configured.
	if err != nil {
		return AutoCaptureVRChatMetadata{}
	}
	return parseVRChatWorldMetadata(string(data))
}

func parseVRChatWorldMetadata(logText string) AutoCaptureVRChatMetadata {
	worldRe := regexp.MustCompile(`wrld_[0-9A-Za-z-]+(?::[0-9A-Za-z~._:()%-]+)?`)
	matches := worldRe.FindAllString(logText, -1)
	if len(matches) == 0 {
		return AutoCaptureVRChatMetadata{}
	}
	last := trimVRChatWorldToken(matches[len(matches)-1])
	parts := strings.SplitN(last, ":", 2)
	meta := AutoCaptureVRChatMetadata{WorldID: parts[0]}
	if len(parts) == 2 {
		meta.InstanceID = parts[1]
	}
	return meta
}

func trimVRChatWorldToken(value string) string {
	return strings.TrimRight(value, ".,;]")
}

func presenceUsersWithoutIDs(users []PresenceUser) []PresenceUser {
	out := make([]PresenceUser, len(users))
	copy(out, users)
	for i := range out {
		out[i].UserID = ""
		if out[i].Confidence == "confirmed" {
			out[i].Confidence = "partial"
		}
	}
	return out
}

func autoCaptureSidecarUsers(cfg AutoCaptureConfig, users []PresenceUser) []PresenceUser {
	if cfg.Presence.IncludeUserIDsInSidecar {
		return users
	}
	return presenceUsersWithoutIDs(users)
}

func autoCaptureVRChatMetadata(world AutoCaptureVRChatMetadata, confidence string) AutoCaptureVRChatMetadata {
	world.UsersSource = "output_log"
	world.UsersConfidence = confidence
	return world
}

func autoCaptureDiscordContent(cfg AutoCaptureConfig, view CameraViewConfig, users []PresenceUser) string {
	lines := []string{
		"VRChat自動撮影",
		"構図: " + view.Name,
		"撮影方式: " + cfg.Capture.Mode,
	}
	if cfg.Presence.IncludeDisplayNamesInDiscord || cfg.Presence.IncludeUserIDsInDiscord {
		parts := make([]string, 0, len(users))
		for _, user := range users {
			name := strings.TrimSpace(user.DisplayName)
			if name == "" {
				name = "unknown"
			}
			if cfg.Presence.IncludeUserIDsInDiscord && user.UserID != "" {
				name += " (" + user.UserID + ")"
			}
			parts = append(parts, name)
		}
		if len(parts) > 0 {
			lines = append(lines, "同席ユーザー: "+strings.Join(parts, ", "))
		}
	}
	return strings.Join(lines, "\n")
}

func uploadAutoCaptureDiscord(webhookURL string, imagePath string, content string) (DiscordUpload, error) {
	var uploaded DiscordUpload
	data, err := os.ReadFile(imagePath) // #nosec G304 -- captured image path comes from configured VRChat photo directory.
	if err != nil {
		return uploaded, err
	}
	ext := strings.ToLower(filepath.Ext(imagePath))
	mime := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".webp":
		mime = "image/webp"
	case ".png":
	default:
		ext = ".png"
	}
	return UploadDiscordWithContent(webhookURL, filepath.Base(imagePath), EncodedImage{Extension: ext, Mime: mime, Data: data}, content)
}

func autoCaptureStreamMetadata(result SpoutCaptureResult) *AutoCaptureStreamMetadata {
	if strings.TrimSpace(result.SenderName) == "" && result.Width == 0 && result.Height == 0 && result.Frame == 0 {
		return nil
	}
	return &AutoCaptureStreamMetadata{
		Backend:    "spout",
		SenderName: result.SenderName,
		Width:      result.Width,
		Height:     result.Height,
		Frame:      result.Frame,
		CapturedAt: result.CapturedAt,
	}
}

func convertAutoCaptureImage(sourcePath string, outputPath string, outputFormat string) error {
	img, _, err := DecodeImageFile(sourcePath)
	if err != nil {
		return err
	}
	encoded, err := EncodeImage(img, outputFormat, 92)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return WritePrivateFile(outputPath, encoded.Data)
}

func latestVRChatOutputLog(dir string) string {
	matches := make([]string, 0)
	for _, candidate := range vrchatLogDirectoryCandidates(dir) {
		found, err := filepath.Glob(filepath.Join(candidate, "output_log_*.txt"))
		if err == nil {
			matches = append(matches, found...)
		}
	}
	if len(matches) == 0 {
		return ""
	}
	sort.Slice(matches, func(i, j int) bool {
		ii, _ := os.Stat(matches[i])
		jj, _ := os.Stat(matches[j])
		if ii == nil || jj == nil {
			return matches[i] > matches[j]
		}
		return ii.ModTime().After(jj.ModTime())
	})
	return matches[0]
}

var (
	usrIDPattern           = regexp.MustCompile(`usr_[0-9a-fA-F-]{36}`)
	joinPattern            = regexp.MustCompile(`(?i)(?:OnPlayerJoined|Joining|joined).*?((?:usr_[0-9a-fA-F-]{36})|$)`)
	leavePattern           = regexp.MustCompile(`(?i)(?:OnPlayerLeft|Leaving|left).*?((?:usr_[0-9a-fA-F-]{36})|$)`)
	namePattern            = regexp.MustCompile(`\b(?:displayName|userName|name)[:=]\s*"?([^",\]]+)`)
	playerEventNamePattern = regexp.MustCompile(`(?i)OnPlayer(?:Joined|Left)\s+(.+?)\s+\(usr_[0-9a-fA-F-]{36}\)`)
)

func vrchatLogDirectoryCandidates(dir string) []string {
	dir = strings.Trim(strings.TrimSpace(dir), `"`)
	if dir == "" {
		dir = DefaultVRChatLogDirectory()
	}
	candidates := []string{}
	add := func(value string) {
		value = strings.Trim(strings.TrimSpace(value), `"`)
		if value == "" {
			return
		}
		cleaned := filepath.Clean(value)
		for _, existing := range candidates {
			if strings.EqualFold(existing, cleaned) {
				return
			}
		}
		candidates = append(candidates, cleaned)
	}
	add(dir)
	normalized := strings.ReplaceAll(filepath.ToSlash(dir), "/Local/Low/", "/LocalLow/")
	add(filepath.FromSlash(normalized))
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		add(filepath.Join(userProfile, "AppData", "LocalLow", "VRChat", "VRChat"))
	}
	add(DefaultVRChatLogDirectory())
	return candidates
}

func parseVRChatPresenceLog(path string) (map[string]PresenceUser, bool) {
	data, err := os.ReadFile(path) // #nosec G304 -- VRChat output log path is configured by the local user.
	if err != nil {
		return nil, false
	}
	users := map[string]PresenceUser{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Entering Room") || strings.Contains(line, "Joining wrld_") {
			users = map[string]PresenceUser{}
			continue
		}
		id := usrIDPattern.FindString(line)
		name := extractPresenceName(line)
		switch {
		case joinPattern.MatchString(line):
			key := id
			if key == "" {
				key = name
			}
			if key != "" {
				users[key] = PresenceUser{DisplayName: name, UserID: id, Status: "present", Source: "output_log", Confidence: presenceConfidence(id), JoinedAt: extractLogTime(line)}
			}
		case leavePattern.MatchString(line):
			key := id
			if key == "" {
				key = name
			}
			delete(users, key)
		}
	}
	return users, true
}

func extractPresenceName(line string) string {
	if match := playerEventNamePattern.FindStringSubmatch(line); len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	if match := namePattern.FindStringSubmatch(line); len(match) == 2 {
		name := strings.TrimSpace(match[1])
		if idx := strings.Index(name, " usr_"); idx >= 0 {
			name = strings.TrimSpace(name[:idx])
		}
		return name
	}
	if idx := strings.LastIndex(line, "usr_"); idx > 0 {
		prefix := strings.TrimSpace(line[:idx])
		fields := strings.Fields(prefix)
		if len(fields) > 0 {
			return strings.Trim(fields[len(fields)-1], `"'[]():`)
		}
	}
	return ""
}

func presenceConfidence(id string) string {
	if id == "" {
		return "partial"
	}
	return "confirmed"
}

func extractLogTime(line string) string {
	if len(line) >= 19 {
		prefix := line[:19]
		if _, err := time.Parse("2006.01.02 15:04:05", prefix); err == nil {
			return prefix
		}
	}
	return ""
}

func enabledCameraViews(views []CameraViewConfig) []CameraViewConfig {
	out := make([]CameraViewConfig, 0, len(views))
	for _, view := range views {
		if view.Enabled {
			out = append(out, view)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].SortOrder < out[j].SortOrder
	})
	return out
}

func waitForNewPhoto(ctx context.Context, dir string, before map[string]time.Time, timeout time.Duration, notBefore time.Time, logPath string) string {
	if strings.TrimSpace(dir) == "" {
		diagAutoCapture(logPath, "photo wait skipped: dir_empty=true timeout_ms=%d", timeout.Milliseconds())
		return ""
	}
	initial, status := scanPhotoFilesWithStatus(dir)
	diagAutoCapture(logPath, "photo wait begin: dir=%q timeout_ms=%d not_before=%s before_files=%d current_files=%d current_latest=%s new_candidates=%d scan_error=%q limit_reached=%t",
		dir,
		timeout.Milliseconds(),
		notBefore.Format(time.RFC3339Nano),
		len(before),
		len(initial),
		photoFileSummary(initial),
		len(newPhotoCandidates(initial, before, notBefore)),
		status.Error,
		status.LimitReached,
	)
	deadline := time.Now().Add(timeout)
	var latestCandidate string
	for time.Now().Before(deadline) {
		current := scanPhotoFiles(dir)
		paths := newPhotoCandidates(current, before, notBefore)
		if len(paths) > 0 {
			latestCandidate = paths[0]
		}
		for _, path := range paths {
			if fileLooksStable(path) {
				diagAutoCapture(logPath, "photo wait found: dir=%q path=%q current_files=%d new_candidates=%d reason=%q", dir, path, len(current), len(paths), "stable")
				return path
			}
		}
		if !sleepContext(ctx, 500*time.Millisecond) {
			diagAutoCapture(logPath, "photo wait cancelled: dir=%q err=%v current_files=%d current_latest=%s new_candidates=%d latest_candidate=%q", dir, ctx.Err(), len(current), photoFileSummary(current), len(paths), latestCandidate)
			return ""
		}
	}
	current, status := scanPhotoFilesWithStatus(dir)
	paths := newPhotoCandidates(current, before, notBefore)
	if len(paths) > 0 {
		for _, path := range paths {
			if fileHasContent(path) {
				diagAutoCapture(logPath, "photo wait found: dir=%q path=%q current_files=%d new_candidates=%d reason=%q", dir, path, len(current), len(paths), "timeout_candidate")
				return path
			}
		}
	}
	diagAutoCapture(logPath, "photo wait timeout: dir=%q timeout_ms=%d not_before=%s before_files=%d current_files=%d current_latest=%s new_candidates=%d latest_candidate=%q scan_error=%q limit_reached=%t",
		dir,
		timeout.Milliseconds(),
		notBefore.Format(time.RFC3339Nano),
		len(before),
		len(current),
		photoFileSummary(current),
		len(paths),
		latestCandidate,
		status.Error,
		status.LimitReached,
	)
	return ""
}

func scanAutoCapturePhotoFiles(photoDir string, outputDir string) map[string]time.Time {
	files := scanPhotoFiles(photoDir)
	for path, modTime := range scanPhotoFiles(outputDir) {
		files[path] = modTime
	}
	return files
}

func photoFileSummary(files map[string]time.Time) string {
	paths := photoPathsByModTimeDesc(files)
	if len(paths) == 0 {
		return "none"
	}
	path := paths[0]
	return fmt.Sprintf("%q@%s", path, files[path].Format(time.RFC3339))
}

func newPhotoCandidates(files map[string]time.Time, before map[string]time.Time, notBefore time.Time) []string {
	paths := make([]string, 0)
	for path, modTime := range files {
		if _, ok := before[path]; ok {
			continue
		}
		if !notBefore.IsZero() && modTime.Before(notBefore) {
			continue
		}
		paths = append(paths, path)
	}
	sort.SliceStable(paths, func(i, j int) bool {
		left := files[paths[i]]
		right := files[paths[j]]
		if left.Equal(right) {
			return paths[i] > paths[j]
		}
		return left.After(right)
	})
	return paths
}

func photoPathsByModTimeDesc(files map[string]time.Time) []string {
	paths := sortedPhotoPaths(files)
	sort.SliceStable(paths, func(i, j int) bool {
		left := files[paths[i]]
		right := files[paths[j]]
		if left.Equal(right) {
			return paths[i] > paths[j]
		}
		return left.After(right)
	})
	return paths
}

func fileHasContent(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Size() > 0
}

func autoCapturePhotoDirectory(cfg Config) string {
	photoDir := strings.TrimSpace(cfg.AutoPhoto.PhotoDirectory)
	if photoDir == "" {
		return DefaultVRChatPhotoDirectory()
	}
	return photoDir
}

func newBatchID(t time.Time) string {
	return "batch-" + t.Format("20060102-150405")
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- hashing the captured local image selected by the app.
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func diagAutoCapture(path string, format string, args ...any) {
	AppendDiagnosticLog(path, "auto-capture "+format, args...)
}

func sendCameraButton(ctx context.Context, client oscClient, address string, releaseDelayMS int, logPath string, detail string) error {
	if releaseDelayMS < 1 {
		releaseDelayMS = 1
	}
	diagAutoCapture(logPath, "osc button press begin: address=%q detail=%q", address, detail)
	if err := client.sendBool(address, true); err != nil {
		diagAutoCapture(logPath, "osc button press error: address=%q detail=%q err=%v", address, detail, err)
		return err
	}
	diagAutoCapture(logPath, "osc button press success: address=%q detail=%q", address, detail)
	released := false
	defer func() {
		if released {
			return
		}
		diagAutoCapture(logPath, "osc button release begin: address=%q detail=%q reason=%q", address, detail, "cleanup")
		if err := client.sendBool(address, false); err != nil {
			diagAutoCapture(logPath, "osc button release error: address=%q detail=%q reason=%q err=%v", address, detail, "cleanup", err)
			return
		}
		diagAutoCapture(logPath, "osc button release success: address=%q detail=%q reason=%q", address, detail, "cleanup")
	}()
	diagAutoCapture(logPath, "button_release wait begin: address=%q detail=%q duration_ms=%d", address, detail, releaseDelayMS)
	if !sleepContext(ctx, time.Duration(releaseDelayMS)*time.Millisecond) {
		diagAutoCapture(logPath, "button_release wait cancelled: address=%q detail=%q err=%v", address, detail, ctx.Err())
		return ctx.Err()
	}
	diagAutoCapture(logPath, "osc button release begin: address=%q detail=%q", address, detail)
	if err := client.sendBool(address, false); err != nil {
		diagAutoCapture(logPath, "osc button release error: address=%q detail=%q err=%v", address, detail, err)
		return err
	}
	diagAutoCapture(logPath, "osc button release success: address=%q detail=%q", address, detail)
	released = true
	return nil
}

func sendCameraBoolCompat(client oscClient, logPath string, address string, value bool, detail string) error {
	diagAutoCapture(logPath, "osc bool compat send begin: address=%q value=%t detail=%q", address, value, detail)
	if err := client.sendBool(address, value); err != nil {
		diagAutoCapture(logPath, "osc bool compat send error: address=%q value=%t detail=%q err=%v", address, value, detail, err)
		return err
	}
	diagAutoCapture(logPath, "osc bool compat send success: address=%q value=%t detail=%q", address, value, detail)
	intValue := boolOSCInt(value)
	diagAutoCapture(logPath, "osc bool compat int send begin: address=%q value=%d detail=%q", address, intValue, detail)
	if err := client.sendInt(address, intValue); err != nil {
		diagAutoCapture(logPath, "osc bool compat int send error: address=%q value=%d detail=%q err=%v", address, intValue, detail, err)
		return err
	}
	diagAutoCapture(logPath, "osc bool compat int send success: address=%q value=%d detail=%q", address, intValue, detail)
	return nil
}

func boolOSCInt(value bool) int32 {
	if value {
		return 1
	}
	return 0
}

func sendOptionalFloat(client oscClient, address string, value *float64) int {
	if value != nil {
		_ = client.sendFloat(address, float32(*value))
		return 1
	}
	return 0
}

func sendOptionalBool(client oscClient, address string, value *bool) int {
	if value != nil {
		_ = client.sendBool(address, *value)
		return 1
	}
	return 0
}

func ParseOSCPacket(packet []byte) (string, string, []byte, bool) {
	address, next, ok := readOSCString(packet, 0)
	if !ok {
		return "", "", nil, false
	}
	typeTags, next, ok := readOSCString(packet, next)
	if !ok {
		return "", "", nil, false
	}
	return address, typeTags, packet[next:], true
}

func ParseOSCPose(packet []byte) (CameraPoseConfig, bool) {
	address, typeTags, payload, ok := ParseOSCPacket(packet)
	if !ok || address != "/usercamera/Pose" {
		return CameraPoseConfig{}, false
	}
	typeTags = strings.TrimPrefix(typeTags, ",")
	if len(typeTags) < 6 {
		return CameraPoseConfig{}, false
	}
	values := make([]float32, 0, 6)
	offset := 0
	for _, tag := range typeTags {
		if len(values) == 6 {
			break
		}
		switch tag {
		case 'f':
			if offset+4 > len(payload) {
				return CameraPoseConfig{}, false
			}
			values = append(values, math.Float32frombits(binary.BigEndian.Uint32(payload[offset:offset+4])))
			offset += 4
		case 'i':
			if offset+4 > len(payload) {
				return CameraPoseConfig{}, false
			}
			values = append(values, float32(int32(binary.BigEndian.Uint32(payload[offset:offset+4]))))
			offset += 4
		default:
			return CameraPoseConfig{}, false
		}
	}
	if len(values) != 6 {
		return CameraPoseConfig{}, false
	}
	return CameraPoseConfig{
		Position: CameraVector3Config{X: float64(values[0]), Y: float64(values[1]), Z: float64(values[2])},
		Rotation: CameraVector3Config{X: float64(values[3]), Y: float64(values[4]), Z: float64(values[5])},
	}, true
}

func ParseOSCUserCameraSample(packet []byte) (UserCameraOSCSample, bool) {
	address, typeTags, payload, ok := ParseOSCPacket(packet)
	if !ok {
		return UserCameraOSCSample{}, false
	}
	return DecodeOSCUserCameraSample(address, typeTags, payload)
}

func DecodeOSCUserCameraSample(address string, typeTags string, payload []byte) (UserCameraOSCSample, bool) {
	if !strings.HasPrefix(address, "/usercamera/") {
		return UserCameraOSCSample{}, false
	}
	if address == "/usercamera/Pose" {
		pose, ok := decodeUserCameraPose(typeTags, payload)
		if !ok {
			return UserCameraOSCSample{}, false
		}
		return UserCameraOSCSample{Address: address, Pose: pose, HasPose: true}, true
	}
	if address == "/usercamera/Mode" {
		value, ok := decodeOSCFirstInt(typeTags, payload)
		if !ok {
			return UserCameraOSCSample{}, false
		}
		return UserCameraOSCSample{Address: address, Int: value, HasInt: true}, true
	}
	if userCameraBoolAddress(address) {
		value, ok := decodeOSCFirstBool(typeTags, payload)
		if !ok {
			return UserCameraOSCSample{}, false
		}
		return UserCameraOSCSample{Address: address, Bool: value, HasBool: true}, true
	}
	if userCameraFloatAddress(address) {
		value, ok := decodeOSCFirstFloat(typeTags, payload)
		if !ok {
			return UserCameraOSCSample{}, false
		}
		return UserCameraOSCSample{Address: address, Float: value, HasFloat: true}, true
	}
	return UserCameraOSCSample{}, false
}

func decodeUserCameraPose(typeTags string, payload []byte) (CameraPoseConfig, bool) {
	typeTags = strings.TrimPrefix(typeTags, ",")
	if len(typeTags) < 6 {
		return CameraPoseConfig{}, false
	}
	values := make([]float32, 0, 6)
	offset := 0
	for _, tag := range typeTags {
		if len(values) == 6 {
			break
		}
		switch tag {
		case 'f':
			if offset+4 > len(payload) {
				return CameraPoseConfig{}, false
			}
			values = append(values, math.Float32frombits(binary.BigEndian.Uint32(payload[offset:offset+4])))
			offset += 4
		case 'i':
			if offset+4 > len(payload) {
				return CameraPoseConfig{}, false
			}
			values = append(values, float32(int32(binary.BigEndian.Uint32(payload[offset:offset+4]))))
			offset += 4
		default:
			return CameraPoseConfig{}, false
		}
	}
	if len(values) != 6 {
		return CameraPoseConfig{}, false
	}
	return CameraPoseConfig{
		Position: CameraVector3Config{X: float64(values[0]), Y: float64(values[1]), Z: float64(values[2])},
		Rotation: CameraVector3Config{X: float64(values[3]), Y: float64(values[4]), Z: float64(values[5])},
	}, true
}

func decodeOSCFirstInt(typeTags string, payload []byte) (int, bool) {
	tag, ok := firstOSCTag(typeTags)
	if !ok {
		return 0, false
	}
	switch tag {
	case 'i':
		if len(payload) < 4 {
			return 0, false
		}
		return int(int32(binary.BigEndian.Uint32(payload[:4]))), true
	case 'f':
		if len(payload) < 4 {
			return 0, false
		}
		return int(math.Round(float64(math.Float32frombits(binary.BigEndian.Uint32(payload[:4]))))), true
	default:
		return 0, false
	}
}

func decodeOSCFirstFloat(typeTags string, payload []byte) (float64, bool) {
	tag, ok := firstOSCTag(typeTags)
	if !ok {
		return 0, false
	}
	switch tag {
	case 'f':
		if len(payload) < 4 {
			return 0, false
		}
		return float64(math.Float32frombits(binary.BigEndian.Uint32(payload[:4]))), true
	case 'i':
		if len(payload) < 4 {
			return 0, false
		}
		return float64(int32(binary.BigEndian.Uint32(payload[:4]))), true
	default:
		return 0, false
	}
}

func decodeOSCFirstBool(typeTags string, payload []byte) (bool, bool) {
	tag, ok := firstOSCTag(typeTags)
	if !ok {
		return false, false
	}
	switch tag {
	case 'T':
		return true, true
	case 'F':
		return false, true
	case 'i':
		if len(payload) < 4 {
			return false, false
		}
		return int32(binary.BigEndian.Uint32(payload[:4])) != 0, true
	case 'f':
		if len(payload) < 4 {
			return false, false
		}
		return math.Float32frombits(binary.BigEndian.Uint32(payload[:4])) != 0, true
	default:
		return false, false
	}
}

func firstOSCTag(typeTags string) (rune, bool) {
	typeTags = strings.TrimPrefix(typeTags, ",")
	for _, tag := range typeTags {
		return tag, true
	}
	return 0, false
}

func userCameraBoolAddress(address string) bool {
	switch address {
	case "/usercamera/ShowUIInCamera",
		"/usercamera/Lock",
		"/usercamera/LocalPlayer",
		"/usercamera/RemotePlayer",
		"/usercamera/Environment",
		"/usercamera/GreenScreen",
		"/usercamera/SmoothMovement",
		"/usercamera/LookAtMe",
		"/usercamera/AutoLevelRoll",
		"/usercamera/AutoLevelPitch",
		"/usercamera/Flying",
		"/usercamera/TriggerTakesPhotos",
		"/usercamera/DollyPathsStayVisible",
		"/usercamera/CameraEars",
		"/usercamera/ShowFocus",
		"/usercamera/Streaming",
		"/usercamera/RollWhileFlying",
		"/usercamera/OrientationIsLandscape":
		return true
	default:
		return false
	}
}

func userCameraFloatAddress(address string) bool {
	switch address {
	case "/usercamera/Zoom",
		"/usercamera/Exposure",
		"/usercamera/FocalDistance",
		"/usercamera/Aperture",
		"/usercamera/Hue",
		"/usercamera/Saturation",
		"/usercamera/Lightness",
		"/usercamera/LookAtMeXOffset",
		"/usercamera/LookAtMeYOffset",
		"/usercamera/FlySpeed",
		"/usercamera/TurnSpeed",
		"/usercamera/SmoothingStrength",
		"/usercamera/PhotoRate",
		"/usercamera/Duration":
		return true
	default:
		return false
	}
}

type oscClient struct {
	host string
	port int
	conn net.Conn
}

func (c *oscClient) open() error {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", c.host, c.port), 3*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *oscClient) close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c oscClient) sendInt(address string, value int32) error {
	return c.send(address, ",i", func(buf []byte) []byte {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], uint32(value))
		return append(buf, raw[:]...)
	})
}

func (c oscClient) sendFloat(address string, value float32) error {
	return c.send(address, ",f", func(buf []byte) []byte {
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], math.Float32bits(value))
		return append(buf, raw[:]...)
	})
}

func (c oscClient) sendFloats(address string, values []float32) error {
	types := "," + strings.Repeat("f", len(values))
	return c.send(address, types, func(buf []byte) []byte {
		for _, value := range values {
			var raw [4]byte
			binary.BigEndian.PutUint32(raw[:], math.Float32bits(value))
			buf = append(buf, raw[:]...)
		}
		return buf
	})
}

func (c oscClient) sendBool(address string, value bool) error {
	tag := ",F"
	if value {
		tag = ",T"
	}
	return c.send(address, tag, func(buf []byte) []byte { return buf })
}

func SendDebugOSCLine(cfg AutoCaptureOSCConfig, line string) (DebugOSCSendResult, error) {
	address, typeTags, appendArgs, err := parseDebugOSCLine(line)
	target := fmt.Sprintf("%s:%d", strings.TrimSpace(cfg.Host), cfg.SendPort)
	result := DebugOSCSendResult{
		OK:       false,
		Target:   target,
		Address:  address,
		TypeTags: typeTags,
	}
	if err != nil {
		result.Message = err.Error()
		return result, err
	}
	client := &oscClient{host: strings.TrimSpace(cfg.Host), port: cfg.SendPort}
	if client.host == "" {
		client.host = "127.0.0.1"
	}
	if client.port <= 0 || client.port > 65535 {
		err := fmt.Errorf("OSC送信先portが不正です: %d", client.port)
		result.Target = fmt.Sprintf("%s:%d", client.host, client.port)
		result.Message = err.Error()
		return result, err
	}
	result.Target = fmt.Sprintf("%s:%d", client.host, client.port)
	if err := client.open(); err != nil {
		result.Message = err.Error()
		return result, err
	}
	defer client.close()
	if err := client.send(address, typeTags, appendArgs); err != nil {
		result.Message = err.Error()
		return result, err
	}
	result.OK = true
	result.Message = "OSCを送信しました"
	return result, nil
}

func parseDebugOSCLine(line string) (string, string, func([]byte) []byte, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return "", "", nil, fmt.Errorf("OSC入力が空です")
	}
	address := fields[0]
	if !strings.HasPrefix(address, "/") {
		return "", "", nil, fmt.Errorf("OSC addressは / から始めてください")
	}
	args := make([]debugOSCArg, 0, len(fields)-1)
	typeTags := ","
	for _, raw := range fields[1:] {
		arg, err := parseDebugOSCArg(raw)
		if err != nil {
			return "", "", nil, err
		}
		args = append(args, arg)
		typeTags += string(arg.tag)
	}
	return address, typeTags, func(buf []byte) []byte {
		for _, arg := range args {
			switch arg.tag {
			case 'i':
				var raw [4]byte
				binary.BigEndian.PutUint32(raw[:], uint32(arg.intV))
				buf = append(buf, raw[:]...)
			case 'f':
				var raw [4]byte
				binary.BigEndian.PutUint32(raw[:], math.Float32bits(arg.float))
				buf = append(buf, raw[:]...)
			case 's':
				buf = appendOSCString(buf, arg.str)
			case 'T', 'F':
			}
		}
		return buf
	}, nil
}

func parseDebugOSCArg(raw string) (debugOSCArg, error) {
	value := strings.TrimSpace(raw)
	lower := strings.ToLower(value)
	var arg debugOSCArg
	switch {
	case strings.HasPrefix(lower, "i:"):
		parsed, err := strconv.ParseInt(value[2:], 10, 32)
		if err != nil {
			return arg, fmt.Errorf("int引数が不正です: %s", value)
		}
		arg.tag = 'i'
		arg.intV = int32(parsed)
		return arg, nil
	case strings.HasPrefix(lower, "f:"):
		parsed, err := strconv.ParseFloat(value[2:], 32)
		if err != nil {
			return arg, fmt.Errorf("float引数が不正です: %s", value)
		}
		arg.tag = 'f'
		arg.float = float32(parsed)
		return arg, nil
	case strings.HasPrefix(lower, "s:"):
		arg.tag = 's'
		arg.str = value[2:]
		return arg, nil
	case strings.HasPrefix(lower, "b:"):
		boolValue, ok := parseDebugOSCBool(value[2:])
		if !ok {
			return arg, fmt.Errorf("bool引数が不正です: %s", value)
		}
		if boolValue {
			arg.tag = 'T'
		} else {
			arg.tag = 'F'
		}
		return arg, nil
	}
	if boolValue, ok := parseDebugOSCBool(value); ok {
		if boolValue {
			arg.tag = 'T'
		} else {
			arg.tag = 'F'
		}
		return arg, nil
	}
	if parsed, err := strconv.ParseInt(value, 10, 32); err == nil {
		arg.tag = 'i'
		arg.intV = int32(parsed)
		return arg, nil
	}
	if parsed, err := strconv.ParseFloat(value, 32); err == nil {
		arg.tag = 'f'
		arg.float = float32(parsed)
		return arg, nil
	}
	arg.tag = 's'
	arg.str = value
	return arg, nil
}

func parseDebugOSCBool(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "t", "1", "on":
		return true, true
	case "false", "f", "0", "off":
		return false, true
	default:
		return false, false
	}
}

func (c oscClient) send(address string, typeTags string, appendArgs func([]byte) []byte) error {
	if c.conn == nil {
		return fmt.Errorf("OSC接続が開かれていません")
	}
	packet := buildOSCPacket(address, typeTags, appendArgs)
	_, err := c.conn.Write(packet)
	status := "ok"
	errText := ""
	if err != nil {
		status = "error"
		errText = err.Error()
	}
	_, _, payload, _ := ParseOSCPacket(packet)
	emitOSCTrace(OSCTraceEvent{
		Direction: "send",
		Address:   address,
		TypeTags:  typeTags,
		Payload:   append([]byte(nil), payload...),
		Target:    fmt.Sprintf("%s:%d", c.host, c.port),
		Status:    status,
		Error:     errText,
	})
	return err
}

func buildOSCPacket(address string, typeTags string, appendArgs func([]byte) []byte) []byte {
	packet := appendOSCString(nil, address)
	packet = appendOSCString(packet, typeTags)
	return appendArgs(packet)
}

func readOSCString(packet []byte, offset int) (string, int, bool) {
	if offset < 0 || offset >= len(packet) {
		return "", 0, false
	}
	end := offset
	for end < len(packet) && packet[end] != 0 {
		end++
	}
	if end >= len(packet) {
		return "", 0, false
	}
	next := end + 1
	for next%4 != 0 {
		next++
	}
	if next > len(packet) {
		return "", 0, false
	}
	return string(packet[offset:end]), next, true
}

func appendOSCString(buf []byte, value string) []byte {
	buf = append(buf, []byte(value)...)
	buf = append(buf, 0)
	for len(buf)%4 != 0 {
		buf = append(buf, 0)
	}
	return buf
}
