# Spout senderが空のときStreamingをOFF/ONして再確認する

## 問題

実機確認で、ゲーム内カメラとSpout/StreamingがONに見える状態でも、`spout-capture.exe --list-senders` が一時的に `{"ok":true,"senders":[]}` を返すことがある。

この状態のまま撮影へ進むと、helper側のsender待機だけではVRChat側のStreaming状態を揺り戻せず、sender未検出または空フレームで失敗する可能性がある。

## 期待する挙動

Stream方式のSpout取得直前にsender一覧が0件の場合、ClipForVRChatがOSCで `/usercamera/Streaming=false` と `/usercamera/Streaming=true` を一度送ってからsenderを再確認する。

## 受け入れ条件

- [x] Spout取得直前のsender一覧が0件の場合だけ、StreamingのOFF/ON回復を行う。
- [x] `OpenCameraBeforeBatch=false` でも、既にカメラを手動起動している場合のsender空回復として動く。
- [x] helper未検出やsender一覧取得失敗では、回復OSCを送らず既存のSpout取得エラーへ進む。
- [x] 回復前後のsender数とOSC送信結果が診断ログに残る。
- [x] Go testが通る。

## 対応メモ

- 2026-07-07: Spout取得直前に `ListSpoutSenders` を実行し、senderが0件のときだけ `/usercamera/Streaming=false` と `true` を互換int付きで送信する回復処理を追加した。
- 2026-07-07: sender一覧取得に失敗した場合は回復OSCを送らず、従来通り後続のSpout取得で詳細エラーを出す。
- ローカル検証: `go test ./internal/appcore -run 'TestRecoverEmptySpoutSenderList|TestAutoCaptureRunnerRunOnceSkipsCameraAutoOpenWhenDisabled'`、`go test ./...`。
