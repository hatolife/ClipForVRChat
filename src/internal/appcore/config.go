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
	Image     ImageConfig     `json:"image"`
	Output    OutputConfig    `json:"output"`
	Discord   DiscordConfig   `json:"discord"`
	AutoPhoto AutoPhotoConfig `json:"autoPhoto"`
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
}

type DiscordConfig struct {
	WebhookURL string `json:"webhookUrl"`
}

type AutoPhotoConfig struct {
	Enabled        bool   `json:"enabled"`
	PhotoDirectory string `json:"photoDirectory"`
	WebhookURL     string `json:"webhookUrl"`
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
			UploadDiscord:              true,
			ShowUI:                     "auto",
			CopySingleURLToClipboard:   true,
			DeleteOutputOnHistoryPurge: true,
		},
		AutoPhoto: AutoPhotoConfig{
			Enabled:        false,
			PhotoDirectory: DefaultVRChatPhotoDirectory(),
		},
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
	data, err := os.ReadFile(path) // #nosec G304 -- config path is the app config path or an explicitly opened local config file.
	if errors.Is(err, os.ErrNotExist) {
		return cfg, SaveConfig(path, cfg)
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
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
