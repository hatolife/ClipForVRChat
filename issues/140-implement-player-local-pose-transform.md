# `player_local` からUser Camera Poseへの変換処理を実装する

## 問題

現行の `AutoCaptureRunner.applyCameraView()` は保存済みPoseをそのまま `/usercamera/Pose` に送っている。
このままでは `CoordinateSpace` を `player_local` にしても、実際にはワールド座標として動作してしまう。

## 期待する挙動

有効なローカルプレイヤー基準Poseがある場合、`player_local` の構図値をワールド座標のUser Camera Poseへ変換してからOSC送信する。
変換不能な場合は、ワールド座標を送って成功扱いにせず、明確なエラーにする。

## 受け入れ条件

- `PlayerBasis` 相当の入力型を定義し、位置、Yaw、アンカー高さ/注視点を表現できる。
- `CameraViewConfig` の `CoordinateSpace == "player_local"` を検出し、送信前にワールドPoseへ変換する。
- Yaw 0/90/180/-90度、正面/背後/斜めoffset、Pitch/Rollを含む単体テストを追加する。
- `world` 既存構図は従来通りPoseをそのまま送る。
- 基準Poseが取得できない場合は、結果/診断ログに「プレイヤー基準Poseを取得できないため撮影できない」と分かるエラーを出す。
- 変換処理は副作用のない関数として切り出し、OSC送信処理と分離してテストできる。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- 現行の `TransformPlayerLocalPose(basis, local)` は撮影時のworld pose算出として維持し、逆変換 `InverseTransformPlayerLocalPose(basis, world)` を追加する。
- 逆変換は現在Pose保存/追加で、world poseを `player_local` viewのローカルPoseへ変換するために使う。
- `applyCameraView()` は送信したresolved world poseを返す形に変更し、sidecar/埋め込みメタデータへ「設定Pose」と「送信Pose」を区別して残す。
- PlayerBasis型を新設するより、v0.1.8では `AutoCapturePlayerLocalConfig` を基準Pose保持型として使う。将来、自動player root取得を入れる場合に `PlayerBasis` へ分離する。

対象ファイル:

- `src/internal/appcore/player_local.go`
- `src/internal/appcore/player_local_test.go`
- `src/internal/appcore/autocapture.go`
- `src/app.go`

小タスク:

- `InverseTransformPlayerLocalPose()` を追加する。
- Yaw 0/90/180/-90 の forward/inverse round-trip testを追加する。
- 正面/背後/斜めの初期offsetが基準Yawに追従するテストを追加する。
- Pitch/Roll/Yawの加算/減算テストを追加する。
- `ResolveCameraViewPose()` の未キャリブレーションエラーは現行文言を維持しつつ、診断ログにview IDと座標系を残す。
- `AutoCaptureSidecar` と `AutoCaptureEmbeddedMetadata` に `ResolvedPose` または同等フィールドを追加する。

確認方法:

- `TransformPlayerLocalPose(basis, InverseTransformPlayerLocalPose(basis, world))` が元のworld poseへ戻る。
- `world` 構図は従来通り保存値をそのまま送信する。
- `player_local` 未キャリブレーションでは `/usercamera/Pose` を送信せず失敗する。
