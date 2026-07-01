# Spout sender設定と診断UIを追加する

## 問題

VRChat Stream CameraのSpout sender名は環境やVRChat側状態によって変わる可能性がある。
senderを固定名で決め打ちすると、Stream Cameraが出ているのに取得できない、または別アプリのSpout senderを誤取得する可能性がある。

## 期待する挙動

設定画面の「自動撮影」タブからSpout受信ヘルパーの状態確認、sender一覧取得、sender選択/自動選択設定ができる。
取得できない場合は、VRChat側でStream Cameraが起動していないのか、Spout senderがないのか、複数senderで選択が必要なのかを区別して表示できる。

## 受け入れ条件

- 自動撮影タブにSpout helper確認ボタンを追加し、ヘルパーexeの検出結果とバージョンまたは実行可否を表示する。
- 自動撮影タブにSpout sender一覧更新ボタンを追加し、`--list-senders` の結果を表示できる。
- Stream方式の設定として、sender自動選択ON/OFF、sender名、取得タイムアウト、起動後待機時間を保存できる。
- senderが0件の場合、VRChatでStream Cameraを起動し `/usercamera/Streaming=true` が必要であることを表示する。
- senderが複数件あり自動選択できない場合、候補一覧と選択を促す文を表示する。
- 設定保存/読み込みで既存configを壊さず、未設定時は自動選択ONの妥当な初期値になる。
- 既存のPhoto方式、Presence、Discord、通常投稿設定を壊さない。

## 実装メモ

- configにはFFmpeg前提の `ffmpegPath` / `inputArgs` とは別に、Spout用の設定フィールドを追加する。
- FFmpeg設定を残す場合も、Stream Camera直接取得の主設定として見せない。
- UI文言は「Stream Camera(Spout)」を基準にし、画面キャプチャと誤解される説明を避ける。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- `AutoCaptureStreamConfig.StartDelayMS` はconfigと撮影処理に存在するが、設定画面にないため自動撮影タブへ追加する。
- helper確認は `--version` と `--list-senders` の結果を表示し、sender 0件、複数件、helper未検出を区別する。
- Spout helperの失敗JSONに `senders` が含まれるようにしたうえで、Go側の `SpoutHelperStatus.Senders` とfrontendの `spoutStatus.senders` へ流す。

対象ファイル:

- `src/internal/appcore/config.go`
- `src/internal/appcore/spout.go`
- `src/app.go`
- `src/frontend/src/main.js`

小タスク:

- Stream設定UIに「Stream起動後待機」数値入力を追加し、`autoCaptureSettings.stream.startDelayMs` へbindする。
- `CheckSpoutHelper()` のメッセージを `helper実行可 / sender 0件 / 複数候補 / sender選択必要` に分ける。
- sender一覧表示にhostPathを追加し、候補が複数の場合は自動選択OFFとsender名選択を促す。
- UIのdisabled条件を `capture.mode !== 'stream'` とSpout操作中に揃える。

確認方法:

- 既存configに `startDelayMs` がない場合は初期値1000msへNormalizeされる。
- 設定画面で変更した `startDelayMs` が保存後の `config.json` に反映される。
- `npm run build` でfrontend buildが通る。
