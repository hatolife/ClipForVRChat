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
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	URL              string   `json:"url"`
	Thumbnail        string   `json:"thumbnail"`
	SourcePath       string   `json:"sourcePath"`
	OutputPath       string   `json:"outputPath"`
	QRURLs           []string `json:"qrUrls,omitempty"`
	DiscordMessageID string   `json:"discordMessageId"`
	DiscordWebhookID string   `json:"discordWebhookId"`
	DiscordToken     string   `json:"discordToken"`
	Cleared          bool     `json:"cleared"`
	Pinned           bool     `json:"pinned,omitempty"`
	DiscordDeleted   bool     `json:"discordDeleted"`
	LocalExists      bool     `json:"localExists,omitempty"`
	LocalDeleted     bool     `json:"localDeleted,omitempty"`
	CreatedAt        string   `json:"createdAt"`
	ClearedAt        string   `json:"clearedAt,omitempty"`
	DeletedAt        string   `json:"deletedAt,omitempty"`
	LocalDeletedAt   string   `json:"localDeletedAt,omitempty"`
}

func HistoryPath(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return "history.json"
	}
	return filepath.Join(filepath.Dir(configPath), "history.json")
}

func LoadHistory(path string) ([]HistoryEntry, error) {
	return LoadHistoryWithBaseDir(path, filepath.Dir(path))
}

func LoadHistoryWithBaseDir(path string, baseDir string) ([]HistoryEntry, error) {
	return LoadHistoryWithManagedOutputDir(path, baseDir, filepath.Join(baseDir, "output"))
}

func LoadHistoryWithManagedOutputDir(path string, baseDir string, outputDir string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- history path is derived from the active app config path.
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
	enrichHistoryStatus(history, baseDir, outputDir)
	return history, nil
}

func SaveHistory(path string, history []HistoryEntry) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), privateDirMode); err != nil {
		return err
	}
	return WritePrivateFile(path, append(data, '\n'))
}

func AddResultsToHistory(path string, results []Result) ([]HistoryEntry, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, err
	}
	now := time.Now().Format(time.RFC3339)
	for i := range results {
		if results[i].Error != "" || !ResultHasUserVisibleWork(results[i]) {
			continue
		}
		trustedDiscordURL := strings.TrimSpace(results[i].URL) != "" && IsTrustedDiscordImageURL(results[i].URL)
		entry := HistoryEntry{
			ID:         historyID(now, len(history)),
			Name:       results[i].Name,
			Thumbnail:  results[i].Thumbnail,
			SourcePath: results[i].SourcePath,
			OutputPath: results[i].OutputPath,
			QRURLs:     results[i].QRURLs,
			CreatedAt:  now,
		}
		if trustedDiscordURL {
			entry.URL = results[i].URL
			entry.DiscordMessageID = results[i].DiscordMessageID
			entry.DiscordWebhookID = results[i].DiscordWebhookID
			entry.DiscordToken = results[i].DiscordToken
		}
		history = append([]HistoryEntry{entry}, history...)
		results[i].HistoryID = entry.ID
	}
	enrichHistoryStatus(history, filepath.Dir(path), filepath.Join(filepath.Dir(path), "output"))
	return history, SaveHistory(path, history)
}

func ResultHasUserVisibleWork(result Result) bool {
	if strings.TrimSpace(result.Error) != "" {
		return true
	}
	if strings.TrimSpace(result.OutputPath) != "" {
		return true
	}
	if len(result.QRURLs) > 0 {
		return true
	}
	return strings.TrimSpace(result.URL) != "" && IsTrustedDiscordImageURL(result.URL)
}

func MarkHistoryCleared(path string, ids []string) ([]HistoryEntry, error) {
	history, err := LoadHistoryWithBaseDir(path, filepath.Dir(path))
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

func SetHistoryPinned(path string, id string, pinned bool) ([]HistoryEntry, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, err
	}
	for i := range history {
		if history[i].ID == id {
			history[i].Pinned = pinned
			break
		}
	}
	return history, SaveHistory(path, history)
}

func DeleteHistoryEntries(path string, ids []string) ([]HistoryEntry, int, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, 0, err
	}
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	kept := make([]HistoryEntry, 0, len(history))
	removed := 0
	for _, entry := range history {
		if idSet[entry.ID] && !entry.Pinned {
			removed++
			continue
		}
		kept = append(kept, entry)
	}
	return kept, removed, SaveHistory(path, kept)
}

func DeleteLocalHistoryFiles(path string, ids []string) ([]HistoryEntry, int, error) {
	return DeleteLocalHistoryFilesWithManagedOutputDir(path, ids, filepath.Join(filepath.Dir(path), "output"))
}

func DeleteLocalHistoryFilesWithManagedOutputDir(path string, ids []string, outputDir string) ([]HistoryEntry, int, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, 0, err
	}
	now := time.Now().Format(time.RFC3339)
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	deleted := 0
	for i := range history {
		if !idSet[history[i].ID] || history[i].Pinned || strings.TrimSpace(history[i].OutputPath) == "" || history[i].LocalDeleted {
			continue
		}
		if err := removeHistoryOutputFile(history[i].OutputPath, filepath.Dir(path), outputDir); err != nil {
			return history, deleted, err
		}
		_ = removeHistorySidecarFile(history[i].OutputPath, filepath.Dir(path), outputDir)
		history[i].LocalDeleted = true
		history[i].LocalDeletedAt = now
		deleted++
	}
	enrichHistoryStatus(history, filepath.Dir(path), outputDir)
	return history, deleted, SaveHistory(path, history)
}

func PurgeUnavailableHistory(path string) ([]HistoryEntry, int, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, 0, err
	}
	kept := make([]HistoryEntry, 0, len(history))
	removed := 0
	for _, entry := range history {
		if entry.Pinned {
			kept = append(kept, entry)
			continue
		}
		if entry.URL != "" && !RemoteURLAvailable(entry.URL) {
			removed++
			continue
		}
		kept = append(kept, entry)
	}
	return kept, removed, SaveHistory(path, kept)
}

func PurgeDiscordDeletedHistory(path string, deleteOutput bool) ([]HistoryEntry, int, error) {
	history, err := LoadHistory(path)
	if err != nil {
		return history, 0, err
	}
	kept := make([]HistoryEntry, 0, len(history))
	removed := 0
	for _, entry := range history {
		if entry.Pinned || !entry.DiscordDeleted {
			kept = append(kept, entry)
			continue
		}
		if deleteOutput {
			_ = removeHistoryOutputFile(entry.OutputPath, filepath.Dir(path), filepath.Join(filepath.Dir(path), "output"))
			_ = removeHistorySidecarFile(entry.OutputPath, filepath.Dir(path), filepath.Join(filepath.Dir(path), "output"))
		}
		removed++
	}
	return kept, removed, SaveHistory(path, kept)
}

func removeHistoryOutputFile(path string, baseDir string, outputDir string) error {
	path = ResolveManagedHistoryOutputPath(path, baseDir, outputDir)
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func removeHistorySidecarFile(path string, baseDir string, outputDir string) error {
	path = ResolveManagedHistoryOutputPath(path, baseDir, outputDir)
	if path == "" {
		return nil
	}
	sidecarPath := path + ".json"
	if err := os.Remove(sidecarPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func enrichHistoryStatus(history []HistoryEntry, baseDir string, outputDir string) {
	for i := range history {
		history[i].LocalExists = localHistoryFileExists(history[i], baseDir, outputDir)
	}
}

func localHistoryFileExists(entry HistoryEntry, baseDir string, outputDir string) bool {
	path := ResolveManagedHistoryOutputPath(entry.OutputPath, baseDir, outputDir)
	if path == "" || entry.LocalDeleted {
		return false
	}
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func ResolveHistoryOutputPath(path string, baseDir string) string {
	cleaned := strings.Trim(strings.TrimSpace(path), `"`)
	if cleaned == "" {
		return ""
	}
	if filepath.IsAbs(cleaned) {
		return cleaned
	}
	if strings.TrimSpace(baseDir) == "" {
		return ""
	}
	return filepath.Join(baseDir, cleaned)
}

func ResolveManagedHistoryOutputPath(path string, baseDir string, outputDir string) string {
	resolved := ResolveHistoryOutputPath(path, baseDir)
	if resolved == "" {
		return ""
	}
	managedDir := ResolveHistoryOutputPath(outputDir, baseDir)
	if managedDir == "" {
		return ""
	}
	absPath, err := filepath.Abs(resolved)
	if err != nil {
		return ""
	}
	absManagedDir, err := filepath.Abs(managedDir)
	if err != nil {
		return ""
	}
	rel, err := filepath.Rel(absManagedDir, absPath)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." || filepath.IsAbs(rel) {
		return ""
	}
	return absPath
}

func RemoteURLAvailable(url string) bool {
	if !IsTrustedDiscordImageURL(url) {
		return false
	}
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
