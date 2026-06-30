# プレイヤー中心ローカル座標系の仕様を確定する

## 問題

現在の自動撮影構図は、実機確認でワールド原点を基準にしたような挙動になり、ユーザーが想定する「プレイヤーを中心としたローカル座標系」になっていない。
正面、背面、斜めなどの構図を、プレイヤーの位置と向きに対する相対座標として扱う仕様が未確定である。

## 期待する挙動

自動撮影の各構図が、撮影時点のローカルプレイヤー位置/向きを基準にして決まり、正面は顔中心、背面は三人称視点のような後頭部側から少し見下ろす視点、斜めは正面側かつ少し横から少し見下ろす視点になる。

## 受け入れ条件

- [x] VRChat User Camera OSCのPose仕様と、現行コードのPose送信/受信処理を確認する。
- [x] プレイヤー中心ローカル座標系として採用する座標軸、回転、初期構図値を明文化する。
- [x] 現行実装との差分と実装方法を決める。
- [x] 実装作業を細かい日本語issueへ分割する。
- [x] この調査では実装コード変更を行わない。

## 調査結果

- 現行の `CameraViewConfig.CoordinateSpace` は設定として保存されるだけで、`AutoCaptureRunner.applyCameraView()` は `view.Pose` をそのまま `/usercamera/Pose` へ送っている。
- 初期構図は `defaultCameraView()` で `CoordinateSpace: "world"` 固定であり、`SaveCurrentCameraPoseToView()` と `AddCurrentCameraPoseAsView()` も保存時に `world` を明示している。
- VRChat公式のCamera OSC endpointでは `/usercamera/Pose` はカメラ位置/回転のGet/Setであり、プレイヤーroot位置/向きを返すendpointではない。
- VRChatのOSC Trackersは外部プログラムからVRChatへトラッカー位置を送る入力用途であり、ローカルプレイヤーroot poseの取得元としては扱えない。
- 組み込みAvatar Parametersには `VelocityX/Y/Z`、`AngularY`、`EyeHeightAsMeters` などはあるが、ローカルプレイヤーのワールド位置やroot yawそのものはない。

## 採用方針

v0.1.8では、ワールド座標をプレイヤー中心ローカル座標として扱うダミー実装は禁止する。
`player_local` は「撮影時点のローカルプレイヤー基準Poseが取得できる場合だけ、構図PoseをワールドPoseへ変換して `/usercamera/Pose` に送る」仕様として実装する。

標準OSC/OSCQueryだけでプレイヤー基準Poseを取得できないことが確定した場合は、機能を未対応として明示するか、手動基準Poseなどの実際に成立する方式へ仕様を切り替える。
その場合も、UIやsidecarで `player_local` と表示してワールド座標を送る実装にはしない。

## 分割issue

- [#138](138-define-player-local-coordinate-spec.md) `player_local` 座標仕様を確定する。
- [#139](139-investigate-player-basis-source.md) ローカルプレイヤー基準Poseの取得可否を実機調査する。
- [#140](140-implement-player-local-pose-transform.md) `player_local` からUser Camera Poseへの変換処理を実装する。
- [#141](141-integrate-player-local-coordinate-ui.md) 自動撮影構図の設定/保存/リセット/移動UIへ `player_local` を統合する。
- [#142](142-verify-player-local-camera-compositions.md) プレイヤー中心構図の実機確認手順を整備する。

## 参照

- VRChat OSC Overview: https://docs.vrchat.com/docs/osc-overview
- VRChat 2025.3.3 Open Beta Camera Endpoints: https://docs.vrchat.com/docs/vrchat-202533-openbeta
- VRChat OSC Trackers: https://docs.vrchat.com/docs/osc-trackers
- VRChat Animator Parameters: https://creators.vrchat.com/avatars/animator-parameters/
