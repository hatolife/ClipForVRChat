package appcore

import (
	"math"
	"testing"
)

func TestTransformPlayerLocalPoseYaw(t *testing.T) {
	basis := CameraPoseConfig{
		Position: CameraVector3Config{X: 10, Y: 1, Z: 20},
		Rotation: CameraVector3Config{Y: 90},
	}
	local := CameraPoseConfig{
		Position: CameraVector3Config{X: 1, Y: 2, Z: 3},
		Rotation: CameraVector3Config{X: 4, Y: 5, Z: 6},
	}
	got := TransformPlayerLocalPose(basis, local)
	assertClose(t, got.Position.X, 13)
	assertClose(t, got.Position.Y, 3)
	assertClose(t, got.Position.Z, 19)
	assertClose(t, got.Rotation.X, 4)
	assertClose(t, got.Rotation.Y, 95)
	assertClose(t, got.Rotation.Z, 6)
}

func TestResolveCameraViewPoseRequiresPlayerBasis(t *testing.T) {
	cfg := DefaultAutoCaptureConfig()
	view := cfg.Views[0]
	view.CoordinateSpace = "player_local"
	_, err := ResolveCameraViewPose(cfg, view)
	if err == nil {
		t.Fatal("ResolveCameraViewPose should require calibrated player basis")
	}
	cfg.PlayerLocal.Calibrated = true
	got, err := ResolveCameraViewPose(cfg, view)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, got.Position.Z, view.Pose.Position.Z)
}

func assertClose(t *testing.T, got float64, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.000001 {
		t.Fatalf("got %f, want %f", got, want)
	}
}
