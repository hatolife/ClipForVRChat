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

const (
	oldDesktopFFmpegInputArgs = "-f gdigrab -framerate 30 -i desktop"
	oldTitleFFmpegInputArgs   = "-f gdigrab -framerate 30 -i title=VRChat"
)

const (
	userCameraZoomMin     = 20
	userCameraZoomMax     = 150
	userCameraZoomDefault = 45
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
	OSC         AutoCaptureOSCConfig         `json:"osc"`
	PlayerLocal AutoCapturePlayerLocalConfig `json:"playerLocal"`
	Schedule    AutoCaptureScheduleConfig    `json:"schedule"`
	Capture     AutoCaptureCaptureConfig     `json:"capture"`
	Stream      AutoCaptureStreamConfig      `json:"stream"`
	Restore     AutoCaptureRestoreConfig     `json:"restore"`
	Output      AutoCaptureOutputConfig      `json:"output"`
	Presence    AutoCapturePresenceConfig    `json:"presence"`
	Discord     AutoCaptureDiscordConfig     `json:"discord"`
	Views       []CameraViewConfig           `json:"views"`
}

type AutoCaptureOSCConfig struct {
	Host             string           `json:"vrcHost"`
	SendPort         int              `json:"vrcInPort"`
	ReceivePort      int              `json:"appOutPort"`
	PoseFreshnessSec int              `json:"poseFreshnessSec"`
	Forward          OSCForwardConfig `json:"forward,omitempty"`
}

type OSCForwardConfig struct {
	Enabled bool               `json:"enabled"`
	Mode    string             `json:"mode"`
	Targets []OSCForwardTarget `json:"targets"`
}

type OSCForwardTarget struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type AutoCapturePlayerLocalConfig struct {
	BasisPose   CameraPoseConfig                      `json:"basisPose"`
	BasisSource string                                `json:"basisSource,omitempty"`
	AvatarOSC   AutoCapturePlayerLocalAvatarOSCConfig `json:"avatarOsc,omitempty"`
	Calibrated  bool                                  `json:"calibrated"`
	UpdatedAt   string                                `json:"updatedAt,omitempty"`
}

type AutoCapturePlayerLocalAvatarOSCConfig struct {
	ParameterPrefix       string  `json:"parameterPrefix"`
	PositionScale         float64 `json:"positionScale"`
	InvertMagnitude       bool    `json:"invertMagnitude"`
	PositiveFlagThreshold float64 `json:"positiveFlagThreshold"`
	MaxAbsPosition        float64 `json:"maxAbsPosition"`
	MaxAbsForward         float64 `json:"maxAbsForward"`
	FreshnessSec          int     `json:"freshnessSec"`
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
	PreplacedLocalAnchor  bool   `json:"preplacedLocalAnchor"`
	OpenCameraBeforeBatch bool   `json:"openCameraBeforeBatch"`
	CloseCameraAfterBatch bool   `json:"closeCameraAfterBatch"`
	SettleDelayMS         int    `json:"settleDelayMs"`
	ButtonReleaseDelayMS  int    `json:"buttonReleaseDelayMs"`
}

type AutoCaptureStreamConfig struct {
	SpoutHelperPath         string `json:"spoutHelperPath"`
	SpoutSenderName         string `json:"spoutSenderName"`
	SpoutAutoSelect         bool   `json:"spoutAutoSelect"`
	CaptureTimeoutMS        int    `json:"captureTimeoutMs"`
	StartDelayMS            int    `json:"startDelayMs"`
	DebugRecordingEnabled   bool   `json:"debugRecordingEnabled"`
	DebugFrameCount         int    `json:"debugFrameCount"`
	DebugRecordingDirectory string `json:"-"`
	LegacyFFmpegPath        string `json:"legacyFfmpegPath,omitempty"`
	LegacyInputArgs         string `json:"legacyInputArgs,omitempty"`
}

type AutoCaptureRestoreConfig struct {
	Enabled              bool                                `json:"enabled"`
	PreferSnapshot       bool                                `json:"preferSnapshot"`
	SnapshotFreshnessSec int                                 `json:"snapshotFreshnessSec"`
	Fallback             AutoCaptureUserCameraFallbackConfig `json:"fallback"`
	Snapshot             AutoCaptureUserCameraState          `json:"-"`
}

type AutoCaptureUserCameraFallbackConfig struct {
	Mode                   int              `json:"mode"`
	Streaming              bool             `json:"streaming"`
	SmoothMovement         bool             `json:"smoothMovement"`
	RestorePose            bool             `json:"restorePose"`
	Pose                   CameraPoseConfig `json:"pose"`
	Zoom                   float64          `json:"zoom"`
	Exposure               float64          `json:"exposure"`
	FocalDistance          float64          `json:"focalDistance"`
	Aperture               float64          `json:"aperture"`
	Hue                    float64          `json:"hue"`
	Saturation             float64          `json:"saturation"`
	Lightness              float64          `json:"lightness"`
	LookAtMeXOffset        float64          `json:"lookAtMeXOffset"`
	LookAtMeYOffset        float64          `json:"lookAtMeYOffset"`
	FlySpeed               float64          `json:"flySpeed"`
	TurnSpeed              float64          `json:"turnSpeed"`
	SmoothingStrength      float64          `json:"smoothingStrength"`
	PhotoRate              float64          `json:"photoRate"`
	Duration               float64          `json:"duration"`
	ShowUIInCamera         bool             `json:"showUiInCamera"`
	Lock                   bool             `json:"lock"`
	LocalPlayer            bool             `json:"localPlayer"`
	RemotePlayer           bool             `json:"remotePlayer"`
	Environment            bool             `json:"environment"`
	GreenScreen            bool             `json:"greenScreen"`
	LookAtMe               bool             `json:"lookAtMe"`
	AutoLevelRoll          bool             `json:"autoLevelRoll"`
	AutoLevelPitch         bool             `json:"autoLevelPitch"`
	Flying                 bool             `json:"flying"`
	TriggerTakesPhotos     bool             `json:"triggerTakesPhotos"`
	DollyPathsStayVisible  bool             `json:"dollyPathsStayVisible"`
	CameraEars             bool             `json:"cameraEars"`
	ShowFocus              bool             `json:"showFocus"`
	RollWhileFlying        bool             `json:"rollWhileFlying"`
	OrientationIsLandscape bool             `json:"orientationIsLandscape"`
}

type AutoCaptureUserCameraState struct {
	Mode                   *int              `json:"mode,omitempty"`
	Pose                   *CameraPoseConfig `json:"pose,omitempty"`
	Streaming              *bool             `json:"streaming,omitempty"`
	SmoothMovement         *bool             `json:"smoothMovement,omitempty"`
	Zoom                   *float64          `json:"zoom,omitempty"`
	Exposure               *float64          `json:"exposure,omitempty"`
	FocalDistance          *float64          `json:"focalDistance,omitempty"`
	Aperture               *float64          `json:"aperture,omitempty"`
	Hue                    *float64          `json:"hue,omitempty"`
	Saturation             *float64          `json:"saturation,omitempty"`
	Lightness              *float64          `json:"lightness,omitempty"`
	LookAtMeXOffset        *float64          `json:"lookAtMeXOffset,omitempty"`
	LookAtMeYOffset        *float64          `json:"lookAtMeYOffset,omitempty"`
	FlySpeed               *float64          `json:"flySpeed,omitempty"`
	TurnSpeed              *float64          `json:"turnSpeed,omitempty"`
	SmoothingStrength      *float64          `json:"smoothingStrength,omitempty"`
	PhotoRate              *float64          `json:"photoRate,omitempty"`
	Duration               *float64          `json:"duration,omitempty"`
	ShowUIInCamera         *bool             `json:"showUiInCamera,omitempty"`
	Lock                   *bool             `json:"lock,omitempty"`
	LocalPlayer            *bool             `json:"localPlayer,omitempty"`
	RemotePlayer           *bool             `json:"remotePlayer,omitempty"`
	Environment            *bool             `json:"environment,omitempty"`
	GreenScreen            *bool             `json:"greenScreen,omitempty"`
	LookAtMe               *bool             `json:"lookAtMe,omitempty"`
	AutoLevelRoll          *bool             `json:"autoLevelRoll,omitempty"`
	AutoLevelPitch         *bool             `json:"autoLevelPitch,omitempty"`
	Flying                 *bool             `json:"flying,omitempty"`
	TriggerTakesPhotos     *bool             `json:"triggerTakesPhotos,omitempty"`
	DollyPathsStayVisible  *bool             `json:"dollyPathsStayVisible,omitempty"`
	CameraEars             *bool             `json:"cameraEars,omitempty"`
	ShowFocus              *bool             `json:"showFocus,omitempty"`
	RollWhileFlying        *bool             `json:"rollWhileFlying,omitempty"`
	OrientationIsLandscape *bool             `json:"orientationIsLandscape,omitempty"`
}

type AutoCaptureOutputConfig struct {
	Directory           string `json:"directory"`
	ImageFormat         string `json:"imageFormat"`
	FilenameTemplate    string `json:"filenameTemplate"`
	WriteSidecarJSON    bool   `json:"writeSidecarJson"`
	WriteEXIF           bool   `json:"writeExif"`
	WriteUserListToEXIF bool   `json:"writeUserListToExif"`
	WriteUserIDsToEXIF  bool   `json:"writeUserIdsToExif"`
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
			Forward: OSCForwardConfig{
				Enabled: false,
				Mode:    OSCForwardModeAll,
			},
		},
		PlayerLocal: AutoCapturePlayerLocalConfig{
			BasisSource: PlayerLocalBasisSourceAvatarOSC,
			AvatarOSC:   defaultAutoCapturePlayerLocalAvatarOSCConfig(),
		},
		Schedule: AutoCaptureScheduleConfig{
			Enabled:                    false,
			CaptureIntervalSec:         300,
			InitialDelaySec:            0,
			SkipIfPreviousBatchRunning: true,
			CaptureOnStart:             true,
		},
		Capture: AutoCaptureCaptureConfig{
			Mode:                  "stream",
			ConcurrentMode:        "sequential",
			RequestedCameraCount:  1,
			MultiBackend:          "dolly_multi",
			FallbackToSequential:  true,
			OpenCameraBeforeBatch: false,
			CloseCameraAfterBatch: false,
			SettleDelayMS:         1500,
			ButtonReleaseDelayMS:  200,
		},
		Stream: AutoCaptureStreamConfig{
			SpoutHelperPath:  "spout-capture.exe",
			SpoutAutoSelect:  true,
			CaptureTimeoutMS: 10000,
			StartDelayMS:     1000,
			DebugFrameCount:  8,
		},
		Restore: defaultAutoCaptureRestoreConfig(),
		Output: AutoCaptureOutputConfig{
			Directory:           DefaultAutoCaptureDirectory(),
			ImageFormat:         "png",
			FilenameTemplate:    "{timestamp_local}_{batch_id}_{shot_index}_{view_name}_{mode}.{ext}",
			WriteSidecarJSON:    true,
			WriteEXIF:           true,
			WriteUserListToEXIF: true,
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
	c.OSC.Forward.Normalize()
	c.PlayerLocal.Normalize()
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
	c.Stream.SpoutHelperPath = strings.Trim(strings.TrimSpace(c.Stream.SpoutHelperPath), `"`)
	if c.Stream.SpoutHelperPath == "" {
		c.Stream.SpoutHelperPath = "spout-capture.exe"
	}
	c.Stream.SpoutSenderName = strings.TrimSpace(c.Stream.SpoutSenderName)
	if c.Stream.SpoutSenderName == "" {
		c.Stream.SpoutAutoSelect = true
	}
	c.Stream.LegacyFFmpegPath = strings.Trim(strings.TrimSpace(c.Stream.LegacyFFmpegPath), `"`)
	if c.Stream.LegacyFFmpegPath == "" {
		c.Stream.LegacyFFmpegPath = "ffmpeg"
	}
	c.Stream.LegacyInputArgs = strings.TrimSpace(c.Stream.LegacyInputArgs)
	if c.Stream.LegacyInputArgs == oldDesktopFFmpegInputArgs || c.Stream.LegacyInputArgs == oldTitleFFmpegInputArgs {
		c.Stream.LegacyInputArgs = DefaultAutoCaptureFFmpegInputArgs()
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
	if c.Stream.DebugFrameCount <= 0 {
		c.Stream.DebugFrameCount = 8
	}
	if c.Stream.DebugFrameCount > 120 {
		c.Stream.DebugFrameCount = 120
	}
	c.Stream.DebugRecordingDirectory = strings.Trim(strings.TrimSpace(c.Stream.DebugRecordingDirectory), `"`)
	c.Restore.Normalize()
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
	if !c.Output.WriteEXIF {
		c.Output.WriteUserListToEXIF = false
		c.Output.WriteUserIDsToEXIF = false
	}
	if !c.Output.WriteUserListToEXIF {
		c.Output.WriteUserIDsToEXIF = false
	}
	c.Presence.OutputLogDirectory = strings.Trim(strings.TrimSpace(c.Presence.OutputLogDirectory), `"`)
	if c.Presence.OutputLogDirectory == "" {
		c.Presence.OutputLogDirectory = DefaultVRChatLogDirectory()
	}
	c.Discord.WebhookURL = strings.Trim(strings.TrimSpace(c.Discord.WebhookURL), `"`)
	c.Discord.PostMode = "shot"
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

const (
	OSCForwardModeAll           = "all"
	OSCForwardModeUnhandledOnly = "unhandled_only"
)

func (c *OSCForwardConfig) Normalize() {
	switch c.Mode {
	case OSCForwardModeAll, OSCForwardModeUnhandledOnly:
	default:
		c.Mode = OSCForwardModeAll
	}
	if c.Targets == nil {
		c.Targets = []OSCForwardTarget{}
	}
	for i := range c.Targets {
		c.Targets[i].Normalize()
	}
}

func (t *OSCForwardTarget) Normalize() {
	t.Host = strings.TrimSpace(strings.Trim(strings.TrimSpace(t.Host), `"`))
	if t.Host == "" {
		t.Host = "127.0.0.1"
	}
}

func defaultAutoCaptureRestoreConfig() AutoCaptureRestoreConfig {
	return AutoCaptureRestoreConfig{
		Enabled:              true,
		PreferSnapshot:       true,
		SnapshotFreshnessSec: 10,
		Fallback:             defaultAutoCaptureUserCameraFallbackConfig(),
	}
}

func defaultAutoCaptureUserCameraFallbackConfig() AutoCaptureUserCameraFallbackConfig {
	return AutoCaptureUserCameraFallbackConfig{
		Mode:                   0,
		Streaming:              false,
		SmoothMovement:         true,
		RestorePose:            false,
		Zoom:                   45,
		Exposure:               4,
		FocalDistance:          1.5,
		Aperture:               15,
		Hue:                    120,
		Saturation:             100,
		Lightness:              60,
		LookAtMeXOffset:        0,
		LookAtMeYOffset:        0,
		FlySpeed:               3,
		TurnSpeed:              1,
		SmoothingStrength:      5,
		PhotoRate:              1,
		Duration:               2,
		ShowUIInCamera:         false,
		Lock:                   false,
		LocalPlayer:            true,
		RemotePlayer:           true,
		Environment:            true,
		GreenScreen:            false,
		LookAtMe:               false,
		AutoLevelRoll:          false,
		AutoLevelPitch:         false,
		Flying:                 false,
		TriggerTakesPhotos:     false,
		DollyPathsStayVisible:  false,
		CameraEars:             false,
		ShowFocus:              false,
		RollWhileFlying:        false,
		OrientationIsLandscape: true,
	}
}

func (c *AutoCaptureRestoreConfig) Normalize() {
	if c.isZero() {
		*c = defaultAutoCaptureRestoreConfig()
	}
	if c.SnapshotFreshnessSec <= 0 {
		c.SnapshotFreshnessSec = 10
	}
	if c.SnapshotFreshnessSec > 300 {
		c.SnapshotFreshnessSec = 300
	}
	c.Fallback.Normalize()
}

func (c AutoCaptureRestoreConfig) isZero() bool {
	return !c.Enabled &&
		!c.PreferSnapshot &&
		c.SnapshotFreshnessSec == 0 &&
		c.Fallback.isZero()
}

func (c *AutoCaptureUserCameraFallbackConfig) Normalize() {
	if c.isZero() {
		*c = defaultAutoCaptureUserCameraFallbackConfig()
	}
	if c.Mode < 0 || c.Mode > 6 {
		c.Mode = 0
	}
	c.Zoom = clampFiniteDefault(c.Zoom, userCameraZoomMin, userCameraZoomMax, userCameraZoomDefault)
	c.Exposure = clampFiniteDefault(c.Exposure, 0, 10, 4)
	c.FocalDistance = clampFiniteDefault(c.FocalDistance, 0, 10, 1.5)
	c.Aperture = clampFiniteDefault(c.Aperture, 1.4, 32, 15)
	c.Hue = clampFiniteDefault(c.Hue, 0, 360, 120)
	c.Saturation = clampFiniteDefault(c.Saturation, 0, 100, 100)
	c.Lightness = clampFiniteDefault(c.Lightness, 0, 50, 50)
	c.LookAtMeXOffset = clampFiniteDefault(c.LookAtMeXOffset, -25, 25, 0)
	c.LookAtMeYOffset = clampFiniteDefault(c.LookAtMeYOffset, -25, 25, 0)
	c.FlySpeed = clampFiniteDefault(c.FlySpeed, 0.1, 15, 3)
	c.TurnSpeed = clampFiniteDefault(c.TurnSpeed, 0.1, 5, 1)
	c.SmoothingStrength = clampFiniteDefault(c.SmoothingStrength, 0.1, 10, 5)
	c.PhotoRate = clampFiniteDefault(c.PhotoRate, 0.1, 2, 1)
	c.Duration = clampFiniteDefault(c.Duration, 0.1, 60, 2)
}

func (c AutoCaptureUserCameraFallbackConfig) isZero() bool {
	return c.Mode == 0 &&
		!c.Streaming &&
		!c.SmoothMovement &&
		!c.RestorePose &&
		c.Pose.isZero() &&
		c.Zoom == 0 &&
		c.Exposure == 0 &&
		c.FocalDistance == 0 &&
		c.Aperture == 0 &&
		c.Hue == 0 &&
		c.Saturation == 0 &&
		c.Lightness == 0 &&
		c.LookAtMeXOffset == 0 &&
		c.LookAtMeYOffset == 0 &&
		c.FlySpeed == 0 &&
		c.TurnSpeed == 0 &&
		c.SmoothingStrength == 0 &&
		c.PhotoRate == 0 &&
		c.Duration == 0 &&
		!c.ShowUIInCamera &&
		!c.Lock &&
		!c.LocalPlayer &&
		!c.RemotePlayer &&
		!c.Environment &&
		!c.GreenScreen &&
		!c.LookAtMe &&
		!c.AutoLevelRoll &&
		!c.AutoLevelPitch &&
		!c.Flying &&
		!c.TriggerTakesPhotos &&
		!c.DollyPathsStayVisible &&
		!c.CameraEars &&
		!c.ShowFocus &&
		!c.RollWhileFlying &&
		!c.OrientationIsLandscape
}

func clampFiniteDefault(value float64, min float64, max float64, fallback float64) float64 {
	if !isFiniteFloat64(value) || value < min || value > max {
		return fallback
	}
	return value
}

const (
	PlayerLocalBasisSourceManual    = "manual"
	PlayerLocalBasisSourceAvatarOSC = "avatar_osc"
)

func (c *AutoCapturePlayerLocalConfig) Normalize() {
	c.BasisSource = normalizePlayerLocalBasisSource(c.BasisSource)
	c.AvatarOSC.Normalize()
}

func normalizePlayerLocalBasisSource(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "":
		return PlayerLocalBasisSourceAvatarOSC
	case PlayerLocalBasisSourceManual:
		return PlayerLocalBasisSourceManual
	case PlayerLocalBasisSourceAvatarOSC:
		return PlayerLocalBasisSourceAvatarOSC
	default:
		return PlayerLocalBasisSourceManual
	}
}

func (c *AutoCapturePlayerLocalAvatarOSCConfig) Normalize() {
	c.ParameterPrefix = strings.TrimSpace(strings.Trim(strings.TrimSpace(c.ParameterPrefix), `"`))
	if c.ParameterPrefix == "" || c.ParameterPrefix == "CFVRC/basis" {
		c.ParameterPrefix = "coord"
	}
	if c.PositionScale <= 0 {
		c.PositionScale = 1000
	}
	if !isFiniteFloat64(c.PositionScale) {
		c.PositionScale = 1000
	}
	if !isFiniteFloat64(c.PositiveFlagThreshold) {
		c.PositiveFlagThreshold = 0
	}
	if c.MaxAbsPosition <= 0 || !isFiniteFloat64(c.MaxAbsPosition) {
		c.MaxAbsPosition = 10000
	}
	if c.MaxAbsForward <= 0 || !isFiniteFloat64(c.MaxAbsForward) {
		c.MaxAbsForward = 2000
	}
	if c.FreshnessSec <= 0 {
		c.FreshnessSec = 3
	}
	if c.FreshnessSec > 60 {
		c.FreshnessSec = 60
	}
}

func defaultAutoCapturePlayerLocalAvatarOSCConfig() AutoCapturePlayerLocalAvatarOSCConfig {
	cfg := AutoCapturePlayerLocalAvatarOSCConfig{
		ParameterPrefix:       "coord",
		PositionScale:         1000,
		InvertMagnitude:       true,
		PositiveFlagThreshold: 0,
		MaxAbsPosition:        10000,
		MaxAbsForward:         2000,
		FreshnessSec:          3,
	}
	cfg.Normalize()
	return cfg
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
	case "world", "player_local":
	default:
		v.CoordinateSpace = "world"
	}
	if v.SettleDelayMS < 1500 {
		v.SettleDelayMS = 1500
	}
	if v.CaptureDelayMS < 0 {
		v.CaptureDelayMS = 0
	}
	if v.Zoom != nil {
		zoom := clampFiniteDefault(*v.Zoom, userCameraZoomMin, userCameraZoomMax, userCameraZoomDefault)
		v.Zoom = float64ConfigPtr(zoom)
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

func DefaultAutoCaptureFFmpegInputArgs() string {
	return "-f gdigrab -framerate 30 -offset_x {window_x} -offset_y {window_y} -video_size {window_width}x{window_height} -i desktop"
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
		defaultCameraView("front", "正面", 0, "player_local", CameraPoseConfig{
			Position: CameraVector3Config{X: 0, Y: 0, Z: 1.0},
			Rotation: CameraVector3Config{X: 0, Y: 180, Z: 0},
		}, userCameraZoomDefault),
		defaultCameraView("back", "背後", 1, "player_local", CameraPoseConfig{
			Position: CameraVector3Config{X: 0, Y: 0.35, Z: -1.6},
			Rotation: CameraVector3Config{X: 12, Y: 0, Z: 0},
		}, userCameraZoomDefault),
		defaultCameraView("diagonal", "斜め", 2, "player_local", CameraPoseConfig{
			Position: CameraVector3Config{X: 0.8, Y: 0.2, Z: 1.1},
			Rotation: CameraVector3Config{X: 8, Y: -145, Z: 0},
		}, userCameraZoomDefault),
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

func defaultCameraView(id string, name string, order int, coordinateSpace string, pose CameraPoseConfig, zoom float64) CameraViewConfig {
	return CameraViewConfig{
		ID:              id,
		Name:            name,
		Enabled:         true,
		SortOrder:       order,
		CoordinateSpace: coordinateSpace,
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
