package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
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
	if _, err := explorerRevealPath("", ""); err == nil {
		t.Fatal("expected empty path error")
	}
	if _, err := explorerRevealPath(filepath.Join(t.TempDir(), "missing.png"), ""); err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestExplorerSelectArgsRejectsDirectory(t *testing.T) {
	if _, err := explorerRevealPath(t.TempDir(), ""); err == nil {
		t.Fatal("expected directory error")
	}
}

func TestExplorerSelectArgsResolvesRelativePathFromBaseDir(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "out.png"), []byte("image"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := explorerRevealPath("output/out.png", dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(outputDir, "out.png") {
		t.Fatalf("path = %q, want resolved relative output path", got)
	}
}

func TestAppLogUserActionWritesDiagnosticLog(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults})

	app.LogUserAction("button_click", "history_delete_entries selected=1")

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, `ui action="button_click"`) || !strings.Contains(text, `history_delete_entries`) {
		t.Fatalf("diagnostic log = %q", text)
	}
	if filepath.Dir(appcore.DiagnosticLogPath(configPath)) != filepath.Join(dir, "logs") {
		t.Fatalf("diagnostic log path = %q, want dated log under log dir", appcore.DiagnosticLogPath(configPath))
	}
}

func TestTrustedExternalURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		ok   bool
	}{
		{name: "github", raw: "https://github.com/hatolife/ClipForVRChat/releases", ok: true},
		{name: "booth", raw: "https://hatolife.booth.pm/items/8531663", ok: true},
		{name: "discord help", raw: "https://support.discord.com/hc/ja/articles/228383668", ok: true},
		{name: "twitter", raw: "https://x.com/hato_poppo_life", ok: true},
		{name: "http", raw: "http://github.com/hatolife/ClipForVRChat", ok: false},
		{name: "file", raw: "file:///C:/Windows/win.ini", ok: false},
		{name: "unknown host", raw: "https://example.com", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := trustedExternalURL(tt.raw)
			if tt.ok && err != nil {
				t.Fatalf("trustedExternalURL() error = %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("trustedExternalURL() should reject URL")
			}
		})
	}
}

func TestIsGoTestBinary(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "/tmp/app.test", want: true},
		{path: `C:\Users\runner\AppData\Local\Temp\app.test.exe`, want: true},
		{path: `C:\app\ClipForVRChat.exe`, want: false},
		{path: "/tmp/not-test", want: false},
	}
	for _, tt := range tests {
		if got := isGoTestBinary(tt.path); got != tt.want {
			t.Fatalf("isGoTestBinary(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestAppStartupWritesVersionHashAndRedactedConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Discord.WebhookURL = "https://discord.com/api/webhooks/secret"
	cfg.AutoPhoto.WebhookURL = "https://discord.com/api/webhooks/auto-secret"
	app := NewApp(configPath, appcore.UIState{Config: cfg})

	app.startup(context.Background())

	data, err := os.ReadFile(appcore.DiagnosticLogPath(configPath))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "startup app_version=") || !strings.Contains(text, "exe_sha256=") {
		t.Fatalf("diagnostic log = %q, want startup version and hash", text)
	}
	if !strings.Contains(text, `"webhookConfigured":true`) {
		t.Fatalf("diagnostic log = %q, want redacted webhook configured flag", text)
	}
	if strings.Contains(text, "secret") {
		t.Fatalf("diagnostic log leaked webhook URL: %q", text)
	}
}

func TestCreateEncryptedDiagnosticPackageEncryptsZip(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	userProfile := filepath.Join(dir, "Users", "alice")
	t.Setenv("USERPROFILE", userProfile)
	cfg := appcore.DefaultConfig()
	cfg.Discord.WebhookURL = "https://discord.com/api/webhooks/secret"
	cfg.AutoPhoto.PhotoDirectory = filepath.Join(userProfile, "Pictures", "VRChat")
	cfg.Image.OutputDirectory = filepath.Join(userProfile, "Pictures", "ClipForVRChatOutput")
	if err := appcore.SaveConfig(configPath, cfg); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(cfg.Image.OutputDirectory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.Image.OutputDirectory, "result.png"), []byte("real output image"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), []appcore.HistoryEntry{{ID: "1", URL: "https://example.com", SourcePath: filepath.Join(userProfile, "Pictures", "VRChat", "photo.png")}}); err != nil {
		t.Fatal(err)
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "test log path=%q", filepath.Join(userProfile, "Pictures", "VRChat", "photo.png"))

	app := NewApp(configPath, appcore.UIState{Config: cfg})
	path, err := app.CreateEncryptedDiagnosticPackage()
	if err != nil {
		t.Fatal(err)
	}
	workDir := filepath.Dir(path)
	if filepath.Base(filepath.Dir(workDir)) != "diagnostics" || filepath.Ext(path) != ".gpg" {
		t.Fatalf("path = %q, want encrypted package under diagnostics timestamp dir", path)
	}
	zipPath := strings.TrimSuffix(path, ".gpg")
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("plain zip should remain for user review: %v", err)
	}
	dataDir := filepath.Join(workDir, "data")
	if _, err := os.Stat(dataDir); !os.IsNotExist(err) {
		t.Fatalf("data dir should be removed after zip creation, err=%v", err)
	}
	plainZip, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer plainZip.Close()
	logData, err := readZipFile(t, plainZip.File, "logs/"+filepath.Base(appcore.DiagnosticLogPath(configPath)))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(logData), userProfile) || !strings.Contains(string(logData), `%USERPROFILE%`) {
		t.Fatalf("redacted log = %q, want USERPROFILE placeholder and no raw profile path", string(logData))
	}
	if !strings.Contains(string(logData), "diagnostic output directory=") || !strings.Contains(string(logData), "result.png") {
		t.Fatalf("redacted log = %q, want output directory listing", string(logData))
	}
	configData, err := readZipFile(t, plainZip.File, "config.json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(configData), userProfile) || !strings.Contains(string(configData), `%USERPROFILE%`) {
		t.Fatalf("redacted config = %q, want USERPROFILE placeholder and no raw profile path", string(configData))
	}
	historyData, err := readZipFile(t, plainZip.File, "history.json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(historyData), userProfile) || !strings.Contains(string(historyData), `%USERPROFILE%`) {
		t.Fatalf("redacted history = %q, want USERPROFILE placeholder and no raw profile path", string(historyData))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, leaked := range []string{"config.json", "history.json", "manifest.json", "secret"} {
		if strings.Contains(string(data), leaked) {
			t.Fatalf("encrypted package leaked %q", leaked)
		}
	}
}

func TestEncryptZipFileWithPublicKey(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "sample.zip")
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	writer, err := zipWriter.Create("inside.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("secret text")); err != nil {
		t.Fatal(err)
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zipPath, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}

	outputPath, err := encryptZipFileWithPublicKey(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	if outputPath != zipPath+".gpg" {
		t.Fatalf("outputPath = %q, want %q", outputPath, zipPath+".gpg")
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, leaked := range []string{"inside.txt", "secret text"} {
		if strings.Contains(string(data), leaked) {
			t.Fatalf("encrypted zip leaked %q", leaked)
		}
	}
}

func TestDiagnosticZipDoesNotIncludeOutputImages(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Image.OutputDirectory = "./output"
	if err := appcore.SaveConfig(configPath, cfg); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "real.png"), []byte("not the real image in diagnostics"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "real.jpg"), []byte("not the real jpeg in diagnostics"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := appcore.SaveHistory(appcore.HistoryPath(configPath), []appcore.HistoryEntry{{ID: "1", OutputPath: "output/history.png"}}); err != nil {
		t.Fatal(err)
	}
	appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "test log")
	zipData, _, err := buildDiagnosticZip(configPath, cfg)
	if err != nil {
		t.Fatal(err)
	}
	plainZip, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		t.Fatalf("plain diagnostic zip is invalid: %v", err)
	}
	entries := zipEntryNames(plainZip.File)
	for _, want := range []string{"logs/" + filepath.Base(appcore.DiagnosticLogPath(configPath))} {
		if !entries[want] {
			t.Fatalf("zip entries = %+v, missing %s", entries, want)
		}
	}
	if entries["log/"+filepath.Base(appcore.DiagnosticLogPath(configPath))] {
		t.Fatalf("zip entries = %+v, should use logs/ not log/", entries)
	}
	for _, unwanted := range []string{"output/history.png", "output/real.jpg", "output/real.png"} {
		if entries[unwanted] {
			t.Fatalf("zip entries = %+v, should not include output images", entries)
		}
	}

	entity, err := openpgp.NewEntity("Diagnostic Test", "", "diagnostic@example.test", nil)
	if err != nil {
		t.Fatal(err)
	}
	encrypted, err := encryptDiagnosticZip(zipData, openpgp.EntityList{entity})
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(encrypted), "config.json") {
		t.Fatal("encrypted package leaked zip entry name")
	}
	message, err := openpgp.ReadMessage(bytes.NewReader(encrypted), openpgp.EntityList{entity}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if message.LiteralData == nil || !message.LiteralData.IsBinary {
		t.Fatalf("literal data = %+v, want binary", message.LiteralData)
	}
	decrypted, err := io.ReadAll(message.UnverifiedBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := zip.NewReader(bytes.NewReader(decrypted), int64(len(decrypted))); err != nil {
		t.Fatalf("decrypted diagnostic zip is invalid: %v", err)
	}
}

func zipEntryNames(files []*zip.File) map[string]bool {
	names := make(map[string]bool, len(files))
	for _, file := range files {
		names[file.Name] = true
	}
	return names
}

func readZipFile(t *testing.T, files []*zip.File, name string) ([]byte, error) {
	t.Helper()
	for _, file := range files {
		if file.Name != name {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	return nil, fmt.Errorf("zip entry %s not found", name)
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
	cfg.Output.CopySingleURLToClipboard = true
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
	if !strings.Contains(msg, "表示できる処理結果はありません") {
		t.Fatalf("message = %q", msg)
	}

	cfg.Output.UploadDiscord = false
	cfg.Output.SaveLocal = false
	cfg.Output.DetectQRCodeURLs = false
	msg = resultMessage(cfg, nil, nil)
	if !strings.Contains(msg, "すべてOFF") || !strings.Contains(msg, "設定を確認") {
		t.Fatalf("all disabled message = %q", msg)
	}

	cfg.Output.DetectQRCodeURLs = true
	msg = resultMessage(cfg, nil, nil)
	if !strings.Contains(msg, "URLを取得できません") || !strings.Contains(msg, "Discord投稿またはローカル保存") {
		t.Fatalf("qr-only no result message = %q", msg)
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
