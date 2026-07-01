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
	"strings"
	"time"
)

type AutoCaptureEvent struct {
	BatchID string `json:"batchId"`
	ShotID  string `json:"shotId"`
	Path    string `json:"path"`
	Error   string `json:"error"`
	Message string `json:"message"`
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
	for _, address := range []string{"/usercamera/Capture", "/usercamera/Close", "/usercamera/Streaming"} {
		if err := client.sendBool(address, false); err != nil {
			diagAutoCapture(logPath, "osc recovery bool error: address=%q err=%v", address, err)
			return err
		}
		diagAutoCapture(logPath, "osc recovery bool success: address=%q value=false", address)
		if !sleepContext(ctx, 100*time.Millisecond) {
			return ctx.Err()
		}
	}
	if err := client.sendInt("/usercamera/Mode", 0); err != nil {
		diagAutoCapture(logPath, "osc recovery mode error: err=%v", err)
		return err
	}
	diagAutoCapture(logPath, "osc recovery mode success: value=0")
	return nil
}

func (r AutoCaptureRunner) RunOnce(ctx context.Context) ([]Result, error) {
	cfg := r.Config
	cfg.Normalize()
	ac := cfg.AutoCapture
	logPath := cfg.DiagnosticLogPath
	diagAutoCapture(logPath, "run_once begin: mode=%q schedule_enabled=%t capture_on_start=%t interval_sec=%d close_after_batch=%t button_release_ms=%d settle_ms=%d",
		ac.Capture.Mode,
		ac.Schedule.Enabled,
		ac.Schedule.CaptureOnStart,
		ac.Schedule.CaptureIntervalSec,
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
	diagAutoCapture(logPath, "osc open success: target=%s:%d", ac.OSC.Host, ac.OSC.SendPort)
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
		diagAutoCapture(logPath, "osc button press begin: address=%q detail=%q", "/usercamera/Streaming", "stream_start")
		if err := client.sendBool("/usercamera/Streaming", true); err != nil {
			diagAutoCapture(logPath, "osc button press error: address=%q detail=%q err=%v", "/usercamera/Streaming", "stream_start", err)
			return nil, err
		}
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
	results := make([]Result, 0, len(views))
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
	if ac.Capture.CloseCameraAfterBatch && (successCount > 0 || ac.Capture.Mode == "stream") {
		if ac.Capture.Mode == "stream" {
			diagAutoCapture(logPath, "osc button release begin: address=%q detail=%q", "/usercamera/Streaming", "stream_stop")
			if err := client.sendBool("/usercamera/Streaming", false); err != nil {
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
	if err := r.applyCameraView(client, view); err != nil {
		return Result{Name: name, Error: err.Error()}
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
	diagAutoCapture(logPath, "shot optional_params sent: view_id=%q count=%d", view.ID, sentOptions)
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
		diagAutoCapture(logPath, "camera pose resolve error: view_id=%q coordinate_space=%q err=%v", view.ID, view.CoordinateSpace, err)
		return err
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

func (r AutoCaptureRunner) captureStreamShot(ctx context.Context, client oscClient, batchID string, shotID string, index int, view CameraViewConfig, sidecarUsers []PresenceUser, discordUsers []PresenceUser, confidence string, world AutoCaptureVRChatMetadata) Result {
	cfg := r.Config.AutoCapture
	logPath := r.Config.DiagnosticLogPath
	name := view.Name
	if name == "" {
		name = view.ID
	}
	diagAutoCapture(logPath, "stream shot begin: batch_id=%q shot_id=%q index=%d view_id=%q view_name=%q", batchID, shotID, index, view.ID, name)
	if err := r.applyCameraView(client, view); err != nil {
		return Result{Name: name, Error: err.Error()}
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
	diagAutoCapture(logPath, "shot optional_params sent: view_id=%q count=%d", view.ID, sentOptions)
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
		sidecar := AutoCaptureSidecar{
			SchemaVersion:    1,
			BatchID:          batchID,
			ShotID:           shotID,
			CapturedAtLocal:  time.Now().Format(time.RFC3339),
			CapturedAtUTC:    time.Now().UTC().Format(time.RFC3339),
			CaptureMode:      cfg.Capture.Mode,
			View:             view,
			Stream:           autoCaptureStreamMetadata(streamInfo),
			VRChat:           autoCaptureVRChatMetadata(world, confidence),
			Users:            sidecarUsers,
			MetadataWarnings: metadataWarnings,
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
	SchemaVersion    int                        `json:"schema_version"`
	BatchID          string                     `json:"batch_id"`
	ShotID           string                     `json:"shot_id"`
	CapturedAtLocal  string                     `json:"captured_at_local"`
	CapturedAtUTC    string                     `json:"captured_at_utc"`
	CaptureMode      string                     `json:"capture_mode"`
	View             CameraViewConfig           `json:"view"`
	ResolvedPose     *CameraPoseConfig          `json:"resolved_pose,omitempty"`
	Stream           *AutoCaptureStreamMetadata `json:"stream,omitempty"`
	VRChat           AutoCaptureVRChatMetadata  `json:"vrchat"`
	Users            []PresenceUser             `json:"users"`
	Files            AutoCaptureFileMetadata    `json:"files"`
	MetadataWarnings []string                   `json:"metadata_warnings,omitempty"`
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
	return nil
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

func (c oscClient) send(address string, typeTags string, appendArgs func([]byte) []byte) error {
	if c.conn == nil {
		return fmt.Errorf("OSC接続が開かれていません")
	}
	_, err := c.conn.Write(buildOSCPacket(address, typeTags, appendArgs))
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
