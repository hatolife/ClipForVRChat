package appcore

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCapturedImageRejectsBlankFrames(t *testing.T) {
	tests := []struct {
		name string
		fill color.RGBA
		want string
	}{
		{name: "white", fill: color.RGBA{R: 255, G: 255, B: 255, A: 255}, want: "ほぼ白"},
		{name: "black", fill: color.RGBA{A: 255}, want: "ほぼ黒"},
		{name: "transparent", fill: color.RGBA{}, want: "ほぼ透明"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "capture.png")
			writeSolidPNG(t, path, tt.fill)
			_, err := validateCapturedImage(path)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateCapturedImageAcceptsVariedFrame(t *testing.T) {
	path := filepath.Join(t.TempDir(), "capture.png")
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 8), G: uint8(y * 8), B: 80, A: 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := validateCapturedImage(path); err != nil {
		t.Fatal(err)
	}
}

func TestClassifySpoutCaptureFailure(t *testing.T) {
	tests := []struct {
		name     string
		result   SpoutCaptureResult
		wantKind spoutCaptureFailureKind
		wantText string
	}{
		{
			name:     "sender missing",
			result:   SpoutCaptureResult{Code: "sender_not_found"},
			wantKind: spoutCaptureFailureKindSenderMissing,
			wantText: "Spout senderがありません。",
		},
		{
			name:     "old timeout is no new frame",
			result:   SpoutCaptureResult{Code: "capture_timeout"},
			wantKind: spoutCaptureFailureKindNoNewFrame,
			wantText: "新しいフレームを受信できませんでした",
		},
		{
			name:     "no new frame",
			result:   SpoutCaptureResult{Code: "capture_no_new_frame"},
			wantKind: spoutCaptureFailureKindNoNewFrame,
			wantText: "新しいフレームを受信できませんでした",
		},
		{
			name:     "receive stalled",
			result:   SpoutCaptureResult{Code: "capture_receive_stalled"},
			wantKind: spoutCaptureFailureKindReceiveStalled,
			wantText: "画像を受信できませんでした",
		},
		{
			name:     "transparent frame",
			result:   SpoutCaptureResult{Code: "capture_blank_frame", FrameStats: &SpoutFrameStats{Samples: 1024, TransparentRatio: 1}},
			wantKind: spoutCaptureFailureKindTransparentFrame,
			wantText: "ほぼ透明",
		},
		{
			name: "transparent frame with stuck sender frame",
			result: SpoutCaptureResult{
				Code:              "capture_blank_frame",
				Frame:             0,
				FirstFrame:        0,
				LastReceivedFrame: 0,
				ReceiveSuccesses:  264,
				FrameStats:        &SpoutFrameStats{Samples: 1024, TransparentRatio: 1},
			},
			wantKind: spoutCaptureFailureKindTransparentFrame,
			wantText: "frame番号が進まず",
		},
		{
			name:     "black frame",
			result:   SpoutCaptureResult{Code: "capture_blank_frame", FrameStats: &SpoutFrameStats{Samples: 1024, Mean: 0, Stddev: 0, NearBlackRatio: 1}},
			wantKind: spoutCaptureFailureKindBlackFrame,
			wantText: "ほぼ黒一色",
		},
		{
			name:     "white frame",
			result:   SpoutCaptureResult{Code: "capture_blank_frame", FrameStats: &SpoutFrameStats{Samples: 1024, Mean: 255, Stddev: 0, NearWhiteRatio: 1}},
			wantKind: spoutCaptureFailureKindWhiteFrame,
			wantText: "ほぼ白一色",
		},
		{
			name:     "fallback message",
			result:   SpoutCaptureResult{Code: "custom_error", Message: "custom detail"},
			wantKind: spoutCaptureFailureKindUnknown,
			wantText: "custom detail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, message := classifySpoutCaptureFailure(tt.result)
			if kind != tt.wantKind {
				t.Fatalf("kind = %q, want %q", kind, tt.wantKind)
			}
			if !strings.Contains(message, tt.wantText) {
				t.Fatalf("message = %q, want substring %q", message, tt.wantText)
			}
		})
	}
}

func TestClassifyCapturedImageValidationFailure(t *testing.T) {
	tests := []struct {
		name     string
		stats    capturedImageStats
		wantKind spoutCaptureFailureKind
		wantText string
	}{
		{
			name:     "transparent",
			stats:    capturedImageStats{TransparentRatio: 1},
			wantKind: spoutCaptureFailureKindTransparentFrame,
			wantText: "ほぼ透明",
		},
		{
			name:     "black",
			stats:    capturedImageStats{NearBlackRatio: 1, Stddev: 0},
			wantKind: spoutCaptureFailureKindBlackFrame,
			wantText: "ほぼ黒一色",
		},
		{
			name:     "white",
			stats:    capturedImageStats{NearWhiteRatio: 1, Stddev: 0},
			wantKind: spoutCaptureFailureKindWhiteFrame,
			wantText: "ほぼ白一色",
		},
		{
			name:     "decode fallback",
			stats:    capturedImageStats{},
			wantKind: spoutCaptureFailureKindInvalidImage,
			wantText: "decode failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, message := classifyCapturedImageValidationFailure(tt.stats, fmt.Errorf("decode failed"))
			if kind != tt.wantKind {
				t.Fatalf("kind = %q, want %q", kind, tt.wantKind)
			}
			if !strings.Contains(message, tt.wantText) {
				t.Fatalf("message = %q, want substring %q", message, tt.wantText)
			}
		})
	}
}

func writeSolidPNG(t *testing.T, path string, fill color.RGBA) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.SetRGBA(x, y, fill)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
