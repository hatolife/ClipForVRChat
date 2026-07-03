# 自動撮影後にUser Camera状態をできるだけ元へ戻す

## 問題

自動撮影は `/usercamera/Mode`、`/usercamera/Streaming`、`/usercamera/Pose`、Zoom/Exposure/各種maskなどを変更するが、撮影後に元の状態へ戻す仕組みが不足している。
現状は成功時中心に `Streaming=false` と `Close` を送るだけで、撮影失敗、キャンセル、Spout取得失敗でもCamera Mode、Streaming、SmoothMovement、Pose、Zoom、Exposure、mask類がユーザーの意図しない状態で残る可能性がある。

ユーザー確認メモ:

- `/usercamera/Mode value=2` でStreamモード
- `/usercamera/SmoothMovement value=1` でSmoothをON
- `/usercamera/Streaming value=1` でSpoutをON

VRChat公式の2025.3.3 Open Beta notesでは、User Camera OSC endpoints はMode/Pose/Toggles/Slidersがread/write accessと説明されているため、受信できる値は撮影前スナップショットとして使える可能性がある。

## 期待する挙動

自動撮影開始前に、可能な範囲でUser Cameraの現在状態を記録する。
撮影後は成功/失敗/キャンセルに関係なく、可能な範囲で記録した状態へ戻す。

撮影前の現在値を取得できない項目については、設定画面の末尾付近で「撮影後に書き戻す値」を設定できるようにし、その値をフォールバックとして使う。

最低限、次を復元対象にする。

- Camera Mode
- Streaming
- SmoothMovement
- Pose
- Zoom
- Exposure
- FocalDistance
- Aperture
- LookAtMe
- ShowUIInCamera
- LocalPlayer / RemotePlayer / Environment mask
- GreenScreen

## 受け入れ条件

- [x] User Camera OSC受信器が `/usercamera/Pose` 以外のMode/Toggles/Slidersも保持できる。
- [x] 自動撮影開始前のUser Camera状態が新鮮なら、その値を撮影後に復元する。
- [x] 受信できない項目は設定画面の復元用デフォルト値を使って書き戻せる。
- [x] 成功、失敗、キャンセルのいずれでも復元処理ができるだけ実行される。
- [x] Stream方式では `/usercamera/Mode=2`、`/usercamera/Streaming=true`、必要に応じて `/usercamera/SmoothMovement=true` を送る既存経路を維持する。
- [x] 復元処理の成否と、スナップショット由来かフォールバック由来かを診断ログで追える。
- [x] 設定画面末尾付近に、復元用デフォルト値の設定を追加する。
- [x] Go test、frontend build、template literal check、Wails API surface checkが通る。

## 対応メモ

- `AutoCaptureRestoreConfig` を追加し、撮影前スナップショット優先/フォールバック設定/受信値有効秒数を設定化した。
- OSC受信器で `/usercamera/Mode`、`/usercamera/Pose`、User Camera toggle、sliderを保持し、撮影実行時に新鮮な値だけをランナーへ注入する。
- `RunOnce` でOSC接続が開いた後に `defer` 復元を登録し、撮影失敗やキャンセルでもMode、Streaming、SmoothMovement、Pose、Zoom、Exposure、mask類などを書き戻す。
- Stream方式開始時は `/usercamera/Mode=2`、`/usercamera/SmoothMovement=true`、`/usercamera/Streaming=true` を送る。
- 復元開始ログに snapshot/fallback/target の値一覧を出し、各OSC送信の成否を記録する。
- 検証: `go test ./...`、`npm run build`、`check-frontend-template-literals`、`check-wails-api-surface`。

## 実装メモ

- 受信器はすでに `/usercamera/Pose` を保持しているため、同じUDP受信ループで `/usercamera/Mode`、bool toggle、float sliderをparseして保持する。
- 自動撮影実行前に `prepareAutoCaptureConfigForRunLocked` で最新スナップショットを `Config.AutoCapture` へ注入し、`AutoCaptureRunner.RunOnce` 側で `defer` 復元できる形にする。
- VRChat側から全項目が即時送られる保証はないため、スナップショットがない/古い項目は設定値フォールバックで扱う。
- まずv0.1.8では「完全な初期状態取得」より、「撮影後に安全に戻す」「ログから戻せなかった項目が分かる」ことを優先する。
