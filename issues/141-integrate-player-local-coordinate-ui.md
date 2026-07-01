# 自動撮影構図UIへ `player_local` を統合する

## 問題

現在Pose保存、現在Pose追加、初期ポーズリセット、設定Poseへカメラ移動、テスト撮影のUIは、実質的にワールド座標前提で動く。
ユーザーがプレイヤー中心ローカル構図として設定しても、保存値や表示文言が実際の動作と一致しない可能性がある。

## 期待する挙動

自動撮影タブで各構図が `player_local` か `world` かを誤解なく設定でき、現在Pose保存/追加/リセット/カメラ移動/テスト撮影が座標系に応じた動作をする。

## 受け入れ条件

- `CameraViewConfig.Normalize()` に `player_local` を正式な許可値として追加し、未接続の値をユーザーに見せない。
- 初期の正面/背後/斜め構図を `player_local` の妥当な値へ更新する。ただし基準Pose取得が未成立の場合は未対応表示にする。
- 「現在Poseから追加」「現在Poseを保存」は、基準Poseがある場合だけ `player_local` へ変換して保存し、ない場合は理由を表示する。
- 「設定値にカメラを移動」は、`player_local` では基準PoseからワールドPoseを算出して移動し、変換不能なら送信しない。
- 「初期ポーズにリセット」は、構図ごとに初期 `player_local` 値、ズーム、撮影トグルを復元する。
- 構図カード内の表示文言から「Pose鮮度」のような意味が伝わりにくい表現をなくし、OSC受信値の有効期限として説明する。
- sidecar JSONと診断ログに、送信したワールドPoseと元の `player_local` 設定を区別して残す。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- `SaveCurrentCameraPoseToView()` は選択中viewの `CoordinateSpace` を尊重する。`player_local` の場合は基準Poseがあるときだけ逆変換して保存し、基準Poseがない場合は保存しない。
- `AddCurrentCameraPoseAsView()` は各構図カード内ボタンから呼ばれるため、参照元view IDを受け取り、参照元の座標系/ズーム/表示マスクを引き継いだ新規構図を作る。
- 既存のWails公開メソッド名を変える場合は [#156](156-add-wails-api-surface-check.md) のAPI surface検査で漏れを検出する。
- 「初期Poseへ戻す」は [#138](138-define-player-local-coordinate-spec.md) の `player_local` 初期値へ戻す。
- 「このPoseへカメラ移動」と「テスト撮影」は `ResolveCameraViewPose()` を通す現行設計を維持する。

対象ファイル:

- `src/app.go`
- `src/internal/appcore/player_local.go`
- `src/internal/appcore/config.go`
- `src/frontend/src/main.js`
- `src/app_test.go` または `src/internal/appcore/player_local_test.go`

小タスク:

- `SaveCurrentCameraPoseToView()` の `CoordinateSpace = "world"` 強制をやめる。
- `player_local` 保存時は `InverseTransformPlayerLocalPose(cfg.AutoCapture.PlayerLocal.BasisPose, pose)` を使う。
- `AddCurrentCameraPoseAsView(viewID string)` または新APIを追加し、frontendの `addCurrentCameraPoseAsView(cameraView)` から参照view IDを渡す。
- 新規構図の `SortOrder`、`Enabled`、`Calibrated`、`Zoom`、`LookAtMe`、`LocalPlayer`、`RemotePlayer`、`Environment` を参照viewから引き継ぐ。
- frontendは `player_local` で基準未保存の場合、現在Pose保存/追加ボタンの失敗理由をそのまま表示する。
- sidecar JSONに元のview poseと送信済みworld poseが残ることを確認する。

確認方法:

- `world` viewで現在Pose保存するとworld poseがそのまま保存される。
- `player_local` viewで現在Pose保存するとローカルoffsetとして保存され、カメラ移動時に元のworld poseへ戻る。
- 基準Pose未保存で `player_local` 保存/追加は明確に失敗し、設定を壊さない。
