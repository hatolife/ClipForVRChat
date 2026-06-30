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
