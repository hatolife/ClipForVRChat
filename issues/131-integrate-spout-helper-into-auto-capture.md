# 自動撮影Stream方式をSpoutヘルパーへ統合する

## 問題

現在の `captureStreamShot` は構図適用後に `captureStreamFrameWithFFmpeg` を呼び、VRChatウィンドウ周辺の画面キャプチャを保存している。
この経路ではVRChat Stream CameraのSpout映像そのものを保存できず、白画像やデスクトップ周辺取得が発生する。

## 期待する挙動

自動撮影のStream方式では、OSCでUser CameraをStreamモードにし、Streamingを有効化した後、Spoutヘルパーから実際のStream Cameraフレームを取得して保存する。
保存後は既存のsidecar JSON、Discord投稿、履歴処理へ同じように流れる。

## 受け入れ条件

- Stream方式開始時に `/usercamera/Mode=2` と `/usercamera/Streaming=true` を送る既存制御を維持する。
- 構図ごとにPose、Zoom、表示対象などのOSC設定を送った後、Spoutヘルパーで1フレームを取得する。
- Stream方式ではVRChat標準の `/usercamera/Capture` を使わず、シャッター音を発生させない。
- Spout取得成功時は既存の `finalizeAutoCaptureImage` 相当の経路でsidecar JSON、Discord投稿、履歴へ連携する。
- Spout取得失敗時は成功履歴やDiscord投稿を作らず、画面上の自動撮影エラーと診断ログへ詳細を出す。
- context cancellationとタイムアウトでヘルパープロセスを確実に終了できる。
- 診断ログにhelper path、sender選択方式、sender名、出力先、タイムアウト、処理時間、解像度を記録する。
- Photo方式と通常画像処理の挙動を変えない。

## 実装メモ

- FFmpeg経路は主経路から外す。
- 互換目的で残す場合は、明示的なlegacy/debug fallbackとして扱い、Stream Camera直接取得の代替とは表示しない。
- helperのstdout JSONをGo側で構造体として解析し、画像メタデータへ渡す。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- 現行のStream経路は `RunOnce()` で `/usercamera/Mode=2`、`/usercamera/Streaming=true` を送った後、`captureStreamFrameWithSpout()` を呼ぶ構造になっているため、この経路を維持する。
- `captureStreamFrameWithSpout()` に処理時間計測を追加し、開始/成功/失敗/timeoutすべてで `duration_ms` を診断ログに残す。
- helperのstdout JSONを構造体で解析し、sender情報、フレーム番号、取得時刻、候補一覧をsidecar/ログへ渡す。
- context cancellation時は `exec.CommandContext` でプロセス終了させる現行設計を維持し、timeout時のログにhelper path、args、output path、sender、自動選択、durationを入れる。

対象ファイル:

- `src/internal/appcore/autocapture.go`
- `src/internal/appcore/spout.go`
- `src/internal/appcore/autocapture_test.go`

小タスク:

- `captureStreamFrameWithSpout()` の開始直後に `started := time.Now()` を置く。
- helper実行失敗時に `trimmed stdout/stderr` と `duration_ms` をログへ出す。
- 成功時に `duration_ms`、`result.OutputPath`、`senderName`、`width/height`、`frame` をログへ出す。
- Stream方式のDiscord本文に `撮影方式: stream/spout` とsender名を含める処理は [#132](132-validate-spout-capture-output-and-metadata.md) 側で実装する。

確認方法:

- helper pathを存在しない値にして、画面エラーと診断ログにhelper path/原因が出る。
- timeoutを短くして、プロセスが残らずtimeout理由が出る。
- Stream方式成功時にPhotoの `/usercamera/Capture` が呼ばれない。
