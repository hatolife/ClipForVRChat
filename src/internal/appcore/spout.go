package appcore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	spoutHelperFileName = "spout-capture.exe"
	spoutLibraryDLLName = "SpoutLibrary.dll"
)

var spoutHelperCacheDirForTest string

type SpoutSenderInfo struct {
	Name     string `json:"name"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	HostPath string `json:"hostPath,omitempty"`
}

type SpoutListResult struct {
	OK      bool              `json:"ok"`
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Senders []SpoutSenderInfo `json:"senders,omitempty"`
}

type SpoutCaptureResult struct {
	OK         bool              `json:"ok"`
	Code       string            `json:"code,omitempty"`
	Message    string            `json:"message,omitempty"`
	SenderName string            `json:"senderName,omitempty"`
	Width      int               `json:"width,omitempty"`
	Height     int               `json:"height,omitempty"`
	Frame      int64             `json:"frame,omitempty"`
	CapturedAt string            `json:"capturedAt,omitempty"`
	OutputPath string            `json:"outputPath,omitempty"`
	Senders    []SpoutSenderInfo `json:"senders,omitempty"`
}

type SpoutHelperStatus struct {
	Available bool              `json:"available"`
	Path      string            `json:"path,omitempty"`
	Version   string            `json:"version,omitempty"`
	Message   string            `json:"message"`
	Senders   []SpoutSenderInfo `json:"senders,omitempty"`
}

type SpoutHelperVersionResult struct {
	OK      bool   `json:"ok"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type capturedImageStats struct {
	Format           string
	Width            int
	Height           int
	Samples          int
	Mean             float64
	Stddev           float64
	NearWhiteRatio   float64
	NearBlackRatio   float64
	TransparentRatio float64
}

func ResolveSpoutHelperPath(helperPath string) (string, error) {
	helperPath = strings.Trim(strings.TrimSpace(helperPath), `"`)
	useEmbeddedFallback := helperPath == "" || helperPath == spoutHelperFileName
	if helperPath == "" {
		helperPath = spoutHelperFileName
	}
	if strings.ContainsAny(helperPath, `\/`) || filepath.IsAbs(helperPath) {
		info, err := os.Stat(helperPath)
		if err != nil {
			return "", fmt.Errorf("Spout helperが見つかりません: %w", err)
		}
		if info.IsDir() {
			return "", fmt.Errorf("Spout helperパスがフォルダを指しています。spout-capture.exeを指定してください")
		}
		return helperPath, nil
	}
	if candidate, ok := spoutHelperNextToExecutable(helperPath); ok {
		return candidate, nil
	}
	if useEmbeddedFallback && embeddedSpoutHelperAvailable() {
		embedded, err := resolveEmbeddedSpoutHelper()
		if err != nil {
			return "", err
		}
		return embedded, nil
	}
	resolved, err := exec.LookPath(helperPath)
	if err != nil {
		return "", fmt.Errorf("Spout helperが見つかりません。単一exe版では埋め込みhelperを利用します。分離版zipを使う場合は%sをアプリと同じフォルダに置くか、設定でhelperパスを指定してください", spoutHelperFileName)
	}
	return resolved, nil
}

func spoutHelperNextToExecutable(helperPath string) (string, bool) {
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	candidate := filepath.Join(filepath.Dir(exe), helperPath)
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		return candidate, true
	}
	return "", false
}

func embeddedSpoutHelperAvailable() bool {
	return embeddedSpoutAvailable && len(embeddedSpoutHelperEXE) > 0 && len(embeddedSpoutLibraryDLL) > 0
}

func resolveEmbeddedSpoutHelper() (string, error) {
	dir, err := embeddedSpoutHelperDir(embeddedSpoutHelperEXE, embeddedSpoutLibraryDLL)
	if err != nil {
		return "", err
	}
	helperPath := filepath.Join(dir, spoutHelperFileName)
	dllPath := filepath.Join(dir, spoutLibraryDLLName)
	if err := ensureEmbeddedSpoutFile(helperPath, embeddedSpoutHelperEXE, true); err != nil {
		return "", err
	}
	if err := ensureEmbeddedSpoutFile(dllPath, embeddedSpoutLibraryDLL, false); err != nil {
		return "", err
	}
	return helperPath, nil
}

func embeddedSpoutHelperDir(helperData []byte, dllData []byte) (string, error) {
	root := spoutHelperCacheRoot()
	if root == "" {
		return "", fmt.Errorf("埋め込みSpout helperの展開先を決定できません")
	}
	hasher := sha256.New()
	writeEmbeddedSpoutHashPart(hasher, spoutHelperFileName, helperData)
	writeEmbeddedSpoutHashPart(hasher, spoutLibraryDLLName, dllData)
	sum := hasher.Sum(nil)
	return filepath.Join(root, hex.EncodeToString(sum[:])[:16]), nil
}

func writeEmbeddedSpoutHashPart(hasher interface {
	Write([]byte) (int, error)
}, name string, data []byte) {
	var size [8]byte
	binary.LittleEndian.PutUint64(size[:], uint64(len(data)))
	_, _ = hasher.Write([]byte(name))
	_, _ = hasher.Write([]byte{0})
	_, _ = hasher.Write(size[:])
	_, _ = hasher.Write(data)
}

func spoutHelperCacheRoot() string {
	if spoutHelperCacheDirForTest != "" {
		return spoutHelperCacheDirForTest
	}
	if cacheDir, err := os.UserCacheDir(); err == nil && cacheDir != "" {
		return filepath.Join(cacheDir, "ClipForVRChat", "spout-helper")
	}
	return filepath.Join(os.TempDir(), "ClipForVRChat", "spout-helper")
}

func ensureEmbeddedSpoutFile(path string, data []byte, executable bool) error {
	if len(data) == 0 {
		return fmt.Errorf("埋め込みSpout helper資産が空です: %s", filepath.Base(path))
	}
	expected := sha256.Sum256(data)
	if existing, err := os.ReadFile(path); err == nil {
		actual := sha256.Sum256(existing)
		if actual == expected {
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), privateDirMode); err != nil {
		return fmt.Errorf("埋め込みSpout helperの展開フォルダを作成できません: %w", err)
	}
	mode := privateFileMode
	if executable {
		mode = 0700
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		return fmt.Errorf("埋め込みSpout helperを展開できません: %w", err)
	}
	if err := os.Chmod(path, mode); err != nil {
		return fmt.Errorf("埋め込みSpout helperの権限を設定できません: %w", err)
	}
	written, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("展開したSpout helperを検証できません: %w", err)
	}
	actual := sha256.Sum256(written)
	if actual != expected {
		return fmt.Errorf("展開したSpout helperのSHA-256が一致しません: %s", filepath.Base(path))
	}
	return nil
}

func CheckSpoutHelper(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) SpoutHelperStatus {
	helper, err := ResolveSpoutHelperPath(cfg.SpoutHelperPath)
	if err != nil {
		diagAutoCapture(logPath, "spout helper check missing: path=%q err=%v", cfg.SpoutHelperPath, err)
		return SpoutHelperStatus{Available: false, Message: err.Error()}
	}
	version, versionErr := checkSpoutHelperVersion(ctx, helper, logPath)
	if versionErr != nil {
		diagAutoCapture(logPath, "spout helper version error: helper=%q err=%v", helper, versionErr)
		return SpoutHelperStatus{Available: false, Path: helper, Message: "Spout helperは見つかりましたが実行確認に失敗しました: " + versionErr.Error()}
	}
	list, err := ListSpoutSenders(ctx, cfg, logPath)
	if err != nil {
		diagAutoCapture(logPath, "spout helper check list error: helper=%q err=%v", helper, err)
		return SpoutHelperStatus{Available: true, Path: helper, Version: version, Message: "Spout helperは実行できますが、sender一覧取得に失敗しました: " + err.Error(), Senders: list.Senders}
	}
	msg := "Spout helperを利用できます。"
	if len(list.Senders) == 0 {
		msg = "Spout helperを利用できます。senderがありません。VRChatでStream Cameraを起動してください。"
	}
	return SpoutHelperStatus{Available: true, Path: helper, Version: version, Message: msg, Senders: list.Senders}
}

func checkSpoutHelperVersion(ctx context.Context, helper string, logPath string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	commandCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(commandCtx, helper, "--version") // #nosec G204 -- local helper path from config or release folder.
	output, err := cmd.CombinedOutput()
	trimmed := trimCommandOutput(output)
	if commandCtx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("Spout helper --versionがタイムアウトしました")
	}
	if err != nil {
		return "", fmt.Errorf("Spout helper --versionに失敗しました: %v %s", err, trimmed)
	}
	var result SpoutHelperVersionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("Spout helper version JSONを解析できません: %w %s", err, trimmed)
	}
	if !result.OK {
		if result.Message == "" {
			result.Message = result.Code
		}
		return "", errors.New(result.Message)
	}
	diagAutoCapture(logPath, "spout helper version success: helper=%q version=%q", helper, result.Version)
	return result.Version, nil
}

func ListSpoutSenders(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) (SpoutListResult, error) {
	helper, err := ResolveSpoutHelperPath(cfg.SpoutHelperPath)
	if err != nil {
		return SpoutListResult{}, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	timeout := time.Duration(cfg.CaptureTimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	commandCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	diagAutoCapture(logPath, "spout list begin: helper=%q timeout_ms=%d", helper, timeout.Milliseconds())
	cmd := exec.CommandContext(commandCtx, helper, "--list-senders") // #nosec G204 -- local helper path from config or release folder.
	output, err := cmd.CombinedOutput()
	trimmed := trimCommandOutput(output)
	if commandCtx.Err() == context.DeadlineExceeded {
		return SpoutListResult{}, fmt.Errorf("Spout sender一覧取得がタイムアウトしました")
	}
	var result SpoutListResult
	parseErr := json.Unmarshal(output, &result)
	if parseErr != nil {
		if err != nil {
			return SpoutListResult{}, fmt.Errorf("Spout sender一覧取得に失敗しました: %v %s", err, trimmed)
		}
		return SpoutListResult{}, fmt.Errorf("Spout helperのsender一覧JSONを解析できません: %w %s", parseErr, trimmed)
	}
	if err != nil {
		if result.Message == "" {
			result.Message = result.Code
		}
		return result, fmt.Errorf("Spout sender一覧取得に失敗しました: %v %s", err, result.Message)
	}
	if !result.OK {
		if result.Message == "" {
			result.Message = result.Code
		}
		return result, errors.New(result.Message)
	}
	diagAutoCapture(logPath, "spout list success: helper=%q senders=%d", helper, len(result.Senders))
	return result, nil
}

func captureStreamFrameWithSpout(ctx context.Context, cfg AutoCaptureStreamConfig, outputPath string, logPath string) (SpoutCaptureResult, error) {
	helper, err := ResolveSpoutHelperPath(cfg.SpoutHelperPath)
	if err != nil {
		diagAutoCapture(logPath, "spout helper missing: path=%q err=%v", cfg.SpoutHelperPath, err)
		return SpoutCaptureResult{}, err
	}
	timeout := time.Duration(cfg.CaptureTimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return SpoutCaptureResult{}, err
	}
	commandCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	args := []string{"--capture", "--output", outputPath, "--timeout-ms", fmt.Sprintf("%d", timeout.Milliseconds())}
	if strings.TrimSpace(cfg.SpoutSenderName) != "" && !cfg.SpoutAutoSelect {
		args = append(args, "--sender", cfg.SpoutSenderName)
	}
	diagAutoCapture(logPath, "spout capture begin: helper=%q args=%q output=%q sender=%q auto_select=%t timeout_ms=%d os=%s", helper, strings.Join(args, " "), outputPath, cfg.SpoutSenderName, cfg.SpoutAutoSelect, timeout.Milliseconds(), runtime.GOOS)
	cmd := exec.CommandContext(commandCtx, helper, args...) // #nosec G204 -- local helper path from config or release folder.
	output, err := cmd.CombinedOutput()
	trimmed := trimCommandOutput(output)
	if commandCtx.Err() == context.DeadlineExceeded {
		return SpoutCaptureResult{}, fmt.Errorf("Spout取得がタイムアウトしました。VRChatでStream Cameraが起動しているか、sender設定を確認してください")
	}
	var result SpoutCaptureResult
	parseErr := json.Unmarshal(output, &result)
	if parseErr != nil {
		if err != nil {
			return SpoutCaptureResult{}, fmt.Errorf("Spout取得に失敗しました: %v %s", err, trimmed)
		}
		return SpoutCaptureResult{}, fmt.Errorf("Spout helperの取得結果JSONを解析できません: %w %s", parseErr, trimmed)
	}
	if err != nil {
		if result.Message == "" {
			result.Message = result.Code
		}
		diagAutoCapture(logPath, "spout capture helper error: code=%q message=%q senders=%d output=%q", result.Code, result.Message, len(result.Senders), trimmed)
		return result, fmt.Errorf("Spout取得に失敗しました: %s", result.Message)
	}
	if !result.OK {
		if result.Message == "" {
			result.Message = result.Code
		}
		return result, errors.New(result.Message)
	}
	if result.OutputPath == "" {
		result.OutputPath = outputPath
	}
	stats, err := validateCapturedImage(outputPath)
	if err != nil {
		diagAutoCapture(logPath, "spout capture invalid image: output=%q err=%v", outputPath, err)
		return result, err
	}
	diagAutoCapture(logPath, "spout capture success: output=%q sender=%q width=%d height=%d frame=%d captured_at=%q stats_format=%q stats_width=%d stats_height=%d samples=%d mean=%.2f stddev=%.2f near_white=%.4f near_black=%.4f transparent=%.4f", outputPath, result.SenderName, result.Width, result.Height, result.Frame, result.CapturedAt, stats.Format, stats.Width, stats.Height, stats.Samples, stats.Mean, stats.Stddev, stats.NearWhiteRatio, stats.NearBlackRatio, stats.TransparentRatio)
	return result, nil
}

func validateCapturedImage(path string) (capturedImageStats, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is generated by auto-capture output path.
	if err != nil {
		return capturedImageStats{}, fmt.Errorf("取得画像を読み込めません: %w", err)
	}
	if len(data) == 0 {
		return capturedImageStats{}, fmt.Errorf("取得画像が0バイトです")
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return capturedImageStats{}, fmt.Errorf("取得画像を画像として解析できません: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return capturedImageStats{}, fmt.Errorf("取得画像の解像度が不正です: %dx%d", cfg.Width, cfg.Height)
	}
	if format != "png" && format != "jpeg" {
		return capturedImageStats{}, fmt.Errorf("取得画像形式が不正です: %s", format)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return capturedImageStats{}, fmt.Errorf("取得画像を画像としてデコードできません: %w", err)
	}
	stats := capturedImageStats{Format: format, Width: cfg.Width, Height: cfg.Height}
	bounds := img.Bounds()
	stepX := 1
	stepY := 1
	maxSamples := 16384
	if pixels := bounds.Dx() * bounds.Dy(); pixels > maxSamples {
		for (bounds.Dx()/stepX)*(bounds.Dy()/stepY) > maxSamples {
			if bounds.Dx()/stepX >= bounds.Dy()/stepY {
				stepX++
			} else {
				stepY++
			}
		}
	}
	var sum float64
	var sumSq float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			r, g, b, a := img.At(x, y).RGBA()
			luma := (0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8))
			alpha := float64(a >> 8)
			sum += luma
			sumSq += luma * luma
			if luma >= 250 && alpha >= 250 {
				stats.NearWhiteRatio++
			}
			if luma <= 5 && alpha >= 250 {
				stats.NearBlackRatio++
			}
			if alpha <= 5 {
				stats.TransparentRatio++
			}
			stats.Samples++
		}
	}
	if stats.Samples == 0 {
		return stats, fmt.Errorf("取得画像からサンプルを取得できません")
	}
	samples := float64(stats.Samples)
	stats.Mean = sum / samples
	variance := (sumSq / samples) - (stats.Mean * stats.Mean)
	if variance < 0 {
		variance = 0
	}
	stats.Stddev = math.Sqrt(variance)
	stats.NearWhiteRatio /= samples
	stats.NearBlackRatio /= samples
	stats.TransparentRatio /= samples
	if stats.TransparentRatio > 0.99 {
		return stats, fmt.Errorf("取得画像がほぼ透明です")
	}
	if stats.NearWhiteRatio > 0.99 && stats.Stddev < 1.5 {
		return stats, fmt.Errorf("取得画像がほぼ白一色です。VRChat Stream Cameraの映像ではない可能性があります")
	}
	if stats.NearBlackRatio > 0.99 && stats.Stddev < 1.5 {
		return stats, fmt.Errorf("取得画像がほぼ黒一色です。VRChat Stream Cameraの映像ではない可能性があります")
	}
	return stats, nil
}

func trimCommandOutput(output []byte) string {
	trimmed := strings.TrimSpace(string(output))
	if len(trimmed) > 2000 {
		trimmed = trimmed[len(trimmed)-2000:]
	}
	return trimmed
}
