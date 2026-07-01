package appcore

import (
	"fmt"
	"math"
	"time"
)

type AvatarOSCBasisAxisSample struct {
	Magnitude float64
	SignFlag  float64
	Present   bool
}

type AvatarOSCBasisVectorSample struct {
	X AvatarOSCBasisAxisSample
	Y AvatarOSCBasisAxisSample
	Z AvatarOSCBasisAxisSample
}

type AvatarOSCBasisSample struct {
	Position   AvatarOSCBasisVectorSample
	Forward    AvatarOSCBasisVectorSample
	ReceivedAt time.Time
}

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

func ResolvePlayerLocalBasisPose(cfg AutoCaptureConfig, avatarOSCBasis AvatarOSCBasisSample, now time.Time) (CameraPoseConfig, error) {
	cfg.Normalize()
	switch cfg.PlayerLocal.BasisSource {
	case PlayerLocalBasisSourceAvatarOSC:
		return RestoreAvatarOSCBasisPose(avatarOSCBasis, cfg.PlayerLocal.AvatarOSC, now)
	default:
		if !cfg.PlayerLocal.Calibrated {
			return CameraPoseConfig{}, fmt.Errorf("プレイヤー基準Poseが未設定のため、player_local構図を撮影できません。自動撮影タブで現在Poseをプレイヤー基準として保存してください")
		}
		return cfg.PlayerLocal.BasisPose, nil
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

func RestoreAvatarOSCBasisPose(sample AvatarOSCBasisSample, cfg AutoCapturePlayerLocalAvatarOSCConfig, now time.Time) (CameraPoseConfig, error) {
	cfg.Normalize()
	if err := ValidateAvatarOSCBasisSample(sample, cfg, now); err != nil {
		return CameraPoseConfig{}, err
	}
	position, err := RestoreAvatarOSCBasisVector(sample.Position, cfg, cfg.MaxAbsPosition)
	if err != nil {
		return CameraPoseConfig{}, err
	}
	forward, err := RestoreAvatarOSCBasisVector(sample.Forward, cfg, cfg.MaxAbsForward)
	if err != nil {
		return CameraPoseConfig{}, err
	}
	yaw, err := YawFromForwardVector(forward)
	if err != nil {
		return CameraPoseConfig{}, err
	}
	return CameraPoseConfig{
		Position: position,
		Rotation: CameraVector3Config{Y: yaw},
	}, nil
}

func RestoreAvatarOSCBasisVector(sample AvatarOSCBasisVectorSample, cfg AutoCapturePlayerLocalAvatarOSCConfig, maxAbs float64) (CameraVector3Config, error) {
	x, err := RestoreAvatarOSCBasisAxis(sample.X, cfg)
	if err != nil {
		return CameraVector3Config{}, fmt.Errorf("x axis: %w", err)
	}
	y, err := RestoreAvatarOSCBasisAxis(sample.Y, cfg)
	if err != nil {
		return CameraVector3Config{}, fmt.Errorf("y axis: %w", err)
	}
	z, err := RestoreAvatarOSCBasisAxis(sample.Z, cfg)
	if err != nil {
		return CameraVector3Config{}, fmt.Errorf("z axis: %w", err)
	}
	vec := CameraVector3Config{X: x, Y: y, Z: z}
	if !isFiniteFloat64(vec.X) || !isFiniteFloat64(vec.Y) || !isFiniteFloat64(vec.Z) {
		return CameraVector3Config{}, fmt.Errorf("vector contains non-finite value")
	}
	if math.Abs(vec.X) > maxAbs || math.Abs(vec.Y) > maxAbs || math.Abs(vec.Z) > maxAbs {
		return CameraVector3Config{}, fmt.Errorf("vector out of range")
	}
	return vec, nil
}

func RestoreAvatarOSCBasisAxis(sample AvatarOSCBasisAxisSample, cfg AutoCapturePlayerLocalAvatarOSCConfig) (float64, error) {
	if !sample.Present {
		return 0, fmt.Errorf("missing axis sample")
	}
	if !isFiniteFloat64(sample.Magnitude) || !isFiniteFloat64(sample.SignFlag) {
		return 0, fmt.Errorf("axis contains non-finite value")
	}
	magnitude := sample.Magnitude
	if cfg.InvertMagnitude {
		magnitude = 1 - magnitude
	}
	value := magnitude * cfg.PositionScale
	if sample.SignFlag > cfg.PositiveFlagThreshold {
		return value, nil
	}
	return -value, nil
}

func YawFromForwardVector(forward CameraVector3Config) (float64, error) {
	if !isFiniteFloat64(forward.X) || !isFiniteFloat64(forward.Z) {
		return 0, fmt.Errorf("forward vector contains non-finite value")
	}
	horizontalMagnitude := math.Hypot(forward.X, forward.Z)
	if horizontalMagnitude == 0 {
		return 0, fmt.Errorf("forward vector has no horizontal component")
	}
	return math.Atan2(forward.X, forward.Z) * 180 / math.Pi, nil
}

func ValidateAvatarOSCBasisSample(sample AvatarOSCBasisSample, cfg AutoCapturePlayerLocalAvatarOSCConfig, now time.Time) error {
	cfg.Normalize()
	if sample.ReceivedAt.IsZero() {
		return fmt.Errorf("missing avatar OSC basis timestamp")
	}
	if sample.Position.X.Present != sample.Position.Y.Present ||
		sample.Position.X.Present != sample.Position.Z.Present ||
		sample.Forward.X.Present != sample.Forward.Y.Present ||
		sample.Forward.X.Present != sample.Forward.Z.Present {
		return fmt.Errorf("partial avatar OSC basis sample")
	}
	if !sample.Position.X.Present || !sample.Forward.X.Present {
		return fmt.Errorf("missing avatar OSC basis sample")
	}
	if now.IsZero() {
		now = time.Now()
	}
	if sample.ReceivedAt.After(now.Add(10 * time.Second)) {
		return fmt.Errorf("avatar OSC basis timestamp is in the future")
	}
	age := now.Sub(sample.ReceivedAt)
	if age < 0 {
		age = -age
	}
	if age > time.Duration(cfg.FreshnessSec)*time.Second {
		return fmt.Errorf("stale avatar OSC basis")
	}
	if !avatarOSCBasisVectorSampleValid(sample.Position) || !avatarOSCBasisVectorSampleValid(sample.Forward) {
		return fmt.Errorf("invalid avatar OSC basis values")
	}
	return nil
}

func avatarOSCBasisVectorSampleValid(sample AvatarOSCBasisVectorSample) bool {
	return avatarOSCBasisAxisSampleValid(sample.X) &&
		avatarOSCBasisAxisSampleValid(sample.Y) &&
		avatarOSCBasisAxisSampleValid(sample.Z)
}

func avatarOSCBasisAxisSampleValid(sample AvatarOSCBasisAxisSample) bool {
	if !sample.Present {
		return false
	}
	if !isFiniteFloat64(sample.Magnitude) || !isFiniteFloat64(sample.SignFlag) {
		return false
	}
	if math.Abs(sample.Magnitude) > 1 {
		return false
	}
	if math.Abs(sample.SignFlag) > 1 {
		return false
	}
	return true
}

func InverseTransformPlayerLocalPose(basis CameraPoseConfig, world CameraPoseConfig) CameraPoseConfig {
	yaw := basis.Rotation.Y * math.Pi / 180
	sinY := math.Sin(yaw)
	cosY := math.Cos(yaw)
	dx := world.Position.X - basis.Position.X
	dz := world.Position.Z - basis.Position.Z
	x := dx*cosY - dz*sinY
	z := dx*sinY + dz*cosY
	return CameraPoseConfig{
		Position: CameraVector3Config{
			X: x,
			Y: world.Position.Y - basis.Position.Y,
			Z: z,
		},
		Rotation: CameraVector3Config{
			X: world.Rotation.X - basis.Rotation.X,
			Y: world.Rotation.Y - basis.Rotation.Y,
			Z: world.Rotation.Z - basis.Rotation.Z,
		},
	}
}

func isFiniteFloat64(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
