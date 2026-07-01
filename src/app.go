package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	copySingleURLIfNeeded = appcore.CopySingleURLIfNeeded
	deleteDiscordMessage  = appcore.DeleteDiscordMessage
	runningGoTest         = isGoTestBinary(os.Args[0])
)

func isGoTestBinary(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	return strings.HasSuffix(name, ".test") || strings.HasSuffix(name, ".test.exe")
}

func emitWailsEvent(ctx context.Context, name string, data any) {
	if ctx == nil || runningGoTest {
		return
	}
	runtime.EventsEmit(ctx, name, data)
}

type App struct {
	ctx        context.Context
	configPath string
	state      appcore.UIState
	autoCancel context.CancelFunc
	oscCancel  context.CancelFunc
	latestPose appcore.CameraPoseConfig
	poseAt     time.Time
	mu         sync.Mutex
}

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	GitHub  string `json:"github"`
}

type OSSLicense struct {
	Name      string `json:"name"`
	License   string `json:"license"`
	Copyright string `json:"copyright"`
	URL       string `json:"url"`
}

type FFmpegStatus struct {
	Available bool   `json:"available"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Message   string `json:"message"`
}

func NewApp(configPath string, initial appcore.UIState) *App {
	return &App{
		configPath: configPath,
		state:      initial,
	}
}

func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ctx = ctx
	a.logStartupLocked()
	if a.state.Mode == appcore.ModeResults {
		a.restartAutoPhotoWatcher(a.state.Config)
	}
}

func (a *App) shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.stopBackgroundTasksLocked()
}

func (a *App) stopBackgroundTasksLocked() {
	if a.autoCancel != nil {
		a.autoCancel()
		a.autoCancel = nil
	}
	if a.oscCancel != nil {
		a.oscCancel()
		a.oscCancel = nil
	}
}

func (a *App) GetInitialState() appcore.UIState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.refreshHistory()
	return a.state
}

func (a *App) GetAppInfo() AppInfo {
	return AppInfo{
		Name:    "ClipForVRChat",
		Version: appVersion(),
		GitHub:  githubURL,
	}
}

func (a *App) OpenURL(rawURL string) error {
	trustedURL, err := trustedExternalURL(rawURL)
	if err != nil {
		return err
	}
	runtime.BrowserOpenURL(a.ctx, trustedURL)
	return nil
}

func (a *App) LogUserAction(action string, detail string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logUserActionLocked(action, detail)
}

func trustedExternalURL(rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme != "https" {
		return "", fmt.Errorf("開けないURLです")
	}
	switch strings.ToLower(parsed.Hostname()) {
	case "github.com", "hatolife.booth.pm", "support.discord.com", "x.com":
	default:
		return "", fmt.Errorf("許可されていないURLです")
	}
	return parsed.String(), nil
}

func (a *App) CreateEncryptedDiagnosticPackage() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	path, err := createEncryptedDiagnosticPackage(a.configPath, a.state.Config)
	if err != nil {
		a.logUserActionLocked("diagnostic_package_failed", err.Error())
		return "", err
	}
	a.logUserActionLocked("diagnostic_package_created", path)
	return path, nil
}

func (a *App) logUserActionLocked(action string, detail string) {
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "ui action=%q detail=%q", strings.TrimSpace(action), strings.TrimSpace(detail))
}

func (a *App) RevealFileInExplorer(path string) error {
	target, err := explorerRevealPath(path, filepath.Dir(a.configPath))
	if err != nil {
		return err
	}
	return revealFileInExplorer(target)
}

func explorerRevealPath(path string, baseDir string) (string, error) {
	cleaned := appcore.ResolveHistoryOutputPath(path, baseDir)
	if cleaned == "" {
		return "", fmt.Errorf("保存済みファイルがありません")
	}
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}
	stat, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("保存済みファイルを確認できません: %w", err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("保存済み画像ファイルではありません: %s", abs)
	}
	return abs, nil
}

func (a *App) CheckForUpdate() (appcore.UpdateInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	return appcore.CheckLatestRelease(ctx, nil, appVersion(), appReleaseTime())
}

func (a *App) GetOSSLicenses() []OSSLicense {
	return []OSSLicense{
		{Name: "Wails", License: "MIT", Copyright: "Copyright (c) 2018-Present Lea Anthony", URL: "https://github.com/wailsapp/wails"},
		{Name: "go-arg", License: "MIT", Copyright: "Copyright (c) 2015, Alex Flint", URL: "https://github.com/alexflint/go-arg"},
		{Name: "Vue.js", License: "MIT", Copyright: "Copyright (c) 2018-present, Yuxi (Evan) You", URL: "https://github.com/vuejs/core"},
		{Name: "Vite", License: "MIT", Copyright: "Copyright (c) 2019-present, VoidZero Inc. and Vite contributors", URL: "https://github.com/vitejs/vite"},
		{Name: "ProtonMail/go-crypto", License: "MIT", Copyright: "Copyright (c) 2022 Proton AG", URL: "https://github.com/ProtonMail/go-crypto"},
		{Name: "imaging", License: "MIT", Copyright: "Copyright (c) 2012 Grigory Dryapak", URL: "https://github.com/disintegration/imaging"},
		{Name: "flock", License: "BSD-3-Clause", Copyright: "Copyright (c) 2018-2025, The Gofrs; Copyright (c) 2015-2020, Tim Heckman", URL: "https://github.com/gofrs/flock"},
		{Name: "gozxing", License: "MIT", Copyright: "Copyright (c) 2018 Daisuke MAKIUCHI", URL: "https://github.com/makiuchi-d/gozxing"},
		{Name: "go-qrcode", License: "MIT", Copyright: "Copyright (c) 2014 Tom Harwood", URL: "https://github.com/skip2/go-qrcode"},
		{Name: "golang.design/x/clipboard", License: "MIT", Copyright: "Copyright (c) 2021 Changkun Ou", URL: "https://github.com/golang-design/clipboard"},
		{Name: "golang.org/x/crypto", License: "BSD-3-Clause", Copyright: "Copyright (c) The Go Authors", URL: "https://cs.opensource.google/go/x/crypto"},
		{Name: "golang.org/x/image", License: "BSD-3-Clause", Copyright: "Copyright (c) The Go Authors", URL: "https://cs.opensource.google/go/x/image"},
		{Name: "Spout2", License: "BSD-2-Clause", Copyright: "Copyright (c) 2016-2025, Lynn Jarvis", URL: "https://github.com/leadedge/Spout2"},
	}
}

func (a *App) SelectOutputDirectory(current string) (string, error) {
	return a.selectDirectory("出力先フォルダを選択", current)
}

func (a *App) SelectPhotoDirectory(current string) (string, error) {
	return a.selectDirectory("VRChat写真フォルダを選択", current)
}

func (a *App) SelectScreenshotDirectory(current string) (string, error) {
	return a.selectDirectory("スクリーンショットフォルダを選択", current)
}

func (a *App) selectDirectory(title string, current string) (string, error) {
	dir := strings.Trim(strings.TrimSpace(current), `"`)
	if dir != "" && !filepath.IsAbs(dir) {
		exe, err := os.Executable()
		if err == nil {
			dir = filepath.Join(filepath.Dir(exe), dir)
		}
	}
	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		dir = ""
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: dir,
	})
}

func (a *App) ClearResults() appcore.UIState {
	a.mu.Lock()
	defer a.mu.Unlock()
	var ids []string
	for _, result := range a.state.Results {
		if result.HistoryID != "" {
			ids = append(ids, result.HistoryID)
		}
	}
	if len(ids) > 0 {
		if history, err := appcore.MarkHistoryCleared(appcore.HistoryPath(a.configPath), ids); err == nil {
			a.state.History = history
		}
	}
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state
}

func (a *App) GetHistory() ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.LoadHistoryWithManagedOutputDir(appcore.HistoryPath(a.configPath), filepath.Dir(a.configPath), managedOutputDir(a.configPath, a.state.Config))
	if err != nil {
		return nil, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) MarkHistoryCleared(ids []string) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.MarkHistoryCleared(appcore.HistoryPath(a.configPath), ids)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) SetHistoryPinned(id string, pinned bool) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.SetHistoryPinned(appcore.HistoryPath(a.configPath), id, pinned)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) DeleteDiscordHistoryEntries(ids []string) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, err := appcore.LoadHistory(appcore.HistoryPath(a.configPath))
	if err != nil {
		return history, err
	}
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	for i := range history {
		if !idSet[history[i].ID] || history[i].Pinned {
			continue
		}
		if history[i].DiscordDeleted || !appcore.IsTrustedDiscordImageURL(history[i].URL) || history[i].DiscordWebhookID == "" || history[i].DiscordToken == "" || history[i].DiscordMessageID == "" {
			continue
		}
		if err := deleteDiscordMessage(history[i].DiscordWebhookID, history[i].DiscordToken, history[i].DiscordMessageID); err != nil {
			if saveErr := appcore.SaveHistory(appcore.HistoryPath(a.configPath), history); saveErr != nil {
				return history, fmt.Errorf("%v; 履歴保存にも失敗しました: %w", err, saveErr)
			}
			a.state.History = history
			return history, err
		}
		history[i].DiscordDeleted = true
		history[i].DeletedAt = nowRFC3339()
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(a.configPath), history); err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) DeleteLocalHistoryFiles(ids []string) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, _, err := appcore.DeleteLocalHistoryFilesWithManagedOutputDir(appcore.HistoryPath(a.configPath), ids, managedOutputDir(a.configPath, a.state.Config))
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) DeleteHistoryEntries(ids []string) ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	history, _, err := appcore.DeleteHistoryEntries(appcore.HistoryPath(a.configPath), ids)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) PurgeDeletedHistoryEntries() ([]appcore.HistoryEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return a.state.History, err
	}
	history, _, err := appcore.PurgeDiscordDeletedHistory(appcore.HistoryPath(a.configPath), cfg.Output.DeleteOutputOnHistoryPurge)
	if err != nil {
		return history, err
	}
	a.state.History = history
	return history, nil
}

func (a *App) LoadConfig() (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return appcore.LoadConfig(a.configPath)
}

func (a *App) SaveConfig(cfg appcore.Config) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.saveConfigLocked(cfg)
}

func (a *App) saveConfigLocked(cfg appcore.Config) error {
	if err := appcore.SaveConfig(a.configPath, cfg); err != nil {
		return err
	}
	a.state.Config = cfg
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.restartCameraPoseReceiverLocked(cfg)
	a.restartAutoPhotoWatcher(cfg)
	return nil
}

func (a *App) GetLatestCameraPose() appcore.CameraPoseSnapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.latestCameraPoseLocked(a.state.Config)
}

func (a *App) SaveCurrentCameraPoseToView(viewID string) (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg := a.state.Config
	cfg.Normalize()
	pose, err := a.freshCameraPoseLocked(cfg)
	if err != nil {
		return cfg, err
	}
	viewID = strings.TrimSpace(viewID)
	found := false
	for i := range cfg.AutoCapture.Views {
		if cfg.AutoCapture.Views[i].ID != viewID {
			continue
		}
		savedPose := pose
		if cfg.AutoCapture.Views[i].CoordinateSpace == "player_local" {
			if !cfg.AutoCapture.PlayerLocal.Calibrated {
				return cfg, fmt.Errorf("プレイヤー基準Poseが未設定のため、player_local構図を保存できません。自動撮影タブで現在Poseをプレイヤー基準として保存してください")
			}
			savedPose = appcore.InverseTransformPlayerLocalPose(cfg.AutoCapture.PlayerLocal.BasisPose, pose)
		}
		cfg.AutoCapture.Views[i].Pose = savedPose
		cfg.AutoCapture.Views[i].Calibrated = true
		found = true
		break
	}
	if !found {
		return cfg, fmt.Errorf("構図が見つかりません: %s", viewID)
	}
	if err := a.saveAutoCaptureConfigFromSettingsLocked(cfg); err != nil {
		return cfg, err
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "auto-capture pose saved to view=%q", viewID)
	return a.state.Config, nil
}

func (a *App) AddCurrentCameraPoseAsView(viewID string) (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg := a.state.Config
	cfg.Normalize()
	pose, err := a.freshCameraPoseLocked(cfg)
	if err != nil {
		return cfg, err
	}
	viewID = strings.TrimSpace(viewID)
	var sourceView appcore.CameraViewConfig
	found := false
	for _, candidate := range cfg.AutoCapture.Views {
		if candidate.ID == viewID {
			sourceView = candidate
			found = true
			break
		}
	}
	if !found {
		return cfg, fmt.Errorf("構図が見つかりません: %s", viewID)
	}
	savedPose := pose
	if sourceView.CoordinateSpace == "player_local" {
		if !cfg.AutoCapture.PlayerLocal.Calibrated {
			return cfg, fmt.Errorf("プレイヤー基準Poseが未設定のため、player_local構図を追加できません。自動撮影タブで現在Poseをプレイヤー基準として保存してください")
		}
		savedPose = appcore.InverseTransformPlayerLocalPose(cfg.AutoCapture.PlayerLocal.BasisPose, pose)
	}
	id := newCameraViewID(cfg.AutoCapture.Views)
	var zoom *float64
	if sourceView.Zoom != nil {
		value := *sourceView.Zoom
		zoom = &value
	}
	var exposure *float64
	if sourceView.Exposure != nil {
		value := *sourceView.Exposure
		exposure = &value
	}
	var focalDistance *float64
	if sourceView.FocalDistance != nil {
		value := *sourceView.FocalDistance
		focalDistance = &value
	}
	var aperture *float64
	if sourceView.Aperture != nil {
		value := *sourceView.Aperture
		aperture = &value
	}
	var lookAtMe *bool
	if sourceView.LookAtMe != nil {
		value := *sourceView.LookAtMe
		lookAtMe = &value
	}
	var showUIInCamera *bool
	if sourceView.ShowUIInCamera != nil {
		value := *sourceView.ShowUIInCamera
		showUIInCamera = &value
	}
	var localPlayer *bool
	if sourceView.LocalPlayer != nil {
		value := *sourceView.LocalPlayer
		localPlayer = &value
	}
	var remotePlayer *bool
	if sourceView.RemotePlayer != nil {
		value := *sourceView.RemotePlayer
		remotePlayer = &value
	}
	var environment *bool
	if sourceView.Environment != nil {
		value := *sourceView.Environment
		environment = &value
	}
	cfg.AutoCapture.Views = append(cfg.AutoCapture.Views, appcore.CameraViewConfig{
		ID:              id,
		Name:            fmt.Sprintf("構図 %d", len(cfg.AutoCapture.Views)+1),
		Enabled:         sourceView.Enabled,
		SortOrder:       nextCameraViewSortOrder(cfg.AutoCapture.Views),
		CoordinateSpace: sourceView.CoordinateSpace,
		Pose:            savedPose,
		Zoom:            zoom,
		Exposure:        exposure,
		FocalDistance:   focalDistance,
		Aperture:        aperture,
		LookAtMe:        lookAtMe,
		ShowUIInCamera:  showUIInCamera,
		LocalPlayer:     localPlayer,
		RemotePlayer:    remotePlayer,
		Environment:     environment,
		SettleDelayMS:   sourceView.SettleDelayMS,
		CaptureDelayMS:  sourceView.CaptureDelayMS,
		Calibrated:      true,
	})
	if err := a.saveAutoCaptureConfigFromSettingsLocked(cfg); err != nil {
		return cfg, err
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "auto-capture pose saved as new view=%q source_view_id=%q coordinate_space=%q", id, viewID, sourceView.CoordinateSpace)
	return a.state.Config, nil
}

func (a *App) SaveCurrentCameraPoseAsPlayerBasis() (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg := a.state.Config
	cfg.Normalize()
	pose, err := a.freshCameraPoseLocked(cfg)
	if err != nil {
		return cfg, err
	}
	cfg.AutoCapture.PlayerLocal.BasisPose = pose
	cfg.AutoCapture.PlayerLocal.Calibrated = true
	cfg.AutoCapture.PlayerLocal.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := a.saveAutoCaptureConfigFromSettingsLocked(cfg); err != nil {
		return cfg, err
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "auto-capture player_local basis saved")
	return a.state.Config, nil
}

func (a *App) ResetCameraPoseToDefault(viewID string) (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg := a.state.Config
	cfg.Normalize()
	viewID = strings.TrimSpace(viewID)
	defaultView, ok := appcore.DefaultCameraViewByID(viewID)
	if !ok {
		return cfg, fmt.Errorf("初期Poseが定義されていない構図です: %s", viewID)
	}
	found := false
	for i := range cfg.AutoCapture.Views {
		if cfg.AutoCapture.Views[i].ID != viewID {
			continue
		}
		cfg.AutoCapture.Views[i].Pose = defaultView.Pose
		cfg.AutoCapture.Views[i].Zoom = defaultView.Zoom
		cfg.AutoCapture.Views[i].CoordinateSpace = defaultView.CoordinateSpace
		cfg.AutoCapture.Views[i].Calibrated = false
		found = true
		break
	}
	if !found {
		return cfg, fmt.Errorf("構図が見つかりません: %s", viewID)
	}
	if err := a.saveAutoCaptureConfigFromSettingsLocked(cfg); err != nil {
		return cfg, err
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "auto-capture pose reset to default view=%q", viewID)
	return a.state.Config, nil
}

func (a *App) ResetCameraViewsToDefaults() (appcore.Config, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg := a.state.Config
	cfg.Normalize()
	cfg.AutoCapture.Views = appcore.DefaultCameraViews()
	if err := a.saveAutoCaptureConfigFromSettingsLocked(cfg); err != nil {
		return cfg, err
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "auto-capture views reset to defaults")
	return a.state.Config, nil
}

func (a *App) MoveCameraToView(viewID string) error {
	a.mu.Lock()
	cfg := a.state.Config
	cfg.Normalize()
	cfg.DiagnosticLogPath = appcore.DiagnosticLogPath(a.configPath)
	a.mu.Unlock()
	if err := appcore.MoveUserCameraToView(context.Background(), cfg, viewID); err != nil {
		a.mu.Lock()
		a.state.Message = "カメラ移動に失敗しました: " + err.Error()
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.state.Message = "カメラを構図のPoseへ移動しました。"
	a.mu.Unlock()
	return nil
}

func (a *App) ResetCameraOSC() error {
	a.mu.Lock()
	cfg := a.state.Config
	cfg.Normalize()
	cfg.DiagnosticLogPath = appcore.DiagnosticLogPath(a.configPath)
	a.mu.Unlock()
	if err := appcore.ResetUserCameraOSC(context.Background(), cfg); err != nil {
		a.mu.Lock()
		a.state.Message = "カメラOSCリセットに失敗しました: " + err.Error()
		a.mu.Unlock()
		return err
	}
	a.mu.Lock()
	a.state.Message = "カメラOSCをリセットしました。"
	a.mu.Unlock()
	return nil
}

func (a *App) CheckSpoutHelper() appcore.SpoutHelperStatus {
	a.mu.Lock()
	cfg := a.state.Config
	cfg.Normalize()
	cfg.DiagnosticLogPath = appcore.DiagnosticLogPath(a.configPath)
	a.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return appcore.CheckSpoutHelper(ctx, cfg.AutoCapture.Stream, cfg.DiagnosticLogPath)
}

func (a *App) ListSpoutSenders() appcore.SpoutHelperStatus {
	a.mu.Lock()
	cfg := a.state.Config
	cfg.Normalize()
	cfg.DiagnosticLogPath = appcore.DiagnosticLogPath(a.configPath)
	a.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := appcore.ListSpoutSenders(ctx, cfg.AutoCapture.Stream, cfg.DiagnosticLogPath)
	if err != nil {
		return appcore.SpoutHelperStatus{Available: true, Message: err.Error(), Senders: result.Senders}
	}
	message := fmt.Sprintf("%d件のSpout senderを検出しました。", len(result.Senders))
	if len(result.Senders) == 0 {
		message = "Spout senderがありません。VRChatでStream Cameraを起動し、OSCでStreamingを有効にしてください。"
	}
	return appcore.SpoutHelperStatus{Available: true, Message: message, Senders: result.Senders}
}

func (a *App) CheckFFmpeg(ffmpegPath string) FFmpegStatus {
	logPath := appcore.DiagnosticLogPath(a.configPath)
	resolved, err := appcore.ResolveFFmpegPath(ffmpegPath)
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "ffmpeg check missing: path=%q err=%v", ffmpegPath, err)
		return FFmpegStatus{
			Available: false,
			Message:   "ffmpegがインストールされていないかPATHにありません。Stream方式を使うにはffmpegをインストールするか、ffmpeg.exeのパスを指定してください。",
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, resolved, "-version") // #nosec G204 -- user-configured local ffmpeg path used for availability check.
	output, runErr := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		appcore.AppendDiagnosticLog(logPath, "ffmpeg check timeout: path=%q resolved=%q", ffmpegPath, resolved)
		return FFmpegStatus{Available: false, Path: resolved, Message: "ffmpeg -versionがタイムアウトしました。ffmpegパスを確認してください。"}
	}
	if runErr != nil {
		appcore.AppendDiagnosticLog(logPath, "ffmpeg check error: path=%q resolved=%q err=%v output=%q", ffmpegPath, resolved, runErr, strings.TrimSpace(string(output)))
		return FFmpegStatus{Available: false, Path: resolved, Message: "ffmpegは見つかりましたが実行できませんでした: " + runErr.Error()}
	}
	version := firstNonEmptyLine(string(output))
	appcore.AppendDiagnosticLog(logPath, "ffmpeg check success: path=%q resolved=%q version=%q", ffmpegPath, resolved, version)
	return FFmpegStatus{
		Available: true,
		Path:      resolved,
		Version:   version,
		Message:   "ffmpegを利用できます。",
	}
}

func (a *App) InstallFFmpegWithWinget() FFmpegStatus {
	logPath := appcore.DiagnosticLogPath(a.configPath)
	if _, err := exec.LookPath("winget"); err != nil {
		appcore.AppendDiagnosticLog(logPath, "ffmpeg install skipped: winget_missing err=%v", err)
		return FFmpegStatus{Available: false, Message: "wingetが見つかりません。Microsoft Storeのアプリ インストーラーを更新するか、手動でffmpegをインストールしてください。"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	appcore.AppendDiagnosticLog(logPath, "ffmpeg install begin: command=%q", "winget install ffmpeg")
	cmd := exec.CommandContext(ctx, "winget", "install", "ffmpeg") // #nosec G204 -- fixed installer command requested by user.
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if len(trimmed) > 1000 {
		trimmed = trimmed[len(trimmed)-1000:]
	}
	if ctx.Err() == context.DeadlineExceeded {
		appcore.AppendDiagnosticLog(logPath, "ffmpeg install timeout: output=%q", trimmed)
		return FFmpegStatus{Available: false, Message: "winget install ffmpegがタイムアウトしました。PowerShellで同じコマンドを確認してください。"}
	}
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "ffmpeg install error: err=%v output=%q", err, trimmed)
		return FFmpegStatus{Available: false, Message: "winget install ffmpegに失敗しました: " + err.Error() + installOutputSuffix(trimmed)}
	}
	appcore.AppendDiagnosticLog(logPath, "ffmpeg install success: output=%q", trimmed)
	status := a.CheckFFmpeg("ffmpeg")
	if status.Available {
		status.Message = "ffmpegをインストールし、利用可能なことを確認しました。"
		return status
	}
	status.Message = "winget install ffmpegは完了しましたが、現在のアプリからffmpegを見つけられません。アプリを再起動するか、ffmpeg.exeのパスを指定してください。"
	return status
}

func firstNonEmptyLine(value string) string {
	for _, line := range strings.Split(value, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
}

func installOutputSuffix(output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}
	return " / " + output
}

func (a *App) TestAutoCaptureView(viewID string) ([]appcore.Result, error) {
	a.mu.Lock()
	cfg := a.state.Config
	cfg.Normalize()
	viewID = strings.TrimSpace(viewID)
	found := false
	for i := range cfg.AutoCapture.Views {
		if cfg.AutoCapture.Views[i].ID == viewID {
			cfg.AutoCapture.Views[i].Enabled = true
			found = true
			continue
		}
		cfg.AutoCapture.Views[i].Enabled = false
	}
	if !found {
		a.mu.Unlock()
		return nil, fmt.Errorf("構図が見つかりません: %s", viewID)
	}
	cfg.AutoCapture.Schedule.Enabled = false
	cfg.DiagnosticLogPath = appcore.DiagnosticLogPath(a.configPath)
	a.mu.Unlock()

	results, err := appcore.AutoCaptureRunner{Config: cfg}.RunOnce(context.Background())

	a.mu.Lock()
	defer a.mu.Unlock()
	if err != nil {
		a.state.Message = "自動撮影テストに失敗しました: " + err.Error()
		return results, err
	}
	visible := autoCaptureVisibleResults(results)
	if len(visible) > 0 {
		a.state.Results = append(visible, a.state.Results...)
		if history, historyErr := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), visible); historyErr == nil {
			a.state.History = history
		}
		a.state.Message = fmt.Sprintf("自動撮影テストで%d件撮影しました。", len(visible))
	} else {
		a.state.Message = summarizeAutoCaptureErrors(results)
	}
	return results, nil
}

func (a *App) saveAutoCaptureConfigFromSettingsLocked(cfg appcore.Config) error {
	if err := appcore.SaveConfig(a.configPath, cfg); err != nil {
		return err
	}
	a.state.Config = cfg
	a.state.ConfigPath = a.configPath
	a.restartCameraPoseReceiverLocked(cfg)
	a.restartAutoPhotoWatcher(cfg)
	return nil
}

func (a *App) CloseSettings() appcore.UIState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state.Mode = appcore.ModeResults
	a.state.Message = ""
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	return a.state
}

func (a *App) OpenSettings(path string) (appcore.UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	configPath := a.configPath
	if strings.TrimSpace(path) != "" {
		configPath = path
	}
	cfg, err := appcore.LoadConfig(configPath)
	if err != nil {
		return a.state, err
	}
	a.configPath = configPath
	a.state.Mode = appcore.ModeSettings
	a.state.Message = ""
	a.state.ConfigPath = a.configPath
	a.state.Config = cfg
	a.state.Results = nil
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	a.restartCameraPoseReceiverLocked(cfg)
	return a.state, nil
}

func (a *App) SaveConfigAndProcess(cfg appcore.Config, paths []string) (appcore.UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.saveConfigLocked(cfg); err != nil {
		return a.state, err
	}
	a.logProcessingStartLocked("save_config_and_process", cfg, len(paths))
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
	if err != nil {
		a.logUserActionLocked("process_error", err.Error())
		a.state.Mode = appcore.ModeError
		a.state.Message = err.Error()
		return a.state, nil
	}
	copyErr := copySingleURLIfNeeded(cfg, results)
	if history, historyErr := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results); historyErr == nil {
		a.state.History = history
	}
	visibleResults := userVisibleResults(results)
	a.state.Results = visibleResults
	a.state.PendingPaths = nil
	a.state.ProcessOnSave = false
	if hasResultErrors(results) {
		a.state.Mode = appcore.ModeError
		a.state.Message = "処理中にエラーが発生しました。内容を確認してください。"
	} else {
		a.state.Mode = appcore.ModeResults
		a.state.Message = resultMessage(cfg, visibleResults, copyErr)
	}
	a.logProcessingResultLocked("save_config_and_process", results, visibleResults, a.state.Message)
	return a.state, nil
}

func (a *App) ProcessToState(paths []string) (appcore.UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if hasJSONPath(paths) {
		a.state.Mode = appcore.ModeError
		a.state.Message = "画像ファイルと設定ファイルが混在しています。設定編集と画像処理は別々に実行してください。"
		a.state.Results = nil
		return a.state, nil
	}
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return a.state, err
	}
	source := "drop"
	if len(paths) == 0 {
		source = "clipboard"
	}
	a.logProcessingStartLocked(source, cfg, len(paths))
	results, err := appcore.Processor{Config: cfg}.ProcessPathsWithProgress(paths, func(event appcore.ProgressEvent) {
		emitWailsEvent(a.ctx, "process:progress", event)
	})
	if err != nil {
		a.logUserActionLocked("process_error", err.Error())
		a.state.Mode = appcore.ModeError
		a.state.Message = err.Error()
		a.state.Results = nil
		return a.state, nil
	}
	copyErr := copySingleURLIfNeeded(cfg, results)
	if history, historyErr := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results); historyErr == nil {
		a.state.History = history
	}
	a.state.Config = cfg
	visibleResults := userVisibleResults(results)
	a.state.Results = visibleResults
	if hasResultErrors(results) {
		a.state.Mode = appcore.ModeError
		a.state.Message = "処理中にエラーが発生しました。内容を確認してください。"
	} else {
		a.state.Mode = appcore.ModeResults
		a.state.Message = resultMessage(cfg, visibleResults, copyErr)
	}
	a.logProcessingResultLocked(source, results, visibleResults, a.state.Message)
	return a.state, nil
}

func (a *App) Process(paths []string) ([]appcore.Result, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg, err := appcore.LoadConfig(a.configPath)
	if err != nil {
		return nil, err
	}
	a.logProcessingStartLocked("process_api", cfg, len(paths))
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(paths)
	if err != nil {
		a.logUserActionLocked("process_error", err.Error())
		return nil, err
	}
	_ = copySingleURLIfNeeded(cfg, results)
	_, _ = appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results)
	a.logProcessingResultLocked("process_api", results, userVisibleResults(results), "")
	return results, nil
}

func (a *App) CopyText(text string) error {
	return appcore.WriteClipboardText(text)
}

func defaultConfigPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	return appcore.ConfigPath(exe)
}

func hasResultErrors(results []appcore.Result) bool {
	for _, result := range results {
		if result.Error != "" {
			return true
		}
	}
	return false
}

func userVisibleResults(results []appcore.Result) []appcore.Result {
	visible := make([]appcore.Result, 0, len(results))
	for _, result := range results {
		if appcore.ResultHasUserVisibleWork(result) {
			visible = append(visible, result)
		}
	}
	return visible
}

func (a *App) logProcessingStartLocked(source string, cfg appcore.Config, inputCount int) {
	appcore.AppendDiagnosticLog(
		appcore.DiagnosticLogPath(a.configPath),
		"process start source=%q input_count=%d save_local=%t upload_discord=%t detect_qr=%t copy_single_url=%t output_format=%q",
		source,
		inputCount,
		cfg.Output.SaveLocal,
		cfg.Output.UploadDiscord,
		cfg.Output.DetectQRCodeURLs,
		cfg.Output.CopySingleURLToClipboard,
		cfg.Image.OutputFormat,
	)
}

func (a *App) logProcessingResultLocked(source string, results []appcore.Result, visibleResults []appcore.Result, message string) {
	discordCount := 0
	localCount := 0
	qrCount := 0
	errorCount := 0
	for _, result := range results {
		if result.URL != "" {
			discordCount++
		}
		if result.OutputPath != "" {
			localCount++
		}
		if len(result.QRURLs) > 0 {
			qrCount++
		}
		if result.Error != "" {
			errorCount++
		}
	}
	appcore.AppendDiagnosticLog(
		appcore.DiagnosticLogPath(a.configPath),
		"process result source=%q total=%d visible=%d discord=%d local=%d qr=%d errors=%d message=%q",
		source,
		len(results),
		len(visibleResults),
		discordCount,
		localCount,
		qrCount,
		errorCount,
		message,
	)
}

func resultMessage(cfg appcore.Config, results []appcore.Result, copyErr error) string {
	if len(results) == 0 {
		return noVisibleResultMessage(cfg)
	}
	if len(results) == 1 && results[0].URL != "" && cfg.Output.CopySingleURLToClipboard {
		if copyErr != nil {
			return fmt.Sprintf("画像URLを取得しましたが、クリップボードへコピーできませんでした: %v サムネイルをクリックすると再コピーできます。", copyErr)
		}
		return "画像URLをクリップボードにコピーしました。"
	}
	if len(results) == 1 && results[0].URL == "" && results[0].Error == "" {
		var parts []string
		if results[0].OutputPath != "" {
			parts = append(parts, "縮小画像を保存しました")
		}
		if len(results[0].QRURLs) > 0 {
			parts = append(parts, "QRコードURLを取得しました")
		}
		if len(parts) == 0 {
			parts = append(parts, "処理しました")
		}
		return strings.Join(parts, "。") + "。"
	}
	actions := map[string]bool{}
	for _, result := range results {
		if result.URL != "" {
			actions["Discordへ投稿"] = true
		}
		if result.OutputPath != "" {
			actions["ローカル保存"] = true
		}
		if len(result.QRURLs) > 0 {
			actions["QRコードURL取得"] = true
		}
	}
	var parts []string
	for _, label := range []string{"Discordへ投稿", "ローカル保存", "QRコードURL取得"} {
		if actions[label] {
			parts = append(parts, label)
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, "、") + "を行いました。"
	}
	return ""
}

func noVisibleResultMessage(cfg appcore.Config) string {
	if !cfg.Output.UploadDiscord && !cfg.Output.SaveLocal && !cfg.Output.DetectQRCodeURLs {
		return "実行された処理はありません。Discord投稿、ローカル保存、QRコードURL取得がすべてOFFのため、表示できる処理結果がありません。設定を確認してください。"
	}
	if cfg.Output.DetectQRCodeURLs && !cfg.Output.UploadDiscord && !cfg.Output.SaveLocal {
		return "表示できる処理結果はありません。QRコードURL取得はONですが、この画像からURLを取得できませんでした。Discord投稿またはローカル保存をONにすると、画像の処理結果も残せます。"
	}
	var disabled []string
	if !cfg.Output.UploadDiscord {
		disabled = append(disabled, "Discord投稿")
	}
	if !cfg.Output.SaveLocal {
		disabled = append(disabled, "ローカル保存")
	}
	if !cfg.Output.DetectQRCodeURLs {
		disabled = append(disabled, "QRコードURL取得")
	}
	if len(disabled) > 0 {
		return fmt.Sprintf("表示できる処理結果はありません。%sがOFFのため、今回の画像では結果として残せる処理がありませんでした。設定を確認してください。", strings.Join(disabled, "、"))
	}
	return "表示できる処理結果はありません。"
}

func hasJSONPath(paths []string) bool {
	for _, path := range paths {
		if strings.EqualFold(filepath.Ext(path), ".json") {
			return true
		}
	}
	return false
}

func (a *App) refreshHistory() {
	if history, err := appcore.LoadHistoryWithManagedOutputDir(appcore.HistoryPath(a.configPath), filepath.Dir(a.configPath), managedOutputDir(a.configPath, a.state.Config)); err == nil {
		a.state.History = history
	}
}

func managedOutputDir(configPath string, cfg appcore.Config) string {
	outputDir := strings.TrimSpace(cfg.Image.OutputDirectory)
	if outputDir == "" {
		outputDir = "./output"
	}
	if filepath.IsAbs(outputDir) {
		return outputDir
	}
	return filepath.Join(filepath.Dir(configPath), outputDir)
}

func nowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func (a *App) restartAutoPhotoWatcher(cfg appcore.Config) {
	cfg.Normalize()
	logPath := appcore.DiagnosticLogPath(a.configPath)
	if a.autoCancel != nil {
		appcore.AppendDiagnosticLog(logPath, "auto-watchers restart: cancelling existing watchers")
		a.autoCancel()
		a.autoCancel = nil
	}
	if a.ctx == nil || (!cfg.AutoPhoto.Enabled && !cfg.ScreenshotAutoPost.Enabled && !cfg.AutoCapture.Schedule.Enabled) {
		appcore.AppendDiagnosticLog(
			logPath,
			"auto-watchers disabled: ctx_nil=%t autoPhoto=%t screenshotAutoPost=%t autoCaptureSchedule=%t",
			a.ctx == nil,
			cfg.AutoPhoto.Enabled,
			cfg.ScreenshotAutoPost.Enabled,
			cfg.AutoCapture.Schedule.Enabled,
		)
		return
	}
	enabledAutoCaptureViews := 0
	for _, view := range cfg.AutoCapture.Views {
		if view.Enabled {
			enabledAutoCaptureViews++
		}
	}
	appcore.AppendDiagnosticLog(
		logPath,
		"auto-watchers start: autoPhoto=%t screenshotAutoPost=%t autoCaptureSchedule=%t autoCapture_captureOnStart=%t autoCapture_interval_sec=%d autoCapture_initial_delay_sec=%d autoCapture_max_batches=%d autoCapture_mode=%q osc=%s:%d photo_dir=%q auto_capture_output_dir=%q auto_capture_views=%d/%d",
		cfg.AutoPhoto.Enabled,
		cfg.ScreenshotAutoPost.Enabled,
		cfg.AutoCapture.Schedule.Enabled,
		cfg.AutoCapture.Schedule.CaptureOnStart,
		cfg.AutoCapture.Schedule.CaptureIntervalSec,
		cfg.AutoCapture.Schedule.InitialDelaySec,
		cfg.AutoCapture.Schedule.MaxBatches,
		cfg.AutoCapture.Capture.Mode,
		cfg.AutoCapture.OSC.Host,
		cfg.AutoCapture.OSC.SendPort,
		cfg.AutoPhoto.PhotoDirectory,
		cfg.AutoCapture.Output.Directory,
		enabledAutoCaptureViews,
		len(cfg.AutoCapture.Views),
	)
	ctx, cancel := context.WithCancel(a.ctx)
	a.autoCancel = cancel
	handler := func(event appcore.AutoPhotoEvent) {
		a.mu.Lock()
		defer a.mu.Unlock()
		if event.Error != "" {
			a.state.Results = append([]appcore.Result{event.Result}, a.state.Results...)
			a.state.Mode = appcore.ModeResults
			emitWailsEvent(a.ctx, "auto-photo:result", event)
			return
		}
		if appcore.ResultHasUserVisibleWork(event.Result) && event.Result.Error == "" {
			results := []appcore.Result{event.Result}
			if history, err := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), results); err == nil {
				a.state.History = history
				event.Result = results[0]
			}
		}
		if appcore.ResultHasUserVisibleWork(event.Result) {
			a.state.Results = append([]appcore.Result{event.Result}, a.state.Results...)
		}
		a.state.Mode = appcore.ModeResults
		emitWailsEvent(a.ctx, "auto-photo:result", event)
	}
	if cfg.AutoPhoto.Enabled {
		watcher := appcore.AutoPhotoWatcher{
			Config:     cfg,
			Directory:  cfg.AutoPhoto.PhotoDirectory,
			WebhookURL: cfg.AutoPhoto.WebhookURL,
			Interval:   time.Duration(cfg.AutoPhoto.ScanIntervalSeconds) * time.Second,
			Handler:    handler,
		}
		go watcher.Run(ctx)
	}
	if cfg.ScreenshotAutoPost.Enabled {
		watcher := appcore.AutoPhotoWatcher{
			Config:     cfg,
			Directory:  cfg.ScreenshotAutoPost.ScreenshotDirectory,
			WebhookURL: cfg.ScreenshotAutoPost.WebhookURL,
			Interval:   time.Duration(cfg.ScreenshotAutoPost.ScanIntervalSeconds) * time.Second,
			Handler:    handler,
		}
		go watcher.Run(ctx)
	}
	if cfg.AutoCapture.Schedule.Enabled {
		appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler launch")
		go a.runAutoCaptureScheduler(ctx, cfg)
	}
}

func (a *App) runAutoCaptureScheduler(ctx context.Context, cfg appcore.Config) {
	cfg.Normalize()
	logPath := appcore.DiagnosticLogPath(a.configPath)
	ac := cfg.AutoCapture
	appcore.AppendDiagnosticLog(
		logPath,
		"auto-capture scheduler start: capture_on_start=%t initial_delay_sec=%d interval_sec=%d max_batches=%d",
		ac.Schedule.CaptureOnStart,
		ac.Schedule.InitialDelaySec,
		ac.Schedule.CaptureIntervalSec,
		ac.Schedule.MaxBatches,
	)
	if ac.Schedule.InitialDelaySec > 0 {
		appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler initial_delay begin: seconds=%d", ac.Schedule.InitialDelaySec)
		if !sleepWithContext(ctx, time.Duration(ac.Schedule.InitialDelaySec)*time.Second) {
			appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler initial_delay cancelled: err=%v", ctx.Err())
			return
		}
		appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler initial_delay complete")
	}
	batches := 0
	if ac.Schedule.CaptureOnStart {
		if ac.Schedule.MaxBatches > 0 && batches >= ac.Schedule.MaxBatches {
			appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler skip: reason=%q batches=%d max_batches=%d", "max_batches_before_capture_on_start", batches, ac.Schedule.MaxBatches)
			return
		}
		appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler trigger: reason=%q", "capture_on_start")
		a.runAutoCaptureBatch(ctx, cfg)
		batches++
	}
	interval := time.Duration(ac.Schedule.CaptureIntervalSec) * time.Second
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler overlap_policy: synchronous=true skip_if_previous_batch_running=%t", ac.Schedule.SkipIfPreviousBatchRunning)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler stop: reason=%q err=%v", "context_done", ctx.Err())
			return
		case <-ticker.C:
			if ac.Schedule.MaxBatches > 0 && batches >= ac.Schedule.MaxBatches {
				appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler stop: reason=%q batches=%d max_batches=%d", "max_batches", batches, ac.Schedule.MaxBatches)
				return
			}
			appcore.AppendDiagnosticLog(logPath, "auto-capture scheduler trigger: reason=%q batch_index=%d", "interval", batches+1)
			a.runAutoCaptureBatch(ctx, cfg)
			batches++
		}
	}
}

func (a *App) runAutoCaptureBatch(ctx context.Context, cfg appcore.Config) {
	logPath := appcore.DiagnosticLogPath(a.configPath)
	cfg.DiagnosticLogPath = logPath
	appcore.AppendDiagnosticLog(logPath, "auto-capture batch begin")
	runner := appcore.AutoCaptureRunner{
		Config: cfg,
		Handler: func(event appcore.AutoCaptureEvent) {
			emitWailsEvent(a.ctx, "auto-capture:result", event)
		},
	}
	results, err := runner.RunOnce(ctx)
	a.mu.Lock()
	defer a.mu.Unlock()
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "auto-capture batch error: %v", err)
		a.state.Message = "自動処理でエラーが発生しました: " + err.Error()
		a.state.Mode = appcore.ModeResults
		emitWailsEvent(a.ctx, "auto-photo:result", appcore.AutoPhotoEvent{Result: appcore.Result{Name: "自動撮影", Error: err.Error()}, Error: err.Error()})
		return
	}
	errorCount := 0
	for _, result := range results {
		if result.Error != "" {
			errorCount++
		}
	}
	appcore.AppendDiagnosticLog(logPath, "auto-capture batch complete: results=%d errors=%d", len(results), errorCount)
	visible := autoCaptureVisibleResults(results)
	if len(visible) > 0 {
		a.state.Results = append(visible, a.state.Results...)
		if history, historyErr := appcore.AddResultsToHistory(appcore.HistoryPath(a.configPath), visible); historyErr == nil {
			a.state.History = history
		}
		a.state.Message = fmt.Sprintf("自動撮影で%d件撮影しました。", len(visible))
		for _, result := range results {
			emitWailsEvent(a.ctx, "auto-photo:result", appcore.AutoPhotoEvent{Path: result.SourcePath, Result: result, Error: result.Error})
		}
	} else if len(results) > 0 {
		a.state.Message = summarizeAutoCaptureErrors(results)
	}
	a.state.Mode = appcore.ModeResults
}

func autoCaptureVisibleResults(results []appcore.Result) []appcore.Result {
	visible := make([]appcore.Result, 0, len(results))
	for _, result := range results {
		if appcore.ResultHasUserVisibleWork(result) {
			visible = append(visible, result)
		}
	}
	return visible
}

func summarizeAutoCaptureErrors(results []appcore.Result) string {
	if len(results) == 0 {
		return "自動処理でエラーが発生しました: 撮影結果がありません。"
	}
	for _, result := range results {
		if strings.TrimSpace(result.Error) != "" {
			return "自動処理でエラーが発生しました: " + result.Error
		}
	}
	return "自動処理でエラーが発生しました。"
}

func sleepWithContext(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (a *App) restartCameraPoseReceiverLocked(cfg appcore.Config) {
	cfg.Normalize()
	logPath := appcore.DiagnosticLogPath(a.configPath)
	if a.oscCancel != nil {
		a.oscCancel()
		a.oscCancel = nil
	}
	if a.ctx == nil {
		return
	}
	host := strings.TrimSpace(cfg.AutoCapture.OSC.Host)
	if host == "" {
		host = "127.0.0.1"
	}
	port := cfg.AutoCapture.OSC.ReceivePort
	if port <= 0 || port > 65535 {
		port = 9001
	}
	ctx, cancel := context.WithCancel(a.ctx)
	a.oscCancel = cancel
	go a.runCameraPoseReceiver(ctx, host, port, logPath)
}

func (a *App) runCameraPoseReceiver(ctx context.Context, host string, port int, logPath string) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "auto-capture pose receiver resolve error: %v", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "auto-capture pose receiver listen error: addr=%s err=%v", addr.String(), err)
		return
	}
	defer conn.Close()
	appcore.AppendDiagnosticLog(logPath, "auto-capture pose receiver start: addr=%s", addr.String())
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	buf := make([]byte, 2048)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil {
				appcore.AppendDiagnosticLog(logPath, "auto-capture pose receiver stop: err=%v", ctx.Err())
				return
			}
			appcore.AppendDiagnosticLog(logPath, "auto-capture pose receiver read error: %v", err)
			continue
		}
		pose, ok := appcore.ParseOSCPose(buf[:n])
		if !ok {
			continue
		}
		a.mu.Lock()
		a.latestPose = pose
		a.poseAt = time.Now()
		a.mu.Unlock()
		appcore.AppendDiagnosticLog(logPath, "auto-capture pose received: x=%.3f y=%.3f z=%.3f rx=%.3f ry=%.3f rz=%.3f", pose.Position.X, pose.Position.Y, pose.Position.Z, pose.Rotation.X, pose.Rotation.Y, pose.Rotation.Z)
	}
}

func (a *App) latestCameraPoseLocked(cfg appcore.Config) appcore.CameraPoseSnapshot {
	cfg.Normalize()
	if a.poseAt.IsZero() {
		return appcore.CameraPoseSnapshot{Configured: true, Fresh: false}
	}
	age := time.Since(a.poseAt)
	freshness := time.Duration(cfg.AutoCapture.OSC.PoseFreshnessSec) * time.Second
	if freshness <= 0 {
		freshness = 3 * time.Second
	}
	return appcore.CameraPoseSnapshot{
		Pose:       a.latestPose,
		UpdatedAt:  a.poseAt.Format(time.RFC3339Nano),
		AgeMS:      age.Milliseconds(),
		Fresh:      age <= freshness,
		Configured: true,
	}
}

func (a *App) freshCameraPoseLocked(cfg appcore.Config) (appcore.CameraPoseConfig, error) {
	snapshot := a.latestCameraPoseLocked(cfg)
	if snapshot.UpdatedAt == "" {
		return appcore.CameraPoseConfig{}, fmt.Errorf("VRChatからUser Camera Poseをまだ受信していません。VRChatのOSCを有効にし、User Cameraを表示して少し動かしてください。")
	}
	if !snapshot.Fresh {
		return appcore.CameraPoseConfig{}, fmt.Errorf("User Camera Poseが古いです。VRChat内でUser Cameraを少し動かしてから保存してください。")
	}
	return snapshot.Pose, nil
}

func newCameraViewID(views []appcore.CameraViewConfig) string {
	seen := map[string]bool{}
	for _, view := range views {
		seen[view.ID] = true
	}
	for i := len(views) + 1; ; i++ {
		id := fmt.Sprintf("view-%d", i)
		if !seen[id] {
			return id
		}
	}
}

func nextCameraViewSortOrder(views []appcore.CameraViewConfig) int {
	maxOrder := 0
	for _, view := range views {
		if view.SortOrder > maxOrder {
			maxOrder = view.SortOrder
		}
	}
	return maxOrder + 10
}
