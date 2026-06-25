package main

import (
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

func TestAppGetInitialStateRefreshesHistory(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	historyPath := appcore.HistoryPath(configPath)
	if err := appcore.SaveHistory(historyPath, []appcore.HistoryEntry{{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}); err != nil {
		t.Fatal(err)
	}

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	state := app.GetInitialState()
	if len(state.History) != 1 || state.History[0].ID != "1" {
		t.Fatalf("history = %+v", state.History)
	}
}

func TestAppClearResultsMarksHistoryCleared(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{{ID: "h1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{
		Mode:    appcore.ModeResults,
		Results: []appcore.Result{{HistoryID: "h1", URL: history[0].URL}},
	})

	state := app.ClearResults()
	if len(state.Results) != 0 {
		t.Fatalf("Results = %+v, want empty", state.Results)
	}
	loaded, err := appcore.LoadHistory(appcore.HistoryPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || !loaded[0].Cleared {
		t.Fatalf("history = %+v, want cleared", loaded)
	}
}

func TestAppSaveConfigAndOpenSettings(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	cfg := appcore.DefaultConfig()
	cfg.Image.MaxWidth = 777

	if err := app.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	state, err := app.OpenSettings("")
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != appcore.ModeSettings || state.Config.Image.MaxWidth != 777 {
		t.Fatalf("state = %+v", state)
	}
	closed := app.CloseSettings()
	if closed.Mode != appcore.ModeResults {
		t.Fatalf("closed mode = %s, want results", closed.Mode)
	}
}

func TestAppOpenSettingsKeepsConfigPathOnLoadError(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	badConfigPath := filepath.Join(dir, "bad.json")
	writeTextFile(t, badConfigPath, "{")

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults, ConfigPath: configPath})
	if _, err := app.OpenSettings(badConfigPath); err == nil {
		t.Fatal("expected invalid config error")
	}
	if app.configPath != configPath {
		t.Fatalf("configPath = %q, want %q", app.configPath, configPath)
	}
}

func TestAppSaveConfigAndProcessHandlesDecodeError(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	source := filepath.Join(dir, "bad.png")
	writeTextFile(t, source, "not image")

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeSettings, ProcessOnSave: true})
	cfg := appcore.DefaultConfig()
	cfg.Image.OutputDirectory = filepath.Join(dir, "out")
	cfg.Output.UploadDiscord = false

	state, err := app.SaveConfigAndProcess(cfg, []string{source})
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != appcore.ModeError || !strings.Contains(state.Message, "処理中にエラー") {
		t.Fatalf("state = %+v", state)
	}
}

func TestAppProcessToStateHidesNoWorkResults(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	source := filepath.Join(dir, "image.png")
	writeTestPNG(t, source, 2, 2)
	cfg := appcore.DefaultConfig()
	cfg.Output.UploadDiscord = false
	cfg.Output.SaveLocal = false
	cfg.Output.DetectQRCodeURLs = false
	if err := appcore.SaveConfig(configPath, cfg); err != nil {
		t.Fatal(err)
	}

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	state, err := app.ProcessToState([]string{source})
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Results) != 0 {
		t.Fatalf("results = %+v, want empty", state.Results)
	}
	if !strings.Contains(state.Message, "実行された処理はありません") {
		t.Fatalf("message = %q", state.Message)
	}
	history, err := appcore.LoadHistory(appcore.HistoryPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 0 {
		t.Fatalf("history = %+v, want empty", history)
	}
}

func TestAppStartupStartsAutoPhotoWatcherWithoutDiscordUpload(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Output.UploadDiscord = false
	cfg.AutoPhoto.Enabled = true
	cfg.AutoPhoto.PhotoDirectory = t.TempDir()
	cfg.AutoPhoto.ScanIntervalSeconds = 3600

	app := NewApp(configPath, appcore.UIState{Config: cfg})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.startup(ctx)
	defer func() {
		if app.autoCancel != nil {
			app.autoCancel()
		}
	}()

	if app.autoCancel == nil {
		t.Fatal("auto photo watcher was not started")
	}
}

func TestExplorerSelectArgsRejectsMissingPath(t *testing.T) {
	if _, _, err := explorerSelectArgs(""); err == nil {
		t.Fatal("expected empty path error")
	}
	if _, _, err := explorerSelectArgs(filepath.Join(t.TempDir(), "missing.png")); err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestExplorerSelectArgsRejectsDirectory(t *testing.T) {
	if _, _, err := explorerSelectArgs(t.TempDir()); err == nil {
		t.Fatal("expected directory error")
	}
}

func TestAppProcessToStateRejectsMixedJSON(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})
	state, err := app.ProcessToState([]string{"config.json", "image.png"})
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != appcore.ModeError || !strings.Contains(state.Message, "混在") {
		t.Fatalf("state = %+v", state)
	}
}

func TestAppHistoryOperations(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{
		{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png", DiscordDeleted: true},
		{ID: "2", URL: ""},
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	got, err := app.GetHistory()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len(history) = %d, want 2", len(got))
	}

	got, err = app.MarkHistoryCleared([]string{"2"})
	if err != nil {
		t.Fatal(err)
	}
	if !got[1].Cleared {
		t.Fatalf("entry 2 should be cleared: %+v", got)
	}

	got, err = app.PurgeDeletedHistoryEntries()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "2" {
		t.Fatalf("purged history = %+v", got)
	}
}

func TestAppDeleteDiscordHistoryEntriesSkipsEntriesWithoutStoredWebhookData(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	got, err := app.DeleteDiscordHistoryEntries([]string{"1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].DiscordDeleted {
		t.Fatalf("history = %+v, want unchanged", got)
	}
}

func TestAppDeleteDiscordHistoryEntriesPersistsPartialSuccess(t *testing.T) {
	oldDelete := deleteDiscordMessage
	t.Cleanup(func() {
		deleteDiscordMessage = oldDelete
	})

	calls := 0
	deleteDiscordMessage = func(webhookID, token, messageID string) error {
		calls++
		if messageID == "m2" {
			return errors.New("delete failed")
		}
		return nil
	}

	configPath := filepath.Join(t.TempDir(), "config.json")
	history := []appcore.HistoryEntry{
		{ID: "1", URL: "https://cdn.discordapp.com/attachments/1/2/a.png", DiscordWebhookID: "w", DiscordToken: "t", DiscordMessageID: "m1"},
		{ID: "2", URL: "https://cdn.discordapp.com/attachments/1/2/b.png", DiscordWebhookID: "w", DiscordToken: "t", DiscordMessageID: "m2"},
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), history); err != nil {
		t.Fatal(err)
	}
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	if _, err := app.DeleteDiscordHistoryEntries([]string{"1", "2"}); err == nil {
		t.Fatal("expected partial delete error")
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	loaded, err := appcore.LoadHistory(appcore.HistoryPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	if !loaded[0].DiscordDeleted || loaded[0].DeletedAt == "" {
		t.Fatalf("first entry should be persisted as deleted: %+v", loaded[0])
	}
	if loaded[1].DiscordDeleted {
		t.Fatalf("second entry should not be deleted: %+v", loaded[1])
	}
}

func TestResultMessage(t *testing.T) {
	cfg := appcore.DefaultConfig()
	msg := resultMessage(cfg, []appcore.Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}, nil)
	if !strings.Contains(msg, "コピー") {
		t.Fatalf("message = %q", msg)
	}

	msg = resultMessage(cfg, []appcore.Result{{URL: "https://cdn.discordapp.com/attachments/1/2/a.png"}}, errors.New("clipboard busy"))
	if !strings.Contains(msg, "コピーできません") {
		t.Fatalf("copy error message = %q", msg)
	}

	msg = resultMessage(cfg, []appcore.Result{{OutputPath: "out.png"}}, nil)
	if !strings.Contains(msg, "保存") {
		t.Fatalf("message = %q", msg)
	}

	msg = resultMessage(cfg, []appcore.Result{{QRURLs: []string{"https://example.com/qr"}}}, nil)
	if !strings.Contains(msg, "QRコードURL") {
		t.Fatalf("message = %q", msg)
	}

	msg = resultMessage(cfg, nil, nil)
	if !strings.Contains(msg, "実行された処理はありません") {
		t.Fatalf("message = %q", msg)
	}
}

func writeTextFile(t *testing.T, path string, text string) {
	t.Helper()
	if err := appcore.WritePrivateFile(path, []byte(text)); err != nil {
		t.Fatal(err)
	}
}

func writeTestPNG(t *testing.T, path string, width int, height int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 80, G: 120, B: 180, A: 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
}
