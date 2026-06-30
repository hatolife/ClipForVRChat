# 自動撮影OSCの押下状態を解除するデバッグボタンを追加する

## 問題

自動撮影後にゲーム内でカメラを開けなくなる場合がある。
`/usercamera/Close` や `/usercamera/Streaming` などのOSC操作が押下/有効状態として残り、VRChat側でカメラ操作を阻害している可能性がある。

## 期待する挙動

- 設定画面からUser Camera関連OSCを安全側へ戻すデバッグ操作を実行できる。
- Close/Capture/Streamingなどをfalseへ戻し、Camera ModeもOffへ戻せる。
- 操作結果は画面と診断ログで確認できる。

## 受け入れ条件

- 自動撮影タブに `カメラOSCをリセット` ボタンがある。
- ボタン実行時に `/usercamera/Capture=false`、`/usercamera/Close=false`、`/usercamera/Streaming=false`、`/usercamera/Mode=0` を送る。
- OSC送信失敗時は画面にエラーを表示する。
- 既存の自動撮影、テスト撮影、Pose移動を壊さない。

## 対応内容

- 自動撮影タブに `カメラOSCをリセット` ボタンを追加した。
- `ResetCameraOSC` APIを追加し、`/usercamera/Capture=false`、`/usercamera/Close=false`、`/usercamera/Streaming=false`、`/usercamera/Mode=0` を送るようにした。
- 実行結果を画面メッセージと診断ログへ出すようにした。
