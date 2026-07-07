package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/hatolife/ClipForVRChat/internal/appcore"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

type cliArgs struct {
	Version bool `arg:"--version" help:"バージョンを表示して終了します"`
}

func (cliArgs) Description() string {
	return "ClipForVRChat"
}

func main() {
	stdout, stderr, cleanup := cliOutputWriters(os.Args[1:], os.Stdout, os.Stderr)
	defer cleanup()
	if handled, exitCode := handleCLIArgs(os.Args[1:], stdout, stderr); handled {
		os.Exit(exitCode)
	}
	appcore.SetOSCVersionNotice(appVersion())

	args := os.Args[1:]
	imageArgs, encryptionArgs := splitStartupPaths(args)
	for _, path := range encryptionArgs {
		outputPath, err := encryptPathWithPublicKey(path)
		if err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Fprintf(stdout, "Encrypted: %s\n", outputPath)
	}
	if len(args) > 0 && len(imageArgs) == 0 {
		return
	}
	args = imageArgs

	instance, err := initializeSingleInstance(stderr, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}
	if instance == nil {
		return
	}
	defer instance.Close()

	configPath := defaultConfigPath()
	state := appcore.UIState{
		Mode:       appcore.ModeResults,
		ConfigPath: configPath,
	}

	configExists := appcore.ConfigExists(configPath)
	cfg, err := appcore.LoadConfig(configPath)
	if err != nil {
		state.Mode = appcore.ModeError
		state.Message = fmt.Sprintf("設定を読み込めませんでした: %v", err)
		state.Config = appcore.DefaultConfig()
		runUI(configPath, state, instance)
		return
	}
	state.Config = cfg
	if history, err := appcore.LoadHistoryWithManagedOutputDir(appcore.HistoryPath(configPath), filepath.Dir(configPath), managedOutputDir(configPath, cfg)); err == nil {
		state.History = history
	}
	if draft, ok, reason, diffPaths, err := loadSettingsDraftForConfig(configPath, cfg, configExists); err != nil {
		appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "settings draft load error: %v", err)
	} else if ok {
		baseline := cfg
		state.Mode = appcore.ModeSettings
		state.Message = "保存されていない設定の一時変更を復元しました。保存するまで設定ファイルには反映されません。"
		state.Config = draft
		state.SettingsBaselineConfig = &baseline
		state.UnsavedSettingsDraft = true
		appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "settings draft restored: diff_paths=%s", strings.Join(diffPaths, ","))
	} else if reason != "" {
		appcore.AppendDiagnosticLog(appcore.DiagnosticLogPath(configPath), "settings draft ignored: %s", reason)
	}

	if !configExists {
		state.Mode = appcore.ModeSettings
		state.Message = "初回起動です。設定を確認して保存すると、続けて通常処理を実行します。"
		state.PendingPaths = args
		state.ProcessOnSave = len(args) > 0
		runUI(configPath, state, instance)
		return
	}

	if len(args) == 0 {
		runUI(configPath, state, instance)
		return
	}

	state.Config = cfg
	state.SettingsBaselineConfig = nil
	state.UnsavedSettingsDraft = false
	results, err := appcore.Processor{Config: cfg}.ProcessPaths(args)
	if err != nil {
		state.Mode = appcore.ModeError
		state.Message = err.Error()
		runUI(configPath, state, instance)
		return
	}
	copyErr := copySingleURLIfNeeded(cfg, results)
	if history, err := appcore.AddResultsToHistory(appcore.HistoryPath(configPath), results); err == nil {
		state.History = history
	}

	state.Results = results
	if shouldExitWithoutUI(cfg, results, copyErr) {
		return
	}
	if hasErrors(results) {
		state.Mode = appcore.ModeError
		state.Message = "処理中にエラーが発生しました。内容を確認してください。"
	} else {
		state.Mode = appcore.ModeResults
		state.Message = resultMessage(cfg, results, copyErr)
	}
	runUI(configPath, state, instance)
}

func splitStartupPaths(paths []string) (imagePaths []string, encryptionPaths []string) {
	for _, path := range paths {
		if isStartupImagePath(path) {
			imagePaths = append(imagePaths, path)
		} else {
			encryptionPaths = append(encryptionPaths, path)
		}
	}
	return imagePaths, encryptionPaths
}

func isStartupImagePath(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return false
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".webp":
		return true
	default:
		return false
	}
}

func handleCLIArgs(args []string, stdout io.Writer, stderr io.Writer) (bool, int) {
	var parsed cliArgs
	parser, err := arg.NewParser(arg.Config{Program: "ClipForVRChat"}, &parsed)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return true, 2
	}
	if err := parser.Parse(args); err != nil {
		if err == arg.ErrHelp {
			parser.WriteHelp(stdout)
			return true, 0
		}
		return false, 0
	}
	if parsed.Version {
		fmt.Fprintf(stdout, "ClipForVRChat %s\n", appVersion())
		return true, 0
	}
	return false, 0
}

func shouldExitWithoutUI(cfg appcore.Config, results []appcore.Result, copyErr error) bool {
	if cfg.Output.ShowUI == "always" {
		return false
	}
	if copyErr != nil {
		return false
	}
	if len(results) != 1 || hasErrors(results) {
		return false
	}
	if cfg.Output.ShowUI == "never" {
		return true
	}
	return results[0].URL != "" && cfg.Output.CopySingleURLToClipboard
}

func hasErrors(results []appcore.Result) bool {
	for _, result := range results {
		if result.Error != "" {
			return true
		}
	}
	return false
}

func runUI(configPath string, state appcore.UIState, instance *singleInstance) {
	app := NewApp(configPath, state)
	if instance != nil {
		instance.BindApp(app)
	}
	logPath := appcore.DiagnosticLogPath(configPath)
	appcore.AppendDiagnosticLog(
		logPath,
		"ui wails run begin: app_version=%q revision=%q channel=%q goos=%q goarch=%q args=%q cwd=%q config=%q initial_mode=%q frontend_assets=%q icon_bytes=%d",
		appVersion(),
		revision,
		buildChannel,
		runtime.GOOS,
		runtime.GOARCH,
		os.Args,
		mustGetwd(),
		configPath,
		state.Mode,
		frontendAssetSummary(),
		len(icon),
	)
	err := wails.Run(&options.App{
		Title:  fmt.Sprintf("ClipForVRChat %s", appVersion()),
		Width:  900,
		Height: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.startup,
		OnDomReady:    app.domReady,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:     true,
			DisableWebViewDrop: true,
		},
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		appcore.AppendDiagnosticLog(logPath, "ui wails run error: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	appcore.AppendDiagnosticLog(logPath, "ui wails run returned: err=nil")
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "error:" + err.Error()
	}
	return wd
}

func frontendAssetSummary() string {
	entries, err := assets.ReadDir("frontend/dist")
	if err != nil {
		return fmt.Sprintf("read_dir_error=%v", err)
	}
	indexExists := true
	if _, err := assets.ReadFile("frontend/dist/index.html"); err != nil {
		indexExists = false
	}
	return fmt.Sprintf("root_entries=%d index_html=%t", len(entries), indexExists)
}
