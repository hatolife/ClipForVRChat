# `/usercamera/Close` が送信されないようにする

## 指示

> /usercamera/Closeが送信される 現在の仕様では送信されることはない想定 
> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b8-windows-amd64/logs'

## 文脈

`v0.1.8-b8` のログで、カメラOSCリセット時に `/usercamera/Close` が送信されていた。
同じログでは古い設定値由来と見られる `close_after_batch=true` の時間帯にも、バッチ後Closeが送信されていた。

## 解釈

現仕様ではClipForVRChatから `/usercamera/Close` を送らない。リセット処理のボタン解除目的の `Close=false` も送信しない。旧設定ファイルに `closeCameraAfterBatch: true` が残っていても、正規化で無効化する。

## 問題

- `ResetUserCameraOSC` が `/usercamera/Close false` を送っている。
- 旧設定値 `closeCameraAfterBatch=true` が残っていると、バッチ後に `/usercamera/Close` が送られる。

## 期待する挙動

- カメラOSCリセットで `/usercamera/Close` を送信しない。
- 自動撮影バッチ後にも `/usercamera/Close` を送信しない。
- 旧設定値が残っていても `closeCameraAfterBatch` は無効化される。

## 受け入れ条件

- [x] `/usercamera/Close` 送信経路を止める。
- [x] 既存または追加テストで `/usercamera/Close` が送られないことを確認する。

## 確認

- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./internal/appcore -run 'TestResetUserCameraOSC|TestLoadConfigDefaultsAutoLevelRollBeforeShotWhenMissing|TestAutoCaptureRunnerRunOnceSkipsCameraAutoOpenWhenDisabled'`
- `git diff --check`

初回の `go test ./...` は `TestAppWaitForAutoCaptureStartReadinessAutoFallbacksWithoutAvatarOSCBasis` が `52.5ms` でタイミング閾値を超えて失敗したが、再実行で通過した。
