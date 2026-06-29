package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

type configLogSummary struct {
	Image              imageConfigLogSummary              `json:"image"`
	Output             appcore.OutputConfig               `json:"output"`
	Discord            webhookConfigLogSummary            `json:"discord"`
	AutoPhoto          autoPhotoConfigLogSummary          `json:"autoPhoto"`
	AutoCapture        autoCaptureConfigLogSummary        `json:"autoCapture"`
	ScreenshotAutoPost screenshotAutoPostConfigLogSummary `json:"screenshotAutoPost"`
	Update             appcore.UpdateConfig               `json:"update"`
}

type imageConfigLogSummary struct {
	MaxWidth        int    `json:"maxWidth"`
	MaxHeight       int    `json:"maxHeight"`
	MaxInputMB      int    `json:"maxInputMb"`
	Suffix          string `json:"suffix"`
	OutputFormat    string `json:"outputFormat"`
	Overwrite       bool   `json:"overwrite"`
	JPEGQuality     int    `json:"jpegQuality"`
	OutputDirectory string `json:"outputDirectory"`
}

type webhookConfigLogSummary struct {
	WebhookConfigured bool `json:"webhookConfigured"`
}

type autoPhotoConfigLogSummary struct {
	Enabled             bool   `json:"enabled"`
	PhotoDirectory      string `json:"photoDirectory"`
	WebhookConfigured   bool   `json:"webhookConfigured"`
	ScanIntervalSeconds int    `json:"scanIntervalSeconds"`
}

type screenshotAutoPostConfigLogSummary struct {
	Enabled             bool   `json:"enabled"`
	ScreenshotDirectory string `json:"screenshotDirectory"`
	WebhookConfigured   bool   `json:"webhookConfigured"`
	ScanIntervalSeconds int    `json:"scanIntervalSeconds"`
}

type autoCaptureConfigLogSummary struct {
	Schedule        appcore.AutoCaptureScheduleConfig `json:"schedule"`
	Capture         appcore.AutoCaptureCaptureConfig  `json:"capture"`
	OSC             appcore.AutoCaptureOSCConfig      `json:"osc"`
	Output          autoCaptureOutputLogSummary       `json:"output"`
	Presence        appcore.AutoCapturePresenceConfig `json:"presence"`
	Discord         autoCaptureDiscordLogSummary      `json:"discord"`
	ViewCount       int                               `json:"viewCount"`
	EnabledViews    int                               `json:"enabledViews"`
	CalibratedViews int                               `json:"calibratedViews"`
}

type autoCaptureOutputLogSummary struct {
	Directory           string `json:"directory"`
	ImageFormat         string `json:"imageFormat"`
	WriteSidecarJSON    bool   `json:"writeSidecarJson"`
	WriteEXIF           bool   `json:"writeExif"`
	WriteUserListToEXIF bool   `json:"writeUserListToExif"`
}

type autoCaptureDiscordLogSummary struct {
	Enabled           bool   `json:"enabled"`
	WebhookConfigured bool   `json:"webhookConfigured"`
	PostMode          string `json:"postMode"`
	IncludeImages     bool   `json:"includeImages"`
}

func (a *App) logStartupLocked() {
	exePath, err := os.Executable()
	if err != nil {
		exePath = ""
	}
	exeSHA256 := ""
	if exePath != "" {
		if hash, err := fileSHA256(exePath); err == nil {
			exeSHA256 = hash
		}
	}
	appcore.AppendDiagnosticLog(
		appcore.DiagnosticLogPath(a.configPath),
		"startup app_version=%q version=%q revision=%q release_time=%q exe=%q exe_sha256=%q config_path=%q history_path=%q log_path=%q",
		appVersion(),
		version,
		revision,
		appReleaseTime(),
		exePath,
		exeSHA256,
		a.configPath,
		appcore.HistoryPath(a.configPath),
		appcore.DiagnosticLogPath(a.configPath),
	)
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(a.configPath), "startup config=%s", configSummaryForLog(a.state.Config))
}

func configSummaryForLog(cfg appcore.Config) string {
	cfg.Normalize()
	summary := configLogSummary{
		Image: imageConfigLogSummary{
			MaxWidth:        cfg.Image.MaxWidth,
			MaxHeight:       cfg.Image.MaxHeight,
			MaxInputMB:      cfg.Image.MaxInputMB,
			Suffix:          cfg.Image.Suffix,
			OutputFormat:    cfg.Image.OutputFormat,
			Overwrite:       cfg.Image.Overwrite,
			JPEGQuality:     cfg.Image.JPEGQuality,
			OutputDirectory: cfg.Image.OutputDirectory,
		},
		Output: cfg.Output,
		Discord: webhookConfigLogSummary{
			WebhookConfigured: strings.TrimSpace(cfg.Discord.WebhookURL) != "",
		},
		AutoPhoto: autoPhotoConfigLogSummary{
			Enabled:             cfg.AutoPhoto.Enabled,
			PhotoDirectory:      cfg.AutoPhoto.PhotoDirectory,
			WebhookConfigured:   strings.TrimSpace(cfg.AutoPhoto.WebhookURL) != "",
			ScanIntervalSeconds: cfg.AutoPhoto.ScanIntervalSeconds,
		},
		AutoCapture: autoCaptureSummaryForLog(cfg.AutoCapture),
		ScreenshotAutoPost: screenshotAutoPostConfigLogSummary{
			Enabled:             cfg.ScreenshotAutoPost.Enabled,
			ScreenshotDirectory: cfg.ScreenshotAutoPost.ScreenshotDirectory,
			WebhookConfigured:   strings.TrimSpace(cfg.ScreenshotAutoPost.WebhookURL) != "",
			ScanIntervalSeconds: cfg.ScreenshotAutoPost.ScanIntervalSeconds,
		},
		Update: cfg.Update,
	}
	data, err := json.Marshal(summary)
	if err != nil {
		return fmt.Sprintf("%+v", summary)
	}
	return string(data)
}

func autoCaptureSummaryForLog(cfg appcore.AutoCaptureConfig) autoCaptureConfigLogSummary {
	enabledViews := 0
	calibratedViews := 0
	for _, view := range cfg.Views {
		if view.Enabled {
			enabledViews++
		}
		if view.Calibrated {
			calibratedViews++
		}
	}
	return autoCaptureConfigLogSummary{
		Schedule: cfg.Schedule,
		Capture:  cfg.Capture,
		OSC:      cfg.OSC,
		Output: autoCaptureOutputLogSummary{
			Directory:           cfg.Output.Directory,
			ImageFormat:         cfg.Output.ImageFormat,
			WriteSidecarJSON:    cfg.Output.WriteSidecarJSON,
			WriteEXIF:           cfg.Output.WriteEXIF,
			WriteUserListToEXIF: cfg.Output.WriteUserListToEXIF,
		},
		Presence: cfg.Presence,
		Discord: autoCaptureDiscordLogSummary{
			Enabled:           cfg.Discord.Enabled,
			WebhookConfigured: strings.TrimSpace(cfg.Discord.WebhookURL) != "",
			PostMode:          cfg.Discord.PostMode,
			IncludeImages:     cfg.Discord.IncludeImages,
		},
		ViewCount:       len(cfg.Views),
		EnabledViews:    enabledViews,
		CalibratedViews: calibratedViews,
	}
}

func fileSHA256(path string) (string, error) {
	file, err := os.Open(path) // #nosec G304 -- path is the currently running executable or selected diagnostic file.
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
