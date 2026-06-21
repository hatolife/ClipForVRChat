package appcore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Image   ImageConfig   `json:"image"`
	Output  OutputConfig  `json:"output"`
	Discord DiscordConfig `json:"discord"`
}

type ImageConfig struct {
	MaxWidth        int    `json:"maxWidth"`
	MaxHeight       int    `json:"maxHeight"`
	Suffix          string `json:"suffix"`
	Overwrite       bool   `json:"overwrite"`
	JPEGQuality     int    `json:"jpegQuality"`
	OutputDirectory string `json:"outputDirectory"`
}

type OutputConfig struct {
	SaveLocal                bool   `json:"saveLocal"`
	UploadDiscord            bool   `json:"uploadDiscord"`
	ShowUI                   string `json:"showUi"`
	CopySingleURLToClipboard bool   `json:"copySingleUrlToClipboard"`
}

type DiscordConfig struct {
	WebhookURL string `json:"webhookUrl"`
}

func DefaultConfig() Config {
	return Config{
		Image: ImageConfig{
			MaxWidth:    2048,
			MaxHeight:   2048,
			Suffix:      "_2048",
			Overwrite:   false,
			JPEGQuality: 92,
		},
		Output: OutputConfig{
			SaveLocal:                true,
			UploadDiscord:            true,
			ShowUI:                   "auto",
			CopySingleURLToClipboard: true,
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
	data, err := os.ReadFile(path)
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
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func (c *Config) Normalize() {
	if c.Image.MaxWidth <= 0 {
		c.Image.MaxWidth = 2048
	}
	if c.Image.MaxHeight <= 0 {
		c.Image.MaxHeight = 2048
	}
	if c.Image.Suffix == "" {
		c.Image.Suffix = "_2048"
	}
	if c.Image.JPEGQuality <= 0 || c.Image.JPEGQuality > 100 {
		c.Image.JPEGQuality = 92
	}
	switch c.Output.ShowUI {
	case "auto", "always", "never":
	default:
		c.Output.ShowUI = "auto"
	}
}
