package appcore

import (
	"fmt"
	"math"
)

func ResolveCameraViewPose(cfg AutoCaptureConfig, view CameraViewConfig) (CameraPoseConfig, error) {
	switch view.CoordinateSpace {
	case "", "world":
		return view.Pose, nil
	case "player_local":
		if !cfg.PlayerLocal.Calibrated {
			return CameraPoseConfig{}, fmt.Errorf("プレイヤー基準Poseが未設定のため、player_local構図を撮影できません。自動撮影タブで現在Poseをプレイヤー基準として保存してください")
		}
		return TransformPlayerLocalPose(cfg.PlayerLocal.BasisPose, view.Pose), nil
	default:
		return CameraPoseConfig{}, fmt.Errorf("未対応の座標系です: %s", view.CoordinateSpace)
	}
}

func TransformPlayerLocalPose(basis CameraPoseConfig, local CameraPoseConfig) CameraPoseConfig {
	yaw := basis.Rotation.Y * math.Pi / 180
	sinY := math.Sin(yaw)
	cosY := math.Cos(yaw)
	x := local.Position.X*cosY + local.Position.Z*sinY
	z := -local.Position.X*sinY + local.Position.Z*cosY
	return CameraPoseConfig{
		Position: CameraVector3Config{
			X: basis.Position.X + x,
			Y: basis.Position.Y + local.Position.Y,
			Z: basis.Position.Z + z,
		},
		Rotation: CameraVector3Config{
			X: basis.Rotation.X + local.Rotation.X,
			Y: basis.Rotation.Y + local.Rotation.Y,
			Z: basis.Rotation.Z + local.Rotation.Z,
		},
	}
}
