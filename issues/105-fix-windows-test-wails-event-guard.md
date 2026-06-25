# Windows GoテストでWailsイベント送信ガードが効かない

## 問題

master CIのWindows `go test ./...` で、Wails runtimeの `EventsEmit` が通常のテストcontextに対して呼ばれ、テストが失敗する。

WindowsではGoのテストバイナリ名が `.test.exe` になる場合があり、既存の `.test` suffix判定だけではテスト実行中と判定できない。

## 期待する挙動

Windows CIのGoテストではWails runtimeへイベント送信せず、GUI実行時のみWails lifecycle contextへイベントを送信する。

## 受け入れ条件

- Windows CIの `go test ./...` が成功する。
- GUI実行時の `process:progress` / `auto-photo:result` イベント送信は維持される。
- master CIが成功する。
