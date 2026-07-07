# b10で3枚撮影想定なのに6枚投稿される

## 指示

> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b10-windows-amd64/history.json''/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b10-windows-amd64/logs' 3枚とる想定で6枚投稿される

## 文脈

`v0.1.8-b10` の実行結果で、3枚撮影される想定の自動撮影がDiscordへ6枚投稿された。ユーザーから実行環境の `history.json` と `logs/` が提示された。

## 解釈

自動撮影バッチ内で撮影処理またはDiscord投稿処理が重複している可能性がある。提示された履歴とログから、撮影枚数自体が6枚なのか、3枚の撮影結果が二重投稿されたのかを切り分け、原因を修正する。

## 問題

- 3枚撮影想定に対してDiscord投稿が6枚発生する。
- 重複が撮影側、履歴登録側、Discord投稿側のどこで起きているか未確定。

## 期待する挙動

- 3構図の自動撮影では、Discord投稿も3件だけ発生する。
- 同じ自動撮影結果が重複投稿されない。

## 受け入れ条件

- [x] 提示された `history.json` と `logs/` から重複の発生箇所を特定する。
- [x] 自動撮影のDiscord投稿が重複しないよう修正する。
- [x] 関連テストを追加または更新する。
- [x] 関連チェックが通る。

## 調査

- `history.json` には6件あり、3件は `outputPath` が `VRC-AutoCapture` の元画像、残り3件は通常処理の `output/*_2048.png` だった。
- 通常ログでは自動撮影バッチ `batch-20260707-155241` の成功は3件で、各画像は自動撮影側でDiscord投稿済みだった。
- `autoPhoto.photoDirectory` が `%USERPROFILE%\Pictures\VRChat`、`autoCapture.output.directory` がその配下の `%USERPROFILE%\Pictures\VRChat\VRC-AutoCapture` だったため、VRChat写真自動投稿が自動撮影の出力画像を新規写真として拾い、同じ画像を再投稿していた。

## 確認

- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./internal/appcore -run 'TestAutoPhotoWatcher|TestScanPhotoFiles'`
- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `git diff --check`
