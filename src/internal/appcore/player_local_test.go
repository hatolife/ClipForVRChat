package appcore

import (
	"math"
	"testing"
)

func TestTransformPlayerLocalPoseRoundTrip(t *testing.T) {
	local := CameraPoseConfig{
		Position: CameraVector3Config{X: 1, Y: 2, Z: 3},
		Rotation: CameraVector3Config{X: 4, Y: 5, Z: 6},
	}
	cases := []struct {
		name     string
		yaw      float64
		wantPos  CameraVector3Config
		wantRotY float64
	}{
		{name: "yaw0", yaw: 0, wantPos: CameraVector3Config{X: 11, Y: 3, Z: 23}, wantRotY: 5},
		{name: "yaw90", yaw: 90, wantPos: CameraVector3Config{X: 13, Y: 3, Z: 19}, wantRotY: 95},
		{name: "yaw180", yaw: 180, wantPos: CameraVector3Config{X: 9, Y: 3, Z: 17}, wantRotY: 185},
		{name: "yaw-90", yaw: -90, wantPos: CameraVector3Config{X: 7, Y: 3, Z: 21}, wantRotY: -85},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			basis := CameraPoseConfig{
				Position: CameraVector3Config{X: 10, Y: 1, Z: 20},
				Rotation: CameraVector3Config{X: 11, Y: tt.yaw, Z: 7},
			}
			world := TransformPlayerLocalPose(basis, local)
			assertPoseClose(t, world, CameraPoseConfig{
				Position: tt.wantPos,
				Rotation: CameraVector3Config{X: 15, Y: tt.wantRotY, Z: 13},
			})
			got := InverseTransformPlayerLocalPose(basis, world)
			assertPoseClose(t, got, local)
		})
	}
}

func TestResolveCameraViewPoseRequiresPlayerBasis(t *testing.T) {
	cfg := DefaultAutoCaptureConfig()
	view := cfg.Views[0]
	view.Pose = CameraPoseConfig{
		Position: CameraVector3Config{X: 0.8, Y: 0.2, Z: 1.1},
		Rotation: CameraVector3Config{X: 8, Y: -145, Z: 0},
	}
	view.CoordinateSpace = "player_local"
	_, err := ResolveCameraViewPose(cfg, view)
	if err == nil {
		t.Fatal("ResolveCameraViewPose should require calibrated player basis")
	}
	cfg.PlayerLocal.Calibrated = true
	cfg.PlayerLocal.BasisPose = CameraPoseConfig{
		Position: CameraVector3Config{X: 10, Y: 1, Z: 20},
		Rotation: CameraVector3Config{Y: 90},
	}
	got, err := ResolveCameraViewPose(cfg, view)
	if err != nil {
		t.Fatal(err)
	}
	assertPoseClose(t, got, TransformPlayerLocalPose(cfg.PlayerLocal.BasisPose, view.Pose))
}

func assertPoseClose(t *testing.T, got CameraPoseConfig, want CameraPoseConfig) {
	t.Helper()
	assertClose(t, got.Position.X, want.Position.X)
	assertClose(t, got.Position.Y, want.Position.Y)
	assertClose(t, got.Position.Z, want.Position.Z)
	assertClose(t, got.Rotation.X, want.Rotation.X)
	assertClose(t, got.Rotation.Y, want.Rotation.Y)
	assertClose(t, got.Rotation.Z, want.Rotation.Z)
}

func assertClose(t *testing.T, got float64, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.000001 {
		t.Fatalf("got %f, want %f", got, want)
	}
}
