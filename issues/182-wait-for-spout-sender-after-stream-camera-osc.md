# Stream Camera起動直後にSpout senderが出る前に失敗する

## 問題

`v0.1.8-a26` のStream方式テスト撮影で、アプリは `/usercamera/Mode=2`、`/usercamera/SmoothMovement=true`、`/usercamera/Streaming=true` を送信しているが、Spout helperが起動直後のsender一覧だけを見て `Spout senderがありません` と即失敗する。

実機ログでは `Streaming=true` 送信から約3秒後に helper を起動しているが、helper側は `--timeout-ms 10000` をフレーム取得待ちにしか使っておらず、sender出現待ちには使っていない。そのためVRChat側のStream Camera/Spout sender生成が遅い場合に、カメラ自動起動を送っていても失敗する。

## 期待する挙動

Stream方式では、OSCでStream CameraとSpoutをONにしたあと、Spout senderが出るまで一定時間待ってからフレーム取得する。

また、sender待機中にVRChat側の起動漏れを補うため、必要に応じて `/usercamera/Mode=2` と `/usercamera/Streaming=true` を再送できる。

## 受け入れ条件

- [x] Spout helperがsender 0件の場合に即失敗せず、timeout内でsender一覧を再試行する。
- [x] 指定sender名がある場合も、timeout内で指定sender出現を待つ。
- [x] アプリ側のStream起動待機で、sender待機/再送の状況が診断ログに出る。
- [x] timeout後だけ `Spout senderがありません` を返す。
- [x] 既存の複数sender曖昧エラーは不要に長く待たず、候補を返す。
- [x] Go test、frontend check、Release/CIでSpout helperビルドが通る。

## 実装メモ

- `tools/spout-capture/main.cpp` の `capture()` で `choose_sender()` が `sender_not_found` の場合に、`options.timeout_ms` まで短いsleepで再試行する。
- 複数senderで自動選択不能な `sender_ambiguous` は、待って解決するとは限らないため即エラーでよい。
- Go側は必要に応じて `captureStreamShot()` のSpout取得直前に `/usercamera/Mode=2` と `/usercamera/Streaming=true` を再送し、ログへ出す。

## 対応メモ

- 2026-07-04: `spout-capture.exe --capture` で、sender未検出時に `--timeout-ms` 内で100msごとにsender選択を再試行するようにした。
- 2026-07-04: senderが出た後のフレーム取得も同じdeadline内で行い、timeout後にだけsenderなし/フレームなしのエラーを返す。
- 2026-07-04: Go側のhelperプロセスtimeoutは、helperがJSONエラーを返す余裕を持たせるため、capture timeout + 5秒にした。
- 2026-07-04: Stream方式の各shotでSpout helper起動直前に `/usercamera/Mode=2`、`/usercamera/SmoothMovement=true`、`/usercamera/Streaming=true` を再送し、500ms待つ。
- ローカル検証: `go test ./...`、frontend build、template literal check、Wails API surface check、`tools/spout-capture/main.cpp` object compile。
