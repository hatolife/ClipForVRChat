# rc30 Release workflowのWindowsテストcleanup失敗を安定化する

## 問題

`v0.1.8-rc30` の GitHub Actions Release workflow が、Windows上の `go test ./...` で失敗した。

失敗箇所は `TestAppRestartCameraPoseReceiverRestartsWhenForwardConfigChanges` の `TempDir RemoveAll cleanup` で、テスト本体のassertではなく、一時ディレクトリ削除時に `The directory is not empty` が出ている。

該当テストはOSC受信器の再起動を確認するために実際のUDP listener goroutineを起動しており、Windows上でテスト終了直後に診断ログ/ディレクトリのハンドル解放が間に合わない可能性がある。

## 期待する挙動

Windows CIでもテスト後の一時ディレクトリcleanupが安定して成功する。

OSC受信器の再起動条件、特にforward設定変更時に受信器が再起動されることは維持する。

## 受け入れ条件

- `TestAppRestartCameraPoseReceiverRestartsWhenForwardConfigChanges` がWindows GitHub Actionsで安定して通る。
- テスト終了前にOSC受信器goroutineとUDP listenerが停止するのを待つ。
- テストのためだけに本体の受信器停止挙動を弱めない。
- ローカル `go test ./...` が通る。
- 次回RCのRelease workflowが少なくとも同じcleanup失敗で落ちない。

## 実装メモ

- `stopCameraPoseReceiverLocked()` はcancel後に戻るが、goroutine側の `ListenUDP` close、defer、診断ログ書き込みの完了までは同期していない。
- `runCameraPoseReceiver` に完了通知用の `done` channel を渡し、`defer` の後始末完了後にcloseする。
- テスト側では `stopCameraPoseReceiverLocked()` 後に `done` channel を待ち、Windowsのファイルハンドル解放を待つ。
- 本体に完了待ち機構を追加する場合は、通常終了パスへ悪影響がない範囲に限定する。

## 検証観点

- [x] `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- [x] Windows向けの `GOOS=windows GOARCH=amd64 go test -c`。
- [x] GitHub Actions Release workflowの再実行、または次RCでの成功確認。

## 検証結果

- `v0.1.8-rc34` のRelease workflowとbranch CIでWindows `go test ./...` が通過した。
- rc30で発生した `TestAppRestartCameraPoseReceiverRestartsWhenForwardConfigChanges` の `TempDir RemoveAll cleanup` 失敗は再発していない。
