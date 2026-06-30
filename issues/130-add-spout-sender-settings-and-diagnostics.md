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
