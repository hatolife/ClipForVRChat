package appcore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const LatestReleaseAPIURL = "https://api.github.com/repos/hatolife/ClipForVRChat/releases/latest"

type UpdateInfo struct {
	Available              bool   `json:"available"`
	CurrentVersion         string `json:"currentVersion"`
	CurrentReleaseTime     string `json:"currentReleaseTime"`
	LatestVersion          string `json:"latestVersion"`
	LatestReleasePublished string `json:"latestReleasePublished"`
	URL                    string `json:"url"`
}

type latestReleaseResponse struct {
	TagName     string `json:"tag_name"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

func CheckLatestRelease(ctx context.Context, client *http.Client, currentVersion string, currentReleaseTime string) (UpdateInfo, error) {
	info := UpdateInfo{
		CurrentVersion:     currentVersion,
		CurrentReleaseTime: currentReleaseTime,
	}
	if client == nil {
		client = &http.Client{Timeout: 6 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LatestReleaseAPIURL, nil)
	if err != nil {
		return info, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ClipForVRChat")

	resp, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return info, fmt.Errorf("GitHub Releases の確認に失敗しました: status=%d", resp.StatusCode)
	}

	var latest latestReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&latest); err != nil {
		return info, err
	}
	info.LatestVersion = strings.TrimSpace(latest.TagName)
	info.LatestReleasePublished = strings.TrimSpace(latest.PublishedAt)
	info.URL = strings.TrimSpace(latest.HTMLURL)
	if info.URL == "" && info.LatestVersion != "" {
		info.URL = "https://github.com/hatolife/ClipForVRChat/releases/tag/" + info.LatestVersion
	}
	info.Available = IsNewerRelease(currentVersion, currentReleaseTime, info.LatestVersion, info.LatestReleasePublished)
	return info, nil
}

func IsNewerRelease(currentVersion string, currentReleaseTime string, latestVersion string, latestReleasePublished string) bool {
	current, okCurrent := parseReleaseVersion(currentVersion)
	latest, okLatest := parseReleaseVersion(latestVersion)
	if !okCurrent || !okLatest {
		return false
	}
	for i := range latest {
		if latest[i] > current[i] {
			return true
		}
		if latest[i] < current[i] {
			return false
		}
	}
	return isNewerReleaseTime(currentReleaseTime, latestReleasePublished)
}

func isNewerReleaseTime(currentReleaseTime string, latestReleasePublished string) bool {
	if strings.TrimSpace(currentReleaseTime) == "" || strings.TrimSpace(latestReleasePublished) == "" {
		return false
	}
	current, err := time.Parse(time.RFC3339, currentReleaseTime)
	if err != nil {
		return false
	}
	latest, err := time.Parse(time.RFC3339, latestReleasePublished)
	if err != nil {
		return false
	}
	return latest.After(current.Add(5 * time.Minute))
}

func parseReleaseVersion(value string) ([3]int, bool) {
	var version [3]int
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "v")
	value = strings.TrimPrefix(value, "V")
	value = strings.SplitN(value, "+", 2)[0]
	value = strings.SplitN(value, "-", 2)[0]
	parts := strings.Split(value, ".")
	if len(parts) < 3 {
		return version, false
	}
	for i := 0; i < 3; i++ {
		n, err := strconv.Atoi(parts[i])
		if err != nil || n < 0 {
			return version, false
		}
		version[i] = n
	}
	return version, true
}
