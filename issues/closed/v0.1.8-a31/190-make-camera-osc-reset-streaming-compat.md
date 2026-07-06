# カメラ撮影機能終了時に一時OSC状態を確実に解除する

## 問題

`カメラOSCをリセット` ボタンは `/usercamera/Streaming=false` をOSC boolで送るが、通常のStream起動/停止処理ではVRChat実機互換性のため bool に加えて numeric `0/1` も送っている。
リセット処理だけ送信形式が異なるため、環境によってSpout stream解除が効きにくい可能性がある。

また、ClipForVRChatからの自動撮影、テスト撮影、Stream撮影などで `/usercamera/Capture` や `/usercamera/Streaming` のような一時的なOSC状態を操作した場合、成功、失敗、キャンセルに関係なく終了時に解除される必要がある。
途中エラーで後始末が漏れると、VRChat側で撮影ボタンやStream状態が残ったように見え、以後のカメラ操作を阻害する可能性がある。

## 期待する挙動

リセットボタン実行時も、通常のStream制御と同じ形式で `/usercamera/Streaming` の解除を送る。
ただし、Zoom、Exposure、Pose、mask類などユーザーのカメラ設定まで一律初期化しない。

ClipForVRChatのカメラ撮影系処理は、成否にかかわらず終了時に押下系/一時有効化系OSCを解除する。
一方で、ユーザーがカメラを開いたまま構図確認したい操作まで無条件に閉じないよう、`/usercamera/Mode=0` や `/usercamera/Close` は処理目的と既存の復元/閉じる設定に従って扱う。

## 受け入れ条件

- [x] `ResetUserCameraOSC` が `/usercamera/Streaming=false` を bool false と numeric `0` の互換送信で送る。
- [x] `/usercamera/Capture=false`、`/usercamera/Close=false`、`/usercamera/Mode=0` のリセット動作は維持する。
- [x] リセットボタンが構図、保存済みPose、manual/avatar_osc basis設定、Zoom/Exposure/mask類を変更しない。
- [x] 自動撮影、テスト撮影、Stream撮影などClipForVRChatがカメラ撮影OSCを送る処理では、成功、失敗、キャンセルにかかわらず `/usercamera/Capture=false` と `/usercamera/Streaming=false` が必要に応じて送られる。
- [x] 撮影終了時の一時OSC解除は、通常のUser Camera状態復元処理と矛盾しない順序で実行される。
- [x] `/usercamera/Mode=0` は「カメラを閉じる」意図がある処理、または復元/fallback設定でOffへ戻す場合に限って送る。
- [x] OSC送信失敗時は従来通り画面と診断ログで確認できる。

## 検証

- [x] `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- [x] Windows向けの `GOOS=windows GOARCH=amd64 go test -c`
