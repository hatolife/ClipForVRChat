package appcore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HistoryEntry struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Thumbnail        string `json:"thumbnail"`
	SourcePath       string `json:"sourcePath"`
	OutputPath       string `json:"outputPath"`
	DiscordMessageID string `json:"discordMessageId"`
	DiscordWebhookID string `json:"discordWebhookId"`
	DiscordToken     string `json:"discordToken"`
	Cleared          bool   `json:"cleared"`
	DiscordDeleted   bool   `json:"discordDeleted"`
	CreatedAt        string `json:"createdAt"`
	ClearedAt        string `json:"clearedAt,omitempty"`
	DeletedAt        string `json:"deletedAt,omitempty"`
}

func HistoryPath(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return "history.json"
	}
	return filepath.Join(filepath.Dir(configPath), "history.json")
}

func LoadHistory(path string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var history []HistoryEntry
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func SaveHistory(path string, history []HistoryEntry) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func AddResultsToHistory(path string, results []Result) ([]HistoryEntry, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, err
	}
	now := time.Now().Format(time.RFC3339)
	for i := range results {
		if strings.TrimSpace(results[i].URL) == "" || results[i].Error != "" {
			continue
		}
		entry := HistoryEntry{
			ID:               historyID(now, len(history)),
			Name:             results[i].Name,
			URL:              results[i].URL,
			Thumbnail:        results[i].Thumbnail,
			SourcePath:       results[i].SourcePath,
			OutputPath:       results[i].OutputPath,
			DiscordMessageID: results[i].DiscordMessageID,
			DiscordWebhookID: results[i].DiscordWebhookID,
			DiscordToken:     results[i].DiscordToken,
			CreatedAt:        now,
		}
		history = append([]HistoryEntry{entry}, history...)
		results[i].HistoryID = entry.ID
	}
	return history, SaveHistory(path, history)
}

func MarkHistoryCleared(path string, ids []string) ([]HistoryEntry, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, err
	}
	if len(ids) == 0 {
		return history, nil
	}
	now := time.Now().Format(time.RFC3339)
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	for i := range history {
		if idSet[history[i].ID] {
			history[i].Cleared = true
			if history[i].ClearedAt == "" {
				history[i].ClearedAt = now
			}
		}
	}
	return history, SaveHistory(path, history)
}

func PurgeUnavailableHistory(path string) ([]HistoryEntry, int, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, 0, err
	}
	kept := make([]HistoryEntry, 0, len(history))
	removed := 0
	for _, entry := range history {
		if entry.URL != "" && !RemoteURLAvailable(entry.URL) {
			removed++
			continue
		}
		kept = append(kept, entry)
	}
	return kept, removed, SaveHistory(path, kept)
}

func RemoteURLAvailable(url string) bool {
	status, err := remoteURLStatus(http.MethodHead, url)
	if err != nil {
		return false
	}
	if status == http.StatusMethodNotAllowed || status == http.StatusForbidden {
		status, err = remoteURLStatus(http.MethodGet, url)
		if err != nil {
			return false
		}
	}
	return status >= 200 && status < 400
}

func remoteURLStatus(method, url string) (int, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, err
	}
	if method == http.MethodGet {
		req.Header.Set("Range", "bytes=0-0")
	}
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func historyID(createdAt string, index int) string {
	clean := strings.NewReplacer(":", "", "-", "", ".", "").Replace(createdAt)
	return fmt.Sprintf("%s-%s-%d", clean, time.Now().Format("150405000000"), index)
}
