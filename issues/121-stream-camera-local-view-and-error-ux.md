# 自動撮影をローカル視点基準とStream方式中心に修正する

## 問題

v0.1.8-rc8の初期Poseはワールド原点基準で撮影され、ユーザーが想定するプレーヤー中心の正面、背面、斜め視点にならない。
Photo方式ではシャッター音が複数回鳴り、カメラが閉じられた後の次回自動撮影で写真検出に失敗する。
失敗結果が履歴に残り、設定画面の「Pose鮮度」「自動撮影管理フォルダ」の説明も用途が分かりにくい。
本命のStreamカメラ方式が未実装のままで、Photo方式が実用上のフォールバックに留まっている。

## 期待する挙動

初期構図はプレーヤー中心のローカル視点として扱い、各構図を個別にテスト撮影できる。
自動撮影の失敗は履歴カードとして積まず、画面上のエラー文で理由を具体的に確認できる。
UI文言は、現在Pose保存の鮮度と出力先フォルダの用途が分かる表現にする。
Stream方式はダミーではなく、実際にVRChat Stream Cameraの映像から静止画を保存する実装方針へ進める。

## 受け入れ条件

- 初期構図のPoseはワールド原点ではなく、ユーザー視点を基準にした相対構図として扱える。
- 正面、背面、斜めの各構図にテスト撮影ボタンがある。
- 自動撮影で失敗しただけの結果は履歴に追加しない。
- 自動撮影エラーは「写真ファイル未検出」「カメラ未表示/起動失敗の可能性」など、次の確認につながる文で表示する。
- 「Pose鮮度」は現在Pose保存用の説明に変更する。
- 「自動撮影管理フォルダ」はStream方式の切り出し保存先として分かる説明と初期値に変更する。
- Stream方式を実装する。実装できない制約がある場合は、未実装のUIを出さず、具体的な技術的制約をIssueに記録する。
- 既存のPhoto方式、sidecar JSON、Discord投稿、通常画像処理を壊さない。

## 対応メモ

- rc9実機確認で、Stream方式でも `/usercamera/Mode` にPhoto Camera相当の `1` を送っていたため、Stream方式では `2` を送るよう修正した。

## 対応内容

- v0.1.8実装に合わせて対応済み。Stream方式はSpout主経路へ置き換え、ffmpeg主導線は廃止。詳細はREADME/SPEC/RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- このissueはv0.1.8自動撮影の統合確認issueとして扱う。実コードの主な修正は [#129](129-add-spout-capture-helper.md) から [#135](135-add-spout-stream-camera-verification-guide.md)、[#138](138-define-player-local-coordinate-spec.md)、[#140](140-implement-player-local-pose-transform.md)、[#141](141-integrate-player-local-coordinate-ui.md)、[#163](163-show-autocapture-test-results-in-settings.md) で行う。
- VRChat公式ドキュメント上、Camera OSCは `/usercamera/Mode` で 2=Stream、`/usercamera/Streaming` でSpout stream toggle、`/usercamera/Pose` でカメラ位置/回転を操作する。現行実装のStream起動方向は維持しつつ、Spout取得、検証、UI結果表示、ローカル座標を完了させる。
- `player_local` は標準OSCだけではローカルプレイヤーrootを自動取得できないため、v0.1.8では「手動保存したプレイヤー基準Pose」を基準にする。実装後のREADME/SPECでは、プレイヤー移動へ自動追従する機能ではないことを明記する。

対象ファイル:

- `src/internal/appcore/autocapture.go`
- `src/internal/appcore/spout.go`
- `src/internal/appcore/player_local.go`
- `src/internal/appcore/config.go`
- `src/app.go`
- `src/frontend/src/main.js`
- `README.md`, `src/SPEC.md`, `RELEASE_NOTES.md`

完了確認:

- Stream方式で標準写真撮影を使わず、Spout helperの画像だけが保存される。
- Photo方式の失敗文言からffmpeg主経路の説明が消える。
- 正面/背後/斜めの初期構図が `player_local` の明示仕様に沿う。
- 失敗結果は履歴に成功画像として混入しない。
- `go test ./...`、`npm run build`、Windows CI、Release zip内容確認、実機確認手順を通す。
