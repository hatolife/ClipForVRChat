package appcore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	privateDirMode  os.FileMode = 0700
	privateFileMode os.FileMode = 0600
)

type Config struct {
	Image              ImageConfig              `json:"image"`
	Output             OutputConfig             `json:"output"`
	Discord            DiscordConfig            `json:"discord"`
	AutoPhoto          AutoPhotoConfig          `json:"autoPhoto"`
	AutoCapture        AutoCaptureConfig        `json:"autoCapture"`
	ScreenshotAutoPost ScreenshotAutoPostConfig `json:"screenshotAutoPost"`
	Update             UpdateConfig             `json:"update"`
	DiagnosticLogPath  string                   `json:"-"`
}

type ImageConfig struct {
	MaxWidth        int    `json:"maxWidth"`
	MaxHeight       int    `json:"maxHeight"`
	MaxInputMB      int    `json:"maxInputMb"`
	Suffix          string `json:"suffix"`
	OutputFormat    string `json:"outputFormat"`
	Overwrite       bool   `json:"overwrite"`
	JPEGQuality     int    `json:"jpegQuality"`
	OutputDirectory string `json:"outputDirectory"`
}

type OutputConfig struct {
	SaveLocal                  bool   `json:"saveLocal"`
	UploadDiscord              bool   `json:"uploadDiscord"`
	ShowUI                     string `json:"showUi"`
	CopySingleURLToClipboard   bool   `json:"copySingleUrlToClipboard"`
	DeleteOutputOnHistoryPurge bool   `json:"deleteOutputOnHistoryPurge"`
	DetectQRCodeURLs           bool   `json:"detectQrCodeUrls"`
}

type DiscordConfig struct {
	WebhookURL string `json:"webhookUrl"`
}

type AutoPhotoConfig struct {
	Enabled             bool   `json:"enabled"`
	PhotoDirectory      string `json:"photoDirectory"`
	WebhookURL          string `json:"webhookUrl"`
	ScanIntervalSeconds int    `json:"scanIntervalSeconds"`
}

type ScreenshotAutoPostConfig struct {
	Enabled             bool   `json:"enabled"`
	ScreenshotDirectory string `json:"screenshotDirectory"`
	WebhookURL          string `json:"webhookUrl"`
	ScanIntervalSeconds int    `json:"scanIntervalSeconds"`
}

type UpdateConfig struct {
	CheckEnabled        bool `json:"checkEnabled"`
	NotificationEnabled bool `json:"notificationEnabled"`
}

type AutoCaptureConfig struct {
	OSC      AutoCaptureOSCConfig      `json:"osc"`
	Schedule AutoCaptureScheduleConfig `json:"schedule"`
	Capture  AutoCaptureCaptureConfig  `json:"capture"`
	Stream   AutoCaptureStreamConfig   `json:"stream"`
	Output   AutoCaptureOutputConfig   `json:"output"`
	Presence AutoCapturePresenceConfig `json:"presence"`
	Discord  AutoCaptureDiscordConfig  `json:"discord"`
	Views    []CameraViewConfig        `json:"views"`
}

type AutoCaptureOSCConfig struct {
	Host             string `json:"vrcHost"`
	SendPort         int    `json:"vrcInPort"`
	ReceivePort      int    `json:"appOutPort"`
	PoseFreshnessSec int    `json:"poseFreshnessSec"`
}

type AutoCaptureScheduleConfig struct {
	Enabled                    bool `json:"enabled"`
	CaptureIntervalSec         int  `json:"captureIntervalSec"`
	InitialDelaySec            int  `json:"initialDelaySec"`
	SkipIfPreviousBatchRunning bool `json:"skipIfPreviousBatchRunning"`
	CaptureOnStart             bool `json:"captureOnStart"`
	MaxBatches                 int  `json:"maxBatches"`
}

type AutoCaptureCaptureConfig struct {
	Mode                  string `json:"mode"`
	ConcurrentMode        string `json:"concurrentMode"`
	RequestedCameraCount  int    `json:"requestedCameraCount"`
	MultiBackend          string `json:"multiBackend"`
	FallbackToSequential  bool   `json:"fallbackToSequential"`
	CloseCameraAfterBatch bool   `json:"closeCameraAfterBatch"`
	SettleDelayMS         int    `json:"settleDelayMs"`
	ButtonReleaseDelayMS  int    `json:"buttonReleaseDelayMs"`
}

type AutoCaptureStreamConfig struct {
	FFmpegPath       string `json:"ffmpegPath"`
	InputArgs        string `json:"inputArgs"`
	CaptureTimeoutMS int    `json:"captureTimeoutMs"`
	StartDelayMS     int    `json:"startDelayMs"`
}

type AutoCaptureOutputConfig struct {
	Directory           string `json:"directory"`
	ImageFormat         string `json:"imageFormat"`
	FilenameTemplate    string `json:"filenameTemplate"`
	WriteSidecarJSON    bool   `json:"writeSidecarJson"`
	WriteEXIF           bool   `json:"writeExif"`
	WriteUserListToEXIF bool   `json:"writeUserListToExif"`
}

type AutoCapturePresenceConfig struct {
	WatchOutputLog               bool   `json:"watchOutputLog"`
	OutputLogDirectory           string `json:"outputLogDirectory"`
	IncludeUserIDsInSidecar      bool   `json:"includeUserIdsInSidecar"`
	IncludeUserIDsInDiscord      bool   `json:"includeUserIdsInDiscord"`
	IncludeDisplayNamesInDiscord bool   `json:"includeDisplayNamesInDiscord"`
}

type AutoCaptureDiscordConfig struct {
	Enabled       bool   `json:"enabled"`
	WebhookURL    string `json:"webhookUrl"`
	PostMode      string `json:"postMode"`
	IncludeImages bool   `json:"includeImages"`
}

type CameraViewConfig struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Enabled         bool             `json:"enabled"`
	SortOrder       int              `json:"sortOrder"`
	CoordinateSpace string           `json:"coordinateSpace"`
	Pose            CameraPoseConfig `json:"pose"`
	Zoom            *float64         `json:"zoom,omitempty"`
	Exposure        *float64         `json:"exposure,omitempty"`
	FocalDistance   *float64         `json:"focalDistance,omitempty"`
	Aperture        *float64         `json:"aperture,omitempty"`
	LookAtMe        *bool            `json:"lookAtMe,omitempty"`
	ShowUIInCamera  *bool            `json:"showUiInCamera,omitempty"`
	LocalPlayer     *bool            `json:"localPlayer,omitempty"`
	RemotePlayer    *bool            `json:"remotePlayer,omitempty"`
	Environment     *bool            `json:"environment,omitempty"`
	SettleDelayMS   int              `json:"settleDelayMs"`
	CaptureDelayMS  int              `json:"captureDelayMs"`
	Calibrated      bool             `json:"calibrated"`
}

type CameraPoseConfig struct {
	Position CameraVector3Config `json:"position"`
	Rotation CameraVector3Config `json:"rotation"`
}

type CameraVector3Config struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (p CameraPoseConfig) isZero() bool {
	return p.Position.X == 0 && p.Position.Y == 0 && p.Position.Z == 0 &&
		p.Rotation.X == 0 && p.Rotation.Y == 0 && p.Rotation.Z == 0
}

func DefaultConfig() Config {
	return Config{
		Image: ImageConfig{
			MaxWidth:        2048,
			MaxHeight:       2048,
			MaxInputMB:      32,
			Suffix:          "_2048",
			OutputFormat:    "png",
			Overwrite:       false,
			JPEGQuality:     92,
			OutputDirectory: "./output",
		},
		Output: OutputConfig{
			SaveLocal:                  true,
			UploadDiscord:              false,
			ShowUI:                     "auto",
			CopySingleURLToClipboard:   false,
			DeleteOutputOnHistoryPurge: true,
			DetectQRCodeURLs:           false,
		},
		AutoPhoto: AutoPhotoConfig{
			Enabled:             false,
			PhotoDirectory:      DefaultVRChatPhotoDirectory(),
			ScanIntervalSeconds: 2,
		},
		AutoCapture: DefaultAutoCaptureConfig(),
		ScreenshotAutoPost: ScreenshotAutoPostConfig{
			Enabled:             false,
			ScreenshotDirectory: DefaultScreenshotsDirectory(),
			ScanIntervalSeconds: 2,
		},
		Update: UpdateConfig{
			CheckEnabled:        true,
			NotificationEnabled: true,
		},
	}
}

func DefaultAutoCaptureConfig() AutoCaptureConfig {
	return AutoCaptureConfig{
		OSC: AutoCaptureOSCConfig{
			Host:             "127.0.0.1",
			SendPort:         9000,
			ReceivePort:      9001,
			PoseFreshnessSec: 3,
		},
		Schedule: AutoCaptureScheduleConfig{
			Enabled:                    false,
			CaptureIntervalSec:         300,
			InitialDelaySec:            0,
			SkipIfPreviousBatchRunning: true,
			CaptureOnStart:             false,
		},
		Capture: AutoCaptureCaptureConfig{
			Mode:                  "stream",
			ConcurrentMode:        "sequential",
			RequestedCameraCount:  1,
			MultiBackend:          "dolly_multi",
			FallbackToSequential:  true,
			CloseCameraAfterBatch: true,
			SettleDelayMS:         1500,
			ButtonReleaseDelayMS:  200,
		},
		Stream: AutoCaptureStreamConfig{
			FFmpegPath:       "ffmpeg",
			InputArgs:        "-f gdigrab -framerate 30 -i desktop",
			CaptureTimeoutMS: 10000,
			StartDelayMS:     1000,
		},
		Output: AutoCaptureOutputConfig{
			Directory:        DefaultAutoCaptureDirectory(),
			ImageFormat:      "png",
			FilenameTemplate: "{timestamp_local}_{batch_id}_{shot_index}_{view_name}_{mode}.{ext}",
			WriteSidecarJSON: true,
		},
		Presence: AutoCapturePresenceConfig{
			WatchOutputLog:          true,
			OutputLogDirectory:      DefaultVRChatLogDirectory(),
			IncludeUserIDsInSidecar: true,
		},
		Discord: AutoCaptureDiscordConfig{
			Enabled:       false,
			PostMode:      "shot",
			IncludeImages: true,
		},
		Views: defaultCameraViews(),
	}
}

func ConfigPath(exePath string) string {
	return filepath.Join(filepath.Dir(exePath), "config.json")
}

func ConfigExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	cfg.DiagnosticLogPath = DiagnosticLogPath(path)
	data, err := os.ReadFile(path) // #nosec G304 -- config path is the app config path or an explicitly opened local config file.
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	cfg.DiagnosticLogPath = DiagnosticLogPath(path)
	cfg.Normalize()
	return cfg, nil
}

func SaveConfig(path string, cfg Config) error {
	cfg.Normalize()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), privateDirMode); err != nil {
		return err
	}
	return WritePrivateFile(path, append(data, '\n'))
}

func WritePrivateFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, privateFileMode); err != nil {
		return err
	}
	return os.Chmod(path, privateFileMode)
}

func (c *Config) Normalize() {
	if c.Image.MaxWidth <= 0 {
		c.Image.MaxWidth = 2048
	}
	if c.Image.MaxHeight <= 0 {
		c.Image.MaxHeight = 2048
	}
	if c.Image.MaxInputMB <= 0 {
		c.Image.MaxInputMB = 32
	}
	if c.Image.Suffix == "" {
		c.Image.Suffix = "_2048"
	}
	c.Image.OutputDirectory = strings.Trim(strings.TrimSpace(c.Image.OutputDirectory), `"`)
	if c.Image.OutputDirectory == "" {
		c.Image.OutputDirectory = "./output"
	}
	switch c.Image.OutputFormat {
	case "png", "jpg":
	default:
		c.Image.OutputFormat = "png"
	}
	if c.Image.JPEGQuality <= 0 || c.Image.JPEGQuality > 100 {
		c.Image.JPEGQuality = 92
	}
	switch c.Output.ShowUI {
	case "auto", "always", "never":
	default:
		c.Output.ShowUI = "auto"
	}
	c.AutoPhoto.PhotoDirectory = strings.Trim(strings.TrimSpace(c.AutoPhoto.PhotoDirectory), `"`)
	if c.AutoPhoto.PhotoDirectory == "" {
		c.AutoPhoto.PhotoDirectory = DefaultVRChatPhotoDirectory()
	}
	c.AutoPhoto.WebhookURL = strings.Trim(strings.TrimSpace(c.AutoPhoto.WebhookURL), `"`)
	if c.AutoPhoto.ScanIntervalSeconds <= 0 {
		c.AutoPhoto.ScanIntervalSeconds = 2
	}
	if c.AutoPhoto.ScanIntervalSeconds > 3600 {
		c.AutoPhoto.ScanIntervalSeconds = 3600
	}
	c.ScreenshotAutoPost.ScreenshotDirectory = strings.Trim(strings.TrimSpace(c.ScreenshotAutoPost.ScreenshotDirectory), `"`)
	if c.ScreenshotAutoPost.ScreenshotDirectory == "" {
		c.ScreenshotAutoPost.ScreenshotDirectory = DefaultScreenshotsDirectory()
	}
	c.ScreenshotAutoPost.WebhookURL = strings.Trim(strings.TrimSpace(c.ScreenshotAutoPost.WebhookURL), `"`)
	if c.ScreenshotAutoPost.ScanIntervalSeconds <= 0 {
		c.ScreenshotAutoPost.ScanIntervalSeconds = 2
	}
	if c.ScreenshotAutoPost.ScanIntervalSeconds > 3600 {
		c.ScreenshotAutoPost.ScanIntervalSeconds = 3600
	}
	c.AutoCapture.Normalize()
}

func (c *AutoCaptureConfig) Normalize() {
	c.OSC.Host = strings.TrimSpace(c.OSC.Host)
	if c.OSC.Host == "" {
		c.OSC.Host = "127.0.0.1"
	}
	if c.OSC.SendPort <= 0 || c.OSC.SendPort > 65535 {
		c.OSC.SendPort = 9000
	}
	if c.OSC.ReceivePort <= 0 || c.OSC.ReceivePort > 65535 {
		c.OSC.ReceivePort = 9001
	}
	if c.OSC.PoseFreshnessSec <= 0 {
		c.OSC.PoseFreshnessSec = 3
	}
	if c.Schedule.CaptureIntervalSec <= 0 {
		c.Schedule.CaptureIntervalSec = 300
	}
	if c.Schedule.CaptureIntervalSec < 10 {
		c.Schedule.CaptureIntervalSec = 10
	}
	if c.Schedule.CaptureIntervalSec > 86400 {
		c.Schedule.CaptureIntervalSec = 86400
	}
	if c.Schedule.InitialDelaySec < 0 {
		c.Schedule.InitialDelaySec = 0
	}
	if c.Schedule.MaxBatches < 0 {
		c.Schedule.MaxBatches = 0
	}
	switch c.Capture.Mode {
	case "photo", "stream":
	default:
		c.Capture.Mode = "stream"
	}
	switch c.Capture.ConcurrentMode {
	case "sequential", "multi":
	default:
		c.Capture.ConcurrentMode = "sequential"
	}
	if c.Capture.RequestedCameraCount <= 0 {
		c.Capture.RequestedCameraCount = 1
	}
	if c.Capture.RequestedCameraCount > 4 {
		c.Capture.RequestedCameraCount = 4
	}
	if c.Capture.MultiBackend == "" {
		c.Capture.MultiBackend = "dolly_multi"
	}
	if c.Capture.SettleDelayMS < 1500 {
		c.Capture.SettleDelayMS = 1500
	}
	if c.Capture.ButtonReleaseDelayMS < 200 {
		c.Capture.ButtonReleaseDelayMS = 200
	}
	c.Stream.FFmpegPath = strings.Trim(strings.TrimSpace(c.Stream.FFmpegPath), `"`)
	if c.Stream.FFmpegPath == "" {
		c.Stream.FFmpegPath = "ffmpeg"
	}
	c.Stream.InputArgs = strings.TrimSpace(c.Stream.InputArgs)
	if c.Stream.InputArgs == "" {
		c.Stream.InputArgs = "-f gdigrab -framerate 30 -i desktop"
	}
	if c.Stream.CaptureTimeoutMS <= 0 {
		c.Stream.CaptureTimeoutMS = 10000
	}
	if c.Stream.CaptureTimeoutMS < 1000 {
		c.Stream.CaptureTimeoutMS = 1000
	}
	if c.Stream.CaptureTimeoutMS > 60000 {
		c.Stream.CaptureTimeoutMS = 60000
	}
	if c.Stream.StartDelayMS < 0 {
		c.Stream.StartDelayMS = 0
	}
	if c.Stream.StartDelayMS > 10000 {
		c.Stream.StartDelayMS = 10000
	}
	c.Output.Directory = strings.Trim(strings.TrimSpace(c.Output.Directory), `"`)
	if c.Output.Directory == "" {
		c.Output.Directory = DefaultAutoCaptureDirectory()
	}
	switch c.Output.ImageFormat {
	case "png", "jpg", "jpeg":
	default:
		c.Output.ImageFormat = "png"
	}
	if c.Output.FilenameTemplate == "" {
		c.Output.FilenameTemplate = "{timestamp_local}_{batch_id}_{shot_index}_{view_name}_{mode}.{ext}"
	}
	c.Presence.OutputLogDirectory = strings.Trim(strings.TrimSpace(c.Presence.OutputLogDirectory), `"`)
	if c.Presence.OutputLogDirectory == "" {
		c.Presence.OutputLogDirectory = DefaultVRChatLogDirectory()
	}
	c.Discord.WebhookURL = strings.Trim(strings.TrimSpace(c.Discord.WebhookURL), `"`)
	switch c.Discord.PostMode {
	case "shot", "batch":
	default:
		c.Discord.PostMode = "shot"
	}
	if len(c.Views) == 0 {
		c.Views = defaultCameraViews()
	}
	for i := range c.Views {
		c.Views[i].Normalize(i)
		if defaultView, ok := DefaultCameraViewByID(c.Views[i].ID); ok && !c.Views[i].Calibrated {
			if c.Views[i].CoordinateSpace == "template_relative" || c.Views[i].Pose.isZero() {
				c.Views[i].CoordinateSpace = defaultView.CoordinateSpace
				c.Views[i].Pose = defaultView.Pose
			}
			if c.Views[i].Zoom == nil {
				c.Views[i].Zoom = defaultView.Zoom
			}
		}
	}
}

func (v *CameraViewConfig) Normalize(index int) {
	v.ID = strings.TrimSpace(v.ID)
	if v.ID == "" {
		v.ID = "view"
	}
	v.Name = strings.TrimSpace(v.Name)
	if v.Name == "" {
		v.Name = v.ID
	}
	if v.SortOrder < 0 {
		v.SortOrder = index
	}
	switch v.CoordinateSpace {
	case "world", "dolly_local", "template_relative":
	default:
		v.CoordinateSpace = "world"
	}
	if v.SettleDelayMS < 1500 {
		v.SettleDelayMS = 1500
	}
	if v.CaptureDelayMS < 0 {
		v.CaptureDelayMS = 0
	}
}

func DefaultVRChatPhotoDirectory() string {
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return filepath.Join(userProfile, "Pictures", "VRChat")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "Pictures", "VRChat")
	}
	return ""
}

func DefaultScreenshotsDirectory() string {
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return filepath.Join(userProfile, "Pictures", "Screenshots")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "Pictures", "Screenshots")
	}
	return ""
}

func DefaultAutoCaptureDirectory() string {
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return filepath.Join(userProfile, "Pictures", "VRChat", "VRC-AutoCapture")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "Pictures", "VRChat", "VRC-AutoCapture")
	}
	return ""
}

func DefaultVRChatLogDirectory() string {
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		return filepath.Join(filepath.Dir(localAppData), "LocalLow", "VRChat", "VRChat")
	}
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return filepath.Join(userProfile, "AppData", "LocalLow", "VRChat", "VRChat")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, "AppData", "LocalLow", "VRChat", "VRChat")
	}
	return ""
}

func defaultCameraViews() []CameraViewConfig {
	return []CameraViewConfig{
		defaultCameraView("front", "正面", 0, CameraPoseConfig{
			Position: CameraVector3Config{X: 0, Y: 1.55, Z: 0.9},
			Rotation: CameraVector3Config{X: 0, Y: 180, Z: 0},
		}, 1.0),
		defaultCameraView("back", "背後", 1, CameraPoseConfig{
			Position: CameraVector3Config{X: 0, Y: 1.85, Z: -1.6},
			Rotation: CameraVector3Config{X: 12, Y: 0, Z: 0},
		}, 1.0),
		defaultCameraView("diagonal", "斜め", 2, CameraPoseConfig{
			Position: CameraVector3Config{X: 0.8, Y: 1.75, Z: 1.1},
			Rotation: CameraVector3Config{X: 8, Y: -145, Z: 0},
		}, 1.0),
	}
}

func DefaultCameraViewByID(id string) (CameraViewConfig, bool) {
	for _, view := range defaultCameraViews() {
		if view.ID == id {
			return view, true
		}
	}
	return CameraViewConfig{}, false
}

func DefaultCameraViews() []CameraViewConfig {
	return defaultCameraViews()
}

func defaultCameraView(id string, name string, order int, pose CameraPoseConfig, zoom float64) CameraViewConfig {
	return CameraViewConfig{
		ID:              id,
		Name:            name,
		Enabled:         true,
		SortOrder:       order,
		CoordinateSpace: "world",
		Pose:            pose,
		Zoom:            float64ConfigPtr(zoom),
		LookAtMe:        boolConfigPtr(true),
		LocalPlayer:     boolConfigPtr(true),
		RemotePlayer:    boolConfigPtr(false),
		Environment:     boolConfigPtr(true),
		SettleDelayMS:   1500,
		CaptureDelayMS:  0,
		Calibrated:      false,
	}
}

func float64ConfigPtr(value float64) *float64 {
	return &value
}

func boolConfigPtr(value bool) *bool {
	return &value
}
