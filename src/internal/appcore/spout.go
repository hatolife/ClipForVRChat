package appcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

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
	OK         bool   `json:"ok"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	SenderName string `json:"senderName,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Frame      int64  `json:"frame,omitempty"`
	CapturedAt string `json:"capturedAt,omitempty"`
	OutputPath string `json:"outputPath,omitempty"`
}

type SpoutHelperStatus struct {
	Available bool              `json:"available"`
	Path      string            `json:"path,omitempty"`
	Message   string            `json:"message"`
	Senders   []SpoutSenderInfo `json:"senders,omitempty"`
}

func ResolveSpoutHelperPath(helperPath string) (string, error) {
	helperPath = strings.Trim(strings.TrimSpace(helperPath), `"`)
	if helperPath == "" {
		helperPath = "spout-capture.exe"
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
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), helperPath)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	resolved, err := exec.LookPath(helperPath)
	if err != nil {
		return "", fmt.Errorf("Spout helperが見つかりません。Release zipに含まれるspout-capture.exeをアプリと同じフォルダに置いてください")
	}
	return resolved, nil
}

func CheckSpoutHelper(ctx context.Context, cfg AutoCaptureStreamConfig, logPath string) SpoutHelperStatus {
	helper, err := ResolveSpoutHelperPath(cfg.SpoutHelperPath)
	if err != nil {
		diagAutoCapture(logPath, "spout helper check missing: path=%q err=%v", cfg.SpoutHelperPath, err)
		return SpoutHelperStatus{Available: false, Message: err.Error()}
	}
	list, err := ListSpoutSenders(ctx, cfg, logPath)
	if err != nil {
		diagAutoCapture(logPath, "spout helper check list error: helper=%q err=%v", helper, err)
		return SpoutHelperStatus{Available: true, Path: helper, Message: "Spout helperは実行できますが、sender一覧取得に失敗しました: " + err.Error()}
	}
	msg := "Spout helperを利用できます。"
	if len(list.Senders) == 0 {
		msg = "Spout helperを利用できます。senderがありません。VRChatでStream Cameraを起動してください。"
	}
	return SpoutHelperStatus{Available: true, Path: helper, Message: msg, Senders: list.Senders}
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
	if err != nil {
		return SpoutListResult{}, fmt.Errorf("Spout sender一覧取得に失敗しました: %v %s", err, trimmed)
	}
	var result SpoutListResult
	if err := json.Unmarshal(output, &result); err != nil {
		return SpoutListResult{}, fmt.Errorf("Spout helperのsender一覧JSONを解析できません: %w %s", err, trimmed)
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
	if err != nil {
		return SpoutCaptureResult{}, fmt.Errorf("Spout取得に失敗しました: %v %s", err, trimmed)
	}
	var result SpoutCaptureResult
	if err := json.Unmarshal(output, &result); err != nil {
		return SpoutCaptureResult{}, fmt.Errorf("Spout helperの取得結果JSONを解析できません: %w %s", err, trimmed)
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
	if err := validateCapturedImage(outputPath); err != nil {
		return result, err
	}
	diagAutoCapture(logPath, "spout capture success: output=%q sender=%q width=%d height=%d frame=%d captured_at=%q", outputPath, result.SenderName, result.Width, result.Height, result.Frame, result.CapturedAt)
	return result, nil
}

func validateCapturedImage(path string) error {
	data, err := os.ReadFile(path) // #nosec G304 -- path is generated by auto-capture output path.
	if err != nil {
		return fmt.Errorf("取得画像を読み込めません: %w", err)
	}
	if len(data) == 0 {
		return fmt.Errorf("取得画像が0バイトです")
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("取得画像を画像として解析できません: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return fmt.Errorf("取得画像の解像度が不正です: %dx%d", cfg.Width, cfg.Height)
	}
	if format != "png" && format != "jpeg" {
		return fmt.Errorf("取得画像形式が不正です: %s", format)
	}
	return nil
}

func trimCommandOutput(output []byte) string {
	trimmed := strings.TrimSpace(string(output))
	if len(trimmed) > 2000 {
		trimmed = trimmed[len(trimmed)-2000:]
	}
	return trimmed
}
